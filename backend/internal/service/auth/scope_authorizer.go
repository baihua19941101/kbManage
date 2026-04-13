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
		if targetType == domain.ScopeTypeProject {
			return grantedProjectID != 0 && grantedProjectID == targetProjectID
		}
		if targetType == domain.ScopeTypeWorkspace {
			return grantedWorkspaceID != 0 && grantedWorkspaceID == targetWorkspaceID
		}
		return false
	}

	return false
}

func (a *ScopeAuthorizer) CanAccessObservability(
	grantedType domain.ScopeType,
	grantedWorkspaceID uint64,
	grantedProjectID uint64,
	targetWorkspaceID uint64,
	targetProjectID uint64,
	clusterScoped bool,
) bool {
	if clusterScoped && targetWorkspaceID == 0 && targetProjectID == 0 {
		return grantedType == domain.ScopeTypePlatform
	}

	targetType := domain.ScopeTypeWorkspace
	if targetProjectID != 0 {
		targetType = domain.ScopeTypeProject
	}

	return a.CanAccess(
		grantedType,
		grantedWorkspaceID,
		grantedProjectID,
		targetType,
		targetWorkspaceID,
		targetProjectID,
	)
}

func (a *ScopeAuthorizer) CanAccessObservabilityMapped(
	grantedType domain.ScopeType,
	grantedWorkspaceID uint64,
	grantedProjectID uint64,
	targetWorkspaceIDs []uint64,
	targetProjectIDs []uint64,
	clusterScoped bool,
) bool {
	if grantedType == domain.ScopeTypePlatform {
		return true
	}

	for _, workspaceID := range targetWorkspaceIDs {
		if workspaceID == 0 {
			continue
		}
		if a.CanAccess(
			grantedType,
			grantedWorkspaceID,
			grantedProjectID,
			domain.ScopeTypeWorkspace,
			workspaceID,
			0,
		) {
			return true
		}
	}
	for _, projectID := range targetProjectIDs {
		if projectID == 0 {
			continue
		}
		if a.CanAccess(
			grantedType,
			grantedWorkspaceID,
			grantedProjectID,
			domain.ScopeTypeProject,
			0,
			projectID,
		) {
			return true
		}
	}

	if clusterScoped {
		return false
	}

	return a.CanAccessObservability(
		grantedType,
		grantedWorkspaceID,
		grantedProjectID,
		firstNonZero(targetWorkspaceIDs),
		firstNonZero(targetProjectIDs),
		false,
	)
}

func (a *ScopeAuthorizer) CanAccessWorkloadOps(
	grantedType domain.ScopeType,
	grantedWorkspaceID uint64,
	grantedProjectID uint64,
	targetWorkspaceID uint64,
	targetProjectID uint64,
	clusterScoped bool,
) bool {
	if clusterScoped && targetWorkspaceID == 0 && targetProjectID == 0 {
		return grantedType == domain.ScopeTypePlatform
	}

	targetType := domain.ScopeTypeWorkspace
	if targetProjectID != 0 {
		targetType = domain.ScopeTypeProject
	}

	return a.CanAccess(
		grantedType,
		grantedWorkspaceID,
		grantedProjectID,
		targetType,
		targetWorkspaceID,
		targetProjectID,
	)
}

func (a *ScopeAuthorizer) CanAccessWorkloadOpsMapped(
	grantedType domain.ScopeType,
	grantedWorkspaceID uint64,
	grantedProjectID uint64,
	targetWorkspaceIDs []uint64,
	targetProjectIDs []uint64,
	clusterScoped bool,
) bool {
	if grantedType == domain.ScopeTypePlatform {
		return true
	}

	for _, workspaceID := range targetWorkspaceIDs {
		if workspaceID == 0 {
			continue
		}
		if a.CanAccess(
			grantedType,
			grantedWorkspaceID,
			grantedProjectID,
			domain.ScopeTypeWorkspace,
			workspaceID,
			0,
		) {
			return true
		}
	}
	for _, projectID := range targetProjectIDs {
		if projectID == 0 {
			continue
		}
		if a.CanAccess(
			grantedType,
			grantedWorkspaceID,
			grantedProjectID,
			domain.ScopeTypeProject,
			0,
			projectID,
		) {
			return true
		}
	}

	if clusterScoped {
		return false
	}

	return a.CanAccessWorkloadOps(
		grantedType,
		grantedWorkspaceID,
		grantedProjectID,
		firstNonZero(targetWorkspaceIDs),
		firstNonZero(targetProjectIDs),
		false,
	)
}

func firstNonZero(items []uint64) uint64 {
	for _, item := range items {
		if item != 0 {
			return item
		}
	}
	return 0
}
