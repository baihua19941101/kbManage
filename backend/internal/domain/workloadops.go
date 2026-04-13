package domain

import "time"

type WorkloadActionType string

type BatchOperationStatus string

type BatchOperationItemStatus string

type TerminalSessionStatus string

const (
	WorkloadActionTypeScale           WorkloadActionType = "scale"
	WorkloadActionTypeRestart         WorkloadActionType = "restart"
	WorkloadActionTypeRedeploy        WorkloadActionType = "redeploy"
	WorkloadActionTypeReplaceInstance WorkloadActionType = "replace-instance"
	WorkloadActionTypeRollback        WorkloadActionType = "rollback"

	BatchOperationStatusPending            BatchOperationStatus = "pending"
	BatchOperationStatusRunning            BatchOperationStatus = "running"
	BatchOperationStatusPartiallySucceeded BatchOperationStatus = "partially_succeeded"
	BatchOperationStatusSucceeded          BatchOperationStatus = "succeeded"
	BatchOperationStatusFailed             BatchOperationStatus = "failed"
	BatchOperationStatusCanceled           BatchOperationStatus = "canceled"

	BatchOperationItemStatusPending   BatchOperationItemStatus = "pending"
	BatchOperationItemStatusRunning   BatchOperationItemStatus = "running"
	BatchOperationItemStatusSucceeded BatchOperationItemStatus = "succeeded"
	BatchOperationItemStatusFailed    BatchOperationItemStatus = "failed"
	BatchOperationItemStatusSkipped   BatchOperationItemStatus = "skipped"
	BatchOperationItemStatusCanceled  BatchOperationItemStatus = "canceled"

	TerminalSessionStatusPending TerminalSessionStatus = "pending"
	TerminalSessionStatusActive  TerminalSessionStatus = "active"
	TerminalSessionStatusClosed  TerminalSessionStatus = "closed"
	TerminalSessionStatusExpired TerminalSessionStatus = "expired"
	TerminalSessionStatusDenied  TerminalSessionStatus = "denied"
	TerminalSessionStatusFailed  TerminalSessionStatus = "failed"
)

type WorkloadActionRequest struct {
	ID                uint64             `gorm:"primaryKey"`
	RequestID         string             `gorm:"size:64;uniqueIndex;not null"`
	OperatorID        uint64             `gorm:"index;not null"`
	ClusterID         uint64             `gorm:"index;not null"`
	WorkspaceID       *uint64            `gorm:"index"`
	ProjectID         *uint64            `gorm:"index"`
	Namespace         string             `gorm:"size:255;not null"`
	ResourceKind      string             `gorm:"size:64;not null"`
	ResourceName      string             `gorm:"size:255;not null"`
	TargetInstanceRef string             `gorm:"size:255"`
	ActionType        WorkloadActionType `gorm:"size:64;not null"`
	RiskLevel         RiskLevel          `gorm:"size:32;not null"`
	RiskConfirmed     bool               `gorm:"not null;default:false"`
	PayloadJSON       string             `gorm:"type:longtext"`
	Status            OperationStatus    `gorm:"size:32;not null"`
	ProgressMessage   string             `gorm:"type:text"`
	ResultMessage     string             `gorm:"type:text"`
	FailureReason     string             `gorm:"type:text"`
	BatchID           *uint64            `gorm:"index"`
	StartedAt         *time.Time
	CompletedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type BatchOperationTask struct {
	ID                uint64               `gorm:"primaryKey"`
	RequestID         string               `gorm:"size:64;uniqueIndex;not null"`
	OperatorID        uint64               `gorm:"index;not null"`
	ActionType        WorkloadActionType   `gorm:"size:64;not null"`
	ScopeSnapshotJSON string               `gorm:"type:longtext"`
	RiskLevel         RiskLevel            `gorm:"size:32;not null"`
	RiskConfirmed     bool                 `gorm:"not null;default:false"`
	TotalTargets      int                  `gorm:"not null;default:0"`
	SucceededTargets  int                  `gorm:"not null;default:0"`
	FailedTargets     int                  `gorm:"not null;default:0"`
	CanceledTargets   int                  `gorm:"not null;default:0"`
	Status            BatchOperationStatus `gorm:"size:32;not null"`
	ProgressPercent   int                  `gorm:"not null;default:0"`
	StartedAt         *time.Time
	CompletedAt       *time.Time
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type BatchOperationItem struct {
	ID              uint64                   `gorm:"primaryKey"`
	BatchID         uint64                   `gorm:"index;not null"`
	ClusterID       uint64                   `gorm:"index;not null"`
	WorkspaceID     *uint64                  `gorm:"index"`
	ProjectID       *uint64                  `gorm:"index"`
	Namespace       string                   `gorm:"size:255;not null"`
	ResourceKind    string                   `gorm:"size:64;not null"`
	ResourceName    string                   `gorm:"size:255;not null"`
	Status          BatchOperationItemStatus `gorm:"size:32;not null"`
	ActionRequestID *uint64                  `gorm:"index"`
	FailureReason   string                   `gorm:"type:text"`
	ResultMessage   string                   `gorm:"type:text"`
	StartedAt       *time.Time
	CompletedAt     *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type TerminalSession struct {
	ID                 uint64                `gorm:"primaryKey"`
	SessionKey         string                `gorm:"size:128;uniqueIndex;not null"`
	OperatorID         uint64                `gorm:"index;not null"`
	ClusterID          uint64                `gorm:"index;not null"`
	WorkspaceID        *uint64               `gorm:"index"`
	ProjectID          *uint64               `gorm:"index"`
	Namespace          string                `gorm:"size:255;not null"`
	PodName            string                `gorm:"size:255;not null"`
	ContainerName      string                `gorm:"size:255;not null"`
	WorkloadKind       string                `gorm:"size:64"`
	WorkloadName       string                `gorm:"size:255"`
	Status             TerminalSessionStatus `gorm:"size:32;not null"`
	StartedAt          *time.Time
	EndedAt            *time.Time
	DurationSeconds    int    `gorm:"not null;default:0"`
	CloseReason        string `gorm:"size:255"`
	ClientMetadataJSON string `gorm:"type:longtext"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (WorkloadActionRequest) TableName() string { return "workload_action_requests" }

func (BatchOperationTask) TableName() string { return "workload_batch_tasks" }

func (BatchOperationItem) TableName() string { return "workload_batch_items" }

func (TerminalSession) TableName() string { return "workload_terminal_sessions" }
