package identitytenancy

import (
	"context"
	"sort"
	"strconv"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"
)

const (
	PermissionIdentityRead         = "identity:read"
	PermissionIdentityManageSource = "identity:manage-source"
	PermissionIdentityManageOrg    = "identity:manage-org"
	PermissionIdentityManageRole   = "identity:manage-role"
	PermissionIdentityDelegate     = "identity:delegate"
	PermissionIdentitySession      = "identity:session-govern"
)

var identityRolePermissions = map[string]map[string]struct{}{
	"platform-admin": {
		PermissionIdentityRead:         {},
		PermissionIdentityManageSource: {},
		PermissionIdentityManageOrg:    {},
		PermissionIdentityManageRole:   {},
		PermissionIdentityDelegate:     {},
		PermissionIdentitySession:      {},
	},
	"ops-operator": {
		PermissionIdentityRead:         {},
		PermissionIdentityManageSource: {},
		PermissionIdentityManageOrg:    {},
		PermissionIdentityManageRole:   {},
		PermissionIdentityDelegate:     {},
		PermissionIdentitySession:      {},
	},
	"workspace-owner": {
		PermissionIdentityRead:         {},
		PermissionIdentityManageSource: {},
		PermissionIdentityManageOrg:    {},
		PermissionIdentityManageRole:   {},
		PermissionIdentityDelegate:     {},
		PermissionIdentitySession:      {},
	},
	"project-owner": {
		PermissionIdentityRead:       {},
		PermissionIdentityManageOrg:  {},
		PermissionIdentityManageRole: {},
		PermissionIdentityDelegate:   {},
		PermissionIdentitySession:    {},
	},
	"readonly": {
		PermissionIdentityRead: {},
	},
	"auditor": {
		PermissionIdentityRead: {},
	},
	"audit-reader": {
		PermissionIdentityRead: {},
	},
	"workspace-viewer": {
		PermissionIdentityRead: {},
	},
	"project-viewer": {
		PermissionIdentityRead: {},
	},
}

type ScopeService struct {
	bindings   *repository.ScopeRoleBindingRepository
	projects   *repository.ProjectRepository
	authorizer *authSvc.ScopeAuthorizer
}

func NewScopeService(bindingRepo *repository.ScopeRoleBindingRepository, projectRepo *repository.ProjectRepository) *ScopeService {
	return &ScopeService{
		bindings:   bindingRepo,
		projects:   projectRepo,
		authorizer: authSvc.NewScopeAuthorizer(),
	}
}

func (s *ScopeService) EnsureReadAny(ctx context.Context, userID uint64) error {
	return s.ensureAnyPermission(ctx, userID, PermissionIdentityRead)
}

func (s *ScopeService) EnsureManageSource(ctx context.Context, userID uint64) error {
	return s.ensureAnyPermission(ctx, userID, PermissionIdentityManageSource)
}

func (s *ScopeService) EnsureManageOrg(ctx context.Context, userID uint64, scopeType, scopeRef string) error {
	return s.ensureScopedPermission(ctx, userID, scopeType, scopeRef, PermissionIdentityManageOrg)
}

func (s *ScopeService) EnsureManageRole(ctx context.Context, userID uint64, scopeType, scopeRef string) error {
	return s.ensureScopedPermission(ctx, userID, scopeType, scopeRef, PermissionIdentityManageRole)
}

func (s *ScopeService) EnsureDelegate(ctx context.Context, userID uint64, scopeType, scopeRef string) error {
	return s.ensureScopedPermission(ctx, userID, scopeType, scopeRef, PermissionIdentityDelegate)
}

func (s *ScopeService) EnsureSessionGovern(ctx context.Context, userID uint64) error {
	return s.ensureAnyPermission(ctx, userID, PermissionIdentitySession)
}

func (s *ScopeService) DefaultScope(ctx context.Context, userID uint64) (string, string, error) {
	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return "", "", err
	}
	for _, binding := range bindings {
		if !hasIdentityPermission(binding.RoleKey, PermissionIdentityManageOrg) && !hasIdentityPermission(binding.RoleKey, PermissionIdentityManageRole) {
			continue
		}
		if binding.ScopeType == string(domain.ScopeTypeWorkspace) || binding.ScopeType == string(domain.ScopeTypeProject) {
			return binding.ScopeType, strconv.FormatUint(binding.ScopeID, 10), nil
		}
	}
	return string(domain.ScopeTypePlatform), "platform", nil
}

func (s *ScopeService) ListReadableWorkspaceIDs(ctx context.Context, userID uint64) ([]uint64, error) {
	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return nil, err
	}
	set := map[uint64]struct{}{}
	for _, binding := range bindings {
		if !hasIdentityPermission(binding.RoleKey, PermissionIdentityRead) {
			continue
		}
		switch binding.ScopeType {
		case string(domain.ScopeTypeWorkspace):
			if binding.ScopeID != 0 {
				set[binding.ScopeID] = struct{}{}
			}
		case string(domain.ScopeTypeProject):
			workspaceID, err := s.resolveProjectWorkspace(ctx, binding.ScopeID)
			if err != nil {
				return nil, err
			}
			if workspaceID != 0 {
				set[workspaceID] = struct{}{}
			}
		}
	}
	return sortedScopeIDs(set), nil
}

func (s *ScopeService) ensureAnyPermission(ctx context.Context, userID uint64, permission string) error {
	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return err
	}
	for _, binding := range bindings {
		if hasIdentityPermission(binding.RoleKey, permission) {
			return nil
		}
	}
	return ErrIdentityTenancyForbidden
}

func (s *ScopeService) ensureScopedPermission(ctx context.Context, userID uint64, scopeType, scopeRef, permission string) error {
	scopeType = strings.ToLower(strings.TrimSpace(scopeType))
	scopeRef = strings.TrimSpace(scopeRef)
	if scopeType == "" || scopeType == string(domain.ScopeTypePlatform) || scopeType == "organization" || scopeType == "resource" || scopeRef == "" || scopeRef == "platform" {
		return s.ensureAnyPermission(ctx, userID, permission)
	}
	targetWorkspaceID, targetProjectID := parseScopeNumericRef(scopeType, scopeRef)
	if targetProjectID != 0 && targetWorkspaceID == 0 {
		workspaceID, err := s.resolveProjectWorkspace(ctx, targetProjectID)
		if err != nil {
			return err
		}
		targetWorkspaceID = workspaceID
	}
	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return err
	}
	for _, binding := range bindings {
		if !hasIdentityPermission(binding.RoleKey, permission) {
			continue
		}
		grantedType := parseIdentityScopeType(binding.ScopeType)
		grantedWorkspaceID, grantedProjectID, err := s.resolveGrantedScope(ctx, grantedType, binding.ScopeID)
		if err != nil {
			return err
		}
		if s.authorizer.CanAccess(grantedType, grantedWorkspaceID, grantedProjectID, scopeTargetType(targetProjectID), targetWorkspaceID, targetProjectID) {
			return nil
		}
	}
	return ErrIdentityTenancyForbidden
}

func scopeTargetType(projectID uint64) domain.ScopeType {
	if projectID != 0 {
		return domain.ScopeTypeProject
	}
	return domain.ScopeTypeWorkspace
}

func parseIdentityScopeType(value string) domain.ScopeType {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(domain.ScopeTypeWorkspace):
		return domain.ScopeTypeWorkspace
	case string(domain.ScopeTypeProject):
		return domain.ScopeTypeProject
	default:
		return domain.ScopeTypePlatform
	}
}

func hasIdentityPermission(roleKey, permission string) bool {
	perms, ok := identityRolePermissions[strings.ToLower(strings.TrimSpace(roleKey))]
	if !ok {
		return false
	}
	_, allowed := perms[permission]
	return allowed
}

func (s *ScopeService) listUserBindings(ctx context.Context, userID uint64) ([]repository.ScopeRoleBindingWithRole, error) {
	if s == nil || s.bindings == nil || userID == 0 {
		return []repository.ScopeRoleBindingWithRole{}, nil
	}
	return s.bindings.List(ctx, repository.ScopeRoleBindingFilter{
		SubjectType: "user",
		SubjectID:   userID,
		Limit:       200,
	})
}

func (s *ScopeService) resolveGrantedScope(ctx context.Context, scopeType domain.ScopeType, scopeID uint64) (uint64, uint64, error) {
	switch scopeType {
	case domain.ScopeTypeWorkspace:
		return scopeID, 0, nil
	case domain.ScopeTypeProject:
		workspaceID, err := s.resolveProjectWorkspace(ctx, scopeID)
		return workspaceID, scopeID, err
	default:
		return 0, 0, nil
	}
}

func (s *ScopeService) resolveProjectWorkspace(ctx context.Context, projectID uint64) (uint64, error) {
	if projectID == 0 || s.projects == nil {
		return 0, nil
	}
	project, err := s.projects.GetByID(ctx, projectID)
	if err != nil {
		return 0, err
	}
	return project.WorkspaceID, nil
}

func sortedScopeIDs(set map[uint64]struct{}) []uint64 {
	out := make([]uint64, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
