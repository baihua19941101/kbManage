package backuprestore

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
	PermissionBackupRestoreRead         = "backuprestore:read"
	PermissionBackupRestoreManagePolicy = "backuprestore:manage-policy"
	PermissionBackupRestoreBackup       = "backuprestore:backup"
	PermissionBackupRestoreRestore      = "backuprestore:restore"
	PermissionBackupRestoreMigrate      = "backuprestore:migrate"
	PermissionBackupRestoreDrill        = "backuprestore:drill"
)

var backupRestoreRolePermissions = map[string]map[string]struct{}{
	"platform-admin": {
		PermissionBackupRestoreRead:         {},
		PermissionBackupRestoreManagePolicy: {},
		PermissionBackupRestoreBackup:       {},
		PermissionBackupRestoreRestore:      {},
		PermissionBackupRestoreMigrate:      {},
		PermissionBackupRestoreDrill:        {},
	},
	"ops-operator": {
		PermissionBackupRestoreRead:         {},
		PermissionBackupRestoreManagePolicy: {},
		PermissionBackupRestoreBackup:       {},
		PermissionBackupRestoreRestore:      {},
		PermissionBackupRestoreMigrate:      {},
		PermissionBackupRestoreDrill:        {},
	},
	"workspace-owner": {
		PermissionBackupRestoreRead:         {},
		PermissionBackupRestoreManagePolicy: {},
		PermissionBackupRestoreBackup:       {},
		PermissionBackupRestoreRestore:      {},
		PermissionBackupRestoreMigrate:      {},
		PermissionBackupRestoreDrill:        {},
	},
	"project-owner": {
		PermissionBackupRestoreRead:         {},
		PermissionBackupRestoreManagePolicy: {},
		PermissionBackupRestoreBackup:       {},
		PermissionBackupRestoreRestore:      {},
		PermissionBackupRestoreMigrate:      {},
		PermissionBackupRestoreDrill:        {},
	},
	"readonly": {
		PermissionBackupRestoreRead: {},
	},
	"audit-reader": {
		PermissionBackupRestoreRead: {},
	},
	"workspace-viewer": {
		PermissionBackupRestoreRead: {},
	},
	"project-viewer": {
		PermissionBackupRestoreRead: {},
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
	set := map[uint64]struct{}{}
	for _, binding := range bindings {
		if !hasBackupRestorePermission(binding.RoleKey, PermissionBackupRestoreRead) {
			continue
		}
		switch parseScopeType(binding.ScopeType) {
		case domain.ScopeTypeWorkspace:
			set[binding.ScopeID] = struct{}{}
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
	return sortedIDs(set), true, nil
}

func (s *ScopeService) EnsureRead(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionBackupRestoreRead)
}

func (s *ScopeService) EnsureManagePolicy(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionBackupRestoreManagePolicy)
}

func (s *ScopeService) EnsureBackup(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionBackupRestoreBackup)
}

func (s *ScopeService) EnsureRestore(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionBackupRestoreRestore)
}

func (s *ScopeService) EnsureMigrate(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionBackupRestoreMigrate)
}

func (s *ScopeService) EnsureDrill(ctx context.Context, userID, workspaceID, projectID uint64) error {
	return s.ensurePermission(ctx, userID, workspaceID, projectID, PermissionBackupRestoreDrill)
}

func (s *ScopeService) ensurePermission(ctx context.Context, userID, workspaceID, projectID uint64, permission string) error {
	if s == nil || s.bindings == nil {
		return nil
	}
	bindings, err := s.listUserBindings(ctx, userID)
	if err != nil {
		return err
	}
	for _, binding := range bindings {
		if !hasBackupRestorePermission(binding.RoleKey, permission) {
			continue
		}
		grantedType := parseScopeType(binding.ScopeType)
		grantedWorkspaceID, grantedProjectID, err := s.resolveGrantedScope(ctx, grantedType, binding.ScopeID)
		if err != nil {
			return err
		}
		if s.authorizer.CanAccessClusterLifecycle(grantedType, grantedWorkspaceID, grantedProjectID, workspaceID, projectID, false) {
			return nil
		}
	}
	return ErrBackupRestoreScopeDenied
}

func (s *ScopeService) listUserBindings(ctx context.Context, userID uint64) ([]repository.ScopeRoleBindingWithRole, error) {
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
	case string(domain.ScopeTypeProject):
		return domain.ScopeTypeProject
	case string(domain.ScopeTypeWorkspace):
		return domain.ScopeTypeWorkspace
	default:
		return domain.ScopeTypePlatform
	}
}

func hasBackupRestorePermission(roleKey, permission string) bool {
	perms, ok := backupRestoreRolePermissions[strings.TrimSpace(roleKey)]
	if !ok {
		return false
	}
	_, allowed := perms[permission]
	return allowed
}

func errorsIsNotFound(err error) bool {
	return err != nil && (err == gorm.ErrRecordNotFound || strings.Contains(strings.ToLower(err.Error()), "record not found"))
}

func sortedIDs(set map[uint64]struct{}) []uint64 {
	out := make([]uint64, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
