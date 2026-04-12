package domain

import "time"

type ObservabilityDataSourceType string

const (
	ObservabilityDataSourceTypeMetrics ObservabilityDataSourceType = "metrics"
	ObservabilityDataSourceTypeLogs    ObservabilityDataSourceType = "logs"
	ObservabilityDataSourceTypeAlerts  ObservabilityDataSourceType = "alerts"
	ObservabilityDataSourceTypeEvents  ObservabilityDataSourceType = "events"
)

type ObservabilityProviderKind string

const (
	ObservabilityProviderKindPrometheus   ObservabilityProviderKind = "prometheus-compatible"
	ObservabilityProviderKindLoki         ObservabilityProviderKind = "loki-compatible"
	ObservabilityProviderKindAlertmanager ObservabilityProviderKind = "alertmanager-compatible"
	ObservabilityProviderKindK8sEvents    ObservabilityProviderKind = "k8s-events"
)

type ObservabilityDataSourceStatus string

const (
	ObservabilityDataSourceStatusPending     ObservabilityDataSourceStatus = "pending"
	ObservabilityDataSourceStatusHealthy     ObservabilityDataSourceStatus = "healthy"
	ObservabilityDataSourceStatusDegraded    ObservabilityDataSourceStatus = "degraded"
	ObservabilityDataSourceStatusUnreachable ObservabilityDataSourceStatus = "unreachable"
	ObservabilityDataSourceStatusDisabled    ObservabilityDataSourceStatus = "disabled"
)

type AlertSeverity string

const (
	AlertSeverityInfo     AlertSeverity = "info"
	AlertSeverityWarning  AlertSeverity = "warning"
	AlertSeverityCritical AlertSeverity = "critical"
)

type AlertRuleStatus string

const (
	AlertRuleStatusEnabled  AlertRuleStatus = "enabled"
	AlertRuleStatusDisabled AlertRuleStatus = "disabled"
)

type SilenceWindowStatus string

const (
	SilenceWindowStatusScheduled SilenceWindowStatus = "scheduled"
	SilenceWindowStatusActive    SilenceWindowStatus = "active"
	SilenceWindowStatusExpired   SilenceWindowStatus = "expired"
	SilenceWindowStatusCanceled  SilenceWindowStatus = "canceled"
)

type AlertIncidentStatus string

const (
	AlertIncidentStatusFiring       AlertIncidentStatus = "firing"
	AlertIncidentStatusAcknowledged AlertIncidentStatus = "acknowledged"
	AlertIncidentStatusSilenced     AlertIncidentStatus = "silenced"
	AlertIncidentStatusResolved     AlertIncidentStatus = "resolved"
)

// ObservabilityDataSource stores external observability backend metadata.
type ObservabilityDataSource struct {
	ID             uint64                        `gorm:"primaryKey"`
	ClusterID      *uint64                       `gorm:"index"`
	Type           ObservabilityDataSourceType   `gorm:"size:32;not null"`
	ProviderKind   ObservabilityProviderKind     `gorm:"size:64;not null"`
	Name           string                        `gorm:"size:128;not null"`
	BaseURL        string                        `gorm:"size:1024;not null"`
	AuthSecretRef  string                        `gorm:"size:256"`
	Status         ObservabilityDataSourceStatus `gorm:"size:32;not null;default:pending"`
	LastVerifiedAt *time.Time
	LastError      string `gorm:"size:1024"`
	CreatedBy      uint64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type AlertRule struct {
	ID                   uint64          `gorm:"primaryKey"`
	Name                 string          `gorm:"size:128;not null"`
	Description          string          `gorm:"size:1024"`
	Severity             AlertSeverity   `gorm:"size:16;not null"`
	ScopeSnapshotJSON    string          `gorm:"type:longtext"`
	ConditionExpression  string          `gorm:"type:text;not null"`
	EvaluationWindow     string          `gorm:"size:64"`
	NotificationStrategy string          `gorm:"type:longtext"`
	Status               AlertRuleStatus `gorm:"size:32;not null;default:enabled"`
	CreatedBy            uint64
	CreatedAt            time.Time
	UpdatedAt            time.Time
}

type NotificationTarget struct {
	ID            uint64 `gorm:"primaryKey"`
	Name          string `gorm:"size:128;not null"`
	TargetType    string `gorm:"size:32;not null"`
	ConfigRef     string `gorm:"size:256"`
	ScopeSnapshot string `gorm:"type:longtext"`
	Status        string `gorm:"size:32;not null;default:active"`
	CreatedBy     uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type SilenceWindow struct {
	ID            uint64 `gorm:"primaryKey"`
	Name          string `gorm:"size:128;not null"`
	ScopeSnapshot string `gorm:"type:longtext"`
	Reason        string `gorm:"size:1024"`
	StartsAt      time.Time
	EndsAt        time.Time
	Status        SilenceWindowStatus `gorm:"size:32;not null;default:scheduled"`
	CreatedBy     uint64
	CanceledBy    *uint64
	CreatedAt     time.Time
	UpdatedAt     time.Time
}

type AlertIncidentSnapshot struct {
	ID                uint64              `gorm:"primaryKey"`
	SourceIncidentKey string              `gorm:"size:256;uniqueIndex;not null"`
	RuleID            *uint64             `gorm:"index"`
	ClusterID         *uint64             `gorm:"index"`
	WorkspaceID       *uint64             `gorm:"index"`
	ProjectID         *uint64             `gorm:"index"`
	ResourceKind      string              `gorm:"size:64"`
	ResourceName      string              `gorm:"size:255"`
	Namespace         string              `gorm:"size:255"`
	Severity          AlertSeverity       `gorm:"size:16;not null"`
	Status            AlertIncidentStatus `gorm:"size:32;not null"`
	Summary           string              `gorm:"size:1024"`
	StartsAt          *time.Time
	AcknowledgedAt    *time.Time
	ResolvedAt        *time.Time
	LastSyncedAt      *time.Time
	TimelineJSON      string `gorm:"type:longtext"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
}

type AlertHandlingRecord struct {
	ID         uint64 `gorm:"primaryKey"`
	IncidentID uint64 `gorm:"index;not null"`
	ActionType string `gorm:"size:32;not null"`
	Content    string `gorm:"type:text"`
	ActedBy    uint64
	ActedAt    time.Time
}

func (ObservabilityDataSource) TableName() string {
	return "observability_data_sources"
}

func (AlertRule) TableName() string {
	return "observability_alert_rules"
}

func (NotificationTarget) TableName() string {
	return "observability_notification_targets"
}

func (SilenceWindow) TableName() string {
	return "observability_silence_windows"
}

func (AlertIncidentSnapshot) TableName() string {
	return "observability_alert_incidents"
}

func (AlertHandlingRecord) TableName() string {
	return "observability_alert_handling_records"
}
