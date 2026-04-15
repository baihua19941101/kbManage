package domain

import "time"

type PolicyScopeLevel string

type PolicyCategory string

type PolicyEnforcementMode string

type PolicyStatus string

type PolicyRiskLevel string

type PolicyRolloutStage string

type PolicyAssignmentStatus string

type PolicyDistributionOperation string

type PolicyDistributionStatus string

type PolicyHitResult string

type PolicyRemediationStatus string

type PolicyExceptionStatus string

const (
	PolicyScopeLevelPlatform  PolicyScopeLevel = "platform"
	PolicyScopeLevelWorkspace PolicyScopeLevel = "workspace"
	PolicyScopeLevelProject   PolicyScopeLevel = "project"

	PolicyCategoryPodSecurity PolicyCategory = "pod-security"
	PolicyCategoryImage       PolicyCategory = "image"
	PolicyCategoryResource    PolicyCategory = "resource"
	PolicyCategoryLabel       PolicyCategory = "label"
	PolicyCategoryNetwork     PolicyCategory = "network"
	PolicyCategoryAdmission   PolicyCategory = "admission"

	PolicyEnforcementModeAudit   PolicyEnforcementMode = "audit"
	PolicyEnforcementModeAlert   PolicyEnforcementMode = "alert"
	PolicyEnforcementModeWarn    PolicyEnforcementMode = "warn"
	PolicyEnforcementModeEnforce PolicyEnforcementMode = "enforce"

	PolicyStatusDraft    PolicyStatus = "draft"
	PolicyStatusActive   PolicyStatus = "active"
	PolicyStatusDisabled PolicyStatus = "disabled"
	PolicyStatusArchived PolicyStatus = "archived"

	PolicyRiskLevelLow      PolicyRiskLevel = "low"
	PolicyRiskLevelMedium   PolicyRiskLevel = "medium"
	PolicyRiskLevelHigh     PolicyRiskLevel = "high"
	PolicyRiskLevelCritical PolicyRiskLevel = "critical"

	PolicyRolloutStagePilot  PolicyRolloutStage = "pilot"
	PolicyRolloutStageCanary PolicyRolloutStage = "canary"
	PolicyRolloutStageBroad  PolicyRolloutStage = "broad"
	PolicyRolloutStageFull   PolicyRolloutStage = "full"

	PolicyAssignmentStatusPending PolicyAssignmentStatus = "pending"
	PolicyAssignmentStatusActive  PolicyAssignmentStatus = "active"
	PolicyAssignmentStatusFailed  PolicyAssignmentStatus = "failed"
	PolicyAssignmentStatusPaused  PolicyAssignmentStatus = "paused"

	PolicyDistributionOperationAssign     PolicyDistributionOperation = "assign"
	PolicyDistributionOperationModeSwitch PolicyDistributionOperation = "mode-switch"
	PolicyDistributionOperationPause      PolicyDistributionOperation = "pause"
	PolicyDistributionOperationResume     PolicyDistributionOperation = "resume"
	PolicyDistributionOperationRevoke     PolicyDistributionOperation = "revoke"

	PolicyDistributionStatusPending            PolicyDistributionStatus = "pending"
	PolicyDistributionStatusRunning            PolicyDistributionStatus = "running"
	PolicyDistributionStatusPartiallySucceeded PolicyDistributionStatus = "partially_succeeded"
	PolicyDistributionStatusSucceeded          PolicyDistributionStatus = "succeeded"
	PolicyDistributionStatusFailed             PolicyDistributionStatus = "failed"

	PolicyHitResultPass  PolicyHitResult = "pass"
	PolicyHitResultWarn  PolicyHitResult = "warn"
	PolicyHitResultBlock PolicyHitResult = "block"

	PolicyRemediationOpen       PolicyRemediationStatus = "open"
	PolicyRemediationInProgress PolicyRemediationStatus = "in_progress"
	PolicyRemediationMitigated  PolicyRemediationStatus = "mitigated"
	PolicyRemediationClosed     PolicyRemediationStatus = "closed"

	PolicyExceptionPending  PolicyExceptionStatus = "pending"
	PolicyExceptionApproved PolicyExceptionStatus = "approved"
	PolicyExceptionRejected PolicyExceptionStatus = "rejected"
	PolicyExceptionActive   PolicyExceptionStatus = "active"
	PolicyExceptionExpired  PolicyExceptionStatus = "expired"
	PolicyExceptionRevoked  PolicyExceptionStatus = "revoked"
)

type SecurityPolicy struct {
	ID                     uint64                `gorm:"primaryKey"`
	Name                   string                `gorm:"size:128;not null"`
	WorkspaceID            *uint64               `gorm:"index"`
	ProjectID              *uint64               `gorm:"index"`
	ScopeLevel             PolicyScopeLevel      `gorm:"size:32;not null"`
	Category               PolicyCategory        `gorm:"size:32;not null"`
	RuleTemplateJSON       string                `gorm:"type:longtext"`
	DefaultEnforcementMode PolicyEnforcementMode `gorm:"size:32;not null"`
	RiskLevel              PolicyRiskLevel       `gorm:"size:32;not null;default:medium"`
	Status                 PolicyStatus          `gorm:"size:32;not null;default:draft"`
	CreatedBy              *uint64               `gorm:"index"`
	UpdatedBy              *uint64               `gorm:"index"`
	CreatedAt              time.Time
	UpdatedAt              time.Time
}

type PolicyAssignment struct {
	ID                uint64                 `gorm:"primaryKey"`
	PolicyID          uint64                 `gorm:"index;not null"`
	WorkspaceID       *uint64                `gorm:"index"`
	ProjectID         *uint64                `gorm:"index"`
	ClusterRefsJSON   string                 `gorm:"type:longtext"`
	NamespaceRefsJSON string                 `gorm:"type:longtext"`
	ResourceKindsJSON string                 `gorm:"type:longtext"`
	EnforcementMode   PolicyEnforcementMode  `gorm:"size:32;not null"`
	RolloutStage      PolicyRolloutStage     `gorm:"size:32;not null"`
	Status            PolicyAssignmentStatus `gorm:"size:32;not null;default:pending"`
	EffectiveFrom     *time.Time
	EffectiveTo       *time.Time
	LastTaskID        *uint64 `gorm:"index"`
	CreatedBy         *uint64 `gorm:"index"`
	UpdatedBy         *uint64 `gorm:"index"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type PolicyDistributionTask struct {
	ID             uint64                      `gorm:"primaryKey"`
	PolicyID       uint64                      `gorm:"index;not null"`
	Operation      PolicyDistributionOperation `gorm:"size:32;not null"`
	Status         PolicyDistributionStatus    `gorm:"size:32;not null;default:pending"`
	TargetCount    int                         `gorm:"not null;default:0"`
	SucceededCount int                         `gorm:"not null;default:0"`
	FailedCount    int                         `gorm:"not null;default:0"`
	ResultSummary  string                      `gorm:"type:text"`
	CreatedBy      *uint64                     `gorm:"index"`
	StartedAt      *time.Time
	CompletedAt    *time.Time
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type PolicyHitRecord struct {
	ID                uint64                  `gorm:"primaryKey"`
	PolicyID          uint64                  `gorm:"index;not null"`
	AssignmentID      *uint64                 `gorm:"index"`
	ClusterID         *uint64                 `gorm:"index"`
	Namespace         string                  `gorm:"size:255"`
	ResourceKind      string                  `gorm:"size:64"`
	ResourceName      string                  `gorm:"size:255"`
	HitResult         PolicyHitResult         `gorm:"size:32;not null"`
	RiskLevel         PolicyRiskLevel         `gorm:"size:32;not null"`
	Message           string                  `gorm:"type:text"`
	RemediationStatus PolicyRemediationStatus `gorm:"size:32;not null;default:open"`
	DetectedAt        time.Time               `gorm:"not null"`
	ResolvedAt        *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type PolicyExceptionRequest struct {
	ID            uint64                `gorm:"primaryKey"`
	HitID         *uint64               `gorm:"index"`
	PolicyID      uint64                `gorm:"index;not null"`
	WorkspaceID   *uint64               `gorm:"index"`
	ProjectID     *uint64               `gorm:"index"`
	Status        PolicyExceptionStatus `gorm:"size:32;not null;default:pending"`
	Reason        string                `gorm:"type:text"`
	ExpiresAt     *time.Time
	RequestedBy   *uint64 `gorm:"index"`
	ReviewedBy    *uint64 `gorm:"index"`
	ReviewedAt    *time.Time
	ReviewComment string `gorm:"type:text"`
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

func (SecurityPolicy) TableName() string { return "security_policies" }

func (PolicyAssignment) TableName() string { return "policy_assignments" }

func (PolicyDistributionTask) TableName() string { return "policy_distribution_tasks" }

func (PolicyHitRecord) TableName() string { return "policy_hit_records" }

func (PolicyExceptionRequest) TableName() string { return "policy_exception_requests" }
