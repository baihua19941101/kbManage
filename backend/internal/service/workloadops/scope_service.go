package workloadops

import (
	"context"

	"kbmanage/backend/internal/domain"
	authSvc "kbmanage/backend/internal/service/auth"
)

type ScopeService struct {
	scopeAccess *authSvc.ScopeAccessService
}

func NewScopeService(scopeAccess *authSvc.ScopeAccessService) *ScopeService {
	return &ScopeService{scopeAccess: scopeAccess}
}

func (s *ScopeService) ValidateResourceAccess(ctx context.Context, userID uint64, clusterID uint64, permission string) error {
	return s.ValidateWorkloadAccess(ctx, userID, clusterID, 0, 0, permission)
}

func (s *ScopeService) ValidateWorkloadAccess(
	ctx context.Context,
	userID uint64,
	clusterID uint64,
	workspaceID uint64,
	projectID uint64,
	permission string,
) error {
	if s == nil || s.scopeAccess == nil {
		return nil
	}
	if userID == 0 || clusterID == 0 || permission == "" {
		return ErrWorkloadOpsScopeDenied
	}
	clusterIDs, constrained, err := s.scopeAccess.ListClusterIDsByPermission(ctx, userID, permission)
	if err != nil {
		return err
	}
	if !constrained {
		return ErrWorkloadOpsScopeDenied
	}
	for _, allowedClusterID := range clusterIDs {
		if allowedClusterID == clusterID {
			if workspaceID != 0 {
				allowed, scopeErr := s.scopeAccess.HasScopePermission(
					ctx,
					userID,
					domain.ScopeTypeWorkspace,
					workspaceID,
					0,
					permission,
				)
				if scopeErr != nil {
					return scopeErr
				}
				if !allowed {
					return ErrWorkloadOpsScopeDenied
				}
			}
			if projectID != 0 {
				allowed, scopeErr := s.scopeAccess.HasScopePermission(
					ctx,
					userID,
					domain.ScopeTypeProject,
					workspaceID,
					projectID,
					permission,
				)
				if scopeErr != nil {
					return scopeErr
				}
				if !allowed {
					return ErrWorkloadOpsScopeDenied
				}
			}
			return nil
		}
	}
	if len(clusterIDs) == 0 {
		return ErrWorkloadOpsScopeDenied
	}
	return ErrWorkloadOpsScopeDenied
}
