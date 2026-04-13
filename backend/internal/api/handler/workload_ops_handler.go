package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	workloadops "kbmanage/backend/internal/service/workloadops"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkloadOpsHandler struct {
	svc *workloadops.Service
}

func NewWorkloadOpsHandler(svc *workloadops.Service) *WorkloadOpsHandler {
	if svc == nil {
		svc = workloadops.NewService(nil, nil, nil, nil, nil, nil)
	}
	return &WorkloadOpsHandler{svc: svc}
}

type submitWorkloadActionRequest struct {
	ClusterID         uint64         `json:"clusterId"`
	WorkspaceID       uint64         `json:"workspaceId"`
	ProjectID         uint64         `json:"projectId"`
	Namespace         string         `json:"namespace"`
	ResourceKind      string         `json:"resourceKind"`
	ResourceName      string         `json:"resourceName"`
	TargetInstanceRef string         `json:"targetInstanceRef"`
	ActionType        string         `json:"actionType"`
	RiskConfirmed     bool           `json:"riskConfirmed"`
	Payload           map[string]any `json:"payload"`
}

type submitBatchOperationRequest struct {
	ActionType    string         `json:"actionType"`
	RiskConfirmed bool           `json:"riskConfirmed"`
	Payload       map[string]any `json:"payload"`
	Targets       []struct {
		ClusterID    uint64 `json:"clusterId"`
		WorkspaceID  uint64 `json:"workspaceId"`
		ProjectID    uint64 `json:"projectId"`
		Namespace    string `json:"namespace"`
		ResourceKind string `json:"resourceKind"`
		ResourceName string `json:"resourceName"`
	} `json:"targets"`
}

type createTerminalSessionRequest struct {
	ClusterID     uint64 `json:"clusterId"`
	WorkspaceID   uint64 `json:"workspaceId"`
	ProjectID     uint64 `json:"projectId"`
	Namespace     string `json:"namespace"`
	PodName       string `json:"podName"`
	ContainerName string `json:"containerName"`
	WorkloadKind  string `json:"workloadKind"`
	WorkloadName  string `json:"workloadName"`
	Cols          int    `json:"cols"`
	Rows          int    `json:"rows"`
}

func (h *WorkloadOpsHandler) GetContext(c *gin.Context) {
	target, err := parseWorkloadReferenceFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.GetContext(c.Request.Context(), c.GetUint64(middleware.UserIDKey), target)
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *WorkloadOpsHandler) ListInstances(c *gin.Context) {
	target, err := parseWorkloadReferenceFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListInstances(c.Request.Context(), c.GetUint64(middleware.UserIDKey), target)
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *WorkloadOpsHandler) ListRevisions(c *gin.Context) {
	target, err := parseWorkloadReferenceFromQuery(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListRevisions(c.Request.Context(), c.GetUint64(middleware.UserIDKey), target)
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *WorkloadOpsHandler) SubmitAction(c *gin.Context) {
	var req submitWorkloadActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.SubmitAction(c.Request.Context(), workloadops.SubmitWorkloadActionRequest{
		RequestID:  c.GetString(middleware.RequestIDKey),
		OperatorID: c.GetUint64(middleware.UserIDKey),
		Target: workloadops.WorkloadReference{
			ClusterID:    req.ClusterID,
			WorkspaceID:  req.WorkspaceID,
			ProjectID:    req.ProjectID,
			Namespace:    strings.TrimSpace(req.Namespace),
			ResourceKind: strings.TrimSpace(req.ResourceKind),
			ResourceName: strings.TrimSpace(req.ResourceName),
		},
		TargetInstanceRef: strings.TrimSpace(req.TargetInstanceRef),
		ActionType:        domain.WorkloadActionType(strings.TrimSpace(req.ActionType)),
		RiskLevel:         domain.RiskLevelMedium,
		RiskConfirmed:     req.RiskConfirmed,
		PayloadJSON:       workloadopsEncodePayload(req.Payload),
	})
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, toActionResponse(item))
}

func (h *WorkloadOpsHandler) GetAction(c *gin.Context) {
	actionID, err := strconv.ParseUint(c.Param("actionId"), 10, 64)
	if err != nil || actionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid action id"})
		return
	}
	item, err := h.svc.GetAction(c.Request.Context(), c.GetUint64(middleware.UserIDKey), actionID)
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toActionResponse(item))
}

func (h *WorkloadOpsHandler) SubmitBatch(c *gin.Context) {
	var req submitBatchOperationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	targets := make([]workloadops.WorkloadReference, 0, len(req.Targets))
	for _, t := range req.Targets {
		targets = append(targets, workloadops.WorkloadReference{
			ClusterID:    t.ClusterID,
			WorkspaceID:  t.WorkspaceID,
			ProjectID:    t.ProjectID,
			Namespace:    strings.TrimSpace(t.Namespace),
			ResourceKind: strings.TrimSpace(t.ResourceKind),
			ResourceName: strings.TrimSpace(t.ResourceName),
		})
	}
	task, err := h.svc.SubmitBatch(c.Request.Context(), workloadops.SubmitBatchOperationRequest{
		RequestID:     c.GetString(middleware.RequestIDKey),
		OperatorID:    c.GetUint64(middleware.UserIDKey),
		ActionType:    domain.WorkloadActionType(strings.TrimSpace(req.ActionType)),
		RiskLevel:     domain.RiskLevelHigh,
		RiskConfirmed: req.RiskConfirmed,
		Targets:       targets,
		PayloadJSON:   workloadopsEncodePayload(req.Payload),
	})
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	latestTask, items, getErr := h.svc.GetBatch(c.Request.Context(), c.GetUint64(middleware.UserIDKey), task.ID)
	if getErr != nil {
		writeWorkloadOpsError(c, getErr)
		return
	}
	c.JSON(http.StatusAccepted, toBatchResponse(latestTask, items))
}

func (h *WorkloadOpsHandler) GetBatch(c *gin.Context) {
	batchID, err := strconv.ParseUint(c.Param("batchId"), 10, 64)
	if err != nil || batchID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid batch id"})
		return
	}
	task, items, err := h.svc.GetBatch(c.Request.Context(), c.GetUint64(middleware.UserIDKey), batchID)
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toBatchResponse(task, items))
}

func (h *WorkloadOpsHandler) CreateTerminalSession(c *gin.Context) {
	var req createTerminalSessionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateTerminalSession(c.Request.Context(), workloadops.CreateTerminalSessionRequest{
		OperatorID:    c.GetUint64(middleware.UserIDKey),
		ClusterID:     req.ClusterID,
		WorkspaceID:   req.WorkspaceID,
		ProjectID:     req.ProjectID,
		Namespace:     strings.TrimSpace(req.Namespace),
		PodName:       strings.TrimSpace(req.PodName),
		ContainerName: strings.TrimSpace(req.ContainerName),
		WorkloadKind:  strings.TrimSpace(req.WorkloadKind),
		WorkloadName:  strings.TrimSpace(req.WorkloadName),
		Cols:          req.Cols,
		Rows:          req.Rows,
	})
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"id":            item.ID,
		"status":        item.Status,
		"podName":       item.PodName,
		"containerName": item.ContainerName,
		"workloadKind":  item.WorkloadKind,
		"workloadName":  item.WorkloadName,
		"streamUrl":     "",
		"streamToken":   "",
		"startedAt":     item.StartedAt,
	})
}

func (h *WorkloadOpsHandler) GetTerminalSession(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	if err != nil || sessionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	item, err := h.svc.GetTerminalSession(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sessionID)
	if err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"id":              item.ID,
		"status":          item.Status,
		"podName":         item.PodName,
		"containerName":   item.ContainerName,
		"workloadKind":    item.WorkloadKind,
		"workloadName":    item.WorkloadName,
		"startedAt":       item.StartedAt,
		"endedAt":         item.EndedAt,
		"durationSeconds": item.DurationSeconds,
		"closeReason":     item.CloseReason,
	})
}

func (h *WorkloadOpsHandler) CloseTerminalSession(c *gin.Context) {
	sessionID, err := strconv.ParseUint(c.Param("sessionId"), 10, 64)
	if err != nil || sessionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session id"})
		return
	}
	if err := h.svc.CloseTerminalSession(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sessionID); err != nil {
		writeWorkloadOpsError(c, err)
		return
	}
	c.Status(http.StatusNoContent)
}

func parseWorkloadReferenceFromQuery(c *gin.Context) (workloadops.WorkloadReference, error) {
	clusterID, err := strconv.ParseUint(c.Query("clusterId"), 10, 64)
	if err != nil || clusterID == 0 {
		return workloadops.WorkloadReference{}, errors.New("invalid clusterId")
	}
	target := workloadops.WorkloadReference{
		ClusterID:    clusterID,
		Namespace:    strings.TrimSpace(c.Query("namespace")),
		ResourceKind: strings.TrimSpace(c.Query("resourceKind")),
		ResourceName: strings.TrimSpace(c.Query("resourceName")),
	}
	if workspace := strings.TrimSpace(c.Query("workspaceId")); workspace != "" {
		if n, parseErr := strconv.ParseUint(workspace, 10, 64); parseErr == nil {
			target.WorkspaceID = n
		}
	}
	if project := strings.TrimSpace(c.Query("projectId")); project != "" {
		if n, parseErr := strconv.ParseUint(project, 10, 64); parseErr == nil {
			target.ProjectID = n
		}
	}
	if target.Namespace == "" || target.ResourceKind == "" || target.ResourceName == "" {
		return workloadops.WorkloadReference{}, errors.New("namespace/resourceKind/resourceName are required")
	}
	return target, nil
}

func writeWorkloadOpsError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	msg := err.Error()
	switch {
	case errors.Is(err, workloadops.ErrWorkloadOpsScopeDenied):
		status = http.StatusForbidden
	case errors.Is(err, workloadops.ErrInvalidWorkloadReference):
		status = http.StatusBadRequest
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case strings.Contains(strings.ToLower(msg), "required"), strings.Contains(strings.ToLower(msg), "invalid"):
		status = http.StatusBadRequest
	}
	c.JSON(status, gin.H{"error": msg})
}

func workloadopsEncodePayload(payload map[string]any) string {
	if len(payload) == 0 {
		return "{}"
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return "{}"
	}
	return string(encoded)
}

func toActionResponse(item *domain.WorkloadActionRequest) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":              item.ID,
		"actionType":      item.ActionType,
		"status":          item.Status,
		"riskLevel":       item.RiskLevel,
		"progressMessage": item.ProgressMessage,
		"resultMessage":   item.ResultMessage,
		"failureReason":   item.FailureReason,
		"startedAt":       item.StartedAt,
		"completedAt":     item.CompletedAt,
	}
}

func toBatchResponse(task *domain.BatchOperationTask, items []domain.BatchOperationItem) gin.H {
	if task == nil {
		return gin.H{}
	}
	respItems := make([]gin.H, 0, len(items))
	for _, item := range items {
		respItems = append(respItems, gin.H{
			"resourceRef":   item.ResourceKind + "/" + item.Namespace + "/" + item.ResourceName,
			"status":        item.Status,
			"resultMessage": item.ResultMessage,
			"failureReason": item.FailureReason,
		})
	}
	return gin.H{
		"id":               task.ID,
		"actionType":       task.ActionType,
		"status":           task.Status,
		"totalTargets":     task.TotalTargets,
		"succeededTargets": task.SucceededTargets,
		"failedTargets":    task.FailedTargets,
		"canceledTargets":  task.CanceledTargets,
		"progressPercent":  task.ProgressPercent,
		"items":            respItems,
		"startedAt":        task.StartedAt,
		"completedAt":      task.CompletedAt,
	}
}
