package router

import (
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"

	"gorm.io/gorm"
)

func newScopeAccessService(db *gorm.DB) *authSvc.ScopeAccessService {
	if db == nil {
		return nil
	}

	bindingRepo := repository.NewScopeRoleBindingRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	workspaceClusterRepo := repository.NewWorkspaceClusterRepository(db)
	return authSvc.NewScopeAccessService(
		bindingRepo,
		projectRepo,
		authSvc.NewScopeAuthorizer(),
		authSvc.NewPermissionService(),
		workspaceClusterRepo,
	)
}
