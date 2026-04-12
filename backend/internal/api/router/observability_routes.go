package router

import (
	"context"
	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/api/middleware"
	alertProvider "kbmanage/backend/internal/integration/observability/alerts"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	authSvc "kbmanage/backend/internal/service/auth"
	clusterSvc "kbmanage/backend/internal/service/cluster"
	obsSvc "kbmanage/backend/internal/service/observability"
	"kbmanage/backend/internal/worker"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RegisterObservabilityRoutes(group *gin.RouterGroup, db *gorm.DB, h *handler.ObservabilityHandler) {
	scopeAccess := newScopeAccessService(db)
	clusterScopeService := clusterSvc.NewService(
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		nil,
		repository.NewWorkspaceClusterRepository(db),
		repository.NewProjectRepository(db),
	)
	scopeService := obsSvc.NewScopeService(authSvc.NewScopeAuthorizer(), scopeAccess, clusterScopeService)
	if h == nil {
		h = handler.NewObservabilityHandler(obsSvc.NewService(scopeService))
	}
	var alertHandler *handler.ObservabilityAlertHandler
	var adminHandler *handler.ObservabilityAdminHandler
	var auditWriter *auditSvc.EventWriter
	if db != nil {
		alerts := alertProvider.NewAlertmanagerProvider("")
		incidentRepo := repository.NewAlertIncidentRepository(db)
		auditWriter = auditSvc.NewEventWriter(repository.NewAuditRepository(db))

		handlingRecordSvc := obsSvc.NewHandlingRecordService(incidentRepo)
		alertCenterSvc := obsSvc.NewAlertCenterService(alerts, incidentRepo, handlingRecordSvc)
		alertHandler = handler.NewObservabilityAlertHandler(alertCenterSvc, handlingRecordSvc, auditWriter)

		alertRuleSvc := obsSvc.NewAlertRuleService(repository.NewAlertRuleRepository(db), scopeService)
		notificationTargetSvc := obsSvc.NewNotificationTargetService(repository.NewNotificationTargetRepository(db))
		silenceSvc := obsSvc.NewSilenceService(repository.NewSilenceWindowRepository(db))
		adminHandler = handler.NewObservabilityAdminHandler(
			alertRuleSvc,
			notificationTargetSvc,
			silenceSvc,
			auditWriter,
		)

		if !strings.EqualFold(db.Dialector.Name(), "sqlite") {
			alertSyncSvc := obsSvc.NewAlertSyncService(alerts, incidentRepo)
			alertSyncSvc.SetAuditWriter(auditWriter)
			syncWorker := worker.NewObservabilitySyncWorker(alertSyncSvc, 30*time.Second)
			syncWorker.Start(context.Background())
		}
	}

	obs := group.Group("/observability")
	{
		readGroup := obs.Group("/")
		readGroup.Use(middleware.RequireObservabilityScope(scopeAccess, scopeService, middleware.PermissionObservabilityRead, auditWriter))
		readGroup.GET("/overview", h.Overview)
		readGroup.GET("/logs/query", h.QueryLogs)
		readGroup.GET("/events", h.ListEvents)
		readGroup.GET("/metrics/series", h.QueryMetricSeries)
		readGroup.GET("/resources/context", h.ResourceContext)

		if alertHandler != nil {
			readGroup.GET("/alerts", alertHandler.ListAlerts)
			readGroup.GET("/alerts/:alertId", alertHandler.GetAlert)

			writeGroup := obs.Group("/")
			writeGroup.Use(middleware.RequireObservabilityScope(scopeAccess, scopeService, middleware.PermissionObservabilityWrite))
			writeGroup.POST("/alerts/:alertId/acknowledge", alertHandler.AcknowledgeAlert)
			writeGroup.POST("/alerts/:alertId/handling-records", alertHandler.CreateHandlingRecord)
		}

		if adminHandler != nil {
			readGroup.GET("/alert-rules", adminHandler.ListAlertRules)
			readGroup.GET("/alert-rules/:ruleId", adminHandler.GetAlertRule)
			readGroup.GET("/notification-targets", adminHandler.ListNotificationTargets)
			readGroup.GET("/notification-targets/:targetId", adminHandler.GetNotificationTarget)
			readGroup.GET("/silences", adminHandler.ListSilences)

			writeGroup := obs.Group("/")
			writeGroup.Use(middleware.RequireObservabilityScope(scopeAccess, scopeService, middleware.PermissionObservabilityWrite))
			writeGroup.POST("/alert-rules", adminHandler.CreateAlertRule)
			writeGroup.PUT("/alert-rules/:ruleId", adminHandler.UpdateAlertRule)
			writeGroup.DELETE("/alert-rules/:ruleId", adminHandler.DeleteAlertRule)
			writeGroup.POST("/notification-targets", adminHandler.CreateNotificationTarget)
			writeGroup.PUT("/notification-targets/:targetId", adminHandler.UpdateNotificationTarget)
			writeGroup.DELETE("/notification-targets/:targetId", adminHandler.DeleteNotificationTarget)
			writeGroup.POST("/silences", adminHandler.CreateSilence)
			writeGroup.DELETE("/silences/:silenceId", adminHandler.CancelSilence)
		}
	}
}
