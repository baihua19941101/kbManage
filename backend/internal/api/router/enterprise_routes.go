package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	entint "kbmanage/backend/internal/integration/enterprise"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	entSvc "kbmanage/backend/internal/service/enterprise"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterEnterpriseRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	svc := entSvc.NewService(
		repository.NewPermissionChangeTrailRepository(db),
		repository.NewKeyOperationTraceRepository(db),
		repository.NewCrossTeamAuthorizationSnapshotRepository(db),
		repository.NewGovernanceRiskEventRepository(db),
		repository.NewGovernanceCoverageSnapshotRepository(db),
		repository.NewGovernanceReportPackageRepository(db),
		repository.NewExportRecordRepository(db),
		repository.NewDeliveryArtifactRepository(db),
		repository.NewDeliveryReadinessBundleRepository(db),
		repository.NewDeliveryChecklistItemRepository(db),
		repository.NewGovernanceActionItemRepository(db),
		repository.NewScopeRoleBindingRepository(db),
		repository.NewProjectRepository(db),
		repository.NewWorkspaceClusterRepository(db),
		entint.NewStaticAuditProvider(),
		entint.NewStaticReportBuilder(),
		entint.NewStaticDeliveryCatalog(),
		entSvc.NewReportCache(rdb),
		entSvc.NewExportCoordinator(rdb),
		entSvc.NewTrendCache(rdb),
		auditSvc.NewEventWriter(repository.NewAuditRepository(db)),
	)
	h := handler.NewEnterpriseHandler(svc)
	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.PermissionChangeTrail{},
			&domain.KeyOperationTrace{},
			&domain.CrossTeamAuthorizationSnapshot{},
			&domain.GovernanceRiskEvent{},
			&domain.GovernanceCoverageSnapshot{},
			&domain.GovernanceReportPackage{},
			&domain.ExportRecord{},
			&domain.DeliveryArtifact{},
			&domain.DeliveryReadinessBundle{},
			&domain.DeliveryChecklistItem{},
			&domain.GovernanceActionItem{},
		)
	}

	ent := group.Group("/enterprise")
	{
		ent.GET("/audit/permission-trails", h.ListPermissionTrails)
		ent.GET("/audit/key-operations", h.ListKeyOperations)
		ent.GET("/governance/coverage", h.ListCoverage)
		ent.GET("/governance/action-items", h.ListActionItems)
		ent.GET("/reports", h.ListReports)
		ent.POST("/reports", h.CreateReport)
		ent.POST("/reports/:reportId/exports", h.CreateExportRecord)
		ent.GET("/delivery/artifacts", h.ListDeliveryArtifacts)
		ent.GET("/delivery/bundles", h.ListDeliveryBundles)
		ent.GET("/delivery/bundles/:bundleId/checklists", h.ListDeliveryChecklist)
	}
}
