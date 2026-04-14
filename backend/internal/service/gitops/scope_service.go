package gitops

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"
	authSvc "kbmanage/backend/internal/service/auth"
)

type ScopeService struct {
	scopeAccess *authSvc.ScopeAccessService
}

func NewScopeService(scopeAccess *authSvc.ScopeAccessService) *ScopeService {
	return &ScopeService{scopeAccess: scopeAccess}
}

func (s *ScopeService) ValidateScope(
	ctx context.Context,
	userID uint64,
	workspaceID uint64,
	projectID uint64,
	permission string,
) error {
	if s == nil || s.scopeAccess == nil {
		return nil
	}
	if userID == 0 {
		return ErrGitOpsScopeDenied
	}
	permission = strings.TrimSpace(permission)
	if permission == "" {
		permission = PermissionGitOpsRead
	}

	if projectID != 0 {
		allowed, err := s.scopeAccess.HasScopePermission(
			ctx,
			userID,
			domain.ScopeTypeProject,
			workspaceID,
			projectID,
			permission,
		)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrGitOpsScopeDenied
		}
		return nil
	}

	if workspaceID != 0 {
		allowed, err := s.scopeAccess.HasScopePermission(
			ctx,
			userID,
			domain.ScopeTypeWorkspace,
			workspaceID,
			0,
			permission,
		)
		if err != nil {
			return err
		}
		if !allowed {
			return ErrGitOpsScopeDenied
		}
	}
	return nil
}

func (s *ScopeService) FilterProjects(
	ctx context.Context,
	userID uint64,
	workspaceID uint64,
	permission string,
	projectIDs []uint64,
) ([]uint64, error) {
	if len(projectIDs) == 0 {
		return projectIDs, nil
	}
	allowed := make([]uint64, 0, len(projectIDs))
	for _, projectID := range projectIDs {
		if err := s.ValidateScope(ctx, userID, workspaceID, projectID, permission); err == nil {
			allowed = append(allowed, projectID)
		}
	}
	return allowed, nil
}

func (s *ScopeService) ValidateSourceScope(
	ctx context.Context,
	userID uint64,
	source *domain.DeliverySource,
	permission string,
) error {
	if source == nil {
		return ErrGitOpsScopeDenied
	}
	return s.ValidateScope(ctx, userID, derefScopeUint64(source.WorkspaceID), derefScopeUint64(source.ProjectID), permission)
}

func (s *ScopeService) ValidateTargetGroupScope(
	ctx context.Context,
	userID uint64,
	group *domain.ClusterTargetGroup,
	permission string,
) error {
	if group == nil {
		return ErrGitOpsScopeDenied
	}
	return s.ValidateScope(ctx, userID, group.WorkspaceID, derefScopeUint64(group.ProjectID), permission)
}

func (s *ScopeService) ValidateDeliveryUnitScope(
	ctx context.Context,
	userID uint64,
	unit *domain.ApplicationDeliveryUnit,
	permission string,
) error {
	if unit == nil {
		return ErrGitOpsScopeDenied
	}
	return s.ValidateScope(ctx, userID, unit.WorkspaceID, derefScopeUint64(unit.ProjectID), permission)
}

func (s *ScopeService) ValidateEnvironmentScope(
	ctx context.Context,
	userID uint64,
	unit *domain.ApplicationDeliveryUnit,
	_ *domain.EnvironmentStage,
	permission string,
) error {
	return s.ValidateDeliveryUnitScope(ctx, userID, unit, permission)
}

func derefScopeUint64(value *uint64) uint64 {
	if value == nil {
		return 0
	}
	return *value
}
