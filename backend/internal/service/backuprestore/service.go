package backuprestore

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	executorProvider "kbmanage/backend/internal/integration/backuprestore/executor"
	validatorProvider "kbmanage/backend/internal/integration/backuprestore/validator"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"

	"gorm.io/gorm"
)

const (
	ResourceTypeBackupPolicy = "backup-policy"
	ResourceTypeRestorePoint = "restore-point"
	ResourceTypeRestoreJob   = "restore-job"
	ResourceTypeMigration    = "migration-plan"
	ResourceTypeDrillPlan    = "drill-plan"
	ResourceTypeDrillRecord  = "drill-record"
	ResourceTypeDrillReport  = "drill-report"

	ActionPolicyCreate    = "backuprestore.policy.create"
	ActionPolicyRun       = "backuprestore.policy.run"
	ActionRestoreCreate   = "backuprestore.restore.create"
	ActionRestoreCheck    = "backuprestore.restore.validate"
	ActionMigrationCreate = "backuprestore.migration.create"
	ActionDrillPlanCreate = "backuprestore.drill.plan.create"
	ActionDrillRun        = "backuprestore.drill.run"
	ActionDrillReportGen  = "backuprestore.drill.report.generate"
)

var (
	ErrBackupRestoreScopeDenied = errors.New("backup restore scope access denied")
	ErrBackupRestoreConflict    = errors.New("backup restore operation conflict")
	ErrBackupRestoreInvalid     = errors.New("backup restore invalid request")
	ErrBackupRestoreBlocked     = errors.New("backup restore operation blocked")
)

type Service struct {
	policies      *repository.BackupPolicyRepository
	restorePoints *repository.RestorePointRepository
	restoreJobs   *repository.RestoreJobRepository
	migrations    *repository.MigrationPlanRepository
	drillPlans    *repository.DRDrillPlanRepository
	drillRecords  *repository.DRDrillRecordRepository
	drillReports  *repository.DRDrillReportRepository
	auditRepo     *repository.BackupAuditRepository
	scope         *ScopeService
	progress      *ProgressCache
	prechecks     *PrecheckCache
	lock          *OperationLock
	executor      executorProvider.Provider
	validator     validatorProvider.Provider
	auditWriter   *auditSvc.EventWriter
}

func NewService(
	policyRepo *repository.BackupPolicyRepository,
	restorePointRepo *repository.RestorePointRepository,
	restoreJobRepo *repository.RestoreJobRepository,
	migrationRepo *repository.MigrationPlanRepository,
	drillPlanRepo *repository.DRDrillPlanRepository,
	drillRecordRepo *repository.DRDrillRecordRepository,
	drillReportRepo *repository.DRDrillReportRepository,
	auditRepo *repository.BackupAuditRepository,
	bindingRepo *repository.ScopeRoleBindingRepository,
	projectRepo *repository.ProjectRepository,
	progressCache *ProgressCache,
	precheckCache *PrecheckCache,
	lock *OperationLock,
	executor executorProvider.Provider,
	validator validatorProvider.Provider,
	auditWriter *auditSvc.EventWriter,
) *Service {
	return &Service{
		policies:      policyRepo,
		restorePoints: restorePointRepo,
		restoreJobs:   restoreJobRepo,
		migrations:    migrationRepo,
		drillPlans:    drillPlanRepo,
		drillRecords:  drillRecordRepo,
		drillReports:  drillReportRepo,
		auditRepo:     auditRepo,
		scope:         NewScopeService(bindingRepo, projectRepo),
		progress:      progressCache,
		prechecks:     precheckCache,
		lock:          lock,
		executor:      executor,
		validator:     validator,
		auditWriter:   auditWriter,
	}
}

type CreateBackupPolicyInput struct {
	Name               string `json:"name"`
	Description        string `json:"description"`
	ScopeType          string `json:"scopeType"`
	ScopeRef           string `json:"scopeRef"`
	WorkspaceID        uint64 `json:"workspaceId"`
	ProjectID          uint64 `json:"projectId"`
	ExecutionMode      string `json:"executionMode"`
	ScheduleExpression string `json:"scheduleExpression"`
	RetentionRule      string `json:"retentionRule"`
	ConsistencyLevel   string `json:"consistencyLevel"`
	Status             string `json:"status"`
}

type PolicyListFilter struct {
	ScopeType string
	Status    string
}

type RestorePointListFilter struct {
	PolicyID uint64
	Result   string
}

type CreateRestoreJobInput struct {
	RestorePointID    uint64         `json:"restorePointId"`
	JobType           string         `json:"jobType"`
	SourceEnvironment string         `json:"sourceEnvironment"`
	TargetEnvironment string         `json:"targetEnvironment"`
	ScopeSelection    map[string]any `json:"scopeSelection"`
}

type RestoreJobListFilter struct {
	JobType string
	Status  string
}

type PrecheckResult struct {
	Status            string   `json:"status"`
	Blockers          []string `json:"blockers"`
	Warnings          []string `json:"warnings"`
	ConsistencyNotice string   `json:"consistencyNotice"`
}

type CreateMigrationPlanInput struct {
	Name            string         `json:"name"`
	WorkspaceID     uint64         `json:"workspaceId"`
	ProjectID       uint64         `json:"projectId"`
	SourceClusterID uint64         `json:"sourceClusterId"`
	TargetClusterID uint64         `json:"targetClusterId"`
	ScopeSelection  map[string]any `json:"scopeSelection"`
	MappingRules    map[string]any `json:"mappingRules"`
	CutoverSteps    []string       `json:"cutoverSteps"`
}

type CreateDRDrillPlanInput struct {
	Name                string         `json:"name"`
	Description         string         `json:"description"`
	WorkspaceID         uint64         `json:"workspaceId"`
	ProjectID           uint64         `json:"projectId"`
	ScopeSelection      map[string]any `json:"scopeSelection"`
	RPOTargetMinutes    int            `json:"rpoTargetMinutes"`
	RTOTargetMinutes    int            `json:"rtoTargetMinutes"`
	RoleAssignments     []string       `json:"roleAssignments"`
	CutoverProcedure    []string       `json:"cutoverProcedure"`
	ValidationChecklist []string       `json:"validationChecklist"`
}

func (s *Service) ListPolicies(ctx context.Context, userID uint64, filter PolicyListFilter) ([]domain.BackupPolicy, error) {
	workspaceIDs, constrained, err := s.scope.ListReadableWorkspaceIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if constrained && len(workspaceIDs) == 0 {
		return []domain.BackupPolicy{}, nil
	}
	return s.policies.List(ctx, repository.BackupPolicyListFilter{
		WorkspaceIDs: workspaceIDs,
		ScopeType:    filter.ScopeType,
		Status:       filter.Status,
	})
}

func (s *Service) CreatePolicy(ctx context.Context, userID uint64, input CreateBackupPolicyInput) (*domain.BackupPolicy, error) {
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.ScopeType) == "" || strings.TrimSpace(input.RetentionRule) == "" || strings.TrimSpace(input.ConsistencyLevel) == "" {
		return nil, ErrBackupRestoreInvalid
	}
	if strings.EqualFold(strings.TrimSpace(input.ExecutionMode), "scheduled") && strings.TrimSpace(input.ScheduleExpression) == "" {
		return nil, ErrBackupRestoreInvalid
	}
	workspaceID, projectID := normalizeScopeIDs(input.WorkspaceID, input.ProjectID)
	if workspaceID == 0 {
		var err error
		workspaceID, projectID, err = s.defaultScopeForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	if err := s.scope.EnsureManagePolicy(ctx, userID, workspaceID, projectID); err != nil {
		return nil, err
	}
	if existing, err := s.policies.FindByScopeName(ctx, workspaceID, input.ScopeType, input.ScopeRef, input.Name); err == nil && existing != nil {
		return nil, ErrBackupRestoreConflict
	} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	status := strings.TrimSpace(input.Status)
	if status == "" {
		status = string(domain.BackupPolicyStatusActive)
	}
	item := &domain.BackupPolicy{
		Name:               strings.TrimSpace(input.Name),
		Description:        strings.TrimSpace(input.Description),
		ScopeType:          strings.TrimSpace(input.ScopeType),
		ScopeRef:           strings.TrimSpace(input.ScopeRef),
		WorkspaceID:        workspaceID,
		ProjectID:          uint64PtrIf(projectID),
		ExecutionMode:      strings.TrimSpace(input.ExecutionMode),
		ScheduleExpression: strings.TrimSpace(input.ScheduleExpression),
		RetentionRule:      strings.TrimSpace(input.RetentionRule),
		ConsistencyLevel:   strings.TrimSpace(input.ConsistencyLevel),
		Status:             domain.BackupPolicyStatus(status),
		OwnerUserID:        userID,
	}
	if err := s.policies.Create(ctx, item); err != nil {
		return nil, err
	}
	_ = s.writeAudit(ctx, userID, ActionPolicyCreate, ResourceTypeBackupPolicy, item.ID, workspaceID, projectID, domain.BackupAuditOutcomeSucceeded, map[string]any{
		"name":      item.Name,
		"scopeType": item.ScopeType,
		"scopeRef":  item.ScopeRef,
	})
	return item, nil
}

func (s *Service) RunPolicy(ctx context.Context, userID, policyID uint64) (*domain.RestorePoint, error) {
	policy, err := s.policies.GetByID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureBackup(ctx, userID, policy.WorkspaceID, derefUint64(policy.ProjectID)); err != nil {
		return nil, err
	}
	ok, err := s.lock.Acquire(ctx, "policy", policyID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, ErrBackupRestoreConflict
	}
	defer func() { _ = s.lock.Release(ctx, "policy", policyID) }()

	startedAt := time.Now()
	result, err := s.executor.RunBackup(ctx, executorProvider.BackupExecutionRequest{
		PolicyID:         policy.ID,
		PolicyName:       policy.Name,
		ScopeType:        policy.ScopeType,
		ScopeRef:         policy.ScopeRef,
		ConsistencyLevel: policy.ConsistencyLevel,
	})
	if err != nil {
		return nil, err
	}
	completedAt := time.Now()
	scopeSnapshot := mustJSON(map[string]any{
		"scopeType": policy.ScopeType,
		"scopeRef":  policy.ScopeRef,
	})
	point := &domain.RestorePoint{
		PolicyID:           policy.ID,
		WorkspaceID:        policy.WorkspaceID,
		ProjectID:          policy.ProjectID,
		ScopeSnapshot:      scopeSnapshot,
		BackupStartedAt:    startedAt,
		BackupCompletedAt:  &completedAt,
		DurationSeconds:    durationFrom(startedAt, completedAt, result.DurationSeconds),
		Result:             domain.RestorePointResult(result.Result),
		ConsistencySummary: strings.TrimSpace(result.ConsistencySummary),
		FailureReason:      strings.TrimSpace(result.FailureReason),
		StorageRef:         strings.TrimSpace(result.StorageRef),
		ExpiresAt:          result.ExpiresAt,
		CreatedBy:          userID,
	}
	if err := s.restorePoints.Create(ctx, point); err != nil {
		return nil, err
	}
	_ = s.progress.Set(ctx, "policy", policy.ID, "backup", string(point.Result))
	_ = s.writeAudit(ctx, userID, ActionPolicyRun, ResourceTypeRestorePoint, point.ID, policy.WorkspaceID, derefUint64(policy.ProjectID), outcomeForResult(string(point.Result)), map[string]any{
		"policyId":    policy.ID,
		"storageRef":  point.StorageRef,
		"result":      point.Result,
		"consistency": point.ConsistencySummary,
	})
	return point, nil
}

func (s *Service) ListRestorePoints(ctx context.Context, userID uint64, filter RestorePointListFilter) ([]domain.RestorePoint, error) {
	workspaceIDs, constrained, err := s.scope.ListReadableWorkspaceIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if constrained && len(workspaceIDs) == 0 {
		return []domain.RestorePoint{}, nil
	}
	return s.restorePoints.List(ctx, repository.RestorePointListFilter{
		WorkspaceIDs: workspaceIDs,
		PolicyID:     filter.PolicyID,
		Result:       filter.Result,
	})
}

func (s *Service) GetRestorePoint(ctx context.Context, userID, restorePointID uint64) (*domain.RestorePoint, error) {
	item, err := s.restorePoints.GetByID(ctx, restorePointID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureRead(ctx, userID, item.WorkspaceID, derefUint64(item.ProjectID)); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) CreateRestoreJob(ctx context.Context, userID uint64, input CreateRestoreJobInput) (*domain.RestoreJob, error) {
	point, err := s.restorePoints.GetByID(ctx, input.RestorePointID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureRestore(ctx, userID, point.WorkspaceID, derefUint64(point.ProjectID)); err != nil {
		return nil, err
	}
	if point.Result == domain.RestorePointResultFailed || point.Result == domain.RestorePointResultExpired {
		return nil, ErrBackupRestoreBlocked
	}
	job := &domain.RestoreJob{
		RestorePointID:    point.ID,
		WorkspaceID:       point.WorkspaceID,
		ProjectID:         point.ProjectID,
		JobType:           domain.RestoreJobType(strings.TrimSpace(input.JobType)),
		SourceEnvironment: strings.TrimSpace(input.SourceEnvironment),
		TargetEnvironment: strings.TrimSpace(input.TargetEnvironment),
		ScopeSelection:    mustJSON(input.ScopeSelection),
		Status:            domain.RestoreJobStatusPending,
		RequestedBy:       userID,
	}
	if job.TargetEnvironment == "" || len(input.ScopeSelection) == 0 {
		return nil, ErrBackupRestoreInvalid
	}
	precheck, err := s.ValidateRestoreJob(ctx, userID, input)
	if err != nil {
		return nil, err
	}
	job.ConflictSummary = strings.Join(precheck.Blockers, "; ")
	job.ConsistencyNotice = precheck.ConsistencyNotice
	if len(precheck.Blockers) > 0 {
		job.Status = domain.RestoreJobStatusBlocked
	}
	if err := s.restoreJobs.Create(ctx, job); err != nil {
		return nil, err
	}
	if len(precheck.Blockers) == 0 {
		now := time.Now()
		job.StartedAt = &now
		job.Status = domain.RestoreJobStatusRunning
		execResult, execErr := s.executor.RunRestore(ctx, executorProvider.RestoreExecutionRequest{
			JobID:             job.ID,
			JobType:           string(job.JobType),
			RestorePointID:    point.ID,
			TargetEnvironment: job.TargetEnvironment,
		})
		completedAt := time.Now()
		job.CompletedAt = &completedAt
		if execErr != nil {
			job.Status = domain.RestoreJobStatusFailed
			job.FailureReason = execErr.Error()
		} else {
			job.Status = domain.RestoreJobStatus(execResult.Status)
			job.ResultSummary = execResult.ResultSummary
			if job.ConflictSummary == "" {
				job.ConflictSummary = execResult.ConflictSummary
			}
			if job.ConsistencyNotice == "" {
				job.ConsistencyNotice = execResult.ConsistencyNotice
			}
			job.FailureReason = execResult.FailureReason
		}
		if err := s.restoreJobs.Update(ctx, job); err != nil {
			return nil, err
		}
	}
	_ = s.progress.Set(ctx, "restore-job", job.ID, string(job.JobType), string(job.Status))
	_ = s.writeAudit(ctx, userID, ActionRestoreCreate, ResourceTypeRestoreJob, job.ID, job.WorkspaceID, derefUint64(job.ProjectID), outcomeForResult(string(job.Status)), map[string]any{
		"restorePointId": job.RestorePointID,
		"jobType":        job.JobType,
		"targetEnv":      job.TargetEnvironment,
	})
	return job, nil
}

func (s *Service) ListRestoreJobs(ctx context.Context, userID uint64, filter RestoreJobListFilter) ([]domain.RestoreJob, error) {
	workspaceIDs, constrained, err := s.scope.ListReadableWorkspaceIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if constrained && len(workspaceIDs) == 0 {
		return []domain.RestoreJob{}, nil
	}
	return s.restoreJobs.List(ctx, repository.RestoreJobListFilter{
		WorkspaceIDs: workspaceIDs,
		JobType:      filter.JobType,
		Status:       filter.Status,
	})
}

func (s *Service) ValidateRestoreJob(ctx context.Context, userID uint64, input CreateRestoreJobInput) (*PrecheckResult, error) {
	point, err := s.restorePoints.GetByID(ctx, input.RestorePointID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureRestore(ctx, userID, point.WorkspaceID, derefUint64(point.ProjectID)); err != nil {
		return nil, err
	}
	result, err := s.validator.Validate(ctx, validatorProvider.Request{
		JobType:           strings.TrimSpace(input.JobType),
		ScopeSelection:    input.ScopeSelection,
		TargetEnvironment: strings.TrimSpace(input.TargetEnvironment),
	})
	if err != nil {
		return nil, err
	}
	view := &PrecheckResult{
		Status:            result.Status,
		Blockers:          result.Blockers,
		Warnings:          result.Warnings,
		ConsistencyNotice: result.ConsistencyNotice,
	}
	_ = s.prechecks.Store(ctx, input.RestorePointID, view)
	_ = s.writeAudit(ctx, userID, ActionRestoreCheck, ResourceTypeRestorePoint, point.ID, point.WorkspaceID, derefUint64(point.ProjectID), outcomeForValidation(view.Status), map[string]any{
		"jobType":           input.JobType,
		"targetEnvironment": input.TargetEnvironment,
	})
	return view, nil
}

func (s *Service) ValidateRestoreJobByID(ctx context.Context, userID, jobID uint64) (*PrecheckResult, error) {
	job, err := s.restoreJobs.GetByID(ctx, jobID)
	if err != nil {
		return nil, err
	}
	var selection map[string]any
	_ = json.Unmarshal([]byte(job.ScopeSelection), &selection)
	return s.ValidateRestoreJob(ctx, userID, CreateRestoreJobInput{
		RestorePointID:    job.RestorePointID,
		JobType:           string(job.JobType),
		SourceEnvironment: job.SourceEnvironment,
		TargetEnvironment: job.TargetEnvironment,
		ScopeSelection:    selection,
	})
}

func (s *Service) CreateMigrationPlan(ctx context.Context, userID uint64, input CreateMigrationPlanInput) (*domain.MigrationPlan, error) {
	workspaceID, projectID := normalizeScopeIDs(input.WorkspaceID, input.ProjectID)
	if input.SourceClusterID == 0 || input.TargetClusterID == 0 || input.SourceClusterID == input.TargetClusterID || strings.TrimSpace(input.Name) == "" {
		return nil, ErrBackupRestoreInvalid
	}
	if workspaceID == 0 {
		var err error
		workspaceID, projectID, err = s.defaultScopeForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	if err := s.scope.EnsureMigrate(ctx, userID, workspaceID, projectID); err != nil {
		return nil, err
	}
	item := &domain.MigrationPlan{
		Name:            strings.TrimSpace(input.Name),
		WorkspaceID:     workspaceID,
		ProjectID:       uint64PtrIf(projectID),
		SourceClusterID: input.SourceClusterID,
		TargetClusterID: input.TargetClusterID,
		ScopeSelection:  mustJSON(input.ScopeSelection),
		MappingRules:    mustJSON(input.MappingRules),
		CutoverSteps:    mustJSON(input.CutoverSteps),
		Status:          domain.MigrationPlanStatusDraft,
		CreatedBy:       userID,
	}
	if err := s.migrations.Create(ctx, item); err != nil {
		return nil, err
	}
	if _, err := s.executor.RunMigration(ctx, executorProvider.MigrationExecutionRequest{
		PlanID:          item.ID,
		Name:            item.Name,
		SourceClusterID: item.SourceClusterID,
		TargetClusterID: item.TargetClusterID,
	}); err == nil {
		item.Status = domain.MigrationPlanStatusSucceeded
		_ = s.migrations.Update(ctx, item)
	}
	_ = s.writeAudit(ctx, userID, ActionMigrationCreate, ResourceTypeMigration, item.ID, workspaceID, projectID, outcomeForResult(string(item.Status)), map[string]any{
		"name":            item.Name,
		"sourceClusterId": item.SourceClusterID,
		"targetClusterId": item.TargetClusterID,
	})
	return item, nil
}

func (s *Service) ListDrillPlans(ctx context.Context, userID uint64) ([]domain.DRDrillPlan, error) {
	workspaceIDs, constrained, err := s.scope.ListReadableWorkspaceIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if constrained && len(workspaceIDs) == 0 {
		return []domain.DRDrillPlan{}, nil
	}
	return s.drillPlans.List(ctx, workspaceIDs)
}

func (s *Service) CreateDrillPlan(ctx context.Context, userID uint64, input CreateDRDrillPlanInput) (*domain.DRDrillPlan, error) {
	workspaceID, projectID := normalizeScopeIDs(input.WorkspaceID, input.ProjectID)
	if strings.TrimSpace(input.Name) == "" || input.RPOTargetMinutes <= 0 || input.RTOTargetMinutes <= 0 || len(input.CutoverProcedure) == 0 || len(input.ValidationChecklist) == 0 {
		return nil, ErrBackupRestoreInvalid
	}
	if workspaceID == 0 {
		var err error
		workspaceID, projectID, err = s.defaultScopeForUser(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	if err := s.scope.EnsureDrill(ctx, userID, workspaceID, projectID); err != nil {
		return nil, err
	}
	item := &domain.DRDrillPlan{
		Name:                strings.TrimSpace(input.Name),
		Description:         strings.TrimSpace(input.Description),
		WorkspaceID:         workspaceID,
		ProjectID:           uint64PtrIf(projectID),
		ScopeSelection:      mustJSON(input.ScopeSelection),
		RPOTargetMinutes:    input.RPOTargetMinutes,
		RTOTargetMinutes:    input.RTOTargetMinutes,
		RoleAssignments:     mustJSON(input.RoleAssignments),
		CutoverProcedure:    mustJSON(input.CutoverProcedure),
		ValidationChecklist: mustJSON(input.ValidationChecklist),
		Status:              domain.DRDrillPlanStatusActive,
		CreatedBy:           userID,
	}
	if err := s.drillPlans.Create(ctx, item); err != nil {
		return nil, err
	}
	_ = s.writeAudit(ctx, userID, ActionDrillPlanCreate, ResourceTypeDrillPlan, item.ID, workspaceID, projectID, domain.BackupAuditOutcomeSucceeded, map[string]any{"name": item.Name})
	return item, nil
}

func (s *Service) RunDrillPlan(ctx context.Context, userID, planID uint64) (*domain.DRDrillRecord, error) {
	plan, err := s.drillPlans.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureDrill(ctx, userID, plan.WorkspaceID, derefUint64(plan.ProjectID)); err != nil {
		return nil, err
	}
	result, err := s.executor.RunDrill(ctx, executorProvider.DrillExecutionRequest{PlanID: plan.ID, PlanName: plan.Name})
	if err != nil {
		return nil, err
	}
	started := time.Now().Add(-time.Duration(result.ActualRTOMinutes) * time.Minute)
	completed := time.Now()
	record := &domain.DRDrillRecord{
		PlanID:            plan.ID,
		WorkspaceID:       plan.WorkspaceID,
		ProjectID:         plan.ProjectID,
		StartedAt:         started,
		CompletedAt:       &completed,
		ActualRPOMinutes:  result.ActualRPOMinutes,
		ActualRTOMinutes:  result.ActualRTOMinutes,
		Status:            domain.DRDrillRecordStatus(result.Status),
		StepResults:       mustJSON(result.StepResults),
		ValidationResults: mustJSON(result.ValidationResults),
		IncidentNotes:     result.IncidentNotes,
		ExecutedBy:        userID,
	}
	if err := s.drillRecords.Create(ctx, record); err != nil {
		return nil, err
	}
	_ = s.writeAudit(ctx, userID, ActionDrillRun, ResourceTypeDrillRecord, record.ID, plan.WorkspaceID, derefUint64(plan.ProjectID), outcomeForResult(string(record.Status)), map[string]any{
		"planId": plan.ID,
	})
	return record, nil
}

func (s *Service) GetDrillRecord(ctx context.Context, userID, recordID uint64) (*domain.DRDrillRecord, error) {
	record, err := s.drillRecords.GetByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureRead(ctx, userID, record.WorkspaceID, derefUint64(record.ProjectID)); err != nil {
		return nil, err
	}
	return record, nil
}

func (s *Service) GenerateDrillReport(ctx context.Context, userID, recordID uint64) (*domain.DRDrillReport, error) {
	record, err := s.drillRecords.GetByID(ctx, recordID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureDrill(ctx, userID, record.WorkspaceID, derefUint64(record.ProjectID)); err != nil {
		return nil, err
	}
	plan, err := s.drillPlans.GetByID(ctx, record.PlanID)
	if err != nil {
		return nil, err
	}
	goal, gap := AssessRPORTO(plan.RPOTargetMinutes, plan.RTOTargetMinutes, record.ActualRPOMinutes, record.ActualRTOMinutes)
	report := &domain.DRDrillReport{
		DrillRecordID:      record.ID,
		GoalAssessment:     goal,
		GapSummary:         gap,
		IssuesFound:        mustJSON([]string{"建议补充业务验收脚本", "建议预置跨集群存储连通性巡检"}),
		ImprovementActions: mustJSON([]string{"完善恢复前校验清单", "建立季度灾备演练节奏"}),
		PublishedAt:        time.Now(),
		PublishedBy:        userID,
	}
	if err := s.drillReports.Create(ctx, report); err != nil {
		return nil, err
	}
	_ = s.writeAudit(ctx, userID, ActionDrillReportGen, ResourceTypeDrillReport, report.ID, record.WorkspaceID, derefUint64(record.ProjectID), domain.BackupAuditOutcomeSucceeded, map[string]any{
		"recordId": record.ID,
	})
	return report, nil
}

func (s *Service) ListAuditEvents(ctx context.Context, userID uint64, action, outcome, targetType string) ([]domain.BackupAuditEvent, error) {
	workspaceIDs, constrained, err := s.scope.ListReadableWorkspaceIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if constrained && len(workspaceIDs) == 0 {
		return []domain.BackupAuditEvent{}, nil
	}
	return s.auditRepo.List(ctx, repository.BackupAuditListFilter{
		WorkspaceIDs: workspaceIDs,
		Action:       action,
		Outcome:      outcome,
		TargetType:   targetType,
	})
}

func (s *Service) writeAudit(ctx context.Context, userID uint64, action, targetType string, targetID uint64, workspaceID, projectID uint64, outcome domain.BackupAuditOutcome, details map[string]any) error {
	if details == nil {
		details = map[string]any{}
	}
	details["workspaceId"] = workspaceID
	if projectID != 0 {
		details["projectId"] = projectID
	}
	details["targetRef"] = fmt.Sprintf("workspace:%d/project:%d/%s:%d", workspaceID, projectID, targetType, targetID)
	if s != nil && s.auditRepo != nil {
		if err := s.auditRepo.Create(ctx, &domain.BackupAuditEvent{
			Action:         action,
			ActorUserID:    userID,
			TargetType:     targetType,
			TargetRef:      strconv.FormatUint(targetID, 10),
			WorkspaceID:    workspaceID,
			ProjectID:      uint64PtrIf(projectID),
			ScopeSnapshot:  mustJSON(map[string]any{"workspaceId": workspaceID, "projectId": projectID}),
			Outcome:        outcome,
			DetailSnapshot: mustJSON(details),
			OccurredAt:     time.Now(),
		}); err != nil {
			return err
		}
	}
	if s == nil || s.auditWriter == nil {
		return nil
	}
	actorID := userID
	return s.auditWriter.WriteBackupRestoreEvent(ctx, "", &actorID, action, strconv.FormatUint(targetID, 10), toAuditOutcome(outcome), details)
}

func normalizeScopeIDs(workspaceID, projectID uint64) (uint64, uint64) {
	if projectID != 0 && workspaceID == 0 {
		return 0, projectID
	}
	return workspaceID, projectID
}

func (s *Service) defaultScopeForUser(ctx context.Context, userID uint64) (uint64, uint64, error) {
	workspaceIDs, constrained, err := s.scope.ListReadableWorkspaceIDs(ctx, userID)
	if err != nil {
		return 0, 0, err
	}
	if len(workspaceIDs) > 0 {
		return workspaceIDs[0], 0, nil
	}
	if constrained {
		return 0, 0, ErrBackupRestoreScopeDenied
	}
	return 0, 0, nil
}

func mustJSON(v any) string {
	if v == nil {
		return "{}"
	}
	b, err := json.Marshal(v)
	if err != nil {
		return "{}"
	}
	return string(b)
}

func uint64PtrIf(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	copy := v
	return &copy
}

func derefUint64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func durationFrom(started, completed time.Time, fallback int) int {
	if completed.After(started) {
		if seconds := int(completed.Sub(started).Seconds()); seconds > 0 {
			return seconds
		}
	}
	if fallback > 0 {
		return fallback
	}
	return 1
}

func outcomeForResult(status string) domain.BackupAuditOutcome {
	switch strings.TrimSpace(status) {
	case "succeeded", "approved", "active", "draft":
		return domain.BackupAuditOutcomeSucceeded
	case "blocked":
		return domain.BackupAuditOutcomeBlocked
	case "canceled":
		return domain.BackupAuditOutcomeCanceled
	default:
		return domain.BackupAuditOutcomeFailed
	}
}

func outcomeForValidation(status string) domain.BackupAuditOutcome {
	if strings.EqualFold(strings.TrimSpace(status), "passed") {
		return domain.BackupAuditOutcomeSucceeded
	}
	return domain.BackupAuditOutcomeBlocked
}

func parseUint64(value string) (uint64, error) {
	out, err := strconv.ParseUint(strings.TrimSpace(value), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("parse uint64: %w", err)
	}
	return out, nil
}

func backupRestoreProgressKey(parts ...string) string {
	return backupRestoreRedisKey("backuprestore:progress", parts...)
}

func backupRestorePrecheckKey(parts ...string) string {
	return backupRestoreRedisKey("backuprestore:precheck", parts...)
}

func backupRestoreLockKey(parts ...string) string {
	return backupRestoreRedisKey("backuprestore:lock", parts...)
}

func backupRestoreRedisKey(prefix string, parts ...string) string {
	key := prefix
	for _, part := range parts {
		if strings.TrimSpace(part) == "" {
			continue
		}
		key += ":" + strings.TrimSpace(part)
	}
	return key
}

func toAuditOutcome(outcome domain.BackupAuditOutcome) domain.AuditOutcome {
	switch outcome {
	case domain.BackupAuditOutcomeSucceeded:
		return domain.AuditOutcomeSuccess
	case domain.BackupAuditOutcomeBlocked, domain.BackupAuditOutcomeCanceled:
		return domain.AuditOutcomeDenied
	default:
		return domain.AuditOutcomeFailed
	}
}
