package domain

import "time"

type ClusterLifecycleMode string

const (
	ClusterLifecycleModeImported    ClusterLifecycleMode = "imported"
	ClusterLifecycleModeRegistered  ClusterLifecycleMode = "registered"
	ClusterLifecycleModeProvisioned ClusterLifecycleMode = "provisioned"
)

type ClusterLifecycleStatus string

const (
	ClusterLifecycleStatusPending   ClusterLifecycleStatus = "pending"
	ClusterLifecycleStatusActive    ClusterLifecycleStatus = "active"
	ClusterLifecycleStatusDegraded  ClusterLifecycleStatus = "degraded"
	ClusterLifecycleStatusUpgrading ClusterLifecycleStatus = "upgrading"
	ClusterLifecycleStatusDisabled  ClusterLifecycleStatus = "disabled"
	ClusterLifecycleStatusRetiring  ClusterLifecycleStatus = "retiring"
	ClusterLifecycleStatusRetired   ClusterLifecycleStatus = "retired"
	ClusterLifecycleStatusFailed    ClusterLifecycleStatus = "failed"
)

type ClusterRegistrationStatus string

const (
	ClusterRegistrationNotRequired ClusterRegistrationStatus = "not_required"
	ClusterRegistrationPending     ClusterRegistrationStatus = "pending"
	ClusterRegistrationIssued      ClusterRegistrationStatus = "issued"
	ClusterRegistrationConnected   ClusterRegistrationStatus = "connected"
	ClusterRegistrationFailed      ClusterRegistrationStatus = "failed"
)

type ClusterHealthStatus string

const (
	ClusterHealthHealthy  ClusterHealthStatus = "healthy"
	ClusterHealthWarning  ClusterHealthStatus = "warning"
	ClusterHealthCritical ClusterHealthStatus = "critical"
	ClusterHealthUnknown  ClusterHealthStatus = "unknown"
)

type ValidationStatus string

const (
	ValidationStatusPending ValidationStatus = "pending"
	ValidationStatusPassed  ValidationStatus = "passed"
	ValidationStatusWarning ValidationStatus = "warning"
	ValidationStatusFailed  ValidationStatus = "failed"
)

type DriverStatus string

const (
	DriverStatusDraft      DriverStatus = "draft"
	DriverStatusActive     DriverStatus = "active"
	DriverStatusDeprecated DriverStatus = "deprecated"
	DriverStatusDisabled   DriverStatus = "disabled"
)

type CapabilityOwnerType string

const (
	CapabilityOwnerDriver      CapabilityOwnerType = "driver"
	CapabilityOwnerClusterType CapabilityOwnerType = "cluster-type"
)

type CapabilitySupportLevel string

const (
	CapabilitySupportNative      CapabilitySupportLevel = "native"
	CapabilitySupportExtended    CapabilitySupportLevel = "extended"
	CapabilitySupportPartial     CapabilitySupportLevel = "partial"
	CapabilitySupportUnsupported CapabilitySupportLevel = "unsupported"
)

type CapabilityCompatibilityStatus string

const (
	CapabilityCompatibilityCompatible   CapabilityCompatibilityStatus = "compatible"
	CapabilityCompatibilityConditional  CapabilityCompatibilityStatus = "conditional"
	CapabilityCompatibilityIncompatible CapabilityCompatibilityStatus = "incompatible"
)

type TemplateStatus string

const (
	TemplateStatusDraft      TemplateStatus = "draft"
	TemplateStatusActive     TemplateStatus = "active"
	TemplateStatusDeprecated TemplateStatus = "deprecated"
	TemplateStatusDisabled   TemplateStatus = "disabled"
)

type LifecycleOperationType string

const (
	LifecycleOperationImport        LifecycleOperationType = "import"
	LifecycleOperationRegister      LifecycleOperationType = "register"
	LifecycleOperationCreate        LifecycleOperationType = "create"
	LifecycleOperationValidate      LifecycleOperationType = "validate"
	LifecycleOperationUpgrade       LifecycleOperationType = "upgrade"
	LifecycleOperationScaleNodePool LifecycleOperationType = "scale-node-pool"
	LifecycleOperationDisable       LifecycleOperationType = "disable"
	LifecycleOperationRetire        LifecycleOperationType = "retire"
)

type LifecycleTriggerSource string

const (
	LifecycleTriggerManual    LifecycleTriggerSource = "manual"
	LifecycleTriggerScheduled LifecycleTriggerSource = "scheduled"
	LifecycleTriggerFollowUp  LifecycleTriggerSource = "follow-up"
)

type LifecycleOperationStatus string

const (
	LifecycleOperationPending            LifecycleOperationStatus = "pending"
	LifecycleOperationRunning            LifecycleOperationStatus = "running"
	LifecycleOperationPartiallySucceeded LifecycleOperationStatus = "partially_succeeded"
	LifecycleOperationSucceeded          LifecycleOperationStatus = "succeeded"
	LifecycleOperationFailed             LifecycleOperationStatus = "failed"
	LifecycleOperationCanceled           LifecycleOperationStatus = "canceled"
	LifecycleOperationBlocked            LifecycleOperationStatus = "blocked"
)

type LifecycleRiskLevel string

const (
	LifecycleRiskLow      LifecycleRiskLevel = "low"
	LifecycleRiskMedium   LifecycleRiskLevel = "medium"
	LifecycleRiskHigh     LifecycleRiskLevel = "high"
	LifecycleRiskCritical LifecycleRiskLevel = "critical"
)

type UpgradePlanStatus string

const (
	UpgradePlanDraft     UpgradePlanStatus = "draft"
	UpgradePlanApproved  UpgradePlanStatus = "approved"
	UpgradePlanRunning   UpgradePlanStatus = "running"
	UpgradePlanSucceeded UpgradePlanStatus = "succeeded"
	UpgradePlanFailed    UpgradePlanStatus = "failed"
	UpgradePlanCanceled  UpgradePlanStatus = "canceled"
)

type NodePoolRole string

const (
	NodePoolRoleControlPlane NodePoolRole = "control-plane"
	NodePoolRoleWorker       NodePoolRole = "worker"
	NodePoolRoleMixed        NodePoolRole = "mixed"
)

type NodePoolStatus string

const (
	NodePoolStatusPending   NodePoolStatus = "pending"
	NodePoolStatusActive    NodePoolStatus = "active"
	NodePoolStatusScaling   NodePoolStatus = "scaling"
	NodePoolStatusUpgrading NodePoolStatus = "upgrading"
	NodePoolStatusDegraded  NodePoolStatus = "degraded"
	NodePoolStatusFailed    NodePoolStatus = "failed"
)

type LifecycleAuditOutcome string

const (
	LifecycleAuditSucceeded LifecycleAuditOutcome = "succeeded"
	LifecycleAuditFailed    LifecycleAuditOutcome = "failed"
	LifecycleAuditBlocked   LifecycleAuditOutcome = "blocked"
	LifecycleAuditCanceled  LifecycleAuditOutcome = "canceled"
)

type ClusterLifecycleRecord struct {
	ID                   uint64                    `gorm:"primaryKey" json:"id"`
	Name                 string                    `gorm:"size:128;not null;index:idx_cluster_lifecycle_scope_name,priority:3" json:"name"`
	DisplayName          string                    `gorm:"size:128;not null" json:"displayName"`
	LifecycleMode        ClusterLifecycleMode      `gorm:"size:32;not null" json:"lifecycleMode"`
	InfrastructureType   string                    `gorm:"size:64;not null;index" json:"infrastructureType"`
	DriverRef            string                    `gorm:"size:128;not null;index" json:"driverRef"`
	DriverVersion        string                    `gorm:"size:64;not null" json:"driverVersion"`
	WorkspaceID          uint64                    `gorm:"not null;index;index:idx_cluster_lifecycle_scope_name,priority:1" json:"workspaceId"`
	ProjectID            *uint64                   `gorm:"index;index:idx_cluster_lifecycle_scope_name,priority:2" json:"projectId,omitempty"`
	Status               ClusterLifecycleStatus    `gorm:"size:32;not null;index" json:"status"`
	RegistrationStatus   ClusterRegistrationStatus `gorm:"size:32;not null" json:"registrationStatus"`
	HealthStatus         ClusterHealthStatus       `gorm:"size:32;not null" json:"healthStatus"`
	KubernetesVersion    string                    `gorm:"size:64;not null" json:"kubernetesVersion"`
	TargetVersion        string                    `gorm:"size:64" json:"targetVersion,omitempty"`
	NodePoolSummary      string                    `gorm:"type:text" json:"nodePoolSummary,omitempty"`
	LastValidationStatus ValidationStatus          `gorm:"size:32;not null" json:"lastValidationStatus"`
	LastValidationAt     *time.Time                `json:"lastValidationAt,omitempty"`
	LastOperationID      *uint64                   `gorm:"index" json:"lastOperationId,omitempty"`
	TemplateID           *uint64                   `gorm:"index" json:"templateId,omitempty"`
	RetirementReason     string                    `gorm:"type:text" json:"retirementReason,omitempty"`
	CreatedBy            uint64                    `gorm:"not null;index" json:"createdBy"`
	CreatedAt            time.Time                 `json:"createdAt"`
	UpdatedAt            time.Time                 `json:"updatedAt"`
}

func (ClusterLifecycleRecord) TableName() string {
	return "cluster_lifecycle_records"
}

type ClusterDriverVersion struct {
	ID                       uint64       `gorm:"primaryKey" json:"id"`
	DriverKey                string       `gorm:"size:64;not null;index:idx_cluster_driver_version_unique,unique,priority:1" json:"driverKey"`
	Version                  string       `gorm:"size:64;not null;index:idx_cluster_driver_version_unique,unique,priority:2" json:"version"`
	DisplayName              string       `gorm:"size:128;not null" json:"displayName"`
	ProviderType             string       `gorm:"size:64;not null;index" json:"providerType"`
	Status                   DriverStatus `gorm:"size:32;not null;index" json:"status"`
	CapabilityProfileVersion string       `gorm:"size:64;not null" json:"capabilityProfileVersion"`
	SchemaVersion            string       `gorm:"size:64;not null" json:"schemaVersion"`
	ReleaseNotes             string       `gorm:"type:text" json:"releaseNotes,omitempty"`
	CreatedAt                time.Time    `json:"createdAt"`
	UpdatedAt                time.Time    `json:"updatedAt"`
}

func (ClusterDriverVersion) TableName() string {
	return "cluster_driver_versions"
}

type CapabilityMatrixEntry struct {
	ID                  uint64                        `gorm:"primaryKey" json:"id"`
	OwnerType           CapabilityOwnerType           `gorm:"size:32;not null;index:idx_cluster_capability_unique,unique,priority:1" json:"ownerType"`
	OwnerRef            string                        `gorm:"size:128;not null;index:idx_cluster_capability_unique,unique,priority:2;index" json:"ownerRef"`
	CapabilityDomain    string                        `gorm:"size:64;not null;index:idx_cluster_capability_unique,unique,priority:3" json:"capabilityDomain"`
	SupportLevel        CapabilitySupportLevel        `gorm:"size:32;not null" json:"supportLevel"`
	CompatibilityStatus CapabilityCompatibilityStatus `gorm:"size:32;not null" json:"compatibilityStatus"`
	ConstraintsSummary  string                        `gorm:"type:text" json:"constraintsSummary,omitempty"`
	RecommendedFor      string                        `gorm:"type:text" json:"recommendedFor,omitempty"`
	UpdatedAt           time.Time                     `json:"updatedAt"`
}

func (CapabilityMatrixEntry) TableName() string {
	return "cluster_capability_matrix_entries"
}

type ClusterTemplate struct {
	ID                   uint64         `gorm:"primaryKey" json:"id"`
	Name                 string         `gorm:"size:128;not null;uniqueIndex" json:"name"`
	Description          string         `gorm:"type:text" json:"description,omitempty"`
	InfrastructureType   string         `gorm:"size:64;not null;index" json:"infrastructureType"`
	DriverKey            string         `gorm:"size:64;not null;index" json:"driverKey"`
	DriverVersionRange   string         `gorm:"size:128;not null" json:"driverVersionRange"`
	RequiredCapabilities string         `gorm:"type:text" json:"requiredCapabilities"`
	ParameterSchema      string         `gorm:"type:text" json:"parameterSchema"`
	DefaultValues        string         `gorm:"type:text" json:"defaultValues"`
	Status               TemplateStatus `gorm:"size:32;not null;index" json:"status"`
	CreatedBy            uint64         `gorm:"not null;index" json:"createdBy"`
	CreatedAt            time.Time      `json:"createdAt"`
	UpdatedAt            time.Time      `json:"updatedAt"`
}

func (ClusterTemplate) TableName() string {
	return "cluster_templates"
}

type LifecycleOperation struct {
	ID              uint64                   `gorm:"primaryKey" json:"id"`
	ClusterID       *uint64                  `gorm:"index" json:"clusterId,omitempty"`
	OperationType   LifecycleOperationType   `gorm:"size:32;not null;index" json:"operationType"`
	TriggerSource   LifecycleTriggerSource   `gorm:"size:32;not null" json:"triggerSource"`
	Status          LifecycleOperationStatus `gorm:"size:32;not null;index" json:"status"`
	RiskLevel       LifecycleRiskLevel       `gorm:"size:32;not null" json:"riskLevel"`
	RequestedBy     uint64                   `gorm:"not null;index" json:"requestedBy"`
	RequestSnapshot string                   `gorm:"type:text" json:"requestSnapshot,omitempty"`
	ResultSummary   string                   `gorm:"type:text" json:"resultSummary,omitempty"`
	FailureReason   string                   `gorm:"type:text" json:"failureReason,omitempty"`
	StartedAt       *time.Time               `json:"startedAt,omitempty"`
	CompletedAt     *time.Time               `json:"completedAt,omitempty"`
	CreatedAt       time.Time                `json:"createdAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
}

func (LifecycleOperation) TableName() string {
	return "cluster_lifecycle_operations"
}

type UpgradePlan struct {
	ID              uint64            `gorm:"primaryKey" json:"id"`
	ClusterID       uint64            `gorm:"not null;index" json:"clusterId"`
	FromVersion     string            `gorm:"size:64;not null" json:"fromVersion"`
	ToVersion       string            `gorm:"size:64;not null" json:"toVersion"`
	WindowStart     *time.Time        `json:"windowStart,omitempty"`
	WindowEnd       *time.Time        `json:"windowEnd,omitempty"`
	PrecheckStatus  ValidationStatus  `gorm:"size:32;not null" json:"precheckStatus"`
	ImpactSummary   string            `gorm:"type:text" json:"impactSummary,omitempty"`
	Status          UpgradePlanStatus `gorm:"size:32;not null;index" json:"status"`
	LastOperationID *uint64           `gorm:"index" json:"lastOperationId,omitempty"`
	CreatedBy       uint64            `gorm:"not null;index" json:"createdBy"`
	CreatedAt       time.Time         `json:"createdAt"`
	UpdatedAt       time.Time         `json:"updatedAt"`
}

func (UpgradePlan) TableName() string {
	return "cluster_upgrade_plans"
}

type NodePoolProfile struct {
	ID              uint64         `gorm:"primaryKey" json:"id"`
	ClusterID       uint64         `gorm:"not null;index" json:"clusterId"`
	Name            string         `gorm:"size:128;not null" json:"name"`
	Role            NodePoolRole   `gorm:"size:32;not null" json:"role"`
	DesiredCount    int            `gorm:"not null" json:"desiredCount"`
	CurrentCount    int            `gorm:"not null" json:"currentCount"`
	MinCount        int            `gorm:"not null" json:"minCount"`
	MaxCount        int            `gorm:"not null" json:"maxCount"`
	Version         string         `gorm:"size:64;not null" json:"version"`
	ZoneRefs        string         `gorm:"type:text" json:"zoneRefs,omitempty"`
	Status          NodePoolStatus `gorm:"size:32;not null;index" json:"status"`
	LastOperationID *uint64        `gorm:"index" json:"lastOperationId,omitempty"`
	UpdatedAt       time.Time      `json:"updatedAt"`
	CreatedAt       time.Time      `json:"createdAt"`
}

func (NodePoolProfile) TableName() string {
	return "cluster_node_pool_profiles"
}

type LifecycleAuditEvent struct {
	ID             uint64                `gorm:"primaryKey" json:"id"`
	Action         string                `gorm:"size:128;not null;index" json:"action"`
	ActorUserID    uint64                `gorm:"not null;index" json:"actorUserId"`
	ClusterID      *uint64               `gorm:"index" json:"clusterId,omitempty"`
	TargetType     string                `gorm:"size:32;not null;index" json:"targetType"`
	TargetRef      string                `gorm:"size:128;not null;index" json:"targetRef"`
	Outcome        LifecycleAuditOutcome `gorm:"size:32;not null;index" json:"outcome"`
	DetailSnapshot string                `gorm:"type:text" json:"detailSnapshot,omitempty"`
	OccurredAt     time.Time             `gorm:"not null;index" json:"occurredAt"`
}

func (LifecycleAuditEvent) TableName() string {
	return "cluster_lifecycle_audit_events"
}
