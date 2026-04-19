package domain

import "time"

type BackupPolicyStatus string

const (
	BackupPolicyStatusDraft      BackupPolicyStatus = "draft"
	BackupPolicyStatusActive     BackupPolicyStatus = "active"
	BackupPolicyStatusPaused     BackupPolicyStatus = "paused"
	BackupPolicyStatusDeprecated BackupPolicyStatus = "deprecated"
)

type RestorePointResult string

const (
	RestorePointResultSucceeded          RestorePointResult = "succeeded"
	RestorePointResultPartiallySucceeded RestorePointResult = "partially_succeeded"
	RestorePointResultFailed             RestorePointResult = "failed"
	RestorePointResultExpired            RestorePointResult = "expired"
)

type RestoreJobType string

const (
	RestoreJobTypeInPlace      RestoreJobType = "in-place-restore"
	RestoreJobTypeCrossCluster RestoreJobType = "cross-cluster-restore"
	RestoreJobTypeMigration    RestoreJobType = "environment-migration"
	RestoreJobTypeSelective    RestoreJobType = "selective-restore"
)

type RestoreJobStatus string

const (
	RestoreJobStatusPending            RestoreJobStatus = "pending"
	RestoreJobStatusValidating         RestoreJobStatus = "validating"
	RestoreJobStatusRunning            RestoreJobStatus = "running"
	RestoreJobStatusSucceeded          RestoreJobStatus = "succeeded"
	RestoreJobStatusPartiallySucceeded RestoreJobStatus = "partially_succeeded"
	RestoreJobStatusFailed             RestoreJobStatus = "failed"
	RestoreJobStatusCanceled           RestoreJobStatus = "canceled"
	RestoreJobStatusBlocked            RestoreJobStatus = "blocked"
)

type MigrationPlanStatus string

const (
	MigrationPlanStatusDraft     MigrationPlanStatus = "draft"
	MigrationPlanStatusApproved  MigrationPlanStatus = "approved"
	MigrationPlanStatusRunning   MigrationPlanStatus = "running"
	MigrationPlanStatusSucceeded MigrationPlanStatus = "succeeded"
	MigrationPlanStatusFailed    MigrationPlanStatus = "failed"
	MigrationPlanStatusCanceled  MigrationPlanStatus = "canceled"
)

type DRDrillPlanStatus string

const (
	DRDrillPlanStatusDraft   DRDrillPlanStatus = "draft"
	DRDrillPlanStatusActive  DRDrillPlanStatus = "active"
	DRDrillPlanStatusPaused  DRDrillPlanStatus = "paused"
	DRDrillPlanStatusRetired DRDrillPlanStatus = "retired"
)

type DRDrillRecordStatus string

const (
	DRDrillRecordStatusPending            DRDrillRecordStatus = "pending"
	DRDrillRecordStatusRunning            DRDrillRecordStatus = "running"
	DRDrillRecordStatusSucceeded          DRDrillRecordStatus = "succeeded"
	DRDrillRecordStatusFailed             DRDrillRecordStatus = "failed"
	DRDrillRecordStatusPartiallySucceeded DRDrillRecordStatus = "partially_succeeded"
	DRDrillRecordStatusCanceled           DRDrillRecordStatus = "canceled"
)

type BackupAuditOutcome string

const (
	BackupAuditOutcomeSucceeded BackupAuditOutcome = "succeeded"
	BackupAuditOutcomeFailed    BackupAuditOutcome = "failed"
	BackupAuditOutcomeBlocked   BackupAuditOutcome = "blocked"
	BackupAuditOutcomeCanceled  BackupAuditOutcome = "canceled"
)

type BackupPolicy struct {
	ID                 uint64             `gorm:"primaryKey" json:"id"`
	Name               string             `gorm:"size:128;not null;index:idx_backup_policy_scope_name,priority:3" json:"name"`
	Description        string             `gorm:"type:text" json:"description,omitempty"`
	ScopeType          string             `gorm:"size:64;not null;index:idx_backup_policy_scope_name,priority:1" json:"scopeType"`
	ScopeRef           string             `gorm:"size:128;not null;index:idx_backup_policy_scope_name,priority:2" json:"scopeRef"`
	WorkspaceID        uint64             `gorm:"not null;index" json:"workspaceId"`
	ProjectID          *uint64            `gorm:"index" json:"projectId,omitempty"`
	ExecutionMode      string             `gorm:"size:32;not null" json:"executionMode"`
	ScheduleExpression string             `gorm:"size:256" json:"scheduleExpression,omitempty"`
	RetentionRule      string             `gorm:"size:256;not null" json:"retentionRule"`
	ConsistencyLevel   string             `gorm:"size:64;not null" json:"consistencyLevel"`
	Status             BackupPolicyStatus `gorm:"size:32;not null;index" json:"status"`
	OwnerUserID        uint64             `gorm:"not null;index" json:"ownerUserId"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

func (BackupPolicy) TableName() string { return "backup_policies" }

type RestorePoint struct {
	ID                 uint64             `gorm:"primaryKey" json:"id"`
	PolicyID           uint64             `gorm:"not null;index" json:"policyId"`
	WorkspaceID        uint64             `gorm:"not null;index" json:"workspaceId"`
	ProjectID          *uint64            `gorm:"index" json:"projectId,omitempty"`
	ScopeSnapshot      string             `gorm:"type:text;not null" json:"scopeSnapshot"`
	BackupStartedAt    time.Time          `gorm:"not null" json:"backupStartedAt"`
	BackupCompletedAt  *time.Time         `json:"backupCompletedAt,omitempty"`
	DurationSeconds    int                `gorm:"not null;default:0" json:"durationSeconds"`
	Result             RestorePointResult `gorm:"size:32;not null;index" json:"result"`
	ConsistencySummary string             `gorm:"type:text" json:"consistencySummary,omitempty"`
	FailureReason      string             `gorm:"type:text" json:"failureReason,omitempty"`
	StorageRef         string             `gorm:"size:256" json:"storageRef,omitempty"`
	ExpiresAt          *time.Time         `json:"expiresAt,omitempty"`
	CreatedBy          uint64             `gorm:"not null;index" json:"createdBy"`
	CreatedAt          time.Time          `json:"createdAt"`
	UpdatedAt          time.Time          `json:"updatedAt"`
}

func (RestorePoint) TableName() string { return "restore_points" }

type RestoreJob struct {
	ID                uint64           `gorm:"primaryKey" json:"id"`
	RestorePointID    uint64           `gorm:"not null;index" json:"restorePointId"`
	WorkspaceID       uint64           `gorm:"not null;index" json:"workspaceId"`
	ProjectID         *uint64          `gorm:"index" json:"projectId,omitempty"`
	JobType           RestoreJobType   `gorm:"size:64;not null;index" json:"jobType"`
	SourceEnvironment string           `gorm:"size:128" json:"sourceEnvironment,omitempty"`
	TargetEnvironment string           `gorm:"size:128;not null" json:"targetEnvironment"`
	ScopeSelection    string           `gorm:"type:text;not null" json:"scopeSelection"`
	ConflictSummary   string           `gorm:"type:text" json:"conflictSummary,omitempty"`
	ConsistencyNotice string           `gorm:"type:text" json:"consistencyNotice,omitempty"`
	Status            RestoreJobStatus `gorm:"size:32;not null;index" json:"status"`
	ResultSummary     string           `gorm:"type:text" json:"resultSummary,omitempty"`
	FailureReason     string           `gorm:"type:text" json:"failureReason,omitempty"`
	RequestedBy       uint64           `gorm:"not null;index" json:"requestedBy"`
	StartedAt         *time.Time       `json:"startedAt,omitempty"`
	CompletedAt       *time.Time       `json:"completedAt,omitempty"`
	CreatedAt         time.Time        `json:"createdAt"`
	UpdatedAt         time.Time        `json:"updatedAt"`
}

func (RestoreJob) TableName() string { return "restore_jobs" }

type MigrationPlan struct {
	ID              uint64              `gorm:"primaryKey" json:"id"`
	Name            string              `gorm:"size:128;not null;uniqueIndex" json:"name"`
	WorkspaceID     uint64              `gorm:"not null;index" json:"workspaceId"`
	ProjectID       *uint64             `gorm:"index" json:"projectId,omitempty"`
	SourceClusterID uint64              `gorm:"not null;index" json:"sourceClusterId"`
	TargetClusterID uint64              `gorm:"not null;index" json:"targetClusterId"`
	ScopeSelection  string              `gorm:"type:text;not null" json:"scopeSelection"`
	MappingRules    string              `gorm:"type:text" json:"mappingRules,omitempty"`
	CutoverSteps    string              `gorm:"type:text" json:"cutoverSteps,omitempty"`
	Status          MigrationPlanStatus `gorm:"size:32;not null;index" json:"status"`
	CreatedBy       uint64              `gorm:"not null;index" json:"createdBy"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

func (MigrationPlan) TableName() string { return "migration_plans" }

type DRDrillPlan struct {
	ID                  uint64            `gorm:"primaryKey" json:"id"`
	Name                string            `gorm:"size:128;not null;uniqueIndex" json:"name"`
	Description         string            `gorm:"type:text" json:"description,omitempty"`
	WorkspaceID         uint64            `gorm:"not null;index" json:"workspaceId"`
	ProjectID           *uint64           `gorm:"index" json:"projectId,omitempty"`
	ScopeSelection      string            `gorm:"type:text;not null" json:"scopeSelection"`
	RPOTargetMinutes    int               `gorm:"not null" json:"rpoTargetMinutes"`
	RTOTargetMinutes    int               `gorm:"not null" json:"rtoTargetMinutes"`
	RoleAssignments     string            `gorm:"type:text" json:"roleAssignments,omitempty"`
	CutoverProcedure    string            `gorm:"type:text;not null" json:"cutoverProcedure"`
	ValidationChecklist string            `gorm:"type:text;not null" json:"validationChecklist"`
	Status              DRDrillPlanStatus `gorm:"size:32;not null;index" json:"status"`
	CreatedBy           uint64            `gorm:"not null;index" json:"createdBy"`
	CreatedAt           time.Time         `json:"createdAt"`
	UpdatedAt           time.Time         `json:"updatedAt"`
}

func (DRDrillPlan) TableName() string { return "dr_drill_plans" }

type DRDrillRecord struct {
	ID                uint64              `gorm:"primaryKey" json:"id"`
	PlanID            uint64              `gorm:"not null;index" json:"planId"`
	WorkspaceID       uint64              `gorm:"not null;index" json:"workspaceId"`
	ProjectID         *uint64             `gorm:"index" json:"projectId,omitempty"`
	StartedAt         time.Time           `gorm:"not null" json:"startedAt"`
	CompletedAt       *time.Time          `json:"completedAt,omitempty"`
	ActualRPOMinutes  int                 `gorm:"not null;default:0" json:"actualRpoMinutes"`
	ActualRTOMinutes  int                 `gorm:"not null;default:0" json:"actualRtoMinutes"`
	Status            DRDrillRecordStatus `gorm:"size:32;not null;index" json:"status"`
	StepResults       string              `gorm:"type:text" json:"stepResults,omitempty"`
	ValidationResults string              `gorm:"type:text" json:"validationResults,omitempty"`
	IncidentNotes     string              `gorm:"type:text" json:"incidentNotes,omitempty"`
	ExecutedBy        uint64              `gorm:"not null;index" json:"executedBy"`
	CreatedAt         time.Time           `json:"createdAt"`
	UpdatedAt         time.Time           `json:"updatedAt"`
}

func (DRDrillRecord) TableName() string { return "dr_drill_records" }

type DRDrillReport struct {
	ID                 uint64    `gorm:"primaryKey" json:"id"`
	DrillRecordID      uint64    `gorm:"not null;uniqueIndex" json:"drillRecordId"`
	GoalAssessment     string    `gorm:"type:text;not null" json:"goalAssessment"`
	GapSummary         string    `gorm:"type:text" json:"gapSummary,omitempty"`
	IssuesFound        string    `gorm:"type:text" json:"issuesFound,omitempty"`
	ImprovementActions string    `gorm:"type:text;not null" json:"improvementActions"`
	PublishedAt        time.Time `gorm:"not null" json:"publishedAt"`
	PublishedBy        uint64    `gorm:"not null;index" json:"publishedBy"`
	CreatedAt          time.Time `json:"createdAt"`
	UpdatedAt          time.Time `json:"updatedAt"`
}

func (DRDrillReport) TableName() string { return "dr_drill_reports" }

type BackupAuditEvent struct {
	ID             uint64             `gorm:"primaryKey" json:"id"`
	Action         string             `gorm:"size:128;not null;index" json:"action"`
	ActorUserID    uint64             `gorm:"not null;index" json:"actorUserId"`
	TargetType     string             `gorm:"size:64;not null;index" json:"targetType"`
	TargetRef      string             `gorm:"size:128;not null;index" json:"targetRef"`
	WorkspaceID    uint64             `gorm:"not null;index" json:"workspaceId"`
	ProjectID      *uint64            `gorm:"index" json:"projectId,omitempty"`
	ScopeSnapshot  string             `gorm:"type:text" json:"scopeSnapshot,omitempty"`
	Outcome        BackupAuditOutcome `gorm:"size:32;not null;index" json:"outcome"`
	DetailSnapshot string             `gorm:"type:text" json:"detailSnapshot,omitempty"`
	OccurredAt     time.Time          `gorm:"not null;index" json:"occurredAt"`
	CreatedAt      time.Time          `json:"createdAt"`
	UpdatedAt      time.Time          `json:"updatedAt"`
}

func (BackupAuditEvent) TableName() string { return "backup_audit_events" }
