package domain

import "time"

type HAPolicyStatus string

const (
	HAPolicyStatusDraft      HAPolicyStatus = "draft"
	HAPolicyStatusActive     HAPolicyStatus = "active"
	HAPolicyStatusDegraded   HAPolicyStatus = "degraded"
	HAPolicyStatusRecovering HAPolicyStatus = "recovering"
)

type MaintenanceWindowStatus string

const (
	MaintenanceWindowStatusScheduled MaintenanceWindowStatus = "scheduled"
	MaintenanceWindowStatusActive    MaintenanceWindowStatus = "active"
	MaintenanceWindowStatusCompleted MaintenanceWindowStatus = "completed"
	MaintenanceWindowStatusException MaintenanceWindowStatus = "exception"
)

type PlatformHealthOverallStatus string

const (
	PlatformHealthOverallHealthy     PlatformHealthOverallStatus = "healthy"
	PlatformHealthOverallWarning     PlatformHealthOverallStatus = "warning"
	PlatformHealthOverallCritical    PlatformHealthOverallStatus = "critical"
	PlatformHealthOverallMaintenance PlatformHealthOverallStatus = "maintenance"
)

type CapacityBaselineStatus string

const (
	CapacityBaselineStatusActive   CapacityBaselineStatus = "active"
	CapacityBaselineStatusWarning  CapacityBaselineStatus = "warning"
	CapacityBaselineStatusCritical CapacityBaselineStatus = "critical"
)

type SREUpgradeStatus string

const (
	SREUpgradeStatusDraft            SREUpgradeStatus = "draft"
	SREUpgradeStatusReady            SREUpgradeStatus = "ready"
	SREUpgradeStatusRolling          SREUpgradeStatus = "rolling"
	SREUpgradeStatusAccepted         SREUpgradeStatus = "accepted"
	SREUpgradeStatusRollbackRequired SREUpgradeStatus = "rollback-required"
	SREUpgradeStatusClosed           SREUpgradeStatus = "closed"
	SREUpgradeStatusPrecheckFailed   SREUpgradeStatus = "precheck-failed"
)

type RollbackValidationResult string

const (
	RollbackValidationResultPassed  RollbackValidationResult = "passed"
	RollbackValidationResultWarning RollbackValidationResult = "warning"
	RollbackValidationResultFailed  RollbackValidationResult = "failed"
)

type RunbookStatus string

const (
	RunbookStatusDraft  RunbookStatus = "draft"
	RunbookStatusActive RunbookStatus = "active"
)

type AlertBaselineStatus string

const (
	AlertBaselineStatusActive   AlertBaselineStatus = "active"
	AlertBaselineStatusDisabled AlertBaselineStatus = "disabled"
)

type ScaleEvidenceStatus string

const (
	ScaleEvidenceStatusCaptured ScaleEvidenceStatus = "captured"
	ScaleEvidenceStatusAnalyzed ScaleEvidenceStatus = "analyzed"
	ScaleEvidenceStatusArchived ScaleEvidenceStatus = "archived"
)

type HAPolicy struct {
	ID                    uint64         `gorm:"primaryKey" json:"id"`
	WorkspaceID           uint64         `gorm:"index" json:"workspaceId"`
	ProjectID             *uint64        `gorm:"index" json:"projectId,omitempty"`
	Name                  string         `gorm:"size:128;not null;index:idx_sre_ha_scope_name,priority:3" json:"name"`
	ControlPlaneScope     string         `gorm:"size:128;not null" json:"controlPlaneScope"`
	DeploymentMode        string         `gorm:"size:64;not null" json:"deploymentMode"`
	ReplicaExpectation    int            `gorm:"not null;default:1" json:"replicaExpectation"`
	FailoverTriggerPolicy string         `gorm:"type:text" json:"failoverTriggerPolicy"`
	FailoverCooldown      string         `gorm:"size:64" json:"failoverCooldown"`
	TakeoverStatus        string         `gorm:"size:64;not null" json:"takeoverStatus"`
	LastFailoverAt        *time.Time     `json:"lastFailoverAt,omitempty"`
	LastRecoveryResult    string         `gorm:"type:text" json:"lastRecoveryResult"`
	Status                HAPolicyStatus `gorm:"size:32;not null;index" json:"status"`
	OwnerUserID           uint64         `gorm:"not null;index" json:"ownerUserId"`
	CreatedAt             time.Time      `json:"createdAt"`
	UpdatedAt             time.Time      `json:"updatedAt"`
}

func (HAPolicy) TableName() string { return "sre_ha_policies" }

type MaintenanceWindow struct {
	ID                   uint64                  `gorm:"primaryKey" json:"id"`
	WorkspaceID          uint64                  `gorm:"index" json:"workspaceId"`
	ProjectID            *uint64                 `gorm:"index" json:"projectId,omitempty"`
	Name                 string                  `gorm:"size:128;not null;index:idx_sre_window_scope_name,priority:3" json:"name"`
	WindowType           string                  `gorm:"size:64;not null" json:"windowType"`
	Scope                string                  `gorm:"size:128;not null" json:"scope"`
	StartAt              time.Time               `json:"startAt"`
	EndAt                time.Time               `json:"endAt"`
	AllowedOperations    string                  `gorm:"type:text" json:"allowedOperations"`
	RestrictedOperations string                  `gorm:"type:text" json:"restrictedOperations"`
	Status               MaintenanceWindowStatus `gorm:"size:32;not null;index" json:"status"`
	ExceptionReason      string                  `gorm:"type:text" json:"exceptionReason"`
	ApprovalRecord       string                  `gorm:"type:text" json:"approvalRecord"`
	PostCheckStatus      string                  `gorm:"size:64" json:"postCheckStatus"`
	OwnerUserID          uint64                  `gorm:"not null;index" json:"ownerUserId"`
	CreatedAt            time.Time               `json:"createdAt"`
	UpdatedAt            time.Time               `json:"updatedAt"`
}

func (MaintenanceWindow) TableName() string { return "sre_maintenance_windows" }

type PlatformHealthSnapshot struct {
	ID                      uint64                      `gorm:"primaryKey" json:"id"`
	WorkspaceID             uint64                      `gorm:"index" json:"workspaceId"`
	ProjectID               *uint64                     `gorm:"index" json:"projectId,omitempty"`
	HAPolicyID              *uint64                     `gorm:"index" json:"haPolicyId,omitempty"`
	SnapshotAt              time.Time                   `gorm:"index" json:"snapshotAt"`
	ComponentHealthSummary  string                      `gorm:"type:text" json:"componentHealthSummary"`
	DependencyHealthSummary string                      `gorm:"type:text" json:"dependencyHealthSummary"`
	TaskBacklogSummary      string                      `gorm:"type:text" json:"taskBacklogSummary"`
	CapacityRiskLevel       string                      `gorm:"size:32;not null" json:"capacityRiskLevel"`
	ThrottlingStatus        string                      `gorm:"size:64;not null" json:"throttlingStatus"`
	RecoverySummary         string                      `gorm:"type:text" json:"recoverySummary"`
	MaintenanceStatus       string                      `gorm:"size:64;not null" json:"maintenanceStatus"`
	OverallStatus           PlatformHealthOverallStatus `gorm:"size:32;not null;index" json:"overallStatus"`
	RecommendedActions      string                      `gorm:"type:text" json:"recommendedActions"`
	CreatedAt               time.Time                   `json:"createdAt"`
	UpdatedAt               time.Time                   `json:"updatedAt"`
}

func (PlatformHealthSnapshot) TableName() string { return "sre_platform_health_snapshots" }

type CapacityBaseline struct {
	ID                uint64                 `gorm:"primaryKey" json:"id"`
	WorkspaceID       uint64                 `gorm:"index" json:"workspaceId"`
	ProjectID         *uint64                `gorm:"index" json:"projectId,omitempty"`
	Name              string                 `gorm:"size:128;not null;index:idx_sre_capacity_scope_name,priority:3" json:"name"`
	ResourceDimension string                 `gorm:"size:64;not null" json:"resourceDimension"`
	BaselineRange     string                 `gorm:"type:text" json:"baselineRange"`
	Thresholds        string                 `gorm:"type:text" json:"thresholds"`
	GrowthTrend       string                 `gorm:"type:text" json:"growthTrend"`
	ForecastWindow    string                 `gorm:"size:64;not null" json:"forecastWindow"`
	ForecastResult    string                 `gorm:"type:text" json:"forecastResult"`
	ConfidenceLevel   string                 `gorm:"size:32;not null" json:"confidenceLevel"`
	Status            CapacityBaselineStatus `gorm:"size:32;not null;index" json:"status"`
	OwnerUserID       uint64                 `gorm:"not null;index" json:"ownerUserId"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

func (CapacityBaseline) TableName() string { return "sre_capacity_baselines" }

type SREUpgradePlan struct {
	ID                   uint64           `gorm:"primaryKey" json:"id"`
	WorkspaceID          uint64           `gorm:"index" json:"workspaceId"`
	ProjectID            *uint64          `gorm:"index" json:"projectId,omitempty"`
	MaintenanceWindowID  *uint64          `gorm:"index" json:"maintenanceWindowId,omitempty"`
	Name                 string           `gorm:"size:128;not null;index:idx_sre_upgrade_scope_name,priority:3" json:"name"`
	CurrentVersion       string           `gorm:"size:64;not null" json:"currentVersion"`
	TargetVersion        string           `gorm:"size:64;not null" json:"targetVersion"`
	CompatibilitySummary string           `gorm:"type:text" json:"compatibilitySummary"`
	PrecheckResult       string           `gorm:"type:text" json:"precheckResult"`
	RolloutStrategy      string           `gorm:"type:text" json:"rolloutStrategy"`
	ExecutionStage       string           `gorm:"size:64;not null" json:"executionStage"`
	ExecutionProgress    int              `gorm:"not null;default:0" json:"executionProgress"`
	AcceptanceResult     string           `gorm:"type:text" json:"acceptanceResult"`
	RollbackReadiness    string           `gorm:"type:text" json:"rollbackReadiness"`
	Status               SREUpgradeStatus `gorm:"size:32;not null;index" json:"status"`
	CreatedBy            uint64           `gorm:"not null;index" json:"createdBy"`
	CreatedAt            time.Time        `json:"createdAt"`
	UpdatedAt            time.Time        `json:"updatedAt"`
}

func (SREUpgradePlan) TableName() string { return "sre_upgrade_plans" }

type RollbackValidation struct {
	ID              uint64                   `gorm:"primaryKey" json:"id"`
	UpgradePlanID   uint64                   `gorm:"not null;index" json:"upgradePlanId"`
	ValidationScope string                   `gorm:"size:128;not null" json:"validationScope"`
	Preconditions   string                   `gorm:"type:text" json:"preconditions"`
	Result          RollbackValidationResult `gorm:"size:32;not null;index" json:"result"`
	RemainingRisk   string                   `gorm:"type:text" json:"remainingRisk"`
	ValidatedAt     time.Time                `gorm:"index" json:"validatedAt"`
	ValidatedBy     uint64                   `gorm:"not null;index" json:"validatedBy"`
	CreatedAt       time.Time                `json:"createdAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
}

func (RollbackValidation) TableName() string { return "sre_rollback_validations" }

type RunbookArticle struct {
	ID                   uint64        `gorm:"primaryKey" json:"id"`
	WorkspaceID          uint64        `gorm:"index" json:"workspaceId"`
	ProjectID            *uint64       `gorm:"index" json:"projectId,omitempty"`
	Title                string        `gorm:"size:128;not null;index:idx_sre_runbook_scope_title,priority:3" json:"title"`
	ScenarioType         string        `gorm:"size:64;not null" json:"scenarioType"`
	ApplicableComponents string        `gorm:"type:text" json:"applicableComponents"`
	RiskLevel            string        `gorm:"size:32;not null" json:"riskLevel"`
	ChecklistSummary     string        `gorm:"type:text" json:"checklistSummary"`
	RecoverySteps        string        `gorm:"type:text" json:"recoverySteps"`
	VerificationSummary  string        `gorm:"type:text" json:"verificationSummary"`
	Status               RunbookStatus `gorm:"size:32;not null;index" json:"status"`
	OwnerUserID          uint64        `gorm:"not null;index" json:"ownerUserId"`
	CreatedAt            time.Time     `json:"createdAt"`
	UpdatedAt            time.Time     `json:"updatedAt"`
}

func (RunbookArticle) TableName() string { return "sre_runbook_articles" }

type AlertBaseline struct {
	ID                   uint64              `gorm:"primaryKey" json:"id"`
	WorkspaceID          uint64              `gorm:"index" json:"workspaceId"`
	ProjectID            *uint64             `gorm:"index" json:"projectId,omitempty"`
	Name                 string              `gorm:"size:128;not null" json:"name"`
	ComponentScope       string              `gorm:"size:128;not null" json:"componentScope"`
	SignalType           string              `gorm:"size:64;not null" json:"signalType"`
	BaselineCondition    string              `gorm:"type:text" json:"baselineCondition"`
	Severity             string              `gorm:"size:32;not null" json:"severity"`
	RecommendedRunbookID *uint64             `gorm:"index" json:"recommendedRunbookId,omitempty"`
	Status               AlertBaselineStatus `gorm:"size:32;not null;index" json:"status"`
	OwnerUserID          uint64              `gorm:"not null;index" json:"ownerUserId"`
	CreatedAt            time.Time           `json:"createdAt"`
	UpdatedAt            time.Time           `json:"updatedAt"`
}

func (AlertBaseline) TableName() string { return "sre_alert_baselines" }

type ScaleEvidence struct {
	ID                  uint64              `gorm:"primaryKey" json:"id"`
	WorkspaceID         uint64              `gorm:"index" json:"workspaceId"`
	ProjectID           *uint64             `gorm:"index" json:"projectId,omitempty"`
	CapacityBaselineID  *uint64             `gorm:"index" json:"capacityBaselineId,omitempty"`
	RunbookArticleID    *uint64             `gorm:"index" json:"runbookArticleId,omitempty"`
	EvidenceType        string              `gorm:"size:64;not null;index" json:"evidenceType"`
	Scope               string              `gorm:"size:128;not null" json:"scope"`
	SampleWindow        string              `gorm:"size:64;not null" json:"sampleWindow"`
	Summary             string              `gorm:"type:text" json:"summary"`
	BottleneckSummary   string              `gorm:"type:text" json:"bottleneckSummary"`
	ForecastSummary     string              `gorm:"type:text" json:"forecastSummary"`
	ConfidenceLevel     string              `gorm:"size:32;not null" json:"confidenceLevel"`
	RecoveryObservation string              `gorm:"type:text" json:"recoveryObservation"`
	Status              ScaleEvidenceStatus `gorm:"size:32;not null;index" json:"status"`
	CapturedAt          time.Time           `gorm:"index" json:"capturedAt"`
	OwnerUserID         uint64              `gorm:"not null;index" json:"ownerUserId"`
	CreatedAt           time.Time           `json:"createdAt"`
	UpdatedAt           time.Time           `json:"updatedAt"`
}

func (ScaleEvidence) TableName() string { return "sre_scale_evidence" }
