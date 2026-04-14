package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	gitopsSvc "kbmanage/backend/internal/service/gitops"
	"kbmanage/backend/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterGitOpsRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	sourceRepo := repository.NewDeliverySourceRepository(db)
	targetGroupRepo := repository.NewClusterTargetGroupRepository(db)
	deliveryUnitRepo := repository.NewDeliveryUnitRepository(db)
	revisionRepo := repository.NewReleaseRevisionRepository(db)
	operationRepo := repository.NewDeliveryOperationRepository(db)
	scopeAccess := newScopeAccessService(db)
	progressCache := gitopsSvc.NewProgressCache(rdb, 0)
	diffCache := gitopsSvc.NewDiffCache(rdb, 0)
	lockService := gitopsSvc.NewLockService(rdb)
	operationQueue := gitopsSvc.NewOperationQueue(rdb)

	svc := gitopsSvc.NewService(
		sourceRepo,
		targetGroupRepo,
		deliveryUnitRepo,
		revisionRepo,
		operationRepo,
		gitopsSvc.NewScopeService(scopeAccess),
		progressCache,
		diffCache,
		lockService,
		operationQueue,
	)
	auditWriter := auditSvc.NewEventWriter(repository.NewAuditRepository(db))
	h := handler.NewGitOpsHandler(svc, auditWriter)
	stageExecutionSvc := gitopsSvc.NewStageExecutionService(deliveryUnitRepo, targetGroupRepo)
	promotionSvc := gitopsSvc.NewPromotionService(deliveryUnitRepo, stageExecutionSvc)
	revisionSvc := gitopsSvc.NewRevisionService(revisionRepo, deliveryUnitRepo)
	executor := gitopsSvc.NewExecutor(revisionSvc, promotionSvc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.DeliverySource{},
			&domain.ClusterTargetGroup{},
			&domain.ApplicationDeliveryUnit{},
			&domain.EnvironmentStage{},
			&domain.ConfigurationOverlay{},
			&domain.ReleaseRevision{},
			&domain.DeliveryOperation{},
		)
	}
	operationWorker := worker.NewDeliveryOperationWorker(
		operationRepo,
		deliveryUnitRepo,
		operationQueue,
		executor,
		progressCache,
	)
	operationWorker.Start(context.Background())

	gitops := group.Group("/gitops")
	{
		gitops.GET("/sources", middleware.RequireGitOpsScopeFromRequest(scopeAccess, middleware.PermissionGitOpsRead), h.ListSources)
		gitops.POST("/sources", middleware.RequireGitOpsScopeFromRequest(scopeAccess, middleware.PermissionGitOpsManageSource), h.CreateSource)
		gitops.GET("/sources/:sourceId", middleware.RequireGitOpsSourceScope(scopeAccess, sourceRepo, middleware.PermissionGitOpsRead), h.GetSource)
		gitops.PATCH("/sources/:sourceId", middleware.RequireGitOpsSourceScope(scopeAccess, sourceRepo, middleware.PermissionGitOpsManageSource), h.UpdateSource)
		gitops.PUT("/sources/:sourceId", middleware.RequireGitOpsSourceScope(scopeAccess, sourceRepo, middleware.PermissionGitOpsManageSource), h.UpdateSource)
		gitops.POST("/sources/:sourceId/verify", middleware.RequireGitOpsSourceScope(scopeAccess, sourceRepo, middleware.PermissionGitOpsManageSource), h.VerifySource)

		gitops.GET("/target-groups", middleware.RequireGitOpsScopeFromRequest(scopeAccess, middleware.PermissionGitOpsRead), h.ListTargetGroups)
		gitops.POST("/target-groups", middleware.RequireGitOpsScopeFromRequest(scopeAccess, middleware.PermissionGitOpsOverride), h.CreateTargetGroup)
		gitops.GET("/target-groups/:targetGroupId", middleware.RequireGitOpsTargetGroupScope(scopeAccess, targetGroupRepo, middleware.PermissionGitOpsRead), h.GetTargetGroup)
		gitops.PATCH("/target-groups/:targetGroupId", middleware.RequireGitOpsTargetGroupScope(scopeAccess, targetGroupRepo, middleware.PermissionGitOpsOverride), h.UpdateTargetGroup)
		gitops.PUT("/target-groups/:targetGroupId", middleware.RequireGitOpsTargetGroupScope(scopeAccess, targetGroupRepo, middleware.PermissionGitOpsOverride), h.UpdateTargetGroup)

		gitops.GET("/delivery-units", middleware.RequireGitOpsScopeFromRequest(scopeAccess, middleware.PermissionGitOpsRead), h.ListDeliveryUnits)
		gitops.POST("/delivery-units", middleware.RequireGitOpsScopeFromRequest(scopeAccess, middleware.PermissionGitOpsOverride), h.CreateDeliveryUnit)
		gitops.GET("/delivery-units/:unitId", middleware.RequireGitOpsDeliveryUnitScope(scopeAccess, deliveryUnitRepo, middleware.PermissionGitOpsRead), h.GetDeliveryUnit)
		gitops.PATCH("/delivery-units/:unitId", middleware.RequireGitOpsDeliveryUnitScope(scopeAccess, deliveryUnitRepo, middleware.PermissionGitOpsOverride), h.UpdateDeliveryUnit)
		gitops.PUT("/delivery-units/:unitId", middleware.RequireGitOpsDeliveryUnitScope(scopeAccess, deliveryUnitRepo, middleware.PermissionGitOpsOverride), h.UpdateDeliveryUnit)
		gitops.GET("/delivery-units/:unitId/status", middleware.RequireGitOpsDeliveryUnitScope(scopeAccess, deliveryUnitRepo, middleware.PermissionGitOpsRead), h.GetDeliveryUnitStatus)
		gitops.GET("/delivery-units/:unitId/diff", middleware.RequireGitOpsDeliveryUnitScope(scopeAccess, deliveryUnitRepo, middleware.PermissionGitOpsRead), h.GetDeliveryUnitDiff)
		gitops.POST("/delivery-units/:unitId/actions", middleware.RequireGitOpsActionScope(scopeAccess, deliveryUnitRepo), h.SubmitAction)
		gitops.GET("/delivery-units/:unitId/releases", middleware.RequireGitOpsDeliveryUnitScope(scopeAccess, deliveryUnitRepo, middleware.PermissionGitOpsRead), h.ListReleaseRevisions)

		gitops.GET("/operations/:operationId", middleware.RequireGitOpsOperationScope(scopeAccess, operationRepo, deliveryUnitRepo, middleware.PermissionGitOpsRead), h.GetOperation)
	}
}
