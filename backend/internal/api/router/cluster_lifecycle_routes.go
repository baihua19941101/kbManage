package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	driverProvider "kbmanage/backend/internal/integration/clusterlifecycle/driver"
	validatorProvider "kbmanage/backend/internal/integration/clusterlifecycle/validator"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	clusterLifecycleSvc "kbmanage/backend/internal/service/clusterlifecycle"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterClusterLifecycleRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	clusterRepo := repository.NewClusterLifecycleRepository(db)
	operationRepo := repository.NewClusterLifecycleOperationRepository(db)
	driverRepo := repository.NewClusterDriverRepository(db)
	templateRepo := repository.NewClusterTemplateRepository(db)
	capabilityRepo := repository.NewClusterCapabilityRepository(db)
	upgradeRepo := repository.NewUpgradePlanRepository(db)
	nodePoolRepo := repository.NewNodePoolRepository(db)
	bindingRepo := repository.NewScopeRoleBindingRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	auditWriter := auditSvc.NewEventWriter(repository.NewAuditRepository(db))
	svc := clusterLifecycleSvc.NewService(
		clusterRepo,
		operationRepo,
		driverRepo,
		templateRepo,
		capabilityRepo,
		upgradeRepo,
		nodePoolRepo,
		bindingRepo,
		projectRepo,
		clusterLifecycleSvc.NewProgressCache(rdb),
		clusterLifecycleSvc.NewValidationCache(rdb),
		clusterLifecycleSvc.NewOperationLock(rdb),
		driverProvider.NewStaticProvider(),
		validatorProvider.NewStaticProvider(),
		auditWriter,
	)
	h := handler.NewClusterLifecycleHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.ClusterLifecycleRecord{},
			&domain.ClusterDriverVersion{},
			&domain.CapabilityMatrixEntry{},
			&domain.ClusterTemplate{},
			&domain.LifecycleOperation{},
			&domain.UpgradePlan{},
			&domain.NodePoolProfile{},
			&domain.LifecycleAuditEvent{},
		)
	}

	lifecycle := group.Group("/cluster-lifecycle")
	{
		lifecycle.GET("/clusters", h.ListClusters)
		lifecycle.POST("/clusters/import", h.ImportCluster)
		lifecycle.POST("/clusters/register", h.RegisterCluster)
		lifecycle.GET("/clusters/:clusterId", h.GetCluster)
		lifecycle.POST("/clusters", h.CreateCluster)
		lifecycle.POST("/clusters/:clusterId/validate", h.ValidateClusterChange)
		lifecycle.POST("/clusters/:clusterId/upgrade-plans", h.CreateUpgradePlan)
		lifecycle.POST("/clusters/:clusterId/upgrade-plans/:planId/execute", h.ExecuteUpgradePlan)
		lifecycle.GET("/clusters/:clusterId/node-pools", h.ListNodePools)
		lifecycle.POST("/clusters/:clusterId/node-pools/:nodePoolId/scale", h.ScaleNodePool)
		lifecycle.POST("/clusters/:clusterId/disable", h.DisableCluster)
		lifecycle.POST("/clusters/:clusterId/retire", h.RetireCluster)
		lifecycle.GET("/drivers", h.ListDrivers)
		lifecycle.POST("/drivers", h.CreateDriver)
		lifecycle.GET("/drivers/:driverId/capabilities", h.ListDriverCapabilities)
		lifecycle.GET("/templates", h.ListTemplates)
		lifecycle.POST("/templates", h.CreateTemplate)
		lifecycle.POST("/templates/:templateId/validate", h.ValidateTemplate)
	}
}
