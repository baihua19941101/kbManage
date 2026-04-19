package identitytenancy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	identityint "kbmanage/backend/internal/integration/identity"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"

	"gorm.io/gorm"
)

const (
	ResourceTypeIdentitySource   = "identity-source"
	ResourceTypeOrganizationUnit = "organization-unit"
	ResourceTypeTenantMapping    = "tenant-scope-mapping"
	ResourceTypeRoleDefinition   = "role-definition"
	ResourceTypeRoleAssignment   = "role-assignment"
	ResourceTypeDelegationGrant  = "delegation-grant"
	ResourceTypeSessionRecord    = "session-record"
	ResourceTypeAccessRisk       = "access-risk"

	ActionIdentitySourceCreate  = "identitytenancy.source.create"
	ActionIdentitySourceRead    = "identitytenancy.source.read"
	ActionOrganizationCreate    = "identitytenancy.organization.create"
	ActionTenantMappingCreate   = "identitytenancy.mapping.create"
	ActionRoleDefinitionCreate  = "identitytenancy.role.create"
	ActionRoleAssignmentCreate  = "identitytenancy.assignment.create"
	ActionDelegationGrantCreate = "identitytenancy.delegation.create"
	ActionDelegationGrantRead   = "identitytenancy.delegation.read"
	ActionSessionGovernanceRead = "identitytenancy.session.read"
	ActionSessionRevoke         = "identitytenancy.session.revoke"
	ActionLoginModeUpdate       = "identitytenancy.login-mode.update"
	ActionAccessRiskQuery       = "identitytenancy.risk.query"
	ActionMembershipQuery       = "identitytenancy.membership.read"
)

var (
	ErrIdentityTenancyForbidden = errors.New("identity tenancy permission denied")
	ErrIdentityTenancyConflict  = errors.New("identity tenancy conflict")
	ErrIdentityTenancyInvalid   = errors.New("identity tenancy invalid request")
	ErrIdentityTenancyBlocked   = errors.New("identity tenancy request blocked")
)

type Service struct {
	sources      *repository.IdentitySourceRepository
	accounts     *repository.IdentityAccountRepository
	orgUnits     *repository.OrganizationUnitRepository
	memberships  *repository.OrganizationMembershipRepository
	mappings     *repository.TenantScopeMappingRepository
	roles        *repository.RoleDefinitionRepository
	assignments  *repository.RoleAssignmentRepository
	delegations  *repository.DelegationGrantRepository
	sessions     *repository.SessionRecordRepository
	risks        *repository.AccessRiskRepository
	localAudit   *repository.IdentityAuditRepository
	scope        *ScopeService
	sessionCache *SessionCache
	permCache    *PermissionCache
	revocations  *RevocationCoordinator
	provider     identityint.Provider
	syncProvider identityint.SyncProvider
	auditWriter  *auditSvc.EventWriter
}

func NewService(
	sourceRepo *repository.IdentitySourceRepository,
	accountRepo *repository.IdentityAccountRepository,
	orgRepo *repository.OrganizationUnitRepository,
	membershipRepo *repository.OrganizationMembershipRepository,
	mappingRepo *repository.TenantScopeMappingRepository,
	roleRepo *repository.RoleDefinitionRepository,
	assignmentRepo *repository.RoleAssignmentRepository,
	delegationRepo *repository.DelegationGrantRepository,
	sessionRepo *repository.SessionRecordRepository,
	riskRepo *repository.AccessRiskRepository,
	auditRepo *repository.IdentityAuditRepository,
	bindingRepo *repository.ScopeRoleBindingRepository,
	projectRepo *repository.ProjectRepository,
	sessionCache *SessionCache,
	permissionCache *PermissionCache,
	revocations *RevocationCoordinator,
	provider identityint.Provider,
	syncProvider identityint.SyncProvider,
	auditWriter *auditSvc.EventWriter,
) *Service {
	return &Service{
		sources:      sourceRepo,
		accounts:     accountRepo,
		orgUnits:     orgRepo,
		memberships:  membershipRepo,
		mappings:     mappingRepo,
		roles:        roleRepo,
		assignments:  assignmentRepo,
		delegations:  delegationRepo,
		sessions:     sessionRepo,
		risks:        riskRepo,
		localAudit:   auditRepo,
		scope:        NewScopeService(bindingRepo, projectRepo),
		sessionCache: sessionCache,
		permCache:    permissionCache,
		revocations:  revocations,
		provider:     provider,
		syncProvider: syncProvider,
		auditWriter:  auditWriter,
	}
}

type CreateIdentitySourceInput struct {
	Name       string `json:"name"`
	SourceType string `json:"sourceType"`
	LoginMode  string `json:"loginMode"`
	ScopeMode  string `json:"scopeMode"`
	Status     string `json:"status"`
}

type IdentitySourceListFilter struct {
	SourceType string
	Status     string
}

type SessionListFilter struct {
	Status    string
	RiskLevel string
}

type CreateOrganizationUnitInput struct {
	UnitType         string `json:"unitType"`
	Name             string `json:"name"`
	Description      string `json:"description"`
	ParentUnitID     uint64 `json:"parentUnitId"`
	IdentitySourceID uint64 `json:"identitySourceId"`
	Status           string `json:"status"`
}

type OrganizationUnitListFilter struct {
	UnitType     string
	ParentUnitID uint64
}

type CreateTenantScopeMappingInput struct {
	ScopeType       string `json:"scopeType"`
	ScopeRef        string `json:"scopeRef"`
	InheritanceMode string `json:"inheritanceMode"`
	Status          string `json:"status"`
}

type RoleDefinitionListFilter struct {
	RoleLevel string
	Status    string
}

type CreateRoleDefinitionInput struct {
	Name              string `json:"name"`
	RoleLevel         string `json:"roleLevel"`
	Description       string `json:"description"`
	PermissionSummary string `json:"permissionSummary"`
	InheritancePolicy string `json:"inheritancePolicy"`
	Delegable         bool   `json:"delegable"`
	Status            string `json:"status"`
}

type RoleAssignmentListFilter struct {
	SubjectRef string
	ScopeType  string
	Status     string
}

type CreateRoleAssignmentInput struct {
	SubjectType       string     `json:"subjectType"`
	SubjectRef        string     `json:"subjectRef"`
	RoleDefinitionID  uint64     `json:"roleDefinitionId"`
	ScopeType         string     `json:"scopeType"`
	ScopeRef          string     `json:"scopeRef"`
	SourceType        string     `json:"sourceType"`
	DelegationGrantID uint64     `json:"delegationGrantId"`
	ValidUntil        *time.Time `json:"validUntil"`
}

type CreateDelegationGrantInput struct {
	GrantorRef        string    `json:"grantorRef"`
	DelegateRef       string    `json:"delegateRef"`
	AllowedRoleLevels []string  `json:"allowedRoleLevels"`
	ValidFrom         time.Time `json:"validFrom"`
	ValidUntil        time.Time `json:"validUntil"`
	Reason            string    `json:"reason"`
}

type AccessRiskListFilter struct {
	SubjectType string
	Severity    string
}

func normalizeName(value string) string {
	return strings.TrimSpace(value)
}

func uint64PtrIf(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	return &v
}

func derefUint64(v *uint64) uint64 {
	if v == nil {
		return 0
	}
	return *v
}

func marshalJSON(v any) string {
	payload, _ := json.Marshal(v)
	return string(payload)
}

func parseScopeNumericRef(scopeType, scopeRef string) (uint64, uint64) {
	n, err := strconv.ParseUint(strings.TrimSpace(scopeRef), 10, 64)
	if err != nil {
		return 0, 0
	}
	switch strings.ToLower(strings.TrimSpace(scopeType)) {
	case string(domain.ScopeTypeWorkspace):
		return n, 0
	case string(domain.ScopeTypeProject):
		return 0, n
	default:
		return 0, 0
	}
}

func roleAssignmentStatusAt(now time.Time, validUntil *time.Time) domain.RoleAssignmentStatus {
	if validUntil != nil && !validUntil.After(now) {
		return domain.RoleAssignmentStatusExpired
	}
	return domain.RoleAssignmentStatusActive
}

func riskLevelForAssignment(role *domain.RoleDefinition, assignment *domain.RoleAssignment) domain.IdentityRiskLevel {
	if assignment == nil || role == nil {
		return domain.IdentityRiskLevelLow
	}
	if assignment.SourceType == domain.RoleAssignmentSourceDelegated || assignment.SourceType == domain.RoleAssignmentSourceTemporary {
		return domain.IdentityRiskLevelHigh
	}
	if role.InheritancePolicy == domain.RoleInheritanceDownwardAllowed {
		return domain.IdentityRiskLevelMedium
	}
	return domain.IdentityRiskLevelLow
}

func (s *Service) writeAudit(ctx context.Context, actorID uint64, action, targetType, targetRef string, outcome domain.IdentityAuditOutcome, details map[string]any) {
	if details == nil {
		details = map[string]any{}
	}
	event := &domain.IdentityGovernanceAuditEvent{
		Action:         action,
		ActorUserID:    actorID,
		TargetType:     targetType,
		TargetRef:      targetRef,
		Outcome:        outcome,
		DetailSnapshot: marshalJSON(details),
		OccurredAt:     time.Now(),
	}
	if s.localAudit != nil {
		_ = s.localAudit.Create(ctx, event)
	}
	if s.auditWriter != nil {
		outcomeMap := map[domain.IdentityAuditOutcome]domain.AuditOutcome{
			domain.IdentityAuditOutcomeSucceeded: domain.AuditOutcomeSuccess,
			domain.IdentityAuditOutcomeFailed:    domain.AuditOutcomeFailed,
			domain.IdentityAuditOutcomeBlocked:   domain.AuditOutcomeDenied,
		}
		_ = s.auditWriter.Write(ctx, "", &actorID, action, "identitytenancy", targetRef, outcomeMap[outcome], details)
	}
}

func notFoundOrNil(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func (s *Service) ensureLocalFallback(ctx context.Context, actorID uint64) error {
	if s.sources == nil {
		return nil
	}
	_, err := s.sources.FindLocal(ctx)
	if err == nil {
		return nil
	}
	if !notFoundOrNil(err) {
		return err
	}
	now := time.Now()
	local := &domain.IdentitySource{
		Name:          "本地管理员登录",
		SourceType:    domain.IdentitySourceTypeLocal,
		Status:        domain.IdentitySourceStatusActive,
		LoginMode:     domain.IdentityLoginModeFallback,
		ScopeMode:     domain.IdentityScopeModePlatform,
		SyncState:     domain.IdentitySyncStateSucceeded,
		ConfigSummary: `{"builtIn":true}`,
		OwnerUserID:   actorID,
		LastCheckedAt: &now,
	}
	if err := s.sources.Create(ctx, local); err != nil {
		return err
	}
	return s.ensureIdentityAccount(ctx, actorID, local.ID, fmt.Sprintf("local:%d", actorID))
}
