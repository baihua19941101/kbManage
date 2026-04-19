package sre

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	sreint "kbmanage/backend/internal/integration/sre"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
)

const (
	ResourceTypeSRE = "sre"

	ActionHAPolicyUpsert          = "sre.ha-policy.upsert"
	ActionHealthOverviewRead      = "sre.health-overview.read"
	ActionMaintenanceWindowUpsert = "sre.maintenance-window.upsert"
	ActionUpgradePrecheck         = "sre.upgrade.precheck"
	ActionUpgradeCreate           = "sre.upgrade.create"
	ActionRollbackValidate        = "sre.rollback.validate"
	ActionScaleEvidenceRead       = "sre.scale-evidence.read"
	ActionRunbookRead             = "sre.runbook.read"
)

var (
	ErrSREScopeDenied = errors.New("sre scope access denied")
	ErrSREInvalid     = errors.New("sre invalid request")
)

type Service struct {
	haPolicies         *repository.HAPolicyRepository
	maintenanceWindows *repository.MaintenanceWindowRepository
	healthSnapshots    *repository.PlatformHealthSnapshotRepository
	capacityBaselines  *repository.CapacityBaselineRepository
	upgradePlans       *repository.SREUpgradePlanRepository
	rollbacks          *repository.RollbackValidationRepository
	runbooks           *repository.RunbookArticleRepository
	alertBaselines     *repository.AlertBaselineRepository
	scaleEvidence      *repository.ScaleEvidenceRepository
	scope              *ScopeService
	healthProvider     sreint.HealthProvider
	upgradeValidator   sreint.UpgradeValidator
	scaleAnalyzer      sreint.ScaleAnalyzer
	healthCache        *HealthCache
	upgradeCoordinator *UpgradeCoordinator
	scaleCache         *ScaleCache
	auditWriter        *auditSvc.EventWriter
}

func NewService(
	haPolicyRepo *repository.HAPolicyRepository,
	maintenanceRepo *repository.MaintenanceWindowRepository,
	healthRepo *repository.PlatformHealthSnapshotRepository,
	capacityRepo *repository.CapacityBaselineRepository,
	upgradeRepo *repository.SREUpgradePlanRepository,
	rollbackRepo *repository.RollbackValidationRepository,
	runbookRepo *repository.RunbookArticleRepository,
	alertRepo *repository.AlertBaselineRepository,
	scaleRepo *repository.ScaleEvidenceRepository,
	bindingRepo *repository.ScopeRoleBindingRepository,
	projectRepo *repository.ProjectRepository,
	workspaceClusterRepo *repository.WorkspaceClusterRepository,
	healthProvider sreint.HealthProvider,
	upgradeValidator sreint.UpgradeValidator,
	scaleAnalyzer sreint.ScaleAnalyzer,
	healthCache *HealthCache,
	upgradeCoordinator *UpgradeCoordinator,
	scaleCache *ScaleCache,
	auditWriter *auditSvc.EventWriter,
) *Service {
	if healthProvider == nil {
		healthProvider = sreint.NewStaticHealthProvider()
	}
	if upgradeValidator == nil {
		upgradeValidator = sreint.NewStaticUpgradeValidator()
	}
	if scaleAnalyzer == nil {
		scaleAnalyzer = sreint.NewStaticScaleAnalyzer()
	}
	return &Service{
		haPolicies:         haPolicyRepo,
		maintenanceWindows: maintenanceRepo,
		healthSnapshots:    healthRepo,
		capacityBaselines:  capacityRepo,
		upgradePlans:       upgradeRepo,
		rollbacks:          rollbackRepo,
		runbooks:           runbookRepo,
		alertBaselines:     alertRepo,
		scaleEvidence:      scaleRepo,
		scope:              NewScopeService(bindingRepo, projectRepo, workspaceClusterRepo),
		healthProvider:     healthProvider,
		upgradeValidator:   upgradeValidator,
		scaleAnalyzer:      scaleAnalyzer,
		healthCache:        healthCache,
		upgradeCoordinator: upgradeCoordinator,
		scaleCache:         scaleCache,
		auditWriter:        auditWriter,
	}
}

type HAPolicyListFilter struct {
	Status  string
	Keyword string
}

type HAPolicyInput struct {
	WorkspaceID           uint64  `json:"workspaceId"`
	ProjectID             *uint64 `json:"projectId"`
	Name                  string  `json:"name"`
	ControlPlaneScope     string  `json:"controlPlaneScope"`
	DeploymentMode        string  `json:"deploymentMode"`
	ReplicaExpectation    int     `json:"replicaExpectation"`
	FailoverTriggerPolicy string  `json:"failoverTriggerPolicy"`
	FailoverCooldown      string  `json:"failoverCooldown"`
}

type MaintenanceWindowInput struct {
	WorkspaceID       uint64   `json:"workspaceId"`
	ProjectID         *uint64  `json:"projectId"`
	Name              string   `json:"name"`
	WindowType        string   `json:"windowType"`
	Scope             string   `json:"scope"`
	StartAt           string   `json:"startAt"`
	EndAt             string   `json:"endAt"`
	AllowedOperations []string `json:"allowedOperations"`
	RestrictedOps     []string `json:"restrictedOperations"`
	ApprovalRecord    string   `json:"approvalRecord"`
	ExceptionReason   string   `json:"exceptionReason"`
}

type UpgradePrecheckInput struct {
	WorkspaceID    uint64  `json:"workspaceId"`
	ProjectID      *uint64 `json:"projectId"`
	CurrentVersion string  `json:"currentVersion"`
	TargetVersion  string  `json:"targetVersion"`
	Scope          string  `json:"scope"`
}

type SREUpgradePlanInput struct {
	WorkspaceID         uint64  `json:"workspaceId"`
	ProjectID           *uint64 `json:"projectId"`
	MaintenanceWindowID *uint64 `json:"maintenanceWindowId"`
	Name                string  `json:"name"`
	CurrentVersion      string  `json:"currentVersion"`
	TargetVersion       string  `json:"targetVersion"`
	RolloutStrategy     string  `json:"rolloutStrategy"`
}

type RollbackValidationInput struct {
	ValidationScope string   `json:"validationScope"`
	Preconditions   []string `json:"preconditions"`
	Result          string   `json:"result"`
	RemainingRisk   string   `json:"remainingRisk"`
}

type CapacityBaselineListFilter struct {
	Status string
}

type ScaleEvidenceListFilter struct {
	EvidenceType string
}

func (s *Service) writeAudit(ctx context.Context, actorID uint64, action, targetType, targetRef string, outcome domain.AuditOutcome, details map[string]any) {
	if s.auditWriter == nil {
		return
	}
	actor := actorID
	_ = s.auditWriter.Write(ctx, "", &actor, action, ResourceTypeSRE, firstNonEmptyString(targetRef, targetType), outcome, details)
}

func (s *Service) ListHAPolicies(ctx context.Context, userID uint64, filter HAPolicyListFilter) ([]domain.HAPolicy, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "sre:read"); err != nil {
		return nil, err
	}
	return s.haPolicies.List(ctx, repository.HAPolicyListFilter{Status: strings.TrimSpace(filter.Status), Keyword: strings.TrimSpace(filter.Keyword)})
}

func (s *Service) UpsertHAPolicy(ctx context.Context, userID uint64, input HAPolicyInput) (*domain.HAPolicy, error) {
	if err := s.scope.EnsureScopePermission(ctx, userID, input.WorkspaceID, input.ProjectID, "sre:manage-ha"); err != nil {
		return nil, err
	}
	if strings.TrimSpace(input.Name) == "" || strings.TrimSpace(input.DeploymentMode) == "" || input.ReplicaExpectation < 1 {
		return nil, ErrSREInvalid
	}
	item := &domain.HAPolicy{
		WorkspaceID:           input.WorkspaceID,
		ProjectID:             input.ProjectID,
		Name:                  strings.TrimSpace(input.Name),
		ControlPlaneScope:     firstNonEmptyString(input.ControlPlaneScope, "platform-control-plane"),
		DeploymentMode:        strings.TrimSpace(input.DeploymentMode),
		ReplicaExpectation:    input.ReplicaExpectation,
		FailoverTriggerPolicy: input.FailoverTriggerPolicy,
		FailoverCooldown:      input.FailoverCooldown,
		TakeoverStatus:        "standby",
		LastRecoveryResult:    "最近一次恢复校验通过",
		Status:                domain.HAPolicyStatusActive,
		OwnerUserID:           userID,
	}
	if err := s.haPolicies.Create(ctx, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionHAPolicyUpsert, "ha-policy", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"workspaceId": item.WorkspaceID})
	return item, nil
}

func (s *Service) ListMaintenanceWindows(ctx context.Context, userID uint64) ([]domain.MaintenanceWindow, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "sre:read"); err != nil {
		return nil, err
	}
	return s.maintenanceWindows.List(ctx, repository.MaintenanceWindowListFilter{})
}

func (s *Service) UpsertMaintenanceWindow(ctx context.Context, userID uint64, input MaintenanceWindowInput) (*domain.MaintenanceWindow, error) {
	if err := s.scope.EnsureScopePermission(ctx, userID, input.WorkspaceID, input.ProjectID, "sre:manage-ha"); err != nil {
		return nil, err
	}
	startAt, err := time.Parse(time.RFC3339, strings.TrimSpace(input.StartAt))
	if err != nil {
		return nil, ErrSREInvalid
	}
	endAt, err := time.Parse(time.RFC3339, strings.TrimSpace(input.EndAt))
	if err != nil || !endAt.After(startAt) {
		return nil, ErrSREInvalid
	}
	item := &domain.MaintenanceWindow{
		WorkspaceID:          input.WorkspaceID,
		ProjectID:            input.ProjectID,
		Name:                 strings.TrimSpace(input.Name),
		WindowType:           firstNonEmptyString(input.WindowType, "maintenance"),
		Scope:                firstNonEmptyString(input.Scope, "platform"),
		StartAt:              startAt,
		EndAt:                endAt,
		AllowedOperations:    strings.Join(input.AllowedOperations, ","),
		RestrictedOperations: strings.Join(input.RestrictedOps, ","),
		Status:               domain.MaintenanceWindowStatusScheduled,
		ExceptionReason:      input.ExceptionReason,
		ApprovalRecord:       input.ApprovalRecord,
		PostCheckStatus:      "pending",
		OwnerUserID:          userID,
	}
	if err := s.maintenanceWindows.Create(ctx, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionMaintenanceWindowUpsert, "maintenance-window", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"workspaceId": item.WorkspaceID})
	return item, nil
}

func (s *Service) GetHealthOverview(ctx context.Context, userID uint64, workspaceID uint64, projectID *uint64) (*domain.PlatformHealthSnapshot, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "sre:read"); err != nil {
		return nil, err
	}
	var resolvedProjectID uint64
	if projectID != nil {
		resolvedProjectID = *projectID
	}
	policy, _ := s.haPolicies.FindLatestByScope(ctx, workspaceID, resolvedProjectID)
	window, _ := s.maintenanceWindows.FindActiveByScope(ctx, workspaceID, resolvedProjectID, time.Now())
	baseline, _ := s.capacityBaselines.FindLatestByScope(ctx, workspaceID, resolvedProjectID)
	evidence, _ := s.scaleEvidence.FindLatestByScope(ctx, workspaceID, resolvedProjectID)
	upgrade, _ := s.upgradePlans.FindLatestByScope(ctx, workspaceID, resolvedProjectID)
	item := s.healthProvider.BuildOverview(ctx, sreint.HealthOverviewInput{
		Policy:            policy,
		ActiveWindow:      window,
		LatestBaseline:    baseline,
		LatestEvidence:    evidence,
		LatestUpgradePlan: upgrade,
		WorkspaceID:       workspaceID,
		ProjectID:         resolvedProjectID,
	})
	if err := s.healthSnapshots.Create(ctx, item); err != nil {
		return nil, err
	}
	_ = s.healthCache.Store(ctx, fmt.Sprintf("%d:%d", workspaceID, resolvedProjectID), item)
	s.writeAudit(ctx, userID, ActionHealthOverviewRead, "health-overview", fmt.Sprintf("%d:%d", workspaceID, resolvedProjectID), domain.AuditOutcomeSuccess, nil)
	return item, nil
}

func (s *Service) RunUpgradePrecheck(ctx context.Context, userID uint64, input UpgradePrecheckInput) (sreint.UpgradePrecheckResult, error) {
	if err := s.scope.EnsureScopePermission(ctx, userID, input.WorkspaceID, input.ProjectID, "sre:manage-upgrade"); err != nil {
		return sreint.UpgradePrecheckResult{}, err
	}
	result := s.upgradeValidator.Validate(ctx, sreint.UpgradePrecheckInput{
		CurrentVersion: input.CurrentVersion,
		TargetVersion:  input.TargetVersion,
		WorkspaceID:    input.WorkspaceID,
	})
	outcome := domain.AuditOutcomeSuccess
	if result.Decision == "block" {
		outcome = domain.AuditOutcomeFailed
	}
	s.writeAudit(ctx, userID, ActionUpgradePrecheck, "upgrade-precheck", firstNonEmptyString(input.TargetVersion, "unknown"), outcome, map[string]any{"decision": result.Decision})
	return result, nil
}

func (s *Service) ListUpgradePlans(ctx context.Context, userID uint64) ([]domain.SREUpgradePlan, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "sre:read"); err != nil {
		return nil, err
	}
	return s.upgradePlans.List(ctx, repository.SREUpgradePlanListFilter{})
}

func (s *Service) CreateUpgradePlan(ctx context.Context, userID uint64, input SREUpgradePlanInput) (*domain.SREUpgradePlan, error) {
	if err := s.scope.EnsureScopePermission(ctx, userID, input.WorkspaceID, input.ProjectID, "sre:manage-upgrade"); err != nil {
		return nil, err
	}
	precheck := s.upgradeValidator.Validate(ctx, sreint.UpgradePrecheckInput{
		CurrentVersion: input.CurrentVersion,
		TargetVersion:  input.TargetVersion,
		WorkspaceID:    input.WorkspaceID,
	})
	status := domain.SREUpgradeStatusReady
	if precheck.Decision == "block" {
		status = domain.SREUpgradeStatusPrecheckFailed
	}
	item := &domain.SREUpgradePlan{
		WorkspaceID:          input.WorkspaceID,
		ProjectID:            input.ProjectID,
		MaintenanceWindowID:  input.MaintenanceWindowID,
		Name:                 strings.TrimSpace(input.Name),
		CurrentVersion:       strings.TrimSpace(input.CurrentVersion),
		TargetVersion:        strings.TrimSpace(input.TargetVersion),
		CompatibilitySummary: precheck.CompatibilitySummary,
		PrecheckResult:       strings.Join(append(precheck.Blockers, precheck.Warnings...), "；"),
		RolloutStrategy:      firstNonEmptyString(input.RolloutStrategy, "rolling"),
		ExecutionStage:       "planned",
		ExecutionProgress:    0,
		AcceptanceResult:     "pending",
		RollbackReadiness:    "validated",
		Status:               status,
		CreatedBy:            userID,
	}
	if err := s.upgradePlans.Create(ctx, item); err != nil {
		return nil, err
	}
	_ = s.upgradeCoordinator.Mark(ctx, strconv.FormatUint(item.ID, 10), string(item.Status))
	s.writeAudit(ctx, userID, ActionUpgradeCreate, "upgrade-plan", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"status": item.Status})
	return item, nil
}

func (s *Service) CreateRollbackValidation(ctx context.Context, userID, upgradeID uint64, input RollbackValidationInput) (*domain.RollbackValidation, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "sre:manage-upgrade"); err != nil {
		return nil, err
	}
	item := &domain.RollbackValidation{
		UpgradePlanID:   upgradeID,
		ValidationScope: firstNonEmptyString(input.ValidationScope, "platform"),
		Preconditions:   strings.Join(input.Preconditions, ","),
		Result:          domain.RollbackValidationResult(firstNonEmptyString(input.Result, string(domain.RollbackValidationResultPassed))),
		RemainingRisk:   input.RemainingRisk,
		ValidatedAt:     time.Now(),
		ValidatedBy:     userID,
	}
	if err := s.rollbacks.Create(ctx, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionRollbackValidate, "rollback-validation", strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, nil)
	return item, nil
}

func (s *Service) ListCapacityBaselines(ctx context.Context, userID uint64, filter CapacityBaselineListFilter) ([]domain.CapacityBaseline, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "sre:read"); err != nil {
		return nil, err
	}
	return s.capacityBaselines.List(ctx, repository.CapacityBaselineListFilter{Status: filter.Status})
}

func (s *Service) ListScaleEvidence(ctx context.Context, userID uint64, filter ScaleEvidenceListFilter) ([]domain.ScaleEvidence, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "sre:read"); err != nil {
		return nil, err
	}
	items, err := s.scaleEvidence.List(ctx, repository.ScaleEvidenceListFilter{EvidenceType: filter.EvidenceType})
	if err == nil {
		_ = s.scaleCache.Mark(ctx, firstNonEmptyString(filter.EvidenceType, "all"), strconv.Itoa(len(items)))
	}
	s.writeAudit(ctx, userID, ActionScaleEvidenceRead, "scale-evidence", firstNonEmptyString(filter.EvidenceType, "all"), domain.AuditOutcomeSuccess, nil)
	return items, err
}

func (s *Service) ListRunbooks(ctx context.Context, userID uint64) ([]domain.RunbookArticle, error) {
	if err := s.scope.EnsurePermission(ctx, userID, "sre:read"); err != nil {
		return nil, err
	}
	items, err := s.runbooks.List(ctx, repository.RunbookArticleListFilter{})
	if err == nil {
		s.writeAudit(ctx, userID, ActionRunbookRead, "runbook", "list", domain.AuditOutcomeSuccess, map[string]any{"count": len(items)})
	}
	return items, err
}

func firstNonEmptyString(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
