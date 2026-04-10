package auth

import "kbmanage/backend/internal/domain"

type ScopeAuthorizer struct{}

func NewScopeAuthorizer() *ScopeAuthorizer {
	return &ScopeAuthorizer{}
}

func (a *ScopeAuthorizer) CanAccess(
	grantedType domain.ScopeType,
	grantedWorkspaceID uint64,
	grantedProjectID uint64,
	targetType domain.ScopeType,
	targetWorkspaceID uint64,
	targetProjectID uint64,
) bool {
	if grantedType == domain.ScopeTypePlatform {
		return true
	}

	if grantedType == domain.ScopeTypeWorkspace {
		if targetType == domain.ScopeTypeWorkspace || targetType == domain.ScopeTypeProject {
			return grantedWorkspaceID != 0 && grantedWorkspaceID == targetWorkspaceID
		}
		return false
	}

	if grantedType == domain.ScopeTypeProject {
		return targetType == domain.ScopeTypeProject && grantedProjectID != 0 && grantedProjectID == targetProjectID
	}

	return false
}
