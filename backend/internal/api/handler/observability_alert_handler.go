package handler

import (
	"net/http"
	"strconv"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
	obsSvc "kbmanage/backend/internal/service/observability"

	"github.com/gin-gonic/gin"
)

type ObservabilityAlertHandler struct {
	alertCenter  *obsSvc.AlertCenterService
	handlingRepo *obsSvc.HandlingRecordService
	auditWriter  *auditSvc.EventWriter
}

func NewObservabilityAlertHandler(
	alertCenter *obsSvc.AlertCenterService,
	handlingRepo *obsSvc.HandlingRecordService,
	auditWriter *auditSvc.EventWriter,
) *ObservabilityAlertHandler {
	return &ObservabilityAlertHandler{
		alertCenter:  alertCenter,
		handlingRepo: handlingRepo,
		auditWriter:  auditWriter,
	}
}

func (h *ObservabilityAlertHandler) ListAlerts(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "100"))
	items, err := h.alertCenter.List(
		c.Request.Context(),
		domain.AlertIncidentStatus(c.Query("status")),
		limit,
	)
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "list_alerts_failed", err.Error())
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ObservabilityAlertHandler) GetAlert(c *gin.Context) {
	alertID, err := parsePathUint64(c, "alertId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid alertId")
		return
	}
	item, err := h.alertCenter.Get(c.Request.Context(), alertID)
	if err != nil {
		writeObservabilityError(c, http.StatusNotFound, "alert_not_found", err.Error())
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *ObservabilityAlertHandler) AcknowledgeAlert(c *gin.Context) {
	alertID, err := parsePathUint64(c, "alertId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid alertId")
		return
	}
	var req struct {
		Note string `json:"note"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	item, err := h.alertCenter.Acknowledge(
		c.Request.Context(),
		alertID,
		c.GetUint64(middleware.UserIDKey),
		req.Note,
	)
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "acknowledge_alert_failed", err.Error())
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionAlertAcknowledge,
			strconv.FormatUint(alertID, 10),
			domain.AuditOutcomeSuccess,
			map[string]any{"note": req.Note},
		)
	}
	c.JSON(http.StatusOK, item)
}

func (h *ObservabilityAlertHandler) CreateHandlingRecord(c *gin.Context) {
	alertID, err := parsePathUint64(c, "alertId")
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid alertId")
		return
	}
	var req struct {
		ActionType string `json:"actionType"`
		Content    string `json:"content"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}
	item, err := h.handlingRepo.Create(
		c.Request.Context(),
		alertID,
		c.GetUint64(middleware.UserIDKey),
		req.ActionType,
		req.Content,
	)
	if err != nil {
		writeObservabilityError(c, http.StatusInternalServerError, "create_handling_record_failed", err.Error())
		return
	}
	if h.auditWriter != nil {
		actorID := userIDPointer(c.GetUint64(middleware.UserIDKey))
		_ = h.auditWriter.WriteObservabilityEvent(
			c.Request.Context(),
			c.GetString(middleware.RequestIDKey),
			actorID,
			auditSvc.ObservabilityAuditActionAlertHandlingRecordCreate,
			strconv.FormatUint(alertID, 10),
			domain.AuditOutcomeSuccess,
			map[string]any{"actionType": req.ActionType},
		)
	}
	c.JSON(http.StatusCreated, item)
}

func parsePathUint64(c *gin.Context, key string) (uint64, error) {
	return strconv.ParseUint(c.Param(key), 10, 64)
}

func userIDPointer(id uint64) *uint64 {
	if id == 0 {
		return nil
	}
	return &id
}
