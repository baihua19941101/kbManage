package securitypolicy

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
		return ErrSecurityPolicyScopeDenied
	}
	permission = strings.TrimSpace(permission)
	if permission == "" {
		permission = PermissionSecurityPolicyRead
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
			return ErrSecurityPolicyScopeDenied
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
			return ErrSecurityPolicyScopeDenied
		}
	}
	return nil
}

func (s *ScopeService) ValidatePolicyScope(ctx context.Context, userID uint64, policy *domain.SecurityPolicy, permission string) error {
	if policy == nil {
		return ErrSecurityPolicyScopeDenied
	}
	return s.ValidateScope(ctx, userID, derefScopeUint64(policy.WorkspaceID), derefScopeUint64(policy.ProjectID), permission)
}

func derefScopeUint64(value *uint64) uint64 {
	if value == nil {
		return 0
	}
	return *value
}
