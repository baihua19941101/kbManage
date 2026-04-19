package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	marketplaceint "kbmanage/backend/internal/integration/marketplace"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	marketplaceSvc "kbmanage/backend/internal/service/marketplace"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterMarketplaceRoutes(group *gin.RouterGroup, db *gorm.DB, _ *redis.Client) {
	sourceRepo := repository.NewCatalogSourceRepository(db)
	templateRepo := repository.NewApplicationTemplateRepository(db)
	versionRepo := repository.NewTemplateVersionRepository(db)
	releaseRepo := repository.NewTemplateReleaseScopeRepository(db)
	installationRepo := repository.NewInstallationRecordRepository(db)
	extensionRepo := repository.NewExtensionPackageRepository(db)
	compatibilityRepo := repository.NewCompatibilityStatementRepository(db)
	lifecycleRepo := repository.NewExtensionLifecycleRepository(db)
	localAuditRepo := repository.NewMarketplaceAuditRepository(db)
	bindingRepo := repository.NewScopeRoleBindingRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	workspaceClusterRepo := repository.NewWorkspaceClusterRepository(db)

	svc := marketplaceSvc.NewService(
		sourceRepo,
		templateRepo,
		versionRepo,
		releaseRepo,
		installationRepo,
		extensionRepo,
		compatibilityRepo,
		lifecycleRepo,
		localAuditRepo,
		bindingRepo,
		projectRepo,
		workspaceClusterRepo,
		marketplaceint.NewStaticCatalogProvider(),
		marketplaceint.NewStaticExtensionRegistry(),
		auditSvc.NewEventWriter(repository.NewAuditRepository(db)),
	)
	h := handler.NewMarketplaceHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.CatalogSource{},
			&domain.ApplicationTemplate{},
			&domain.TemplateVersion{},
			&domain.TemplateReleaseScope{},
			&domain.InstallationRecord{},
			&domain.ExtensionPackage{},
			&domain.CompatibilityStatement{},
			&domain.ExtensionLifecycleRecord{},
			&domain.MarketplaceAuditEvent{},
		)
	}

	marketplace := group.Group("/marketplace")
	{
		marketplace.GET("/catalog-sources", h.ListCatalogSources)
		marketplace.POST("/catalog-sources", h.CreateCatalogSource)
		marketplace.POST("/catalog-sources/:sourceId/sync", h.SyncCatalogSource)

		marketplace.GET("/templates", h.ListTemplates)
		marketplace.GET("/templates/:templateId", h.GetTemplateDetail)
		marketplace.GET("/templates/:templateId/releases", h.ListTemplateReleases)
		marketplace.POST("/templates/:templateId/releases", h.CreateTemplateRelease)

		marketplace.GET("/installations", h.ListInstallations)

		marketplace.GET("/extensions", h.ListExtensions)
		marketplace.POST("/extensions", h.RegisterExtension)
		marketplace.POST("/extensions/:extensionId/enable", h.EnableExtension)
		marketplace.POST("/extensions/:extensionId/disable", h.DisableExtension)
		marketplace.GET("/extensions/:extensionId/compatibility", h.GetExtensionCompatibility)
	}
}
