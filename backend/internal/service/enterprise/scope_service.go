package enterprise

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"
)

type ScopeService struct {
	scopeAccess *authSvc.ScopeAccessService
}

func NewScopeService(bindingRepo *repository.ScopeRoleBindingRepository, projectRepo *repository.ProjectRepository, workspaceClusterRepo *repository.WorkspaceClusterRepository) *ScopeService {
	return &ScopeService{
		scopeAccess: authSvc.NewScopeAccessService(
			bindingRepo,
			projectRepo,
			authSvc.NewScopeAuthorizer(),
			authSvc.NewPermissionService(),
			workspaceClusterRepo,
		),
	}
}

func (s *ScopeService) EnsurePermission(ctx context.Context, userID uint64, permission string) error {
	if userID == 0 {
		return ErrEnterpriseScopeDenied
	}
	if s == nil || s.scopeAccess == nil || strings.TrimSpace(permission) == "" {
		return nil
	}
	workspaceIDs, err := s.scopeAccess.ListWorkspaceIDsByPermission(ctx, userID, permission)
	if err != nil {
		return err
	}
	if len(workspaceIDs) == 0 {
		return ErrEnterpriseScopeDenied
	}
	return nil
}

func (s *ScopeService) EnsureScopePermission(ctx context.Context, userID, workspaceID uint64, projectID *uint64, permission string) error {
	if workspaceID == 0 {
		return s.EnsurePermission(ctx, userID, permission)
	}
	if s == nil || s.scopeAccess == nil {
		return nil
	}
	targetType := domain.ScopeTypeWorkspace
	targetProjectID := uint64(0)
	if projectID != nil && *projectID != 0 {
		targetType = domain.ScopeTypeProject
		targetProjectID = *projectID
	}
	allowed, err := s.scopeAccess.HasScopePermission(ctx, userID, targetType, workspaceID, targetProjectID, permission)
	if err != nil {
		return err
	}
	if !allowed {
		return ErrEnterpriseScopeDenied
	}
	return nil
}
