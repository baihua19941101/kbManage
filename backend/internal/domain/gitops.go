package domain

import "time"

type DeliverySourceType string

type DeliverySourceStatus string

type ClusterTargetGroupStatus string

type PromotionMode string

type EnvironmentStageStatus string

type ConfigurationOverlayType string

type DeliverySyncMode string

type DeliveryUnitStatus string

type ReleaseRevisionStatus string

type DeliveryActionType string

type DeliveryOperationStatus string

const (
	DeliverySourceTypeGit     DeliverySourceType = "git"
	DeliverySourceTypePackage DeliverySourceType = "package"

	DeliverySourceStatusPending  DeliverySourceStatus = "pending"
	DeliverySourceStatusReady    DeliverySourceStatus = "ready"
	DeliverySourceStatusFailed   DeliverySourceStatus = "failed"
	DeliverySourceStatusDisabled DeliverySourceStatus = "disabled"

	ClusterTargetGroupStatusActive   ClusterTargetGroupStatus = "active"
	ClusterTargetGroupStatusStale    ClusterTargetGroupStatus = "stale"
	ClusterTargetGroupStatusDisabled ClusterTargetGroupStatus = "disabled"

	PromotionModeManual    PromotionMode = "manual"
	PromotionModeAutomatic PromotionMode = "automatic"

	EnvironmentStageStatusIdle        EnvironmentStageStatus = "idle"
	EnvironmentStageStatusWaiting     EnvironmentStageStatus = "waiting"
	EnvironmentStageStatusProgressing EnvironmentStageStatus = "progressing"
	EnvironmentStageStatusSucceeded   EnvironmentStageStatus = "succeeded"
	EnvironmentStageStatusFailed      EnvironmentStageStatus = "failed"
	EnvironmentStageStatusPaused      EnvironmentStageStatus = "paused"

	ConfigurationOverlayTypeValues          ConfigurationOverlayType = "values"
	ConfigurationOverlayTypePatch           ConfigurationOverlayType = "patch"
	ConfigurationOverlayTypeManifestSnippet ConfigurationOverlayType = "manifest-snippet"

	DeliverySyncModeManual DeliverySyncMode = "manual"
	DeliverySyncModeAuto   DeliverySyncMode = "auto"

	DeliveryUnitStatusReady       DeliveryUnitStatus = "ready"
	DeliveryUnitStatusProgressing DeliveryUnitStatus = "progressing"
	DeliveryUnitStatusDegraded    DeliveryUnitStatus = "degraded"
	DeliveryUnitStatusOutOfSync   DeliveryUnitStatus = "out_of_sync"
	DeliveryUnitStatusPaused      DeliveryUnitStatus = "paused"
	DeliveryUnitStatusUnknown     DeliveryUnitStatus = "unknown"

	ReleaseRevisionStatusActive     ReleaseRevisionStatus = "active"
	ReleaseRevisionStatusHistorical ReleaseRevisionStatus = "historical"
	ReleaseRevisionStatusFailed     ReleaseRevisionStatus = "failed"
	ReleaseRevisionStatusRolledBack ReleaseRevisionStatus = "rolled_back"

	DeliveryActionTypeInstall   DeliveryActionType = "install"
	DeliveryActionTypeSync      DeliveryActionType = "sync"
	DeliveryActionTypeResync    DeliveryActionType = "resync"
	DeliveryActionTypeUpgrade   DeliveryActionType = "upgrade"
	DeliveryActionTypePromote   DeliveryActionType = "promote"
	DeliveryActionTypeRollback  DeliveryActionType = "rollback"
	DeliveryActionTypePause     DeliveryActionType = "pause"
	DeliveryActionTypeResume    DeliveryActionType = "resume"
	DeliveryActionTypeUninstall DeliveryActionType = "uninstall"

	DeliveryOperationStatusPending            DeliveryOperationStatus = "pending"
	DeliveryOperationStatusRunning            DeliveryOperationStatus = "running"
	DeliveryOperationStatusPartiallySucceeded DeliveryOperationStatus = "partially_succeeded"
	DeliveryOperationStatusSucceeded          DeliveryOperationStatus = "succeeded"
	DeliveryOperationStatusFailed             DeliveryOperationStatus = "failed"
	DeliveryOperationStatusCanceled           DeliveryOperationStatus = "canceled"
)

type DeliverySource struct {
	ID               uint64               `gorm:"primaryKey"`
	Name             string               `gorm:"size:128;not null"`
	SourceType       DeliverySourceType   `gorm:"size:32;not null"`
	Endpoint         string               `gorm:"size:1024;not null"`
	DefaultRef       string               `gorm:"size:256"`
	CredentialRef    string               `gorm:"size:256"`
	WorkspaceID      *uint64              `gorm:"index"`
	ProjectID        *uint64              `gorm:"index"`
	Status           DeliverySourceStatus `gorm:"size:32;not null;default:pending"`
	LastVerifiedAt   *time.Time
	LastErrorMessage string `gorm:"type:text"`
	CreatedAt        time.Time
	UpdatedAt        time.Time
}

type ClusterTargetGroup struct {
	ID                      uint64                   `gorm:"primaryKey"`
	Name                    string                   `gorm:"size:128;not null"`
	WorkspaceID             uint64                   `gorm:"index;not null"`
	ProjectID               *uint64                  `gorm:"index"`
	ClusterRefsJSON         string                   `gorm:"type:longtext"`
	ClusterSelectorSnapshot string                   `gorm:"type:longtext"`
	Description             string                   `gorm:"size:1024"`
	Status                  ClusterTargetGroupStatus `gorm:"size:32;not null;default:active"`
	CreatedAt               time.Time
	UpdatedAt               time.Time
}

type EnvironmentStage struct {
	ID              uint64                 `gorm:"primaryKey"`
	DeliveryUnitID  uint64                 `gorm:"index;not null"`
	Name            string                 `gorm:"size:128;not null"`
	OrderIndex      int                    `gorm:"not null"`
	TargetGroupID   uint64                 `gorm:"index;not null"`
	PromotionMode   PromotionMode          `gorm:"size:32;not null;default:manual"`
	Paused          bool                   `gorm:"not null;default:false"`
	Status          EnvironmentStageStatus `gorm:"size:32;not null;default:idle"`
	LastEnteredAt   *time.Time
	LastCompletedAt *time.Time
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

type ConfigurationOverlay struct {
	ID                 uint64                   `gorm:"primaryKey"`
	DeliveryUnitID     uint64                   `gorm:"index;not null"`
	EnvironmentStageID *uint64                  `gorm:"index"`
	OverlayType        ConfigurationOverlayType `gorm:"size:32;not null"`
	OverlayRef         string                   `gorm:"size:1024;not null"`
	Precedence         int                      `gorm:"not null;default:0"`
	EffectiveScopeJSON string                   `gorm:"type:longtext"`
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

type ApplicationDeliveryUnit struct {
	ID                   uint64             `gorm:"primaryKey"`
	Name                 string             `gorm:"size:128;not null"`
	WorkspaceID          uint64             `gorm:"index;not null"`
	ProjectID            *uint64            `gorm:"index"`
	SourceID             uint64             `gorm:"index;not null"`
	SourcePath           string             `gorm:"size:1024"`
	DefaultNamespace     string             `gorm:"size:255"`
	SyncMode             DeliverySyncMode   `gorm:"size:32;not null;default:manual"`
	ReleasePolicyJSON    string             `gorm:"type:longtext"`
	DesiredRevision      string             `gorm:"size:256"`
	DesiredAppVersion    string             `gorm:"size:128"`
	DesiredConfigVersion string             `gorm:"size:128"`
	Paused               bool               `gorm:"not null;default:false"`
	DeliveryStatus       DeliveryUnitStatus `gorm:"size:32;not null;default:unknown"`
	LastSyncedAt         *time.Time
	LastReleaseID        *uint64 `gorm:"index"`
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type ReleaseRevision struct {
	ID                  uint64                `gorm:"primaryKey"`
	DeliveryUnitID      uint64                `gorm:"index;not null"`
	SourceRevision      string                `gorm:"size:256;not null"`
	AppVersion          string                `gorm:"size:128"`
	ConfigVersion       string                `gorm:"size:128"`
	EffectiveScopeJSON  string                `gorm:"type:longtext"`
	ReleaseNotesSummary string                `gorm:"type:text"`
	CreatedBy           uint64                `gorm:"index;not null"`
	RollbackAvailable   bool                  `gorm:"not null;default:false"`
	Status              ReleaseRevisionStatus `gorm:"size:32;not null;default:historical"`
	CreatedAt           time.Time
}

type DeliveryOperation struct {
	ID                 uint64                  `gorm:"primaryKey"`
	RequestID          string                  `gorm:"size:64;uniqueIndex;not null"`
	OperatorID         uint64                  `gorm:"index;not null"`
	DeliveryUnitID     uint64                  `gorm:"index;not null"`
	EnvironmentStageID *uint64                 `gorm:"index"`
	ActionType         DeliveryActionType      `gorm:"size:32;not null"`
	TargetReleaseID    *uint64                 `gorm:"index"`
	PayloadJSON        string                  `gorm:"type:longtext"`
	Status             DeliveryOperationStatus `gorm:"size:32;not null;default:pending"`
	ProgressPercent    int                     `gorm:"not null;default:0"`
	ResultSummary      string                  `gorm:"type:text"`
	FailureReason      string                  `gorm:"type:text"`
	StartedAt          *time.Time
	CompletedAt        *time.Time
	CreatedAt          time.Time
	UpdatedAt          time.Time
}

func (DeliverySource) TableName() string { return "gitops_delivery_sources" }

func (ClusterTargetGroup) TableName() string { return "gitops_cluster_target_groups" }

func (EnvironmentStage) TableName() string { return "gitops_environment_stages" }

func (ConfigurationOverlay) TableName() string { return "gitops_configuration_overlays" }

func (ApplicationDeliveryUnit) TableName() string { return "gitops_delivery_units" }

func (ReleaseRevision) TableName() string { return "gitops_release_revisions" }

func (DeliveryOperation) TableName() string { return "gitops_delivery_operations" }
