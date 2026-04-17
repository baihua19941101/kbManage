package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	baselineProvider "kbmanage/backend/internal/integration/compliance/baseline"
	scannerProvider "kbmanage/backend/internal/integration/compliance/scanner"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	complianceSvc "kbmanage/backend/internal/service/compliance"
	"kbmanage/backend/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterComplianceRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	scopeAccess := newScopeAccessService(db)
	baselineRepo := repository.NewComplianceBaselineRepository(db)
	profileRepo := repository.NewComplianceScanProfileRepository(db)
	scanRepo := repository.NewComplianceScanExecutionRepository(db)
	findingRepo := repository.NewComplianceFindingRepository(db)
	evidenceRepo := repository.NewComplianceEvidenceRepository(db)
	remediationRepo := repository.NewComplianceRemediationRepository(db)
	exceptionRepo := repository.NewComplianceExceptionRepository(db)
	recheckRepo := repository.NewComplianceRecheckRepository(db)
	trendRepo := repository.NewComplianceTrendRepository(db)
	exportRepo := repository.NewComplianceExportRepository(db)
	progressCache := complianceSvc.NewProgressCache(rdb, 0)
	exportCache := complianceSvc.NewExportCache(rdb, 0)
	scheduleCache := complianceSvc.NewScheduleCache(rdb, 0)
	svc := complianceSvc.NewService(baselineRepo, profileRepo, scanRepo, findingRepo, evidenceRepo, remediationRepo, exceptionRepo, recheckRepo, trendRepo, exportRepo, scopeAccess, baselineProvider.NewStaticProvider(), scannerProvider.NewMockProvider(), progressCache, exportCache, scheduleCache)
	auditWriter := auditSvc.NewEventWriter(repository.NewAuditRepository(db))
	remediation := complianceSvc.NewRemediationService(auditWriter)
	exceptions := complianceSvc.NewExceptionService(auditWriter)
	rechecks := complianceSvc.NewRecheckService(auditWriter)
	overview := complianceSvc.NewOverviewService()
	trends := complianceSvc.NewTrendService()
	exports := complianceSvc.NewArchiveExportService(auditWriter)
	h := handler.NewComplianceHandler(svc, remediation, exceptions, rechecks, overview, trends, exports)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.ComplianceBaseline{},
			&domain.ScanProfile{},
			&domain.ScanExecution{},
			&domain.ComplianceFinding{},
			&domain.EvidenceRecord{},
			&domain.RemediationTask{},
			&domain.ComplianceExceptionRequest{},
			&domain.RecheckTask{},
			&domain.ComplianceTrendSnapshot{},
			&domain.ArchiveExportTask{},
		)
	}

	worker.NewComplianceScanWorker(svc.Scans, 30).Start(context.Background())
	worker.NewComplianceExceptionExpiryWorker(exceptions, 0).Start(context.Background())
	worker.NewComplianceRecheckWorker(rechecks, 0, 20).Start(context.Background())
	worker.NewComplianceTrendSnapshotWorker(trends, 0).Start(context.Background())
	worker.NewComplianceExportWorker(exports, 0).Start(context.Background())

	compliance := group.Group("/compliance")
	{
		readGroup := compliance.Group("/")
		readGroup.Use(middleware.RequireComplianceScopeFromRequest(scopeAccess, middleware.PermissionComplianceRead))
		writeBaselineGroup := compliance.Group("/")
		writeBaselineGroup.Use(middleware.RequireComplianceScopeFromRequest(scopeAccess, middleware.PermissionComplianceManageBaseline))
		executeGroup := compliance.Group("/")
		executeGroup.Use(middleware.RequireComplianceScopeFromRequest(scopeAccess, middleware.PermissionComplianceExecuteScan))
		remediationGroup := compliance.Group("/")
		remediationGroup.Use(middleware.RequireComplianceScopeFromRequest(scopeAccess, middleware.PermissionComplianceManageRemediation))
		exceptionGroup := compliance.Group("/")
		exceptionGroup.Use(middleware.RequireComplianceScopeFromRequest(scopeAccess, middleware.PermissionComplianceReviewException))
		exportGroup := compliance.Group("/")
		exportGroup.Use(middleware.RequireComplianceScopeFromRequest(scopeAccess, middleware.PermissionComplianceExportArchive))

		readGroup.GET("/baselines", h.ListBaselines)
		writeBaselineGroup.POST("/baselines", h.CreateBaseline)
		readGroup.GET("/baselines/:baselineId", h.GetBaseline)
		writeBaselineGroup.PATCH("/baselines/:baselineId", h.UpdateBaseline)

		readGroup.GET("/scan-profiles", h.ListScanProfiles)
		executeGroup.POST("/scan-profiles", h.CreateScanProfile)
		readGroup.GET("/scan-profiles/:profileId", h.GetScanProfile)
		executeGroup.PATCH("/scan-profiles/:profileId", h.UpdateScanProfile)
		executeGroup.POST("/scan-profiles/:profileId/execute", h.ExecuteScanProfile)

		readGroup.GET("/scans", h.ListScans)
		readGroup.GET("/scans/:scanId", h.GetScan)
		readGroup.GET("/scans/:scanId/findings", h.ListFindings)
		readGroup.GET("/findings", h.ListFindings)
		readGroup.GET("/findings/:findingId", h.GetFinding)
		readGroup.GET("/findings/:findingId/evidence", h.ListEvidence)

		readGroup.GET("/remediation-tasks", h.ListRemediationTasks)
		remediationGroup.POST("/findings/:findingId/remediation-tasks", h.CreateRemediationTask)
		remediationGroup.PATCH("/remediation-tasks/:taskId", h.UpdateRemediationTask)

		readGroup.GET("/exceptions", h.ListExceptions)
		exceptionGroup.POST("/findings/:findingId/exceptions", h.CreateException)
		exceptionGroup.POST("/exceptions/:exceptionId/review", h.ReviewException)

		readGroup.GET("/rechecks", h.ListRechecks)
		readGroup.GET("/rechecks/:recheckId", h.GetRecheck)
		remediationGroup.POST("/findings/:findingId/rechecks", h.CreateRecheck)
		remediationGroup.POST("/rechecks/:recheckId/complete", h.CompleteRecheck)

		readGroup.GET("/overview", h.GetOverview)
		readGroup.GET("/trends", h.ListTrends)
		exportGroup.GET("/archive-exports", h.ListArchiveExports)
		exportGroup.GET("/archive-exports/:exportId", h.GetArchiveExport)
		exportGroup.POST("/archive-exports", h.CreateArchiveExport)
	}
}
