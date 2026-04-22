package domain

import "time"

type GovernanceRiskStatus string

const (
	GovernanceRiskStatusOpen       GovernanceRiskStatus = "open"
	GovernanceRiskStatusReviewing  GovernanceRiskStatus = "reviewing"
	GovernanceRiskStatusMitigating GovernanceRiskStatus = "mitigating"
	GovernanceRiskStatusClosed     GovernanceRiskStatus = "closed"
)

type GovernanceReportStatus string

const (
	GovernanceReportStatusDraft      GovernanceReportStatus = "draft"
	GovernanceReportStatusGenerating GovernanceReportStatus = "generating"
	GovernanceReportStatusReady      GovernanceReportStatus = "ready"
	GovernanceReportStatusExported   GovernanceReportStatus = "exported"
	GovernanceReportStatusArchived   GovernanceReportStatus = "archived"
)

type DeliveryArtifactStatus string

const (
	DeliveryArtifactStatusDraft      DeliveryArtifactStatus = "draft"
	DeliveryArtifactStatusActive     DeliveryArtifactStatus = "active"
	DeliveryArtifactStatusSuperseded DeliveryArtifactStatus = "superseded"
	DeliveryArtifactStatusArchived   DeliveryArtifactStatus = "archived"
)

type DeliveryReadinessConclusion string

const (
	DeliveryReadinessNotReady      DeliveryReadinessConclusion = "not-ready"
	DeliveryReadinessConditionally DeliveryReadinessConclusion = "conditionally-ready"
	DeliveryReadinessReady         DeliveryReadinessConclusion = "ready"
)

type GovernanceActionStatus string

const (
	GovernanceActionStatusOpen       GovernanceActionStatus = "open"
	GovernanceActionStatusInProgress GovernanceActionStatus = "in-progress"
	GovernanceActionStatusBlocked    GovernanceActionStatus = "blocked"
	GovernanceActionStatusResolved   GovernanceActionStatus = "resolved"
	GovernanceActionStatusClosed     GovernanceActionStatus = "closed"
)

type PermissionChangeTrail struct {
	ID                   uint64    `gorm:"primaryKey" json:"id"`
	WorkspaceID          uint64    `gorm:"index" json:"workspaceId"`
	ProjectID            *uint64   `gorm:"index" json:"projectId,omitempty"`
	SubjectType          string    `gorm:"size:64;not null" json:"subjectType"`
	SubjectRef           string    `gorm:"size:128;not null" json:"subjectRef"`
	SourceIdentity       string    `gorm:"size:128" json:"sourceIdentity"`
	ChangeType           string    `gorm:"size:64;not null" json:"changeType"`
	BeforeState          string    `gorm:"type:text" json:"beforeState"`
	AfterState           string    `gorm:"type:text" json:"afterState"`
	AuthorizationBasis   string    `gorm:"type:text" json:"authorizationBasis"`
	ApprovalReference    string    `gorm:"type:text" json:"approvalReference"`
	ScopeType            string    `gorm:"size:64;not null" json:"scopeType"`
	ScopeRef             string    `gorm:"size:128;not null" json:"scopeRef"`
	EvidenceCompleteness string    `gorm:"size:32;not null" json:"evidenceCompleteness"`
	ChangedAt            time.Time `gorm:"index" json:"changedAt"`
	ChangedBy            uint64    `gorm:"index;not null" json:"changedBy"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

func (PermissionChangeTrail) TableName() string { return "enterprise_permission_change_trails" }

type KeyOperationTrace struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	WorkspaceID    uint64    `gorm:"index" json:"workspaceId"`
	ProjectID      *uint64   `gorm:"index" json:"projectId,omitempty"`
	ActorType      string    `gorm:"size:64;not null" json:"actorType"`
	ActorRef       string    `gorm:"size:128;not null" json:"actorRef"`
	OperationType  string    `gorm:"size:64;not null" json:"operationType"`
	TargetType     string    `gorm:"size:64;not null" json:"targetType"`
	TargetRef      string    `gorm:"size:128;not null" json:"targetRef"`
	ContextSummary string    `gorm:"type:text" json:"contextSummary"`
	RiskLevel      string    `gorm:"size:32;not null" json:"riskLevel"`
	Outcome        string    `gorm:"size:32;not null" json:"outcome"`
	OccurredAt     time.Time `gorm:"index" json:"occurredAt"`
	RelatedTrailID *uint64   `gorm:"index" json:"relatedTrailId,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (KeyOperationTrace) TableName() string { return "enterprise_key_operation_traces" }

type CrossTeamAuthorizationSnapshot struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	WorkspaceID    uint64    `gorm:"index" json:"workspaceId"`
	ProjectID      *uint64   `gorm:"index" json:"projectId,omitempty"`
	SnapshotAt     time.Time `gorm:"index" json:"snapshotAt"`
	SourceTeam     string    `gorm:"size:128;not null" json:"sourceTeam"`
	TargetTeam     string    `gorm:"size:128;not null" json:"targetTeam"`
	GrantType      string    `gorm:"size:64;not null" json:"grantType"`
	ScopeSummary   string    `gorm:"type:text" json:"scopeSummary"`
	Temporality    string    `gorm:"size:64" json:"temporality"`
	DelegationFlag bool      `gorm:"not null;default:false" json:"delegationFlag"`
	RiskHint       string    `gorm:"type:text" json:"riskHint"`
	TrendLabel     string    `gorm:"size:128" json:"trendLabel"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (CrossTeamAuthorizationSnapshot) TableName() string {
	return "enterprise_cross_team_authorization_snapshots"
}

type GovernanceRiskEvent struct {
	ID                uint64               `gorm:"primaryKey" json:"id"`
	WorkspaceID       uint64               `gorm:"index" json:"workspaceId"`
	ProjectID         *uint64              `gorm:"index" json:"projectId,omitempty"`
	RiskType          string               `gorm:"size:64;not null" json:"riskType"`
	Severity          string               `gorm:"size:32;not null" json:"severity"`
	SubjectSummary    string               `gorm:"type:text" json:"subjectSummary"`
	ScopeSummary      string               `gorm:"type:text" json:"scopeSummary"`
	TriggerReason     string               `gorm:"type:text" json:"triggerReason"`
	Status            GovernanceRiskStatus `gorm:"size:32;not null;index" json:"status"`
	RecommendedAction string               `gorm:"type:text" json:"recommendedAction"`
	Owner             string               `gorm:"size:128" json:"owner"`
	FirstSeenAt       time.Time            `gorm:"index" json:"firstSeenAt"`
	LastSeenAt        time.Time            `gorm:"index" json:"lastSeenAt"`
	CreatedAt         time.Time            `json:"createdAt"`
	UpdatedAt         time.Time            `json:"updatedAt"`
}

func (GovernanceRiskEvent) TableName() string { return "enterprise_governance_risk_events" }

type GovernanceCoverageSnapshot struct {
	ID                   uint64    `gorm:"primaryKey" json:"id"`
	WorkspaceID          uint64    `gorm:"index" json:"workspaceId"`
	ProjectID            *uint64   `gorm:"index" json:"projectId,omitempty"`
	SnapshotAt           time.Time `gorm:"index" json:"snapshotAt"`
	CoverageDomain       string    `gorm:"size:64;not null" json:"coverageDomain"`
	CoverageRate         float64   `gorm:"not null;default:0" json:"coverageRate"`
	StatusBreakdown      string    `gorm:"type:text" json:"statusBreakdown"`
	MissingReasonSummary string    `gorm:"type:text" json:"missingReasonSummary"`
	ConfidenceLevel      string    `gorm:"size:32;not null" json:"confidenceLevel"`
	TrendSummary         string    `gorm:"type:text" json:"trendSummary"`
	Owner                string    `gorm:"size:128" json:"owner"`
	CreatedAt            time.Time `json:"createdAt"`
	UpdatedAt            time.Time `json:"updatedAt"`
}

func (GovernanceCoverageSnapshot) TableName() string {
	return "enterprise_governance_coverage_snapshots"
}

type GovernanceReportPackage struct {
	ID                uint64                 `gorm:"primaryKey" json:"id"`
	WorkspaceID       uint64                 `gorm:"index" json:"workspaceId"`
	ProjectID         *uint64                `gorm:"index" json:"projectId,omitempty"`
	ReportType        string                 `gorm:"size:64;not null" json:"reportType"`
	Title             string                 `gorm:"size:128;not null" json:"title"`
	AudienceType      string                 `gorm:"size:64;not null" json:"audienceType"`
	TimeRange         string                 `gorm:"size:128" json:"timeRange"`
	SummarySection    string                 `gorm:"type:text" json:"summarySection"`
	DetailSection     string                 `gorm:"type:text" json:"detailSection"`
	AttachmentCatalog string                 `gorm:"type:text" json:"attachmentCatalog"`
	VisibilityPolicy  string                 `gorm:"type:text" json:"visibilityPolicy"`
	GeneratedAt       time.Time              `gorm:"index" json:"generatedAt"`
	GeneratedBy       uint64                 `gorm:"index;not null" json:"generatedBy"`
	Status            GovernanceReportStatus `gorm:"size:32;not null;index" json:"status"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

func (GovernanceReportPackage) TableName() string { return "enterprise_governance_report_packages" }

type ExportRecord struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	PackageID      uint64    `gorm:"index;not null" json:"packageId"`
	ExportType     string    `gorm:"size:64;not null" json:"exportType"`
	AudienceScope  string    `gorm:"size:128;not null" json:"audienceScope"`
	ContentLevel   string    `gorm:"size:64;not null" json:"contentLevel"`
	Result         string    `gorm:"size:32;not null" json:"result"`
	AuditReference string    `gorm:"size:128" json:"auditReference"`
	ExportedAt     time.Time `gorm:"index" json:"exportedAt"`
	ExportedBy     uint64    `gorm:"index;not null" json:"exportedBy"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (ExportRecord) TableName() string { return "enterprise_export_records" }

type DeliveryArtifact struct {
	ID                uint64                 `gorm:"primaryKey" json:"id"`
	WorkspaceID       uint64                 `gorm:"index" json:"workspaceId"`
	ProjectID         *uint64                `gorm:"index" json:"projectId,omitempty"`
	ArtifactType      string                 `gorm:"size:64;not null" json:"artifactType"`
	Title             string                 `gorm:"size:128;not null" json:"title"`
	VersionScope      string                 `gorm:"size:128;not null" json:"versionScope"`
	EnvironmentScope  string                 `gorm:"size:128;not null" json:"environmentScope"`
	OwnerRole         string                 `gorm:"size:128;not null" json:"ownerRole"`
	ApplicabilityNote string                 `gorm:"type:text" json:"applicabilityNote"`
	Status            DeliveryArtifactStatus `gorm:"size:32;not null;index" json:"status"`
	UpdatedAt         time.Time              `json:"updatedAt"`
	CreatedAt         time.Time              `json:"createdAt"`
}

func (DeliveryArtifact) TableName() string { return "enterprise_delivery_artifacts" }

type DeliveryReadinessBundle struct {
	ID                  uint64                      `gorm:"primaryKey" json:"id"`
	WorkspaceID         uint64                      `gorm:"index" json:"workspaceId"`
	ProjectID           *uint64                     `gorm:"index" json:"projectId,omitempty"`
	Name                string                      `gorm:"size:128;not null" json:"name"`
	TargetEnvironment   string                      `gorm:"size:128;not null" json:"targetEnvironment"`
	TargetAudience      string                      `gorm:"size:128;not null" json:"targetAudience"`
	ArtifactSummary     string                      `gorm:"type:text" json:"artifactSummary"`
	ChecklistStatus     string                      `gorm:"size:64" json:"checklistStatus"`
	MissingItems        string                      `gorm:"type:text" json:"missingItems"`
	ReadinessConclusion DeliveryReadinessConclusion `gorm:"size:32;not null;index" json:"readinessConclusion"`
	UpdatedAt           time.Time                   `json:"updatedAt"`
	CreatedAt           time.Time                   `json:"createdAt"`
}

func (DeliveryReadinessBundle) TableName() string { return "enterprise_delivery_readiness_bundles" }

type DeliveryChecklistItem struct {
	ID                  uint64     `gorm:"primaryKey" json:"id"`
	BundleID            uint64     `gorm:"index;not null" json:"bundleId"`
	CheckItem           string     `gorm:"size:255;not null" json:"checkItem"`
	Category            string     `gorm:"size:64;not null" json:"category"`
	Owner               string     `gorm:"size:128;not null" json:"owner"`
	EvidenceRequirement string     `gorm:"type:text" json:"evidenceRequirement"`
	Status              string     `gorm:"size:32;not null" json:"status"`
	Remark              string     `gorm:"type:text" json:"remark"`
	CompletedAt         *time.Time `json:"completedAt,omitempty"`
	CreatedAt           time.Time  `json:"createdAt"`
	UpdatedAt           time.Time  `json:"updatedAt"`
}

func (DeliveryChecklistItem) TableName() string { return "enterprise_delivery_checklist_items" }

type GovernanceActionItem struct {
	ID                uint64                 `gorm:"primaryKey" json:"id"`
	WorkspaceID       uint64                 `gorm:"index" json:"workspaceId"`
	ProjectID         *uint64                `gorm:"index" json:"projectId,omitempty"`
	SourceType        string                 `gorm:"size:64;not null" json:"sourceType"`
	SourceRef         string                 `gorm:"size:128;not null" json:"sourceRef"`
	Title             string                 `gorm:"size:255;not null" json:"title"`
	Priority          string                 `gorm:"size:32;not null" json:"priority"`
	Owner             string                 `gorm:"size:128" json:"owner"`
	Status            GovernanceActionStatus `gorm:"size:32;not null;index" json:"status"`
	ResolutionSummary string                 `gorm:"type:text" json:"resolutionSummary"`
	DueAt             *time.Time             `json:"dueAt,omitempty"`
	CreatedAt         time.Time              `json:"createdAt"`
	UpdatedAt         time.Time              `json:"updatedAt"`
}

func (GovernanceActionItem) TableName() string { return "enterprise_governance_action_items" }
