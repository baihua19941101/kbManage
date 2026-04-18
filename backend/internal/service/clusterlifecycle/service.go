package clusterlifecycle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	driverProvider "kbmanage/backend/internal/integration/clusterlifecycle/driver"
	validatorProvider "kbmanage/backend/internal/integration/clusterlifecycle/validator"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"

	"gorm.io/gorm"
)

const (
	ResourceTypeClusterLifecycle = "clusterlifecycle"
	ActionClusterImport          = "clusterlifecycle.cluster.import"
	ActionClusterRegister        = "clusterlifecycle.cluster.register"
	ActionClusterCreate          = "clusterlifecycle.cluster.create"
	ActionClusterValidate        = "clusterlifecycle.cluster.validate"
	ActionClusterUpgrade         = "clusterlifecycle.cluster.upgrade"
	ActionClusterScaleNodePool   = "clusterlifecycle.cluster.nodepool.scale"
	ActionClusterDisable         = "clusterlifecycle.cluster.disable"
	ActionClusterRetire          = "clusterlifecycle.cluster.retire"
	ActionDriverCreate           = "clusterlifecycle.driver.create"
	ActionTemplateCreate         = "clusterlifecycle.template.create"
)

var (
	ErrLifecycleScopeDenied = errors.New("cluster lifecycle scope access denied")
	ErrLifecycleConflict    = errors.New("cluster lifecycle operation conflict")
	ErrLifecycleInvalid     = errors.New("cluster lifecycle invalid request")
	ErrLifecycleBlocked     = errors.New("cluster lifecycle operation blocked")
)

type Service struct {
	clusters    *repository.ClusterLifecycleRepository
	operations  *repository.ClusterLifecycleOperationRepository
	drivers     *repository.ClusterDriverRepository
	templates   *repository.ClusterTemplateRepository
	capability  *repository.ClusterCapabilityRepository
	plans       *repository.UpgradePlanRepository
	nodePools   *repository.NodePoolRepository
	scope       *ScopeService
	progress    *ProgressCache
	validations *ValidationCache
	lock        *OperationLock
	driverOps   driverProvider.Provider
	validator   validatorProvider.Provider
	auditWriter *auditSvc.EventWriter
}

func NewService(
	clusterRepo *repository.ClusterLifecycleRepository,
	operationRepo *repository.ClusterLifecycleOperationRepository,
	driverRepo *repository.ClusterDriverRepository,
	templateRepo *repository.ClusterTemplateRepository,
	capabilityRepo *repository.ClusterCapabilityRepository,
	planRepo *repository.UpgradePlanRepository,
	nodePoolRepo *repository.NodePoolRepository,
	bindingRepo *repository.ScopeRoleBindingRepository,
	projectRepo *repository.ProjectRepository,
	progressCache *ProgressCache,
	validationCache *ValidationCache,
	lock *OperationLock,
	driverOps driverProvider.Provider,
	validator validatorProvider.Provider,
	auditWriter *auditSvc.EventWriter,
) *Service {
	return &Service{
		clusters:    clusterRepo,
		operations:  operationRepo,
		drivers:     driverRepo,
		templates:   templateRepo,
		capability:  capabilityRepo,
		plans:       planRepo,
		nodePools:   nodePoolRepo,
		scope:       NewScopeService(bindingRepo, projectRepo),
		progress:    progressCache,
		validations: validationCache,
		lock:        lock,
		driverOps:   driverOps,
		validator:   validator,
		auditWriter: auditWriter,
	}
}

type ClusterListFilter struct {
	Status             string
	InfrastructureType string
	DriverKey          string
	Keyword            string
}

type ImportClusterInput struct {
	Name               string `json:"name"`
	DisplayName        string `json:"displayName"`
	WorkspaceID        uint64 `json:"workspaceId"`
	ProjectID          uint64 `json:"projectId"`
	InfrastructureType string `json:"infrastructureType"`
	DriverKey          string `json:"driverKey"`
	DriverVersion      string `json:"driverVersion"`
	KubernetesVersion  string `json:"kubernetesVersion"`
	APIServer          string `json:"apiServer"`
}

type RegisterClusterInput struct {
	Name               string `json:"name"`
	DisplayName        string `json:"displayName"`
	WorkspaceID        uint64 `json:"workspaceId"`
	ProjectID          uint64 `json:"projectId"`
	InfrastructureType string `json:"infrastructureType"`
	DriverKey          string `json:"driverKey"`
	DriverVersion      string `json:"driverVersion"`
	KubernetesVersion  string `json:"kubernetesVersion"`
}

type CreateClusterInput struct {
	Name               string                `json:"name"`
	DisplayName        string                `json:"displayName"`
	WorkspaceID        uint64                `json:"workspaceId"`
	ProjectID          uint64                `json:"projectId"`
	InfrastructureType string                `json:"infrastructureType"`
	DriverKey          string                `json:"driverKey"`
	DriverVersion      string                `json:"driverVersion"`
	KubernetesVersion  string                `json:"kubernetesVersion"`
	TemplateID         uint64                `json:"templateId"`
	Parameters         map[string]any        `json:"parameters"`
	NodePools          []CreateNodePoolInput `json:"nodePools"`
}

type CreateNodePoolInput struct {
	Name         string   `json:"name"`
	Role         string   `json:"role"`
	DesiredCount int      `json:"desiredCount"`
	MinCount     int      `json:"minCount"`
	MaxCount     int      `json:"maxCount"`
	Version      string   `json:"version"`
	ZoneRefs     []string `json:"zoneRefs"`
}

type ValidationInput struct {
	DriverKey            string         `json:"driverKey"`
	DriverVersion        string         `json:"driverVersion"`
	InfrastructureType   string         `json:"infrastructureType"`
	RequiredCapabilities []string       `json:"requiredCapabilities"`
	Parameters           map[string]any `json:"parameters"`
}

type ValidationResult struct {
	Status      string                          `json:"status"`
	CanContinue bool                            `json:"canContinue"`
	Summary     string                          `json:"summary"`
	Checks      []validatorProvider.CheckResult `json:"checks"`
}

type RegistrationBundle struct {
	ClusterID    uint64 `json:"clusterId"`
	Command      string `json:"command"`
	Status       string `json:"status"`
	Instructions string `json:"instructions"`
}

type CreateUpgradePlanInput struct {
	TargetVersion string `json:"targetVersion"`
	WindowStart   string `json:"windowStart"`
	WindowEnd     string `json:"windowEnd"`
	ImpactSummary string `json:"impactSummary"`
}

type ScaleNodePoolInput struct {
	DesiredCount int `json:"desiredCount"`
}

type DisableClusterInput struct {
	Reason string `json:"reason"`
}

type RetireClusterInput struct {
	Reason            string `json:"reason"`
	ConfirmationScope string `json:"confirmationScope"`
	Conclusion        string `json:"conclusion"`
}

type CreateDriverInput struct {
	DriverKey                string                        `json:"driverKey"`
	Version                  string                        `json:"version"`
	DisplayName              string                        `json:"displayName"`
	ProviderType             string                        `json:"providerType"`
	Status                   string                        `json:"status"`
	CapabilityProfileVersion string                        `json:"capabilityProfileVersion"`
	SchemaVersion            string                        `json:"schemaVersion"`
	ReleaseNotes             string                        `json:"releaseNotes"`
	Capabilities             []CreateCapabilityMatrixInput `json:"capabilities"`
}

type CreateCapabilityMatrixInput struct {
	CapabilityDomain    string `json:"capabilityDomain"`
	SupportLevel        string `json:"supportLevel"`
	CompatibilityStatus string `json:"compatibilityStatus"`
	ConstraintsSummary  string `json:"constraintsSummary"`
	RecommendedFor      string `json:"recommendedFor"`
}

type CreateTemplateInput struct {
	Name                 string         `json:"name"`
	Description          string         `json:"description"`
	InfrastructureType   string         `json:"infrastructureType"`
	DriverKey            string         `json:"driverKey"`
	DriverVersionRange   string         `json:"driverVersionRange"`
	RequiredCapabilities []string       `json:"requiredCapabilities"`
	ParameterSchema      map[string]any `json:"parameterSchema"`
	DefaultValues        map[string]any `json:"defaultValues"`
	Status               string         `json:"status"`
	WorkspaceID          uint64         `json:"workspaceId"`
	ProjectID            uint64         `json:"projectId"`
}

type TemplateValidationInput struct {
	InfrastructureType string         `json:"infrastructureType"`
	DriverVersion      string         `json:"driverVersion"`
	Parameters         map[string]any `json:"parameters"`
}

type ClusterLifecycleDetail struct {
	Cluster      *domain.ClusterLifecycleRecord `json:"cluster"`
	NodePools    []domain.NodePoolProfile       `json:"nodePools"`
	UpgradePlans []domain.UpgradePlan           `json:"upgradePlans"`
}

func (s *Service) ListClusters(ctx context.Context, userID uint64, filter ClusterListFilter) ([]domain.ClusterLifecycleRecord, error) {
	workspaceIDs, constrained, err := s.scope.ListReadableWorkspaceIDs(ctx, userID)
	if err != nil {
		return nil, err
	}
	if constrained && len(workspaceIDs) == 0 {
		return []domain.ClusterLifecycleRecord{}, nil
	}
	return s.clusters.List(ctx, repository.ClusterLifecycleListFilter{
		WorkspaceIDs:       workspaceIDs,
		Status:             filter.Status,
		InfrastructureType: filter.InfrastructureType,
		DriverRef:          filter.DriverKey,
		Keyword:            filter.Keyword,
	})
}

func (s *Service) GetCluster(ctx context.Context, userID, clusterID uint64) (*ClusterLifecycleDetail, error) {
	cluster, err := s.clusters.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureReadableCluster(ctx, userID, cluster.WorkspaceID, derefUint64(cluster.ProjectID)); err != nil {
		return nil, err
	}
	nodePools, err := s.nodePools.ListByClusterID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	plans, err := s.plans.ListByClusterID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	return &ClusterLifecycleDetail{Cluster: cluster, NodePools: nodePools, UpgradePlans: plans}, nil
}

func writeAudit(writer *auditSvc.EventWriter, ctx context.Context, actorID uint64, action, resourceID string, outcome domain.AuditOutcome, details map[string]any) {
	if writer == nil {
		return
	}
	actor := actorID
	_ = writer.Write(ctx, "", &actor, action, ResourceTypeClusterLifecycle, resourceID, outcome, details)
}

func marshalPayload(v any) string {
	if v == nil {
		return ""
	}
	payload, err := json.Marshal(v)
	if err != nil {
		return ""
	}
	return string(payload)
}

func parseOptionalRFC3339(value string) (*time.Time, error) {
	value = strings.TrimSpace(value)
	if value == "" {
		return nil, nil
	}
	out, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, fmt.Errorf("%w: invalid RFC3339 time", ErrLifecycleInvalid)
	}
	return &out, nil
}

func normalizeText(value string) string {
	return strings.TrimSpace(value)
}

func derefUint64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func uint64Ptr(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	return &v
}

func isNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func resourceIDForCluster(clusterID uint64) string {
	return strconv.FormatUint(clusterID, 10)
}

func resourceIDForTemplate(templateID uint64) string {
	return strconv.FormatUint(templateID, 10)
}

func (s *Service) ImportCluster(ctx context.Context, userID uint64, input ImportClusterInput) (*domain.LifecycleOperation, *domain.ClusterLifecycleRecord, error) {
	if err := s.scope.EnsureImportCluster(ctx, userID, input.WorkspaceID, input.ProjectID); err != nil {
		return nil, nil, err
	}
	now := time.Now()
	cluster := &domain.ClusterLifecycleRecord{
		Name:                 normalizeText(input.Name),
		DisplayName:          firstNonEmptyText(input.DisplayName, input.Name),
		LifecycleMode:        domain.ClusterLifecycleModeImported,
		InfrastructureType:   normalizeText(input.InfrastructureType),
		DriverRef:            normalizeText(input.DriverKey),
		DriverVersion:        normalizeText(input.DriverVersion),
		WorkspaceID:          input.WorkspaceID,
		ProjectID:            uint64Ptr(input.ProjectID),
		Status:               domain.ClusterLifecycleStatusActive,
		RegistrationStatus:   domain.ClusterRegistrationNotRequired,
		HealthStatus:         domain.ClusterHealthHealthy,
		KubernetesVersion:    firstNonEmptyText(input.KubernetesVersion, "unknown"),
		LastValidationStatus: domain.ValidationStatusPassed,
		LastValidationAt:     &now,
		CreatedBy:            userID,
	}
	result, driverErr := s.driverOps.ImportCluster(ctx, normalizeText(input.APIServer))
	op := &domain.LifecycleOperation{
		OperationType:   domain.LifecycleOperationImport,
		TriggerSource:   domain.LifecycleTriggerManual,
		Status:          domain.LifecycleOperationSucceeded,
		RiskLevel:       domain.LifecycleRiskLow,
		RequestedBy:     userID,
		RequestSnapshot: marshalPayload(input),
		ResultSummary:   firstNonEmptyText(result.Summary, "cluster import accepted"),
		StartedAt:       &now,
		CompletedAt:     &now,
	}
	if driverErr != nil {
		op.Status = domain.LifecycleOperationFailed
		op.FailureReason = driverErr.Error()
		cluster.Status = domain.ClusterLifecycleStatusFailed
		cluster.HealthStatus = domain.ClusterHealthCritical
		cluster.LastValidationStatus = domain.ValidationStatusFailed
	}
	if err := s.clusters.Create(ctx, cluster); err != nil {
		return nil, nil, err
	}
	op.ClusterID = uint64Ptr(cluster.ID)
	if err := s.operations.Create(ctx, op); err != nil {
		return nil, nil, err
	}
	cluster.LastOperationID = uint64Ptr(op.ID)
	_ = s.clusters.Update(ctx, cluster)
	writeAudit(s.auditWriter, ctx, userID, ActionClusterImport, resourceIDForCluster(cluster.ID), outcomeForOperation(op.Status), map[string]any{
		"clusterId": cluster.ID, "workspaceId": cluster.WorkspaceID, "projectId": input.ProjectID, "apiServer": input.APIServer,
	})
	return op, cluster, nil
}

func (s *Service) RegisterCluster(ctx context.Context, userID uint64, input RegisterClusterInput) (*RegistrationBundle, error) {
	if err := s.scope.EnsureImportCluster(ctx, userID, input.WorkspaceID, input.ProjectID); err != nil {
		return nil, err
	}
	now := time.Now()
	cluster := &domain.ClusterLifecycleRecord{
		Name:                 normalizeText(input.Name),
		DisplayName:          firstNonEmptyText(input.DisplayName, input.Name),
		LifecycleMode:        domain.ClusterLifecycleModeRegistered,
		InfrastructureType:   normalizeText(input.InfrastructureType),
		DriverRef:            normalizeText(input.DriverKey),
		DriverVersion:        normalizeText(input.DriverVersion),
		WorkspaceID:          input.WorkspaceID,
		ProjectID:            uint64Ptr(input.ProjectID),
		Status:               domain.ClusterLifecycleStatusPending,
		RegistrationStatus:   domain.ClusterRegistrationIssued,
		HealthStatus:         domain.ClusterHealthUnknown,
		KubernetesVersion:    firstNonEmptyText(input.KubernetesVersion, "unknown"),
		LastValidationStatus: domain.ValidationStatusPending,
		CreatedBy:            userID,
	}
	if err := s.clusters.Create(ctx, cluster); err != nil {
		return nil, err
	}
	command, err := s.driverOps.IssueRegistration(ctx, cluster.Name)
	if err != nil {
		return nil, err
	}
	op := &domain.LifecycleOperation{
		ClusterID:       uint64Ptr(cluster.ID),
		OperationType:   domain.LifecycleOperationRegister,
		TriggerSource:   domain.LifecycleTriggerManual,
		Status:          domain.LifecycleOperationSucceeded,
		RiskLevel:       domain.LifecycleRiskLow,
		RequestedBy:     userID,
		RequestSnapshot: marshalPayload(input),
		ResultSummary:   "registration guide issued",
		StartedAt:       &now,
		CompletedAt:     &now,
	}
	_ = s.operations.Create(ctx, op)
	cluster.LastOperationID = uint64Ptr(op.ID)
	_ = s.clusters.Update(ctx, cluster)
	writeAudit(s.auditWriter, ctx, userID, ActionClusterRegister, resourceIDForCluster(cluster.ID), domain.AuditOutcomeSuccess, map[string]any{
		"clusterId": cluster.ID, "workspaceId": cluster.WorkspaceID, "projectId": input.ProjectID,
	})
	return &RegistrationBundle{
		ClusterID:    cluster.ID,
		Command:      command,
		Status:       string(cluster.RegistrationStatus),
		Instructions: "在目标集群执行注册命令后，返回生命周期中心确认接入状态。",
	}, nil
}

func (s *Service) ValidateChange(ctx context.Context, userID, clusterID uint64, input ValidationInput) (*ValidationResult, error) {
	if clusterID != 0 {
		cluster, err := s.clusters.GetByID(ctx, clusterID)
		if err != nil {
			return nil, err
		}
		if err := s.scope.EnsureReadableCluster(ctx, userID, cluster.WorkspaceID, derefUint64(cluster.ProjectID)); err != nil {
			return nil, err
		}
	}
	result, err := s.validator.Validate(ctx, validatorProvider.Request{
		InfrastructureType: strings.TrimSpace(input.InfrastructureType),
		DriverKey:          strings.TrimSpace(input.DriverKey),
		DriverVersion:      strings.TrimSpace(input.DriverVersion),
		RequiredDomains:    input.RequiredCapabilities,
		Parameters:         input.Parameters,
	})
	if err != nil {
		return nil, err
	}
	writeAudit(s.auditWriter, ctx, userID, ActionClusterValidate, strconv.FormatUint(clusterID, 10), domain.AuditOutcomeSuccess, map[string]any{
		"clusterId": clusterID, "driverKey": input.DriverKey, "status": result.Status,
	})
	response := &ValidationResult{
		Status:      result.Status,
		CanContinue: result.CanContinue,
		Summary:     result.Summary,
		Checks:      result.Checks,
	}
	_ = s.validations.Store(ctx, clusterID, response)
	return response, nil
}

func (s *Service) CreateCluster(ctx context.Context, userID uint64, input CreateClusterInput) (*domain.LifecycleOperation, *domain.ClusterLifecycleRecord, error) {
	if err := s.scope.EnsureCreateCluster(ctx, userID, input.WorkspaceID, input.ProjectID); err != nil {
		return nil, nil, err
	}
	validation, err := s.ValidateChange(ctx, userID, 0, ValidationInput{
		DriverKey:            input.DriverKey,
		DriverVersion:        input.DriverVersion,
		InfrastructureType:   input.InfrastructureType,
		RequiredCapabilities: []string{},
		Parameters:           input.Parameters,
	})
	if err != nil {
		return nil, nil, err
	}
	if !validation.CanContinue {
		return nil, nil, ErrLifecycleBlocked
	}
	now := time.Now()
	cluster := &domain.ClusterLifecycleRecord{
		Name:                 normalizeText(input.Name),
		DisplayName:          firstNonEmptyText(input.DisplayName, input.Name),
		LifecycleMode:        domain.ClusterLifecycleModeProvisioned,
		InfrastructureType:   normalizeText(input.InfrastructureType),
		DriverRef:            normalizeText(input.DriverKey),
		DriverVersion:        normalizeText(input.DriverVersion),
		WorkspaceID:          input.WorkspaceID,
		ProjectID:            uint64Ptr(input.ProjectID),
		Status:               domain.ClusterLifecycleStatusPending,
		RegistrationStatus:   domain.ClusterRegistrationPending,
		HealthStatus:         domain.ClusterHealthUnknown,
		KubernetesVersion:    firstNonEmptyText(input.KubernetesVersion, "unknown"),
		LastValidationStatus: domain.ValidationStatusPassed,
		LastValidationAt:     &now,
		TemplateID:           uint64Ptr(input.TemplateID),
		CreatedBy:            userID,
	}
	if err := s.clusters.Create(ctx, cluster); err != nil {
		return nil, nil, err
	}
	for _, pool := range input.NodePools {
		_ = s.nodePools.Create(ctx, &domain.NodePoolProfile{
			ClusterID:    cluster.ID,
			Name:         normalizeText(pool.Name),
			Role:         domain.NodePoolRole(firstNonEmptyText(pool.Role, string(domain.NodePoolRoleWorker))),
			DesiredCount: pool.DesiredCount,
			CurrentCount: pool.DesiredCount,
			MinCount:     pool.MinCount,
			MaxCount:     pool.MaxCount,
			Version:      firstNonEmptyText(pool.Version, input.KubernetesVersion),
			ZoneRefs:     marshalPayload(pool.ZoneRefs),
			Status:       domain.NodePoolStatusActive,
		})
	}
	result, driverErr := s.driverOps.ProvisionCluster(ctx, driverProvider.ProvisionRequest{
		ClusterName:        cluster.Name,
		InfrastructureType: cluster.InfrastructureType,
		DriverKey:          cluster.DriverRef,
		DriverVersion:      cluster.DriverVersion,
		KubernetesVersion:  cluster.KubernetesVersion,
		Parameters:         input.Parameters,
	})
	op := &domain.LifecycleOperation{
		ClusterID:       uint64Ptr(cluster.ID),
		OperationType:   domain.LifecycleOperationCreate,
		TriggerSource:   domain.LifecycleTriggerManual,
		Status:          domain.LifecycleOperationSucceeded,
		RiskLevel:       domain.LifecycleRiskMedium,
		RequestedBy:     userID,
		RequestSnapshot: marshalPayload(input),
		ResultSummary:   firstNonEmptyText(result.Summary, "cluster provisioned"),
		StartedAt:       &now,
		CompletedAt:     &now,
	}
	if driverErr != nil {
		op.Status = domain.LifecycleOperationFailed
		op.FailureReason = driverErr.Error()
		cluster.Status = domain.ClusterLifecycleStatusFailed
		cluster.HealthStatus = domain.ClusterHealthCritical
		cluster.LastValidationStatus = domain.ValidationStatusFailed
	} else {
		cluster.Status = domain.ClusterLifecycleStatusActive
		cluster.RegistrationStatus = domain.ClusterRegistrationConnected
		cluster.HealthStatus = domain.ClusterHealthHealthy
	}
	if err := s.operations.Create(ctx, op); err != nil {
		return nil, nil, err
	}
	cluster.LastOperationID = uint64Ptr(op.ID)
	_ = s.clusters.Update(ctx, cluster)
	writeAudit(s.auditWriter, ctx, userID, ActionClusterCreate, resourceIDForCluster(cluster.ID), outcomeForOperation(op.Status), map[string]any{
		"clusterId": cluster.ID, "templateId": input.TemplateID, "workspaceId": cluster.WorkspaceID, "projectId": input.ProjectID,
	})
	return op, cluster, nil
}

func (s *Service) CreateUpgradePlan(ctx context.Context, userID, clusterID uint64, input CreateUpgradePlanInput) (*domain.UpgradePlan, error) {
	cluster, err := s.clusters.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureUpgradeCluster(ctx, userID, cluster.WorkspaceID, derefUint64(cluster.ProjectID)); err != nil {
		return nil, err
	}
	startAt, err := parseOptionalRFC3339(input.WindowStart)
	if err != nil {
		return nil, err
	}
	endAt, err := parseOptionalRFC3339(input.WindowEnd)
	if err != nil {
		return nil, err
	}
	item := &domain.UpgradePlan{
		ClusterID:      clusterID,
		FromVersion:    cluster.KubernetesVersion,
		ToVersion:      normalizeText(input.TargetVersion),
		WindowStart:    startAt,
		WindowEnd:      endAt,
		PrecheckStatus: domain.ValidationStatusPassed,
		ImpactSummary:  input.ImpactSummary,
		Status:         domain.UpgradePlanDraft,
		CreatedBy:      userID,
	}
	if err := s.plans.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) ExecuteUpgradePlan(ctx context.Context, userID, clusterID, planID uint64) (*domain.LifecycleOperation, error) {
	cluster, err := s.clusters.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureUpgradeCluster(ctx, userID, cluster.WorkspaceID, derefUint64(cluster.ProjectID)); err != nil {
		return nil, err
	}
	if err := s.ensureClusterNotBusy(ctx, clusterID); err != nil {
		return nil, err
	}
	locked, err := s.lock.Acquire(ctx, clusterID, "upgrade")
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, ErrLifecycleConflict
	}
	defer func() { _ = s.lock.Release(ctx, clusterID) }()
	plan, err := s.plans.GetByID(ctx, planID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	result, driverErr := s.driverOps.UpgradeCluster(ctx, driverProvider.UpgradeRequest{
		ClusterName:   cluster.Name,
		FromVersion:   plan.FromVersion,
		ToVersion:     plan.ToVersion,
		ImpactSummary: plan.ImpactSummary,
	})
	op := &domain.LifecycleOperation{
		ClusterID:       uint64Ptr(clusterID),
		OperationType:   domain.LifecycleOperationUpgrade,
		TriggerSource:   domain.LifecycleTriggerManual,
		Status:          domain.LifecycleOperationSucceeded,
		RiskLevel:       domain.LifecycleRiskHigh,
		RequestedBy:     userID,
		RequestSnapshot: marshalPayload(plan),
		ResultSummary:   firstNonEmptyText(result.Summary, "cluster upgraded"),
		StartedAt:       &now,
		CompletedAt:     &now,
	}
	if driverErr != nil {
		op.Status = domain.LifecycleOperationFailed
		op.FailureReason = driverErr.Error()
		plan.Status = domain.UpgradePlanFailed
		cluster.Status = domain.ClusterLifecycleStatusFailed
	} else {
		plan.Status = domain.UpgradePlanSucceeded
		cluster.KubernetesVersion = plan.ToVersion
		cluster.Status = domain.ClusterLifecycleStatusActive
	}
	if err := s.operations.Create(ctx, op); err != nil {
		return nil, err
	}
	_ = s.progress.SetOperation(ctx, clusterID, string(domain.LifecycleOperationUpgrade), string(op.Status))
	plan.LastOperationID = uint64Ptr(op.ID)
	cluster.LastOperationID = uint64Ptr(op.ID)
	cluster.TargetVersion = plan.ToVersion
	_ = s.plans.Update(ctx, plan)
	_ = s.clusters.Update(ctx, cluster)
	writeAudit(s.auditWriter, ctx, userID, ActionClusterUpgrade, resourceIDForCluster(clusterID), outcomeForOperation(op.Status), map[string]any{"clusterId": clusterID, "planId": planID, "toVersion": plan.ToVersion})
	return op, nil
}

func (s *Service) ListNodePools(ctx context.Context, userID, clusterID uint64) ([]domain.NodePoolProfile, error) {
	cluster, err := s.clusters.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureReadableCluster(ctx, userID, cluster.WorkspaceID, derefUint64(cluster.ProjectID)); err != nil {
		return nil, err
	}
	return s.nodePools.ListByClusterID(ctx, clusterID)
}

func (s *Service) ScaleNodePool(ctx context.Context, userID, clusterID, nodePoolID uint64, input ScaleNodePoolInput) (*domain.LifecycleOperation, error) {
	cluster, err := s.clusters.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureManageNodePool(ctx, userID, cluster.WorkspaceID, derefUint64(cluster.ProjectID)); err != nil {
		return nil, err
	}
	if err := s.ensureClusterNotBusy(ctx, clusterID); err != nil {
		return nil, err
	}
	locked, err := s.lock.Acquire(ctx, clusterID, "node-pool-scale")
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, ErrLifecycleConflict
	}
	defer func() { _ = s.lock.Release(ctx, clusterID) }()
	nodePool, err := s.nodePools.GetByID(ctx, nodePoolID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	result, driverErr := s.driverOps.ScaleNodePool(ctx, driverProvider.NodePoolScaleRequest{
		ClusterName:  cluster.Name,
		NodePoolName: nodePool.Name,
		DesiredCount: input.DesiredCount,
	})
	op := &domain.LifecycleOperation{
		ClusterID:       uint64Ptr(clusterID),
		OperationType:   domain.LifecycleOperationScaleNodePool,
		TriggerSource:   domain.LifecycleTriggerManual,
		Status:          domain.LifecycleOperationSucceeded,
		RiskLevel:       domain.LifecycleRiskMedium,
		RequestedBy:     userID,
		RequestSnapshot: marshalPayload(input),
		ResultSummary:   firstNonEmptyText(result.Summary, "node pool scaled"),
		StartedAt:       &now,
		CompletedAt:     &now,
	}
	if driverErr != nil {
		op.Status = domain.LifecycleOperationFailed
		op.FailureReason = driverErr.Error()
		nodePool.Status = domain.NodePoolStatusFailed
	} else {
		nodePool.DesiredCount = input.DesiredCount
		nodePool.CurrentCount = input.DesiredCount
		nodePool.Status = domain.NodePoolStatusActive
	}
	if err := s.operations.Create(ctx, op); err != nil {
		return nil, err
	}
	_ = s.progress.SetOperation(ctx, clusterID, string(domain.LifecycleOperationScaleNodePool), string(op.Status))
	nodePool.LastOperationID = uint64Ptr(op.ID)
	_ = s.nodePools.Update(ctx, nodePool)
	writeAudit(s.auditWriter, ctx, userID, ActionClusterScaleNodePool, resourceIDForCluster(clusterID), outcomeForOperation(op.Status), map[string]any{"clusterId": clusterID, "nodePoolId": nodePoolID})
	return op, nil
}

func (s *Service) DisableCluster(ctx context.Context, userID, clusterID uint64, input DisableClusterInput) (*domain.LifecycleOperation, error) {
	cluster, err := s.clusters.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureRetireCluster(ctx, userID, cluster.WorkspaceID, derefUint64(cluster.ProjectID)); err != nil {
		return nil, err
	}
	if err := s.ensureClusterNotBusy(ctx, clusterID); err != nil {
		return nil, err
	}
	locked, err := s.lock.Acquire(ctx, clusterID, "disable")
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, ErrLifecycleConflict
	}
	defer func() { _ = s.lock.Release(ctx, clusterID) }()
	now := time.Now()
	result, driverErr := s.driverOps.DisableCluster(ctx, cluster.Name)
	op := &domain.LifecycleOperation{
		ClusterID:       uint64Ptr(clusterID),
		OperationType:   domain.LifecycleOperationDisable,
		TriggerSource:   domain.LifecycleTriggerManual,
		Status:          domain.LifecycleOperationSucceeded,
		RiskLevel:       domain.LifecycleRiskHigh,
		RequestedBy:     userID,
		RequestSnapshot: marshalPayload(input),
		ResultSummary:   firstNonEmptyText(result.Summary, "cluster disabled"),
		StartedAt:       &now,
		CompletedAt:     &now,
	}
	if driverErr != nil {
		op.Status = domain.LifecycleOperationFailed
		op.FailureReason = driverErr.Error()
	} else {
		cluster.Status = domain.ClusterLifecycleStatusDisabled
		cluster.RetirementReason = input.Reason
	}
	if err := s.operations.Create(ctx, op); err != nil {
		return nil, err
	}
	_ = s.progress.SetOperation(ctx, clusterID, string(domain.LifecycleOperationDisable), string(op.Status))
	cluster.LastOperationID = uint64Ptr(op.ID)
	_ = s.clusters.Update(ctx, cluster)
	writeAudit(s.auditWriter, ctx, userID, ActionClusterDisable, resourceIDForCluster(clusterID), outcomeForOperation(op.Status), map[string]any{"clusterId": clusterID, "reason": input.Reason})
	return op, nil
}

func (s *Service) RetireCluster(ctx context.Context, userID, clusterID uint64, input RetireClusterInput) (*domain.LifecycleOperation, error) {
	cluster, err := s.clusters.GetByID(ctx, clusterID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureRetireCluster(ctx, userID, cluster.WorkspaceID, derefUint64(cluster.ProjectID)); err != nil {
		return nil, err
	}
	if err := s.ensureClusterNotBusy(ctx, clusterID); err != nil {
		return nil, err
	}
	locked, err := s.lock.Acquire(ctx, clusterID, "retire")
	if err != nil {
		return nil, err
	}
	if !locked {
		return nil, ErrLifecycleConflict
	}
	defer func() { _ = s.lock.Release(ctx, clusterID) }()
	now := time.Now()
	result, driverErr := s.driverOps.RetireCluster(ctx, driverProvider.RetirementRequest{
		ClusterName: cluster.Name,
		Reason:      input.Reason,
	})
	op := &domain.LifecycleOperation{
		ClusterID:       uint64Ptr(clusterID),
		OperationType:   domain.LifecycleOperationRetire,
		TriggerSource:   domain.LifecycleTriggerManual,
		Status:          domain.LifecycleOperationSucceeded,
		RiskLevel:       domain.LifecycleRiskCritical,
		RequestedBy:     userID,
		RequestSnapshot: marshalPayload(input),
		ResultSummary:   firstNonEmptyText(result.Summary, "cluster retired"),
		StartedAt:       &now,
		CompletedAt:     &now,
	}
	if driverErr != nil {
		op.Status = domain.LifecycleOperationFailed
		op.FailureReason = driverErr.Error()
		cluster.Status = domain.ClusterLifecycleStatusFailed
	} else {
		cluster.Status = domain.ClusterLifecycleStatusRetired
		cluster.HealthStatus = domain.ClusterHealthUnknown
		cluster.RetirementReason = input.Reason
	}
	if err := s.operations.Create(ctx, op); err != nil {
		return nil, err
	}
	_ = s.progress.SetOperation(ctx, clusterID, string(domain.LifecycleOperationRetire), string(op.Status))
	cluster.LastOperationID = uint64Ptr(op.ID)
	_ = s.clusters.Update(ctx, cluster)
	writeAudit(s.auditWriter, ctx, userID, ActionClusterRetire, resourceIDForCluster(clusterID), outcomeForOperation(op.Status), map[string]any{"clusterId": clusterID, "reason": input.Reason, "conclusion": input.Conclusion})
	return op, nil
}

func (s *Service) ListDrivers(ctx context.Context, userID uint64, providerType string) ([]domain.ClusterDriverVersion, error) {
	if err := s.scope.EnsureManageDriver(ctx, userID, 0, 0); err != nil {
		if err := s.scope.EnsureReadableCluster(ctx, userID, 0, 0); err != nil {
			return nil, err
		}
	}
	return s.drivers.List(ctx, providerType)
}

func (s *Service) UpsertDriver(ctx context.Context, userID uint64, input CreateDriverInput) (*domain.ClusterDriverVersion, error) {
	if err := s.scope.EnsureManageDriver(ctx, userID, 0, 0); err != nil {
		return nil, err
	}
	item, err := s.drivers.FindByKeyVersion(ctx, input.DriverKey, input.Version)
	if err != nil && !isNotFound(err) {
		return nil, err
	}
	if item == nil {
		item = &domain.ClusterDriverVersion{}
	}
	item.DriverKey = normalizeText(input.DriverKey)
	item.Version = normalizeText(input.Version)
	item.DisplayName = firstNonEmptyText(input.DisplayName, input.DriverKey)
	item.ProviderType = normalizeText(input.ProviderType)
	item.Status = domain.DriverStatus(firstNonEmptyText(input.Status, string(domain.DriverStatusActive)))
	item.CapabilityProfileVersion = firstNonEmptyText(input.CapabilityProfileVersion, input.Version)
	item.SchemaVersion = firstNonEmptyText(input.SchemaVersion, "v1")
	item.ReleaseNotes = input.ReleaseNotes
	if item.ID == 0 {
		err = s.drivers.Create(ctx, item)
	} else {
		err = s.drivers.Update(ctx, item)
	}
	if err != nil {
		return nil, err
	}
	entries := make([]domain.CapabilityMatrixEntry, 0, len(input.Capabilities))
	for _, cap := range input.Capabilities {
		entries = append(entries, domain.CapabilityMatrixEntry{
			CapabilityDomain:    cap.CapabilityDomain,
			SupportLevel:        domain.CapabilitySupportLevel(cap.SupportLevel),
			CompatibilityStatus: domain.CapabilityCompatibilityStatus(cap.CompatibilityStatus),
			ConstraintsSummary:  cap.ConstraintsSummary,
			RecommendedFor:      cap.RecommendedFor,
		})
	}
	_ = s.capability.ReplaceForOwner(ctx, domain.CapabilityOwnerDriver, fmt.Sprintf("%s:%s", item.DriverKey, item.Version), entries)
	writeAudit(s.auditWriter, ctx, userID, ActionDriverCreate, strconv.FormatUint(item.ID, 10), domain.AuditOutcomeSuccess, map[string]any{"driverKey": item.DriverKey, "version": item.Version})
	return item, nil
}

func (s *Service) ListCapabilities(ctx context.Context, userID, driverID uint64) ([]domain.CapabilityMatrixEntry, error) {
	if err := s.scope.EnsureManageDriver(ctx, userID, 0, 0); err != nil {
		return nil, err
	}
	driver, err := s.drivers.GetByID(ctx, driverID)
	if err != nil {
		return nil, err
	}
	return s.capability.ListByOwner(ctx, domain.CapabilityOwnerDriver, fmt.Sprintf("%s:%s", driver.DriverKey, driver.Version))
}

func (s *Service) ListTemplates(ctx context.Context, userID uint64, driverKey, infrastructureType string) ([]domain.ClusterTemplate, error) {
	if err := s.scope.EnsureReadableCluster(ctx, userID, 0, 0); err != nil && s.scope.EnsureManageDriver(ctx, userID, 0, 0) != nil {
		return nil, err
	}
	return s.templates.List(ctx, driverKey, infrastructureType)
}

func (s *Service) CreateTemplate(ctx context.Context, userID uint64, input CreateTemplateInput) (*domain.ClusterTemplate, error) {
	if err := s.scope.EnsureManageDriver(ctx, userID, input.WorkspaceID, input.ProjectID); err != nil {
		return nil, err
	}
	item := &domain.ClusterTemplate{
		Name:                 normalizeText(input.Name),
		Description:          input.Description,
		InfrastructureType:   normalizeText(input.InfrastructureType),
		DriverKey:            normalizeText(input.DriverKey),
		DriverVersionRange:   normalizeText(input.DriverVersionRange),
		RequiredCapabilities: marshalPayload(input.RequiredCapabilities),
		ParameterSchema:      marshalPayload(input.ParameterSchema),
		DefaultValues:        marshalPayload(input.DefaultValues),
		Status:               domain.TemplateStatus(firstNonEmptyText(input.Status, string(domain.TemplateStatusActive))),
		CreatedBy:            userID,
	}
	if err := s.templates.Create(ctx, item); err != nil {
		return nil, err
	}
	writeAudit(s.auditWriter, ctx, userID, ActionTemplateCreate, resourceIDForTemplate(item.ID), domain.AuditOutcomeSuccess, map[string]any{"templateId": item.ID, "driverKey": item.DriverKey})
	return item, nil
}

func (s *Service) ValidateTemplate(ctx context.Context, userID, templateID uint64, input TemplateValidationInput) (*ValidationResult, error) {
	if err := s.scope.EnsureReadableCluster(ctx, userID, 0, 0); err != nil && s.scope.EnsureManageDriver(ctx, userID, 0, 0) != nil {
		return nil, err
	}
	template, err := s.templates.GetByID(ctx, templateID)
	if err != nil {
		return nil, err
	}
	result, err := s.validator.Validate(ctx, validatorProvider.Request{
		InfrastructureType: firstNonEmptyText(input.InfrastructureType, template.InfrastructureType),
		DriverKey:          template.DriverKey,
		DriverVersion:      input.DriverVersion,
		RequiredDomains:    unmarshalStringSlice(template.RequiredCapabilities),
		Parameters:         input.Parameters,
	})
	if err != nil {
		return nil, err
	}
	return &ValidationResult{Status: result.Status, CanContinue: result.CanContinue, Summary: result.Summary, Checks: result.Checks}, nil
}

func firstNonEmptyText(values ...string) string {
	for _, value := range values {
		if trimmed := normalizeText(value); trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func unmarshalStringSlice(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}
	var out []string
	if err := json.Unmarshal([]byte(raw), &out); err != nil {
		return []string{}
	}
	return out
}

func outcomeForOperation(status domain.LifecycleOperationStatus) domain.AuditOutcome {
	switch status {
	case domain.LifecycleOperationSucceeded:
		return domain.AuditOutcomeSuccess
	case domain.LifecycleOperationBlocked:
		return domain.AuditOutcomeDenied
	default:
		return domain.AuditOutcomeFailed
	}
}

func (s *Service) ensureClusterNotBusy(ctx context.Context, clusterID uint64) error {
	if s == nil || s.operations == nil || clusterID == 0 {
		return nil
	}
	_, err := s.operations.FindRunningByClusterID(ctx, clusterID)
	if err == nil {
		return ErrLifecycleConflict
	}
	if isNotFound(err) {
		return nil
	}
	return err
}
