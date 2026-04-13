package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	workloadops "kbmanage/backend/internal/service/workloadops"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterWorkloadOpsRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	actionRepo := repository.NewWorkloadActionRepository(db)
	batchRepo := repository.NewBatchOperationRepository(db)
	sessionRepo := repository.NewTerminalSessionRepository(db)
	scopeAccess := newScopeAccessService(db)
	scopeSvc := workloadops.NewScopeService(scopeAccess)
	progressCache := workloadops.NewProgressCache(rdb, 0)
	sessionCache := workloadops.NewSessionCache(rdb, 0)
	svc := workloadops.NewService(actionRepo, batchRepo, sessionRepo, scopeSvc, progressCache, sessionCache)
	svc.SetAuditWriter(auditSvc.NewEventWriter(repository.NewAuditRepository(db)))
	h := handler.NewWorkloadOpsHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.WorkloadActionRequest{},
			&domain.BatchOperationTask{},
			&domain.BatchOperationItem{},
			&domain.TerminalSession{},
		)
	}

	ops := group.Group("/workload-ops")
	{
		ops.GET("/resources/context", middleware.RequireWorkloadOpsClusterScope(scopeAccess, middleware.PermissionWorkloadOpsRead), h.GetContext)
		ops.GET("/resources/instances", middleware.RequireWorkloadOpsClusterScope(scopeAccess, middleware.PermissionWorkloadOpsRead), h.ListInstances)
		ops.GET("/resources/revisions", middleware.RequireWorkloadOpsClusterScope(scopeAccess, middleware.PermissionWorkloadOpsRead), h.ListRevisions)

		ops.POST("/actions", middleware.RequireWorkloadOpsActionScope(scopeAccess), h.SubmitAction)
		ops.GET("/actions/:actionId", h.GetAction)

		ops.POST("/batches", middleware.RequireWorkloadOpsClusterScope(scopeAccess, middleware.PermissionWorkloadOpsBatch), h.SubmitBatch)
		ops.GET("/batches/:batchId", h.GetBatch)

		ops.POST("/terminal/sessions", middleware.RequireWorkloadOpsClusterScope(scopeAccess, middleware.PermissionWorkloadOpsTerminal), h.CreateTerminalSession)
		ops.GET("/terminal/sessions/:sessionId", h.GetTerminalSession)
		ops.DELETE("/terminal/sessions/:sessionId", h.CloseTerminalSession)
	}
}
