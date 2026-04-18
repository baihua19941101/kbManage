package clusterlifecycle

import (
	"context"
	"sort"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"

	"gorm.io/gorm"
)

const (
	PermissionClusterLifecycleRead           = "clusterlifecycle:read"
	PermissionClusterLifecycleImport         = "clusterlifecycle:import"
	PermissionClusterLifecycleCreate         = "clusterlifecycle:create"
	PermissionClusterLifecycleUpgrade        = "clusterlifecycle:upgrade"
	PermissionClusterLifecycleManageNodePool = "clusterlifecycle:manage-nodepool"
	PermissionClusterLifecycleRetire         = "clusterlifecycle:retire"
	PermissionClusterLifecycleManageDriver   = "clusterlifecycle:manage-driver"
)

var clusterLifecycleRolePermissions = map[string]map[string]struct{}{
	"platform-admin": {
		PermissionClusterLifecycleRead:           {},
		PermissionClusterLifecycleImport:         {},
		PermissionClusterLifecycleCreate:         {},
		PermissionClusterLifecycleUpgrade:        {},
		PermissionClusterLifecycleManageNodePool: {},
		PermissionClusterLifecycleRetire:         {},
		PermissionClusterLifecycleManageDriver:   {},
	},
	"ops-operator": {
		PermissionClusterLifecycleRead:           {},
		PermissionClusterLifecycleImport:         {},
		PermissionClusterLifecycleCreate:         {},
		PermissionClusterLifecycleUpgrade:        {},
		PermissionClusterLifecycleManageNodePool: {},
		PermissionClusterLifecycleRetire:         {},
	},
	"workspace-owner": {
		PermissionClusterLifecycleRead:           {},
		PermissionClusterLifecycleImport:         {},
		PermissionClusterLifecycleCreate:         {},
		PermissionClusterLifecycleUpgrade:        {},
		PermissionClusterLifecycleManageNodePool: {},
		PermissionClusterLifecycleRetire:         {},
		PermissionClusterLifecycleManageDriver:   {},
	},
	"project-owner": {
		PermissionClusterLifecycleRead:           {},
		PermissionClusterLifecycleImport:         {},
		PermissionClusterLifecycleCreate:         {},
		PermissionClusterLifecycleUpgrade:        {},
		PermissionClusterLifecycleManageNodePool: {},
		PermissionClusterLifecycleRetire:         {},
	},
	"readonly": {
		PermissionClusterLifecycleRead: {},
	},
	"audit-reader": {
		PermissionClusterLifecycleRead: {},
	},
	"workspace-viewer": {
		PermissionClusterLifecycleRead: {},
	},
	"project-viewer": {
		PermissionClusterLifecycleRead: {},
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

func (s *ScopeService) ListReadableWorkspaceIDs(ctx context.Context, userID uint64) ([]uint64, bool, error) {
	if s == nil || s.bindings == nil {
		return []uint64{}, false, nil
	}
	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return nil, false, err
	}
	set := make(map[uint64]struct{})
	for _, binding := range bindings {
		if !hasClusterLifecyclePermission(binding.RoleKey, PermissionClusterLifecycleRead) {
			continue
		}
		switch parseScopeType(binding.ScopeType) {
		case domain.ScopeTypeWorkspace:
			if binding.ScopeID != 0 {
				set[binding.ScopeID] = struct{}{}
			}
		case domain.ScopeTypeProject:
			workspaceID, err := s.resolveProjectWorkspace(ctx, binding.ScopeID)
			if err != nil {
				return nil, false, err
			}
			if workspaceID != 0 {
				set[workspaceID] = struct{}{}
			}
		}
	}
	return setToSortedIDs(set), true, nil
}

func (s *ScopeService) EnsureReadableCluster(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionClusterLifecycleRead)
}

func (s *ScopeService) EnsureImportCluster(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionClusterLifecycleImport)
}

func (s *ScopeService) EnsureCreateCluster(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionClusterLifecycleCreate)
}

func (s *ScopeService) EnsureUpgradeCluster(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionClusterLifecycleUpgrade)
}

func (s *ScopeService) EnsureManageNodePool(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionClusterLifecycleManageNodePool)
}

func (s *ScopeService) EnsureRetireCluster(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionClusterLifecycleRetire)
}

func (s *ScopeService) EnsureManageDriver(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionClusterLifecycleManageDriver)
}

func (s *ScopeService) ensurePermission(ctx context.Context, userID, workspaceID, projectID uint64, permission string) error {
	if s == nil || s.bindings == nil {
		return nil
	}
	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return err
	}
	targetType := domain.ScopeTypeWorkspace
	if projectID != 0 {
		targetType = domain.ScopeTypeProject
	}
	for _, binding := range bindings {
		if !hasClusterLifecyclePermission(binding.RoleKey, permission) {
			continue
		}
		if workspaceID == 0 && projectID == 0 {
			return nil
		}
		grantedType := parseScopeType(binding.ScopeType)
		grantedWorkspaceID, grantedProjectID, err := s.resolveGrantedScope(ctx, grantedType, binding.ScopeID)
		if err != nil {
			return err
		}
		if s.authorizer.CanAccessClusterLifecycle(grantedType, grantedWorkspaceID, grantedProjectID, workspaceID, projectID, false) {
			return nil
		}
		if targetType == domain.ScopeTypeWorkspace && grantedType == domain.ScopeTypeWorkspace && binding.ScopeID == workspaceID {
			return nil
		}
	}
	return ErrLifecycleScopeDenied
}

func (s *ScopeService) listUserBindings(ctx context.Context, userID uint64) ([]repository.ScopeRoleBindingWithRole, error) {
	if userID == 0 {
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
		if errorsIsNotFound(err) {
			return 0, nil
		}
		return 0, err
	}
	return project.WorkspaceID, nil
}

func parseScopeType(value string) domain.ScopeType {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case string(domain.ScopeTypeWorkspace):
		return domain.ScopeTypeWorkspace
	case string(domain.ScopeTypeProject):
		return domain.ScopeTypeProject
	case string(domain.ScopeTypePlatform):
		return domain.ScopeTypePlatform
	default:
		return ""
	}
}

func hasClusterLifecyclePermission(roleKey, permission string) bool {
	grants, ok := clusterLifecycleRolePermissions[strings.ToLower(strings.TrimSpace(roleKey))]
	if !ok {
		return false
	}
	_, ok = grants[permission]
	return ok
}

func setToSortedIDs(set map[uint64]struct{}) []uint64 {
	out := make([]uint64, 0, len(set))
	for id := range set {
		if id != 0 {
			out = append(out, id)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}

func errorsIsNotFound(err error) bool {
	return err == gorm.ErrRecordNotFound
}
