package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	projectSvc "kbmanage/backend/internal/service/project"
	workspaceSvc "kbmanage/backend/internal/service/workspace"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAccessRoutes mounts US2 workspace/project/role-binding APIs.
func RegisterAccessRoutes(group *gin.RouterGroup, db *gorm.DB) {
	workspaceRepo := repository.NewWorkspaceRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	workspaceClusterRepo := repository.NewWorkspaceClusterRepository(db)
	roleRepo := repository.NewScopeRoleRepository(db)
	bindingRepo := repository.NewScopeRoleBindingRepository(db)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.Workspace{},
			&domain.Project{},
			&repository.WorkspaceClusterBinding{},
			&repository.ScopeRole{},
			&repository.ScopeRoleBinding{},
		)
		_ = roleRepo.EnsureDefaults(context.Background())
	}

	_ = workspaceClusterRepo

	workspaceHandler := handler.NewWorkspaceHandler(workspaceSvc.NewService(workspaceRepo))
	projectHandler := handler.NewProjectHandler(projectSvc.NewService(projectRepo, workspaceRepo))
	roleBindingHandler := handler.NewRoleBindingHandler(roleRepo, bindingRepo)

	group.GET("/workspaces", workspaceHandler.List)
	group.POST("/workspaces", workspaceHandler.Create)

	group.GET("/workspaces/:workspaceId/projects", projectHandler.ListByWorkspace)
	group.POST("/workspaces/:workspaceId/projects", projectHandler.Create)

	group.GET("/role-bindings", roleBindingHandler.List)
	group.POST("/role-bindings", roleBindingHandler.Create)
}
