package domain

import "time"

type IdentitySourceType string

const (
	IdentitySourceTypeLocal IdentitySourceType = "local"
	IdentitySourceTypeSSO   IdentitySourceType = "sso"
	IdentitySourceTypeOIDC  IdentitySourceType = "oidc"
	IdentitySourceTypeLDAP  IdentitySourceType = "ldap"
)

type IdentitySourceStatus string

const (
	IdentitySourceStatusDraft       IdentitySourceStatus = "draft"
	IdentitySourceStatusActive      IdentitySourceStatus = "active"
	IdentitySourceStatusDisabled    IdentitySourceStatus = "disabled"
	IdentitySourceStatusUnavailable IdentitySourceStatus = "unavailable"
)

type IdentityLoginMode string

const (
	IdentityLoginModeExclusive IdentityLoginMode = "exclusive"
	IdentityLoginModeOptional  IdentityLoginMode = "optional"
	IdentityLoginModeFallback  IdentityLoginMode = "fallback"
)

type IdentityScopeMode string

const (
	IdentityScopeModePlatform              IdentityScopeMode = "platform"
	IdentityScopeModeSelectedOrganizations IdentityScopeMode = "selected-organizations"
)

type IdentitySyncState string

const (
	IdentitySyncStateIdle      IdentitySyncState = "idle"
	IdentitySyncStateSyncing   IdentitySyncState = "syncing"
	IdentitySyncStateSucceeded IdentitySyncState = "succeeded"
	IdentitySyncStateFailed    IdentitySyncState = "failed"
)

type IdentityAccountStatus string

const (
	IdentityAccountStatusActive   IdentityAccountStatus = "active"
	IdentityAccountStatusDisabled IdentityAccountStatus = "disabled"
	IdentityAccountStatusLocked   IdentityAccountStatus = "locked"
	IdentityAccountStatusOrphaned IdentityAccountStatus = "orphaned"
)

type IdentityPrincipalType string

const (
	IdentityPrincipalTypeUser           IdentityPrincipalType = "user"
	IdentityPrincipalTypeServiceAccount IdentityPrincipalType = "service-account"
)

type OrganizationUnitType string

const (
	OrganizationUnitTypeOrganization OrganizationUnitType = "organization"
	OrganizationUnitTypeTeam         OrganizationUnitType = "team"
	OrganizationUnitTypeUserGroup    OrganizationUnitType = "user-group"
)

type OrganizationUnitStatus string

const (
	OrganizationUnitStatusActive   OrganizationUnitStatus = "active"
	OrganizationUnitStatusDisabled OrganizationUnitStatus = "disabled"
	OrganizationUnitStatusArchived OrganizationUnitStatus = "archived"
)

type OrganizationMembershipRole string

const (
	OrganizationMembershipRoleOwner      OrganizationMembershipRole = "owner"
	OrganizationMembershipRoleMaintainer OrganizationMembershipRole = "maintainer"
	OrganizationMembershipRoleMember     OrganizationMembershipRole = "member"
	OrganizationMembershipRoleObserver   OrganizationMembershipRole = "observer"
)

type OrganizationMembershipStatus string

const (
	OrganizationMembershipStatusActive    OrganizationMembershipStatus = "active"
	OrganizationMembershipStatusSuspended OrganizationMembershipStatus = "suspended"
	OrganizationMembershipStatusRemoved   OrganizationMembershipStatus = "removed"
)

type TenantScopeType string

const (
	TenantScopeTypeWorkspace TenantScopeType = "workspace"
	TenantScopeTypeProject   TenantScopeType = "project"
	TenantScopeTypeResource  TenantScopeType = "resource"
)

type TenantInheritanceMode string

const (
	TenantInheritanceModeDirect     TenantInheritanceMode = "direct"
	TenantInheritanceModeInherited  TenantInheritanceMode = "inherited"
	TenantInheritanceModeRestricted TenantInheritanceMode = "restricted"
)

type TenantScopeMappingStatus string

const (
	TenantScopeMappingStatusActive  TenantScopeMappingStatus = "active"
	TenantScopeMappingStatusPending TenantScopeMappingStatus = "pending"
	TenantScopeMappingStatusRevoked TenantScopeMappingStatus = "revoked"
)

type RoleLevel string

const (
	RoleLevelPlatform     RoleLevel = "platform"
	RoleLevelOrganization RoleLevel = "organization"
	RoleLevelWorkspace    RoleLevel = "workspace"
	RoleLevelProject      RoleLevel = "project"
	RoleLevelResource     RoleLevel = "resource"
)

type RoleInheritancePolicy string

const (
	RoleInheritanceNone            RoleInheritancePolicy = "none"
	RoleInheritanceUpwardBlocked   RoleInheritancePolicy = "upward-blocked"
	RoleInheritanceDownwardAllowed RoleInheritancePolicy = "downward-allowed"
	RoleInheritanceBounded         RoleInheritancePolicy = "bounded"
)

type RoleDefinitionStatus string

const (
	RoleDefinitionStatusDraft      RoleDefinitionStatus = "draft"
	RoleDefinitionStatusActive     RoleDefinitionStatus = "active"
	RoleDefinitionStatusDeprecated RoleDefinitionStatus = "deprecated"
	RoleDefinitionStatusDisabled   RoleDefinitionStatus = "disabled"
)

type RoleAssignmentSourceType string

const (
	RoleAssignmentSourceDirect    RoleAssignmentSourceType = "direct"
	RoleAssignmentSourceInherited RoleAssignmentSourceType = "inherited"
	RoleAssignmentSourceDelegated RoleAssignmentSourceType = "delegated"
	RoleAssignmentSourceTemporary RoleAssignmentSourceType = "temporary"
)

type RoleAssignmentStatus string

const (
	RoleAssignmentStatusActive  RoleAssignmentStatus = "active"
	RoleAssignmentStatusPending RoleAssignmentStatus = "pending"
	RoleAssignmentStatusExpired RoleAssignmentStatus = "expired"
	RoleAssignmentStatusRevoked RoleAssignmentStatus = "revoked"
)

type DelegationGrantStatus string

const (
	DelegationGrantStatusActive    DelegationGrantStatus = "active"
	DelegationGrantStatusSuspended DelegationGrantStatus = "suspended"
	DelegationGrantStatusExpired   DelegationGrantStatus = "expired"
	DelegationGrantStatusRevoked   DelegationGrantStatus = "revoked"
)

type IdentitySessionStatus string

const (
	IdentitySessionStatusActive      IdentitySessionStatus = "active"
	IdentitySessionStatusIdle        IdentitySessionStatus = "idle"
	IdentitySessionStatusRevoked     IdentitySessionStatus = "revoked"
	IdentitySessionStatusExpired     IdentitySessionStatus = "expired"
	IdentitySessionStatusRiskBlocked IdentitySessionStatus = "risk-blocked"
)

type IdentityRiskLevel string

const (
	IdentityRiskLevelLow      IdentityRiskLevel = "low"
	IdentityRiskLevelMedium   IdentityRiskLevel = "medium"
	IdentityRiskLevelHigh     IdentityRiskLevel = "high"
	IdentityRiskLevelCritical IdentityRiskLevel = "critical"
)

type AccessRiskStatus string

const (
	AccessRiskStatusOpen         AccessRiskStatus = "open"
	AccessRiskStatusAcknowledged AccessRiskStatus = "acknowledged"
	AccessRiskStatusResolved     AccessRiskStatus = "resolved"
)

type IdentityAuditOutcome string

const (
	IdentityAuditOutcomeSucceeded IdentityAuditOutcome = "succeeded"
	IdentityAuditOutcomeFailed    IdentityAuditOutcome = "failed"
	IdentityAuditOutcomeBlocked   IdentityAuditOutcome = "blocked"
)

type IdentitySource struct {
	ID            uint64               `gorm:"primaryKey" json:"id"`
	Name          string               `gorm:"size:128;not null;uniqueIndex" json:"name"`
	SourceType    IdentitySourceType   `gorm:"size:32;not null;index" json:"sourceType"`
	Status        IdentitySourceStatus `gorm:"size:32;not null;index" json:"status"`
	LoginMode     IdentityLoginMode    `gorm:"size:32;not null" json:"loginMode"`
	ScopeMode     IdentityScopeMode    `gorm:"size:48;not null" json:"scopeMode"`
	SyncState     IdentitySyncState    `gorm:"size:32;not null" json:"syncState"`
	ConfigSummary string               `gorm:"type:text" json:"configSummary,omitempty"`
	LastError     string               `gorm:"type:text" json:"lastError,omitempty"`
	OwnerUserID   uint64               `gorm:"not null;index" json:"ownerUserId"`
	LastCheckedAt *time.Time           `json:"lastCheckedAt,omitempty"`
	CreatedAt     time.Time            `json:"createdAt"`
	UpdatedAt     time.Time            `json:"updatedAt"`
}

func (IdentitySource) TableName() string { return "identity_sources" }

type IdentityAccount struct {
	ID               uint64                `gorm:"primaryKey" json:"id"`
	UserID           uint64                `gorm:"not null;index" json:"userId"`
	IdentitySourceID uint64                `gorm:"not null;index" json:"identitySourceId"`
	ExternalRef      string                `gorm:"size:190;not null;uniqueIndex:uk_identity_account_ref" json:"externalRef"`
	PrincipalType    IdentityPrincipalType `gorm:"size:32;not null" json:"principalType"`
	Status           IdentityAccountStatus `gorm:"size:32;not null;index" json:"status"`
	LastLoginAt      *time.Time            `json:"lastLoginAt,omitempty"`
	CreatedAt        time.Time             `json:"createdAt"`
	UpdatedAt        time.Time             `json:"updatedAt"`
}

func (IdentityAccount) TableName() string { return "identity_accounts" }

type OrganizationUnit struct {
	ID               uint64                 `gorm:"primaryKey" json:"id"`
	UnitType         OrganizationUnitType   `gorm:"size:32;not null;index" json:"unitType"`
	Name             string                 `gorm:"size:128;not null;index:idx_org_unit_parent_name,priority:2" json:"name"`
	Description      string                 `gorm:"type:text" json:"description,omitempty"`
	ParentUnitID     *uint64                `gorm:"index:idx_org_unit_parent_name,priority:1" json:"parentUnitId,omitempty"`
	IdentitySourceID *uint64                `gorm:"index" json:"identitySourceId,omitempty"`
	OwnerUserID      uint64                 `gorm:"not null;index" json:"ownerUserId"`
	Status           OrganizationUnitStatus `gorm:"size:32;not null;index" json:"status"`
	CreatedAt        time.Time              `json:"createdAt"`
	UpdatedAt        time.Time              `json:"updatedAt"`
}

func (OrganizationUnit) TableName() string { return "organization_units" }

type OrganizationMembership struct {
	ID             uint64                       `gorm:"primaryKey" json:"id"`
	UnitID         uint64                       `gorm:"not null;index;uniqueIndex:uk_org_membership_active,priority:1" json:"unitId"`
	MemberType     string                       `gorm:"size:32;not null;uniqueIndex:uk_org_membership_active,priority:2" json:"memberType"`
	MemberRef      string                       `gorm:"size:190;not null;uniqueIndex:uk_org_membership_active,priority:3" json:"memberRef"`
	MembershipRole OrganizationMembershipRole   `gorm:"size:32;not null" json:"membershipRole"`
	Status         OrganizationMembershipStatus `gorm:"size:32;not null;index;uniqueIndex:uk_org_membership_active,priority:4" json:"status"`
	JoinedAt       time.Time                    `gorm:"not null" json:"joinedAt"`
	CreatedAt      time.Time                    `json:"createdAt"`
	UpdatedAt      time.Time                    `json:"updatedAt"`
}

func (OrganizationMembership) TableName() string { return "organization_memberships" }

type TenantScopeMapping struct {
	ID              uint64                   `gorm:"primaryKey" json:"id"`
	UnitID          uint64                   `gorm:"not null;index;uniqueIndex:uk_tenant_scope_mapping,priority:1" json:"unitId"`
	ScopeType       TenantScopeType          `gorm:"size:32;not null;uniqueIndex:uk_tenant_scope_mapping,priority:2" json:"scopeType"`
	ScopeRef        string                   `gorm:"size:190;not null;uniqueIndex:uk_tenant_scope_mapping,priority:3" json:"scopeRef"`
	InheritanceMode TenantInheritanceMode    `gorm:"size:32;not null" json:"inheritanceMode"`
	Status          TenantScopeMappingStatus `gorm:"size:32;not null;index" json:"status"`
	CreatedBy       uint64                   `gorm:"not null;index" json:"createdBy"`
	CreatedAt       time.Time                `json:"createdAt"`
	UpdatedAt       time.Time                `json:"updatedAt"`
}

func (TenantScopeMapping) TableName() string { return "tenant_scope_mappings" }

type RoleDefinition struct {
	ID                uint64                `gorm:"primaryKey" json:"id"`
	Name              string                `gorm:"size:128;not null;uniqueIndex:uk_role_definition_level_name,priority:2" json:"name"`
	RoleLevel         RoleLevel             `gorm:"size:32;not null;index;uniqueIndex:uk_role_definition_level_name,priority:1" json:"roleLevel"`
	Description       string                `gorm:"type:text" json:"description,omitempty"`
	PermissionSummary string                `gorm:"type:text;not null" json:"permissionSummary"`
	InheritancePolicy RoleInheritancePolicy `gorm:"size:32;not null" json:"inheritancePolicy"`
	Delegable         bool                  `gorm:"not null;default:false" json:"delegable"`
	Status            RoleDefinitionStatus  `gorm:"size:32;not null;index" json:"status"`
	CreatedBy         uint64                `gorm:"not null;index" json:"createdBy"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
}

func (RoleDefinition) TableName() string { return "role_definitions" }

type RoleAssignment struct {
	ID                uint64                   `gorm:"primaryKey" json:"id"`
	SubjectType       string                   `gorm:"size:32;not null;index" json:"subjectType"`
	SubjectRef        string                   `gorm:"size:190;not null;index" json:"subjectRef"`
	RoleDefinitionID  uint64                   `gorm:"not null;index" json:"roleDefinitionId"`
	ScopeType         string                   `gorm:"size:32;not null;index" json:"scopeType"`
	ScopeRef          string                   `gorm:"size:190;not null;index" json:"scopeRef"`
	SourceType        RoleAssignmentSourceType `gorm:"size:32;not null;index" json:"sourceType"`
	DelegationGrantID *uint64                  `gorm:"index" json:"delegationGrantId,omitempty"`
	ValidFrom         time.Time                `gorm:"not null" json:"validFrom"`
	ValidUntil        *time.Time               `json:"validUntil,omitempty"`
	Status            RoleAssignmentStatus     `gorm:"size:32;not null;index" json:"status"`
	GrantedBy         uint64                   `gorm:"not null;index" json:"grantedBy"`
	CreatedAt         time.Time                `json:"createdAt"`
	UpdatedAt         time.Time                `json:"updatedAt"`
}

func (RoleAssignment) TableName() string { return "role_assignments" }

type DelegationGrant struct {
	ID                   uint64                `gorm:"primaryKey" json:"id"`
	GrantorRef           string                `gorm:"size:190;not null;index" json:"grantorRef"`
	DelegateRef          string                `gorm:"size:190;not null;index" json:"delegateRef"`
	AllowedRoleLevels    string                `gorm:"type:text;not null" json:"allowedRoleLevels"`
	AllowedScopeSnapshot string                `gorm:"type:text;not null" json:"allowedScopeSnapshot"`
	Status               DelegationGrantStatus `gorm:"size:32;not null;index" json:"status"`
	ValidFrom            time.Time             `gorm:"not null" json:"validFrom"`
	ValidUntil           time.Time             `gorm:"not null" json:"validUntil"`
	Reason               string                `gorm:"type:text" json:"reason,omitempty"`
	CreatedBy            uint64                `gorm:"not null;index" json:"createdBy"`
	CreatedAt            time.Time             `json:"createdAt"`
	UpdatedAt            time.Time             `json:"updatedAt"`
}

func (DelegationGrant) TableName() string { return "delegation_grants" }

type SessionRecord struct {
	ID                uint64                `gorm:"primaryKey" json:"id"`
	UserID            uint64                `gorm:"not null;index" json:"userId"`
	IdentitySourceID  uint64                `gorm:"not null;index" json:"identitySourceId"`
	LoginMethod       string                `gorm:"size:64;not null" json:"loginMethod"`
	Status            IdentitySessionStatus `gorm:"size:32;not null;index" json:"status"`
	RiskLevel         IdentityRiskLevel     `gorm:"size:32;not null;index" json:"riskLevel"`
	PermissionVersion string                `gorm:"size:64;not null" json:"permissionVersion"`
	LastSeenAt        *time.Time            `json:"lastSeenAt,omitempty"`
	RevokedAt         *time.Time            `json:"revokedAt,omitempty"`
	CreatedAt         time.Time             `json:"createdAt"`
	UpdatedAt         time.Time             `json:"updatedAt"`
}

func (SessionRecord) TableName() string { return "identity_session_records" }

type AccessRiskSnapshot struct {
	ID                uint64            `gorm:"primaryKey" json:"id"`
	SubjectType       string            `gorm:"size:32;not null;index" json:"subjectType"`
	SubjectRef        string            `gorm:"size:190;not null;index" json:"subjectRef"`
	RiskType          string            `gorm:"size:64;not null;index" json:"riskType"`
	Severity          IdentityRiskLevel `gorm:"size:32;not null;index" json:"severity"`
	Summary           string            `gorm:"type:text;not null" json:"summary"`
	RecommendedAction string            `gorm:"type:text" json:"recommendedAction,omitempty"`
	Status            AccessRiskStatus  `gorm:"size:32;not null;index" json:"status"`
	GeneratedAt       time.Time         `gorm:"not null;index" json:"generatedAt"`
	CreatedAt         time.Time         `json:"createdAt"`
	UpdatedAt         time.Time         `json:"updatedAt"`
}

func (AccessRiskSnapshot) TableName() string { return "access_risk_snapshots" }

type IdentityGovernanceAuditEvent struct {
	ID             uint64               `gorm:"primaryKey" json:"id"`
	Action         string               `gorm:"size:128;not null;index" json:"action"`
	ActorUserID    uint64               `gorm:"not null;index" json:"actorUserId"`
	TargetType     string               `gorm:"size:64;not null;index" json:"targetType"`
	TargetRef      string               `gorm:"size:190;not null;index" json:"targetRef"`
	Outcome        IdentityAuditOutcome `gorm:"size:32;not null;index" json:"outcome"`
	DetailSnapshot string               `gorm:"type:text" json:"detailSnapshot,omitempty"`
	OccurredAt     time.Time            `gorm:"not null;index" json:"occurredAt"`
	CreatedAt      time.Time            `json:"createdAt"`
	UpdatedAt      time.Time            `json:"updatedAt"`
}

func (IdentityGovernanceAuditEvent) TableName() string { return "identity_governance_audit_events" }
