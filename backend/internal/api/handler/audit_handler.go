package handler

import (
	"errors"
	"fmt"
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
	startAt, err := parseOptionalRFC3339(firstNonEmptyQuery(c, "startAt", "timeFrom"), "startAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endAt, err := parseOptionalRFC3339(firstNonEmptyQuery(c, "endAt", "timeTo"), "endAt")
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

	clusterID, err := parseOptionalQueryUint64(c, "clusterId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := parseOptionalQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
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
		StartAt:     startAt,
		EndAt:       endAt,
		ActorID:     actorID,
		ClusterID:   clusterID,
		WorkspaceID: workspaceID,
		ProjectID:   projectID,
		Action:      strings.TrimSpace(c.Query("action")),
		Outcome:     strings.TrimSpace(c.Query("outcome")),
		Result:      strings.TrimSpace(c.Query("result")),
		Resource:    strings.TrimSpace(c.Query("resource")),
		Limit:       limit,
		ViewerID:    c.GetUint64(middleware.UserIDKey),
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

func (h *AuditHandler) ListComplianceEvents(c *gin.Context) {
	startAt, err := parseOptionalRFC3339(firstNonEmptyQuery(c, "startAt", "timeFrom"), "startAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endAt, err := parseOptionalRFC3339(firstNonEmptyQuery(c, "endAt", "timeTo"), "endAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clusterID, err := parseOptionalQueryUint64(c, "clusterId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := parseOptionalQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.QueryComplianceEvents(c.Request.Context(), auditSvc.QueryEventsRequest{StartAt: startAt, EndAt: endAt, ClusterID: clusterID, WorkspaceID: workspaceID, ProjectID: projectID, Action: strings.TrimSpace(c.Query("action")), Outcome: strings.TrimSpace(c.Query("outcome")), Result: strings.TrimSpace(c.Query("result")), Resource: strings.TrimSpace(c.Query("resource")), Limit: 100, ViewerID: c.GetUint64(middleware.UserIDKey)})
	if err != nil {
		writeAuditError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "count": len(items)})
}

func (h *AuditHandler) ListSecurityPolicyEvents(c *gin.Context) {
	startAt, err := parseOptionalRFC3339(firstNonEmptyQuery(c, "startAt", "timeFrom"), "startAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	endAt, err := parseOptionalRFC3339(firstNonEmptyQuery(c, "endAt", "timeTo"), "endAt")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var actorID *uint64
	actorRaw := strings.TrimSpace(c.Query("actorId"))
	if actorRaw != "" {
		parsed, parseErr := strconv.ParseUint(actorRaw, 10, 64)
		if parseErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "actorId must be a positive integer"})
			return
		}
		actorID = &parsed
	}

	clusterID, err := parseOptionalQueryUint64(c, "clusterId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := parseOptionalQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	limit := 100
	if raw := strings.TrimSpace(c.Query("limit")); raw != "" {
		n, parseErr := strconv.Atoi(raw)
		if parseErr != nil || n <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "limit must be a positive integer"})
			return
		}
		limit = n
	}

	items, err := h.svc.QuerySecurityPolicyEvents(c.Request.Context(), auditSvc.QueryEventsRequest{
		StartAt:     startAt,
		EndAt:       endAt,
		ActorID:     actorID,
		ClusterID:   clusterID,
		WorkspaceID: workspaceID,
		ProjectID:   projectID,
		Action:      strings.TrimSpace(c.Query("action")),
		Outcome:     strings.TrimSpace(c.Query("outcome")),
		Result:      strings.TrimSpace(c.Query("result")),
		Resource:    strings.TrimSpace(c.Query("resource")),
		Limit:       limit,
		ViewerID:    c.GetUint64(middleware.UserIDKey),
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
	StartAt     string `json:"startAt"`
	EndAt       string `json:"endAt"`
	ActorID     any    `json:"actorId"`
	ClusterID   any    `json:"clusterId"`
	WorkspaceID any    `json:"workspaceId"`
	ProjectID   any    `json:"projectId"`
	Action      string `json:"action"`
	Outcome     string `json:"outcome"`
	Result      string `json:"result"`
	Resource    string `json:"resource"`
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
	cluster, err := parseUint64(req.ClusterID, false, "clusterId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspace, err := parseUint64(req.WorkspaceID, false, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	project, err := parseUint64(req.ProjectID, false, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var actorID *uint64
	if actor != 0 {
		actorID = &actor
	}
	var clusterID *uint64
	if cluster != 0 {
		clusterID = &cluster
	}
	var workspaceID *uint64
	if workspace != 0 {
		workspaceID = &workspace
	}
	var projectID *uint64
	if project != 0 {
		projectID = &project
	}

	task, err := h.svc.SubmitExport(c.Request.Context(), c.GetUint64(middleware.UserIDKey), auditSvc.SubmitExportRequest{
		StartAt:     startAt,
		EndAt:       endAt,
		ActorID:     actorID,
		ClusterID:   clusterID,
		WorkspaceID: workspaceID,
		ProjectID:   projectID,
		Action:      strings.TrimSpace(req.Action),
		Outcome:     strings.TrimSpace(req.Outcome),
		Result:      strings.TrimSpace(req.Result),
		Resource:    strings.TrimSpace(req.Resource),
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
	task, err := h.svc.GetExportTaskForViewer(c.Request.Context(), c.Param("taskId"), c.GetUint64(middleware.UserIDKey))
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

func (h *AuditHandler) DownloadExport(c *gin.Context) {
	result, err := h.svc.GetExportDownloadForViewer(c.Request.Context(), c.Param("taskId"), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeAuditError(c, err)
		return
	}

	fileName := strings.TrimSpace(result.FileName)
	if fileName == "" {
		fileName = "audit-export.csv"
	}
	contentType := strings.TrimSpace(result.ContentType)
	if contentType == "" {
		contentType = "text/csv; charset=utf-8"
	}

	c.Header("Content-Type", contentType)
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%q", fileName))
	c.Header("Cache-Control", "no-store")
	c.Data(http.StatusOK, contentType, result.Data)
}

func writeAuditError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	lower := strings.ToLower(err.Error())
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case strings.Contains(lower, "not ready"):
		status = http.StatusConflict
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

func parseOptionalQueryUint64(c *gin.Context, field string) (*uint64, error) {
	raw := strings.TrimSpace(c.Query(field))
	if raw == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%s must be a positive integer", field)
	}
	return &parsed, nil
}

func firstNonEmptyQuery(c *gin.Context, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(c.Query(key)); value != "" {
			return value
		}
	}
	return ""
}
