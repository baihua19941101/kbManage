package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/api/middleware"
	auditSvc "kbmanage/backend/internal/service/audit"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AuditHandler struct {
	svc *auditSvc.Service
}

func NewAuditHandler(svc *auditSvc.Service) *AuditHandler {
	return &AuditHandler{svc: svc}
}

func (h *AuditHandler) ListEvents(c *gin.Context) {
	startAt, err := parseOptionalRFC3339(c.Query("startAt"), "startAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endAt, err := parseOptionalRFC3339(c.Query("endAt"), "endAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var actorID *uint64
	actorRaw := strings.TrimSpace(c.Query("actorId"))
	if actorRaw != "" {
		parsed, err := strconv.ParseUint(actorRaw, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "actorId must be a positive integer"})
			return
		}
		actorID = &parsed
	}

	limit := 100
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		n, err := strconv.Atoi(raw)
		if err != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
			return
		}
		limit = n
	}

	items, err := h.svc.QueryEvents(c.Request.Context(), auditSvc.QueryEventsRequest{
		StartAt: startAt,
		EndAt:   endAt,
		ActorID: actorID,
		Action:  strings.TrimSpace(c.Query("action")),
		Outcome: strings.TrimSpace(c.Query("outcome")),
		Limit:   limit,
	})
	if err != nil {
		writeAuditError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"items": items,
		"count": len(items),
	})
}

type submitAuditExportRequest struct {
	StartAt string `json:"startAt"`
	EndAt   string `json:"endAt"`
	ActorID any    `json:"actorId"`
	Action  string `json:"action"`
	Outcome string `json:"outcome"`
}

func (h *AuditHandler) SubmitExport(c *gin.Context) {
	var req submitAuditExportRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	startAt, err := parseOptionalRFC3339(req.StartAt, "startAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endAt, err := parseOptionalRFC3339(req.EndAt, "endAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	actor, err := parseUint64(req.ActorID, false, "actorId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var actorID *uint64
	if actor != 0 {
		actorID = &actor
	}

	task, err := h.svc.SubmitExport(c.Request.Context(), c.GetUint64(middleware.UserIDKey), auditSvc.SubmitExportRequest{
		StartAt: startAt,
		EndAt:   endAt,
		ActorID: actorID,
		Action:  strings.TrimSpace(req.Action),
		Outcome: strings.TrimSpace(req.Outcome),
	})
	if err != nil {
		writeAuditError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{
		"taskId":      task.ID,
		"status":      task.Status,
		"createdAt":   task.CreatedAt,
		"operatorId":  task.OperatorID,
		"resultTotal": task.ResultTotal,
		"downloadUrl": task.DownloadURL,
	})
}

func (h *AuditHandler) GetExportStatus(c *gin.Context) {
	task, err := h.svc.GetExportTask(c.Request.Context(), c.Param("taskId"))
	if err != nil {
		writeAuditError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"taskId":       task.ID,
		"status":       task.Status,
		"operatorId":   task.OperatorID,
		"resultTotal":  task.ResultTotal,
		"downloadUrl":  task.DownloadURL,
		"errorMessage": task.ErrorMessage,
		"createdAt":    task.CreatedAt,
		"updatedAt":    task.UpdatedAt,
		"completedAt":  task.CompletedAt,
	})
}

func writeAuditError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	lower := strings.ToLower(err.Error())
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case strings.Contains(lower, "required"), strings.Contains(lower, "invalid"):
		status = http.StatusBadRequest
	}
	c.JSON(status, gin.H{"error": err.Error()})
}

func parseOptionalRFC3339(raw, field string) (*time.Time, error) {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil, nil
	}
	t, err := time.Parse(time.RFC3339, trimmed)
	if err != nil {
		return nil, errors.New(field + " must be RFC3339 format")
	}
	return &t, nil
}
