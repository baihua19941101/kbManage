package workloadops

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
)

var (
	ErrWorkloadOpsNotConfigured = errors.New("workload operations service is not configured")
	ErrInvalidWorkloadReference = errors.New("invalid workload reference")
	ErrWorkloadOpsScopeDenied   = errors.New("workload operations scope access denied")
)

type WorkloadReference struct {
	ClusterID    uint64
	WorkspaceID  uint64
	ProjectID    uint64
	Namespace    string
	ResourceKind string
	ResourceName string
}

type WorkloadInstance struct {
	PodName           string     `json:"podName"`
	ContainerName     string     `json:"containerName,omitempty"`
	NodeName          string     `json:"nodeName,omitempty"`
	Phase             string     `json:"phase"`
	Ready             bool       `json:"ready"`
	RestartCount      int        `json:"restartCount"`
	StartedAt         *time.Time `json:"startedAt,omitempty"`
	LastTransitionAt  *time.Time `json:"lastTransitionAt,omitempty"`
	LogAvailable      bool       `json:"logAvailable"`
	TerminalAvailable bool       `json:"terminalAvailable"`
}

type ReleaseRevision struct {
	Revision          int        `json:"revision"`
	SourceKind        string     `json:"sourceKind"`
	SourceName        string     `json:"sourceName"`
	ChangeCause       string     `json:"changeCause,omitempty"`
	CreatedAt         *time.Time `json:"createdAt,omitempty"`
	IsCurrent         bool       `json:"isCurrent"`
	RollbackAvailable bool       `json:"rollbackAvailable"`
	Summary           string     `json:"summary,omitempty"`
}

type SubmitWorkloadActionRequest struct {
	RequestID         string
	OperatorID        uint64
	Target            WorkloadReference
	TargetInstanceRef string
	ActionType        domain.WorkloadActionType
	RiskLevel         domain.RiskLevel
	RiskConfirmed     bool
	PayloadJSON       string
	BatchID           *uint64
}

type SubmitBatchOperationRequest struct {
	RequestID     string
	OperatorID    uint64
	ActionType    domain.WorkloadActionType
	RiskLevel     domain.RiskLevel
	RiskConfirmed bool
	Targets       []WorkloadReference
	PayloadJSON   string
}

type CreateTerminalSessionRequest struct {
	OperatorID    uint64
	ClusterID     uint64
	WorkspaceID   uint64
	ProjectID     uint64
	Namespace     string
	PodName       string
	ContainerName string
	WorkloadKind  string
	WorkloadName  string
	Cols          int
	Rows          int
}

type Service struct {
	actions   *repository.WorkloadActionRepository
	batches   *repository.BatchOperationRepository
	sessions  *repository.TerminalSessionRepository
	scope     *ScopeService
	progress  *ProgressCache
	sessionKV *SessionCache
	executor  ActionExecutor
	audit     *auditSvc.EventWriter
}

func NewService(
	actions *repository.WorkloadActionRepository,
	batches *repository.BatchOperationRepository,
	sessions *repository.TerminalSessionRepository,
	scope *ScopeService,
	progress *ProgressCache,
	sessionKV *SessionCache,
) *Service {
	return &Service{
		actions:   actions,
		batches:   batches,
		sessions:  sessions,
		scope:     scope,
		progress:  progress,
		sessionKV: sessionKV,
		executor:  NewActionExecutor(),
	}
}

func (s *Service) SetAuditWriter(writer *auditSvc.EventWriter) {
	if s == nil {
		return
	}
	s.audit = writer
}

func (s *Service) GetContext(ctx context.Context, userID uint64, target WorkloadReference) (map[string]any, error) {
	if err := s.validateTargetAndAccess(ctx, userID, target, "workloadops:read"); err != nil {
		return nil, err
	}
	return s.buildContext(ctx, target), nil
}

func (s *Service) ListInstances(ctx context.Context, userID uint64, target WorkloadReference) ([]WorkloadInstance, error) {
	if err := s.validateTargetAndAccess(ctx, userID, target, "workloadops:read"); err != nil {
		return nil, err
	}
	return s.buildInstances(ctx, target), nil
}

func (s *Service) ListRevisions(ctx context.Context, userID uint64, target WorkloadReference) ([]ReleaseRevision, error) {
	if err := s.validateTargetAndAccess(ctx, userID, target, "workloadops:read"); err != nil {
		return nil, err
	}
	return s.buildRevisions(ctx, target), nil
}

func (s *Service) SubmitAction(ctx context.Context, req SubmitWorkloadActionRequest) (*domain.WorkloadActionRequest, error) {
	return s.createAndExecuteAction(ctx, req)
}

func (s *Service) GetAction(ctx context.Context, userID uint64, actionID uint64) (*domain.WorkloadActionRequest, error) {
	if s.actions == nil {
		return nil, ErrWorkloadOpsNotConfigured
	}
	item, err := s.actions.GetByID(ctx, actionID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureClusterAccess(ctx, userID, item.ClusterID, "workloadops:read"); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) SubmitBatch(ctx context.Context, req SubmitBatchOperationRequest) (*domain.BatchOperationTask, error) {
	if s.batches == nil {
		return nil, ErrWorkloadOpsNotConfigured
	}
	if req.OperatorID == 0 || len(req.Targets) == 0 {
		return nil, ErrInvalidWorkloadReference
	}
	for _, target := range req.Targets {
		if err := s.validateTargetAndAccess(ctx, req.OperatorID, target, "workloadops:batch"); err != nil {
			return nil, err
		}
	}

	task := &domain.BatchOperationTask{
		RequestID:       normalizeRequestID(req.RequestID, req.OperatorID),
		OperatorID:      req.OperatorID,
		ActionType:      req.ActionType,
		RiskLevel:       req.RiskLevel,
		RiskConfirmed:   req.RiskConfirmed,
		TotalTargets:    len(req.Targets),
		Status:          domain.BatchOperationStatusPending,
		ProgressPercent: 0,
	}
	if task.RiskLevel == "" {
		task.RiskLevel = domain.RiskLevelHigh
	}
	if task.ActionType == "" {
		task.ActionType = domain.WorkloadActionTypeRestart
	}
	if err := s.batches.CreateTask(ctx, task); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, req.RequestID, req.OperatorID, auditSvc.WorkloadOpsAuditBatchSubmit, domain.AuditOutcomeSuccess, nil, nil, map[string]any{
		"batchId":       task.ID,
		"actionType":    task.ActionType,
		"totalTargets":  task.TotalTargets,
		"status":        task.Status,
		"riskLevel":     task.RiskLevel,
		"riskConfirmed": task.RiskConfirmed,
	})

	items := make([]domain.BatchOperationItem, 0, len(req.Targets))
	for _, target := range req.Targets {
		items = append(items, domain.BatchOperationItem{
			BatchID:      task.ID,
			ClusterID:    target.ClusterID,
			WorkspaceID:  ptrUint64(target.WorkspaceID),
			ProjectID:    ptrUint64(target.ProjectID),
			Namespace:    strings.TrimSpace(target.Namespace),
			ResourceKind: strings.TrimSpace(target.ResourceKind),
			ResourceName: strings.TrimSpace(target.ResourceName),
			Status:       domain.BatchOperationItemStatusPending,
		})
	}
	if err := s.batches.CreateItems(ctx, task.ID, items); err != nil {
		return nil, err
	}
	_, createdItems, err := s.batches.GetTaskByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	if err := s.executeBatchTask(
		ctx,
		req.OperatorID,
		task,
		createdItems,
		req.Targets,
		req.ActionType,
		req.RiskLevel,
		req.RiskConfirmed,
		req.PayloadJSON,
	); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, req.RequestID, req.OperatorID, auditSvc.WorkloadOpsAuditBatchFinish, domain.AuditOutcomeSuccess, nil, nil, map[string]any{
		"batchId":          task.ID,
		"actionType":       task.ActionType,
		"status":           task.Status,
		"succeededTargets": task.SucceededTargets,
		"failedTargets":    task.FailedTargets,
		"canceledTargets":  task.CanceledTargets,
	})
	task, _, err = s.batches.GetTaskByID(ctx, task.ID)
	if err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) GetBatch(ctx context.Context, userID uint64, batchID uint64) (*domain.BatchOperationTask, []domain.BatchOperationItem, error) {
	if s.batches == nil {
		return nil, nil, ErrWorkloadOpsNotConfigured
	}
	task, items, err := s.batches.GetTaskByID(ctx, batchID)
	if err != nil {
		return nil, nil, err
	}
	for _, item := range items {
		if err := s.ensureClusterAccess(ctx, userID, item.ClusterID, "workloadops:read"); err != nil {
			return nil, nil, err
		}
	}
	return task, items, nil
}

func (s *Service) CreateTerminalSession(ctx context.Context, req CreateTerminalSessionRequest) (*domain.TerminalSession, error) {
	if s.sessions == nil {
		return nil, ErrWorkloadOpsNotConfigured
	}
	if req.OperatorID == 0 || req.ClusterID == 0 || strings.TrimSpace(req.Namespace) == "" || strings.TrimSpace(req.PodName) == "" || strings.TrimSpace(req.ContainerName) == "" {
		return nil, ErrInvalidWorkloadReference
	}
	if err := s.ensureTargetAccess(ctx, req.OperatorID, req.ClusterID, req.WorkspaceID, req.ProjectID, "workloadops:terminal"); err != nil {
		return nil, err
	}
	now := time.Now()
	item := &domain.TerminalSession{
		SessionKey:    fmt.Sprintf("wops-%d-%d", req.OperatorID, now.UnixNano()),
		OperatorID:    req.OperatorID,
		ClusterID:     req.ClusterID,
		WorkspaceID:   ptrUint64(req.WorkspaceID),
		ProjectID:     ptrUint64(req.ProjectID),
		Namespace:     strings.TrimSpace(req.Namespace),
		PodName:       strings.TrimSpace(req.PodName),
		ContainerName: strings.TrimSpace(req.ContainerName),
		WorkloadKind:  strings.TrimSpace(req.WorkloadKind),
		WorkloadName:  strings.TrimSpace(req.WorkloadName),
		Status:        domain.TerminalSessionStatusActive,
		StartedAt:     &now,
	}
	if err := s.sessions.Create(ctx, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, item.SessionKey, req.OperatorID, auditSvc.WorkloadOpsAuditTerminalOpen, domain.AuditOutcomeSuccess, nil, item, nil)
	return item, nil
}

func (s *Service) GetTerminalSession(ctx context.Context, userID uint64, sessionID uint64) (*domain.TerminalSession, error) {
	if s.sessions == nil {
		return nil, ErrWorkloadOpsNotConfigured
	}
	item, err := s.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	if err := s.ensureTargetAccess(ctx, userID, item.ClusterID, derefUint64(item.WorkspaceID), derefUint64(item.ProjectID), "workloadops:terminal"); err != nil {
		return nil, err
	}
	return s.expireSessionIfNeeded(ctx, item)
}

func (s *Service) CloseTerminalSession(ctx context.Context, userID uint64, sessionID uint64) error {
	if s.sessions == nil {
		return ErrWorkloadOpsNotConfigured
	}
	item, err := s.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return err
	}
	if err := s.ensureTargetAccess(ctx, userID, item.ClusterID, derefUint64(item.WorkspaceID), derefUint64(item.ProjectID), "workloadops:terminal"); err != nil {
		return err
	}
	if err := s.sessions.UpdateStatus(ctx, sessionID, domain.TerminalSessionStatusClosed, "closed by user"); err != nil {
		return err
	}
	latest, _ := s.sessions.GetByID(ctx, sessionID)
	s.writeAudit(ctx, normalizeRequestID("", userID), userID, auditSvc.WorkloadOpsAuditTerminalClose, domain.AuditOutcomeSuccess, nil, withTerminalAuditBoundary(latest), nil)
	return nil
}

type ResourceSelector struct {
	ClusterID    uint64
	Namespace    string
	ResourceKind string
	ResourceName string
}

func (s *Service) validateTargetAndAccess(ctx context.Context, userID uint64, target WorkloadReference, permission string) error {
	if userID == 0 || target.ClusterID == 0 || strings.TrimSpace(target.Namespace) == "" || strings.TrimSpace(target.ResourceKind) == "" || strings.TrimSpace(target.ResourceName) == "" {
		return ErrInvalidWorkloadReference
	}
	return s.ensureTargetAccess(ctx, userID, target.ClusterID, target.WorkspaceID, target.ProjectID, permission)
}

func (s *Service) ensureClusterAccess(ctx context.Context, userID uint64, clusterID uint64, permission string) error {
	return s.ensureTargetAccess(ctx, userID, clusterID, 0, 0, permission)
}

func (s *Service) ensureTargetAccess(
	ctx context.Context,
	userID uint64,
	clusterID uint64,
	workspaceID uint64,
	projectID uint64,
	permission string,
) error {
	if s.scope == nil {
		return nil
	}
	if err := s.scope.ValidateWorkloadAccess(ctx, userID, clusterID, workspaceID, projectID, permission); err != nil {
		if errors.Is(err, ErrWorkloadOpsScopeDenied) {
			return ErrWorkloadOpsScopeDenied
		}
		return err
	}
	return nil
}

func ptrUint64(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	out := v
	return &out
}

func derefUint64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func valueOrNil(v uint64) any {
	if v == 0 {
		return nil
	}
	return v
}

func normalizeRequestID(requestID string, operatorID uint64) string {
	if trimmed := strings.TrimSpace(requestID); trimmed != "" {
		return trimmed
	}
	return fmt.Sprintf("wops-%d-%d", operatorID, time.Now().UnixNano())
}

func (s *Service) createAndExecuteAction(ctx context.Context, req SubmitWorkloadActionRequest) (*domain.WorkloadActionRequest, error) {
	if s.actions == nil {
		return nil, ErrWorkloadOpsNotConfigured
	}
	if req.OperatorID == 0 {
		return nil, ErrInvalidWorkloadReference
	}
	if err := s.validateTargetAndAccess(ctx, req.OperatorID, req.Target, requiredPermissionByActionType(req.ActionType)); err != nil {
		return nil, err
	}

	requestID := normalizeRequestID(req.RequestID, req.OperatorID)
	if existing, err := s.actions.GetByRequestID(ctx, requestID); err == nil {
		return existing, nil
	}

	payloadJSON := req.PayloadJSON
	if payloadJSON == "" {
		payloadJSON = "{}"
	}
	item := &domain.WorkloadActionRequest{
		RequestID:         requestID,
		OperatorID:        req.OperatorID,
		ClusterID:         req.Target.ClusterID,
		WorkspaceID:       ptrUint64(req.Target.WorkspaceID),
		ProjectID:         ptrUint64(req.Target.ProjectID),
		Namespace:         strings.TrimSpace(req.Target.Namespace),
		ResourceKind:      strings.TrimSpace(req.Target.ResourceKind),
		ResourceName:      strings.TrimSpace(req.Target.ResourceName),
		TargetInstanceRef: strings.TrimSpace(req.TargetInstanceRef),
		ActionType:        req.ActionType,
		RiskLevel:         req.RiskLevel,
		RiskConfirmed:     req.RiskConfirmed,
		PayloadJSON:       payloadJSON,
		BatchID:           req.BatchID,
		Status:            domain.OperationStatusPending,
		ProgressMessage:   "action submitted",
	}
	if item.RiskLevel == "" {
		item.RiskLevel = domain.RiskLevelMedium
	}
	if item.ActionType == "" {
		item.ActionType = domain.WorkloadActionTypeRestart
	}
	if err := s.actions.Create(ctx, item); err != nil {
		return nil, err
	}
	submitAction := auditSvc.WorkloadOpsAuditActionSubmit
	successAction := auditSvc.WorkloadOpsAuditActionSuccess
	failureAction := auditSvc.WorkloadOpsAuditActionFailure
	if item.ActionType == domain.WorkloadActionTypeRollback {
		submitAction = auditSvc.WorkloadOpsAuditRollbackSubmit
		successAction = auditSvc.WorkloadOpsAuditRollbackFinish
		failureAction = auditSvc.WorkloadOpsAuditRollbackFinish
	}
	s.writeAudit(ctx, item.RequestID, item.OperatorID, submitAction, domain.AuditOutcomeSuccess, item, nil, nil)

	if err := s.actions.UpdateExecutionResult(ctx, item.ID, domain.OperationStatusRunning, "action is running", "", ""); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, item.RequestID, item.OperatorID, auditSvc.WorkloadOpsAuditActionStart, domain.AuditOutcomeSuccess, item, nil, nil)
	result, err := s.executor.Execute(ctx, item)
	if err != nil {
		_ = s.actions.UpdateExecutionResult(ctx, item.ID, domain.OperationStatusFailed, "action failed", result.ResultMessage, err.Error())
		latest, _ := s.actions.GetByID(ctx, item.ID)
		s.writeAudit(ctx, item.RequestID, item.OperatorID, failureAction, domain.AuditOutcomeFailed, latest, nil, nil)
		return s.actions.GetByID(ctx, item.ID)
	}
	if err := s.actions.UpdateExecutionResult(ctx, item.ID, domain.OperationStatusSucceeded, result.ProgressMessage, result.ResultMessage, ""); err != nil {
		return nil, err
	}
	latest, err := s.actions.GetByID(ctx, item.ID)
	if err == nil {
		s.writeAudit(ctx, item.RequestID, item.OperatorID, successAction, domain.AuditOutcomeSuccess, latest, nil, nil)
	}
	return latest, err
}

func encodePayload(payload map[string]any) string {
	if len(payload) == 0 {
		return "{}"
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

func requiredPermissionByActionType(actionType domain.WorkloadActionType) string {
	if actionType == domain.WorkloadActionTypeRollback {
		return "workloadops:rollback"
	}
	return "workloadops:execute"
}

func (s *Service) writeAudit(
	ctx context.Context,
	requestID string,
	operatorID uint64,
	action string,
	outcome domain.AuditOutcome,
	item *domain.WorkloadActionRequest,
	session *domain.TerminalSession,
	extraDetails map[string]any,
) {
	if s == nil || s.audit == nil || operatorID == 0 {
		return
	}
	actorID := operatorID
	details := map[string]any{}
	resourceID := "workloadops"
	if item != nil {
		resourceID = fmt.Sprintf("cluster:%d/ns:%s/kind:%s/name:%s", item.ClusterID, item.Namespace, item.ResourceKind, item.ResourceName)
		details["clusterId"] = item.ClusterID
		details["workspaceId"] = item.WorkspaceID
		details["projectId"] = item.ProjectID
		details["namespace"] = item.Namespace
		details["resourceKind"] = item.ResourceKind
		details["resourceName"] = item.ResourceName
		details["actionType"] = item.ActionType
		details["status"] = item.Status
		details["riskLevel"] = item.RiskLevel
		details["resultMessage"] = item.ResultMessage
		details["failureReason"] = item.FailureReason
	}
	if session != nil {
		resourceID = fmt.Sprintf(
			"cluster:%d/workspace:%d/project:%d/ns:%s/pod:%s/container:%s",
			session.ClusterID,
			derefUint64(session.WorkspaceID),
			derefUint64(session.ProjectID),
			session.Namespace,
			session.PodName,
			session.ContainerName,
		)
		// terminal audit only records session lifecycle metadata, not command/output.
		details["clusterId"] = session.ClusterID
		details["namespace"] = session.Namespace
		details["podName"] = session.PodName
		details["containerName"] = session.ContainerName
		details["durationSeconds"] = session.DurationSeconds
		details["closeReason"] = session.CloseReason
	}
	for k, v := range extraDetails {
		details[k] = v
	}
	_ = s.audit.WriteWorkloadOpsEvent(ctx, requestID, &actorID, action, resourceID, outcome, details)
}

func withTerminalAuditBoundary(item *domain.TerminalSession) *domain.TerminalSession {
	if item == nil {
		return nil
	}
	clone := *item
	if clone.StartedAt != nil {
		endedAt := time.Now()
		if clone.EndedAt != nil {
			endedAt = *clone.EndedAt
		}
		duration := int(endedAt.Sub(*clone.StartedAt).Seconds())
		if duration < 0 {
			duration = 0
		}
		clone.DurationSeconds = duration
	}
	return &clone
}
