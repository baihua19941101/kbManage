package marketplace

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

func NewScopeService(
	bindingRepo *repository.ScopeRoleBindingRepository,
	projectRepo *repository.ProjectRepository,
	workspaceClusterRepo *repository.WorkspaceClusterRepository,
) *ScopeService {
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
		return ErrMarketplaceScopeDenied
	}
	if s == nil || s.scopeAccess == nil || strings.TrimSpace(permission) == "" {
		return nil
	}
	workspaceIDs, err := s.scopeAccess.ListWorkspaceIDsByPermission(ctx, userID, permission)
	if err != nil {
		return err
	}
	if len(workspaceIDs) == 0 {
		return ErrMarketplaceScopeDenied
	}
	return nil
}

func (s *ScopeService) EnsureScopePermission(ctx context.Context, userID uint64, scopeType string, scopeID uint64, permission string) error {
	if userID == 0 || scopeID == 0 {
		return ErrMarketplaceScopeDenied
	}
	if s == nil || s.scopeAccess == nil {
		return nil
	}
	switch strings.TrimSpace(scopeType) {
	case string(domain.ScopeTypeWorkspace):
		allowed, err := s.scopeAccess.HasScopePermission(ctx, userID, domain.ScopeTypeWorkspace, scopeID, 0, permission)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrMarketplaceScopeDenied
		}
	case string(domain.ScopeTypeProject):
		allowed, err := s.scopeAccess.HasScopePermission(ctx, userID, domain.ScopeTypeProject, 0, scopeID, permission)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrMarketplaceScopeDenied
		}
	case "cluster":
		allowed, err := s.scopeAccess.CanAccessClusterByPermission(ctx, userID, scopeID, permission)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrMarketplaceScopeDenied
		}
	default:
		return ErrMarketplaceInvalid
	}
	return nil
}

func (s *ScopeService) CanReadScope(ctx context.Context, userID uint64, scopeType string, scopeID uint64) bool {
	return s.EnsureScopePermission(ctx, userID, scopeType, scopeID, "marketplace:read") == nil
}
