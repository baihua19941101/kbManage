package router

import (
	"context"
	"strings"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	"kbmanage/backend/internal/worker"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAuditRoutes mounts US4 audit APIs.
func RegisterAuditRoutes(group *gin.RouterGroup, db *gorm.DB) {
	auditRepo := repository.NewAuditRepository(db)
	auditExportRepo := repository.NewAuditExportRepository(db)
	scopeAccess := newScopeAccessService(db)
	svc := auditSvc.NewService(auditRepo, auditExportRepo, scopeAccess)
	h := handler.NewAuditHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(&domain.AuditEvent{})
	}

	exportWorker := worker.NewAuditExportWorker(svc, auditExportRepo)
	exportWorker.Start(context.Background())
	if db == nil || !strings.EqualFold(db.Dialector.Name(), "sqlite") {
		retentionWorker := worker.NewAuditRetentionWorker(auditRepo)
		retentionWorker.Start(context.Background())
	}

	group.GET("/audits/events", h.ListEvents)
	group.GET(
		"/audit/security-policies/events",
		middleware.RequireSecurityPolicyScopeFromRequest(scopeAccess, middleware.PermissionSecurityPolicyRead),
		h.ListSecurityPolicyEvents,
	)
	group.GET(
		"/audit/compliance/events",
		h.ListComplianceEvents,
	)
	group.GET(
		"/audit/cluster-lifecycle/events",
		h.ListClusterLifecycleEvents,
	)
	group.GET(
		"/audit/backup-restore/events",
		h.ListBackupRestoreEvents,
	)
	group.GET(
		"/audit/identity/events",
		h.ListIdentityTenancyEvents,
	)
	group.POST("/audits/exports", h.SubmitExport)
	group.GET("/audits/exports/:taskId", h.GetExportStatus)
	group.GET("/audits/exports/:taskId/download", h.DownloadExport)
}
