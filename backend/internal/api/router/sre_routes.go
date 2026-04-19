package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	sreint "kbmanage/backend/internal/integration/sre"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	sreSvc "kbmanage/backend/internal/service/sre"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterSRERoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	haRepo := repository.NewHAPolicyRepository(db)
	windowRepo := repository.NewMaintenanceWindowRepository(db)
	healthRepo := repository.NewPlatformHealthSnapshotRepository(db)
	capacityRepo := repository.NewCapacityBaselineRepository(db)
	upgradeRepo := repository.NewSREUpgradePlanRepository(db)
	rollbackRepo := repository.NewRollbackValidationRepository(db)
	runbookRepo := repository.NewRunbookArticleRepository(db)
	alertRepo := repository.NewAlertBaselineRepository(db)
	scaleRepo := repository.NewScaleEvidenceRepository(db)
	bindingRepo := repository.NewScopeRoleBindingRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	workspaceClusterRepo := repository.NewWorkspaceClusterRepository(db)

	svc := sreSvc.NewService(
		haRepo,
		windowRepo,
		healthRepo,
		capacityRepo,
		upgradeRepo,
		rollbackRepo,
		runbookRepo,
		alertRepo,
		scaleRepo,
		bindingRepo,
		projectRepo,
		workspaceClusterRepo,
		sreint.NewStaticHealthProvider(),
		sreint.NewStaticUpgradeValidator(),
		sreint.NewStaticScaleAnalyzer(),
		sreSvc.NewHealthCache(rdb),
		sreSvc.NewUpgradeCoordinator(rdb),
		sreSvc.NewScaleCache(rdb),
		auditSvc.NewEventWriter(repository.NewAuditRepository(db)),
	)
	h := handler.NewSREHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.HAPolicy{},
			&domain.MaintenanceWindow{},
			&domain.PlatformHealthSnapshot{},
			&domain.CapacityBaseline{},
			&domain.SREUpgradePlan{},
			&domain.RollbackValidation{},
			&domain.RunbookArticle{},
			&domain.AlertBaseline{},
			&domain.ScaleEvidence{},
		)
	}

	sre := group.Group("/sre")
	{
		sre.GET("/ha-policies", h.ListHAPolicies)
		sre.POST("/ha-policies", h.UpsertHAPolicy)
		sre.GET("/health/overview", h.GetHealthOverview)
		sre.GET("/maintenance-windows", h.ListMaintenanceWindows)
		sre.POST("/maintenance-windows", h.UpsertMaintenanceWindow)
		sre.POST("/upgrades/prechecks", h.RunUpgradePrecheck)
		sre.GET("/upgrades", h.ListUpgradePlans)
		sre.POST("/upgrades", h.CreateUpgradePlan)
		sre.POST("/upgrades/:upgradeId/rollback-validations", h.CreateRollbackValidation)
		sre.GET("/capacity/baselines", h.ListCapacityBaselines)
		sre.GET("/scale-evidence", h.ListScaleEvidence)
		sre.GET("/runbooks", h.ListRunbooks)
	}
}
