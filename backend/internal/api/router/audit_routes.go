package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
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
	svc := auditSvc.NewService(auditRepo, auditExportRepo)
	h := handler.NewAuditHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(&domain.AuditEvent{})
	}

	exportWorker := worker.NewAuditExportWorker(svc, auditExportRepo)
	exportWorker.Start(context.Background())

	group.GET("/audits/events", h.ListEvents)
	group.POST("/audits/exports", h.SubmitExport)
	group.GET("/audits/exports/:taskId", h.GetExportStatus)
}
