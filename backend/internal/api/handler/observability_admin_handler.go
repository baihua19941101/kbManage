package handler

import (
	"errors"
	"net/http"
	"strconv"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
	obsSvc "kbmanage/backend/internal/service/observability"

	"github.com/gin-gonic/gin"
)

type ObservabilityAdminHandler struct {
	alertRules          *obsSvc.AlertRuleService
	notificationTargets *obsSvc.NotificationTargetService
	silences            *obsSvc.SilenceService
	auditWriter         *auditSvc.EventWriter
}

func NewObservabilityAdminHandler(
	alertRules *obsSvc.AlertRuleService,
	notificationTargets *obsSvc.NotificationTargetService,
	silences *obsSvc.SilenceService,
	auditWriter *auditSvc.EventWriter,
) *ObservabilityAdminHandler {
	return &ObservabilityAdminHandler{
		alertRules:          alertRules,
		notificationTargets: notificationTargets,
		silences:            silences,
		auditWriter:         auditWriter,
	}
}

func (h *ObservabilityAdminHandler) ListAlertRules(c *gin.Context) {
	items, err := h.alertRules.List(c.Request.Context(), domain.AlertRuleStatus(c.Query("status")))
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "list_alert_rules_failed", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ObservabilityAdminHandler) CreateAlertRule(c *gin.Context) {
	var req obsSvc.UpsertAlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	item, err := h.alertRules.Create(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeObservabilityAdminError(c, "create_alert_rule_failed", err)
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionAlertRuleCreate,
			strconv.FormatUint(item.ID, 10),
			domain.AuditOutcomeSuccess,
			map[string]any{"name": item.Name},
		)
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ObservabilityAdminHandler) GetAlertRule(c *gin.Context) {
	ruleID, err := parsePathUint64(c, "ruleId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid ruleId")
		return
	}
	item, err := h.alertRules.Get(c.Request.Context(), ruleID)
	if err != nil {
		writeObservabilityAdminError(c, "alert_rule_not_found", err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ObservabilityAdminHandler) UpdateAlertRule(c *gin.Context) {
	ruleID, err := parsePathUint64(c, "ruleId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid ruleId")
		return
	}
	var req obsSvc.UpsertAlertRuleRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	item, err := h.alertRules.Update(c.Request.Context(), ruleID, req)
	if err != nil {
		writeObservabilityAdminError(c, "update_alert_rule_failed", err)
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionAlertRuleUpdate,
			strconv.FormatUint(ruleID, 10),
			domain.AuditOutcomeSuccess,
			map[string]any{"name": item.Name},
		)
	}
	c.JSON(http.StatusOK, item)
}

func (h *ObservabilityAdminHandler) DeleteAlertRule(c *gin.Context) {
	ruleID, err := parsePathUint64(c, "ruleId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid ruleId")
		return
	}
	if err := h.alertRules.Delete(c.Request.Context(), ruleID); err != nil {
		writeObservabilityAdminError(c, "delete_alert_rule_failed", err)
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionAlertRuleDelete,
			strconv.FormatUint(ruleID, 10),
			domain.AuditOutcomeSuccess,
			nil,
		)
	}
	c.Status(http.StatusNoContent)
}

func (h *ObservabilityAdminHandler) ListNotificationTargets(c *gin.Context) {
	items, err := h.notificationTargets.List(c.Request.Context())
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "list_notification_targets_failed", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ObservabilityAdminHandler) CreateNotificationTarget(c *gin.Context) {
	var req obsSvc.UpsertNotificationTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	item, err := h.notificationTargets.Create(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "create_notification_target_failed", err.Error())
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionNotificationTargetCreate,
			strconv.FormatUint(item.ID, 10),
			domain.AuditOutcomeSuccess,
			map[string]any{"name": item.Name, "targetType": item.TargetType},
		)
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ObservabilityAdminHandler) GetNotificationTarget(c *gin.Context) {
	targetID, err := parsePathUint64(c, "targetId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid targetId")
		return
	}
	item, err := h.notificationTargets.Get(c.Request.Context(), targetID)
	if err != nil {
		writeObservabilityError(c, http.StatusNotFound, "notification_target_not_found", err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ObservabilityAdminHandler) UpdateNotificationTarget(c *gin.Context) {
	targetID, err := parsePathUint64(c, "targetId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid targetId")
		return
	}
	var req obsSvc.UpsertNotificationTargetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	item, err := h.notificationTargets.Update(c.Request.Context(), targetID, req)
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "update_notification_target_failed", err.Error())
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionNotificationTargetUpdate,
			strconv.FormatUint(targetID, 10),
			domain.AuditOutcomeSuccess,
			map[string]any{"name": item.Name, "targetType": item.TargetType},
		)
	}
	c.JSON(http.StatusOK, item)
}

func (h *ObservabilityAdminHandler) DeleteNotificationTarget(c *gin.Context) {
	targetID, err := parsePathUint64(c, "targetId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid targetId")
		return
	}
	if err := h.notificationTargets.Delete(c.Request.Context(), targetID); err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "delete_notification_target_failed", err.Error())
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionNotificationTargetDelete,
			strconv.FormatUint(targetID, 10),
			domain.AuditOutcomeSuccess,
			nil,
		)
	}
	c.Status(http.StatusNoContent)
}

func (h *ObservabilityAdminHandler) ListSilences(c *gin.Context) {
	items, err := h.silences.List(c.Request.Context(), domain.SilenceWindowStatus(c.Query("status")))
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "list_silences_failed", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ObservabilityAdminHandler) CreateSilence(c *gin.Context) {
	var req obsSvc.CreateSilenceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	item, err := h.silences.Create(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "create_silence_failed", err.Error())
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionSilenceCreate,
			strconv.FormatUint(item.ID, 10),
			domain.AuditOutcomeSuccess,
			map[string]any{"name": item.Name},
		)
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ObservabilityAdminHandler) CancelSilence(c *gin.Context) {
	silenceID, err := parsePathUint64(c, "silenceId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid silenceId")
		return
	}
	item, err := h.silences.Cancel(c.Request.Context(), silenceID, c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "cancel_silence_failed", err.Error())
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionSilenceCancel,
			strconv.FormatUint(silenceID, 10),
			domain.AuditOutcomeSuccess,
			nil,
		)
	}
	c.JSON(http.StatusOK, item)
}

func writeObservabilityAdminError(c *gin.Context, code string, err error) {
	if errors.Is(err, obsSvc.ErrObservabilityScopeDenied) || errors.Is(err, obsSvc.ErrInvalidObservabilityUser) {
		writeObservabilityError(c, http.StatusForbidden, "forbidden", "observability scope access denied")
		return
	}
	if code == "alert_rule_not_found" {
		writeObservabilityError(c, http.StatusNotFound, code, err.Error())
		return
	}
	writeObservabilityError(c, http.StatusInternalServerError, code, err.Error())
}
