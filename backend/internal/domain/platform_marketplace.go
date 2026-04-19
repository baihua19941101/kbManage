package domain

import "time"

type CatalogSourceType string

const (
	CatalogSourceTypeBuiltin CatalogSourceType = "builtin"
	CatalogSourceTypeGit     CatalogSourceType = "git"
	CatalogSourceTypeHelm    CatalogSourceType = "helm"
	CatalogSourceTypeOCI     CatalogSourceType = "oci"
)

type CatalogSourceStatus string

const (
	CatalogSourceStatusDraft    CatalogSourceStatus = "draft"
	CatalogSourceStatusActive   CatalogSourceStatus = "active"
	CatalogSourceStatusDegraded CatalogSourceStatus = "degraded"
	CatalogSourceStatusDisabled CatalogSourceStatus = "disabled"
)

type CatalogSyncState string

const (
	CatalogSyncStateIdle      CatalogSyncState = "idle"
	CatalogSyncStateSyncing   CatalogSyncState = "syncing"
	CatalogSyncStateSucceeded CatalogSyncState = "succeeded"
	CatalogSyncStateFailed    CatalogSyncState = "failed"
)

type TemplatePublishStatus string

const (
	TemplatePublishStatusActive      TemplatePublishStatus = "active"
	TemplatePublishStatusDisabled    TemplatePublishStatus = "disabled"
	TemplatePublishStatusRetired     TemplatePublishStatus = "retired"
	TemplatePublishStatusHistoryOnly TemplatePublishStatus = "history-only"
)

type TemplateVersionStatus string

const (
	TemplateVersionStatusDraft      TemplateVersionStatus = "draft"
	TemplateVersionStatusActive     TemplateVersionStatus = "active"
	TemplateVersionStatusDeprecated TemplateVersionStatus = "deprecated"
	TemplateVersionStatusRetired    TemplateVersionStatus = "retired"
)

type TemplateReleaseStatus string

const (
	TemplateReleaseStatusPublished TemplateReleaseStatus = "published"
	TemplateReleaseStatusWithdrawn TemplateReleaseStatus = "withdrawn"
	TemplateReleaseStatusRetired   TemplateReleaseStatus = "retired"
)

type InstallationLifecycleStatus string

const (
	InstallationLifecycleInstalled        InstallationLifecycleStatus = "installed"
	InstallationLifecycleUpgradeAvailable InstallationLifecycleStatus = "upgrade-available"
	InstallationLifecycleUpgraded         InstallationLifecycleStatus = "upgraded"
	InstallationLifecycleRetired          InstallationLifecycleStatus = "retired"
	InstallationLifecycleOrphaned         InstallationLifecycleStatus = "orphaned"
)

type ExtensionPackageStatus string

const (
	ExtensionPackageStatusDraft      ExtensionPackageStatus = "draft"
	ExtensionPackageStatusRegistered ExtensionPackageStatus = "registered"
	ExtensionPackageStatusEnabled    ExtensionPackageStatus = "enabled"
	ExtensionPackageStatusDisabled   ExtensionPackageStatus = "disabled"
	ExtensionPackageStatusRetired    ExtensionPackageStatus = "retired"
)

type CompatibilityOwnerType string

const (
	CompatibilityOwnerTemplateVersion CompatibilityOwnerType = "template-version"
	CompatibilityOwnerExtension       CompatibilityOwnerType = "extension-package"
)

type CompatibilityResult string

const (
	CompatibilityResultCompatible CompatibilityResult = "compatible"
	CompatibilityResultWarning    CompatibilityResult = "warning"
	CompatibilityResultBlocked    CompatibilityResult = "blocked"
)

type ExtensionLifecycleAction string

const (
	ExtensionLifecycleActionRegister ExtensionLifecycleAction = "register"
	ExtensionLifecycleActionEnable   ExtensionLifecycleAction = "enable"
	ExtensionLifecycleActionDisable  ExtensionLifecycleAction = "disable"
	ExtensionLifecycleActionRetire   ExtensionLifecycleAction = "retire"
)

type CatalogSource struct {
	ID              uint64              `gorm:"primaryKey" json:"id"`
	Name            string              `gorm:"size:128;not null;uniqueIndex" json:"name"`
	SourceType      CatalogSourceType   `gorm:"size:32;not null;index" json:"sourceType"`
	EndpointRef     string              `gorm:"size:255;not null" json:"endpointRef"`
	Status          CatalogSourceStatus `gorm:"size:32;not null;index" json:"status"`
	SyncState       CatalogSyncState    `gorm:"size:32;not null;index" json:"syncState"`
	LastSyncedAt    *time.Time          `json:"lastSyncedAt,omitempty"`
	LastError       string              `gorm:"type:text" json:"lastError,omitempty"`
	OwnerUserID     uint64              `gorm:"not null;index" json:"ownerUserId"`
	VisibilityScope string              `gorm:"size:64;not null;default:'platform'" json:"visibilityScope"`
	ConfigSummary   string              `gorm:"type:text" json:"configSummary,omitempty"`
	CreatedAt       time.Time           `json:"createdAt"`
	UpdatedAt       time.Time           `json:"updatedAt"`
}

func (CatalogSource) TableName() string { return "catalog_sources" }

type ApplicationTemplate struct {
	ID                  uint64                `gorm:"primaryKey" json:"id"`
	CatalogSourceID     uint64                `gorm:"not null;index;uniqueIndex:uk_template_source_slug,priority:1" json:"catalogSourceId"`
	Name                string                `gorm:"size:128;not null" json:"name"`
	Slug                string                `gorm:"size:128;not null;uniqueIndex:uk_template_source_slug,priority:2" json:"slug"`
	Category            string                `gorm:"size:64;not null;index" json:"category"`
	Summary             string                `gorm:"type:text" json:"summary,omitempty"`
	PublishStatus       TemplatePublishStatus `gorm:"size:32;not null;index" json:"publishStatus"`
	DefaultVersionID    *uint64               `gorm:"index" json:"defaultVersionId,omitempty"`
	SupportedScopes     string                `gorm:"type:text" json:"supportedScopes,omitempty"`
	ReleaseNotesSummary string                `gorm:"type:text" json:"releaseNotesSummary,omitempty"`
	CreatedAt           time.Time             `json:"createdAt"`
	UpdatedAt           time.Time             `json:"updatedAt"`
}

func (ApplicationTemplate) TableName() string { return "application_templates" }

type TemplateVersion struct {
	ID                          uint64                `gorm:"primaryKey" json:"id"`
	TemplateID                  uint64                `gorm:"not null;index;uniqueIndex:uk_template_version,priority:1" json:"templateId"`
	Version                     string                `gorm:"size:64;not null;uniqueIndex:uk_template_version,priority:2" json:"version"`
	Status                      TemplateVersionStatus `gorm:"size:32;not null;index" json:"status"`
	DependencySnapshot          string                `gorm:"type:text" json:"dependencySnapshot,omitempty"`
	ParameterSchemaSummary      string                `gorm:"type:text" json:"parameterSchemaSummary,omitempty"`
	DeploymentConstraintSummary string                `gorm:"type:text" json:"deploymentConstraintSummary,omitempty"`
	ReleaseNotes                string                `gorm:"type:text" json:"releaseNotes,omitempty"`
	IsUpgradeable               bool                  `gorm:"not null;default:true" json:"isUpgradeable"`
	SupersedesVersionID         *uint64               `gorm:"index" json:"supersedesVersionId,omitempty"`
	CreatedAt                   time.Time             `json:"createdAt"`
	UpdatedAt                   time.Time             `json:"updatedAt"`
}

func (TemplateVersion) TableName() string { return "template_versions" }

type TemplateReleaseScope struct {
	ID             uint64                `gorm:"primaryKey" json:"id"`
	TemplateID     uint64                `gorm:"not null;index" json:"templateId"`
	VersionID      uint64                `gorm:"not null;index" json:"versionId"`
	ScopeType      string                `gorm:"size:32;not null;index" json:"scopeType"`
	ScopeRef       string                `gorm:"size:128;not null;index" json:"scopeRef"`
	Status         TemplateReleaseStatus `gorm:"size:32;not null;index" json:"status"`
	VisibilityMode string                `gorm:"size:64;not null" json:"visibilityMode"`
	PublishedBy    uint64                `gorm:"not null;index" json:"publishedBy"`
	PublishedAt    time.Time             `json:"publishedAt"`
	WithdrawnAt    *time.Time            `json:"withdrawnAt,omitempty"`
	CreatedAt      time.Time             `json:"createdAt"`
	UpdatedAt      time.Time             `json:"updatedAt"`
}

func (TemplateReleaseScope) TableName() string { return "template_release_scopes" }

type InstallationRecord struct {
	ID                      uint64                      `gorm:"primaryKey" json:"id"`
	TemplateID              uint64                      `gorm:"not null;index" json:"templateId"`
	VersionID               uint64                      `gorm:"not null;index" json:"versionId"`
	ScopeType               string                      `gorm:"size:32;not null;index" json:"scopeType"`
	ScopeRef                string                      `gorm:"size:128;not null;index" json:"scopeRef"`
	ReleaseScopeID          uint64                      `gorm:"not null;index" json:"releaseScopeId"`
	LifecycleStatus         InstallationLifecycleStatus `gorm:"size:32;not null;index" json:"lifecycleStatus"`
	CurrentInstalledVersion string                      `gorm:"size:64;not null" json:"currentInstalledVersion"`
	UpgradeTargetVersion    string                      `gorm:"size:64" json:"upgradeTargetVersion,omitempty"`
	ChangeSummary           string                      `gorm:"type:text" json:"changeSummary,omitempty"`
	InstalledAt             time.Time                   `json:"installedAt"`
	LastChangedAt           time.Time                   `json:"lastChangedAt"`
	CreatedAt               time.Time                   `json:"createdAt"`
	UpdatedAt               time.Time                   `json:"updatedAt"`
}

func (InstallationRecord) TableName() string { return "installation_records" }

type ExtensionPackage struct {
	ID                    uint64                 `gorm:"primaryKey" json:"id"`
	Name                  string                 `gorm:"size:128;not null;uniqueIndex:uk_extension_name_version,priority:1" json:"name"`
	ExtensionType         string                 `gorm:"size:64;not null;index" json:"extensionType"`
	Version               string                 `gorm:"size:64;not null;uniqueIndex:uk_extension_name_version,priority:2" json:"version"`
	Status                ExtensionPackageStatus `gorm:"size:32;not null;index" json:"status"`
	CompatibilityPolicy   string                 `gorm:"type:text" json:"compatibilityPolicy,omitempty"`
	PermissionDeclaration string                 `gorm:"type:text" json:"permissionDeclaration,omitempty"`
	VisibilityScope       string                 `gorm:"size:128;not null" json:"visibilityScope"`
	EntrySummary          string                 `gorm:"type:text" json:"entrySummary,omitempty"`
	OwnerUserID           uint64                 `gorm:"not null;index" json:"ownerUserId"`
	CreatedAt             time.Time              `json:"createdAt"`
	UpdatedAt             time.Time              `json:"updatedAt"`
}

func (ExtensionPackage) TableName() string { return "extension_packages" }

type CompatibilityStatement struct {
	ID          uint64                 `gorm:"primaryKey" json:"id"`
	OwnerType   CompatibilityOwnerType `gorm:"size:32;not null;index;uniqueIndex:uk_compat_owner_target,priority:1" json:"ownerType"`
	OwnerRef    string                 `gorm:"size:128;not null;uniqueIndex:uk_compat_owner_target,priority:2" json:"ownerRef"`
	TargetType  string                 `gorm:"size:64;not null;uniqueIndex:uk_compat_owner_target,priority:3" json:"targetType"`
	TargetRef   string                 `gorm:"size:128;not null;uniqueIndex:uk_compat_owner_target,priority:4" json:"targetRef"`
	Result      CompatibilityResult    `gorm:"size:32;not null;index" json:"result"`
	Summary     string                 `gorm:"type:text" json:"summary"`
	EvaluatedAt time.Time              `json:"evaluatedAt"`
	Evaluator   string                 `gorm:"size:128;not null" json:"evaluator"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}

func (CompatibilityStatement) TableName() string { return "compatibility_statements" }

type ExtensionLifecycleRecord struct {
	ID                 uint64                   `gorm:"primaryKey" json:"id"`
	ExtensionPackageID uint64                   `gorm:"not null;index" json:"extensionPackageId"`
	Action             ExtensionLifecycleAction `gorm:"size:32;not null;index" json:"action"`
	ScopeType          string                   `gorm:"size:32;not null;index" json:"scopeType"`
	ScopeRef           string                   `gorm:"size:128;not null;index" json:"scopeRef"`
	Outcome            string                   `gorm:"size:32;not null" json:"outcome"`
	Reason             string                   `gorm:"type:text" json:"reason,omitempty"`
	ExecutedBy         uint64                   `gorm:"not null;index" json:"executedBy"`
	ExecutedAt         time.Time                `json:"executedAt"`
	CreatedAt          time.Time                `json:"createdAt"`
	UpdatedAt          time.Time                `json:"updatedAt"`
}

func (ExtensionLifecycleRecord) TableName() string { return "extension_lifecycle_records" }

type MarketplaceAuditEvent struct {
	ID             uint64    `gorm:"primaryKey" json:"id"`
	Action         string    `gorm:"size:128;not null;index" json:"action"`
	ActorUserID    uint64    `gorm:"not null;index" json:"actorUserId"`
	TargetType     string    `gorm:"size:64;not null;index" json:"targetType"`
	TargetRef      string    `gorm:"size:128;not null;index" json:"targetRef"`
	Outcome        string    `gorm:"size:32;not null;index" json:"outcome"`
	DetailSnapshot string    `gorm:"type:text" json:"detailSnapshot,omitempty"`
	OccurredAt     time.Time `json:"occurredAt"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
}

func (MarketplaceAuditEvent) TableName() string { return "marketplace_audit_events" }
