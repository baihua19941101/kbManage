package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
	gitopsSvc "kbmanage/backend/internal/service/gitops"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type GitOpsHandler struct {
	svc         *gitopsSvc.Service
	auditWriter *auditSvc.EventWriter
}

func NewGitOpsHandler(svc *gitopsSvc.Service, auditWriter ...*auditSvc.EventWriter) *GitOpsHandler {
	if svc == nil {
		svc = gitopsSvc.NewService(nil, nil, nil, nil, nil, nil, nil, nil, nil, nil)
	}
	var writer *auditSvc.EventWriter
	if len(auditWriter) > 0 {
		writer = auditWriter[0]
	}
	return &GitOpsHandler{
		svc:         svc,
		auditWriter: writer,
	}
}

type createGitOpsSourceRequest struct {
	Name          string `json:"name"`
	SourceType    string `json:"sourceType"`
	Endpoint      string `json:"endpoint"`
	DefaultRef    string `json:"defaultRef"`
	CredentialRef string `json:"credentialRef"`
	WorkspaceID   uint64 `json:"workspaceId"`
	ProjectID     uint64 `json:"projectId"`
}

type updateGitOpsSourceRequest struct {
	Name          *string `json:"name"`
	DefaultRef    *string `json:"defaultRef"`
	CredentialRef *string `json:"credentialRef"`
	Disabled      *bool   `json:"disabled"`
}

type createGitOpsTargetGroupRequest struct {
	Name            string   `json:"name"`
	WorkspaceID     uint64   `json:"workspaceId"`
	ProjectID       uint64   `json:"projectId"`
	ClusterRefs     []uint64 `json:"clusterRefs"`
	SelectorSummary string   `json:"selectorSummary"`
	Description     string   `json:"description"`
}

type updateGitOpsTargetGroupRequest struct {
	Name            *string   `json:"name"`
	ClusterRefs     *[]uint64 `json:"clusterRefs"`
	SelectorSummary *string   `json:"selectorSummary"`
	Description     *string   `json:"description"`
	Disabled        *bool     `json:"disabled"`
}

type gitOpsDeliveryUnitEnvironmentRequest struct {
	Name          string `json:"name"`
	OrderIndex    int    `json:"orderIndex"`
	TargetGroupID uint64 `json:"targetGroupId"`
	PromotionMode string `json:"promotionMode"`
	Paused        bool   `json:"paused"`
}

type gitOpsDeliveryUnitOverlayRequest struct {
	OverlayType    string `json:"overlayType"`
	OverlayRef     string `json:"overlayRef"`
	Precedence     int    `json:"precedence"`
	EffectiveScope string `json:"effectiveScope"`
}

type createGitOpsDeliveryUnitRequest struct {
	Name                 string                                 `json:"name"`
	WorkspaceID          uint64                                 `json:"workspaceId"`
	ProjectID            uint64                                 `json:"projectId"`
	SourceID             uint64                                 `json:"sourceId"`
	SourcePath           string                                 `json:"sourcePath"`
	DefaultNamespace     string                                 `json:"defaultNamespace"`
	SyncMode             string                                 `json:"syncMode"`
	DesiredRevision      string                                 `json:"desiredRevision"`
	DesiredAppVersion    string                                 `json:"desiredAppVersion"`
	DesiredConfigVersion string                                 `json:"desiredConfigVersion"`
	Environments         []gitOpsDeliveryUnitEnvironmentRequest `json:"environments"`
	Overlays             []gitOpsDeliveryUnitOverlayRequest     `json:"overlays"`
}

type updateGitOpsDeliveryUnitRequest struct {
	Name                 *string                                 `json:"name"`
	SourcePath           *string                                 `json:"sourcePath"`
	DefaultNamespace     *string                                 `json:"defaultNamespace"`
	SyncMode             *string                                 `json:"syncMode"`
	DesiredRevision      *string                                 `json:"desiredRevision"`
	DesiredAppVersion    *string                                 `json:"desiredAppVersion"`
	DesiredConfigVersion *string                                 `json:"desiredConfigVersion"`
	Environments         *[]gitOpsDeliveryUnitEnvironmentRequest `json:"environments"`
	Overlays             *[]gitOpsDeliveryUnitOverlayRequest     `json:"overlays"`
}

type submitGitOpsActionRequest struct {
	RequestID          string         `json:"requestId"`
	ActionType         string         `json:"actionType"`
	EnvironmentStageID *uint64        `json:"environmentStageId"`
	TargetReleaseID    *uint64        `json:"targetReleaseId"`
	Payload            map[string]any `json:"payload"`
}

func (h *GitOpsHandler) ListSources(c *gin.Context) {
	workspaceID, err := parseGitOpsOptionalQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseGitOpsOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListSources(
		c.Request.Context(),
		c.GetUint64(middleware.UserIDKey),
		workspaceID,
		projectID,
		gitopsSvc.SourceListFilter{
			SourceType: strings.TrimSpace(c.Query("sourceType")),
			Status:     strings.TrimSpace(c.Query("status")),
		},
	)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	res := make([]gin.H, 0, len(items))
	for i := range items {
		res = append(res, toGitOpsSourceResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

func (h *GitOpsHandler) CreateSource(c *gin.Context) {
	var req createGitOpsSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateSource(c.Request.Context(), c.GetUint64(middleware.UserIDKey), gitopsSvc.CreateSourceInput{
		Name:          strings.TrimSpace(req.Name),
		SourceType:    strings.TrimSpace(req.SourceType),
		Endpoint:      strings.TrimSpace(req.Endpoint),
		DefaultRef:    strings.TrimSpace(req.DefaultRef),
		CredentialRef: strings.TrimSpace(req.CredentialRef),
		WorkspaceID:   req.WorkspaceID,
		ProjectID:     req.ProjectID,
	})
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toGitOpsSourceResponse(item))
}

func (h *GitOpsHandler) GetSource(c *gin.Context) {
	sourceID, err := parseGitOpsRequiredParamUint64(c, "sourceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.GetSource(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sourceID)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toGitOpsSourceResponse(item))
}

func (h *GitOpsHandler) UpdateSource(c *gin.Context) {
	sourceID, err := parseGitOpsRequiredParamUint64(c, "sourceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req updateGitOpsSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.UpdateSource(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sourceID, gitopsSvc.UpdateSourceInput{
		Name:          req.Name,
		DefaultRef:    req.DefaultRef,
		CredentialRef: req.CredentialRef,
		Disabled:      req.Disabled,
	})
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toGitOpsSourceResponse(item))
}

func (h *GitOpsHandler) VerifySource(c *gin.Context) {
	sourceID, err := parseGitOpsRequiredParamUint64(c, "sourceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.VerifySource(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sourceID, c.GetString(middleware.RequestIDKey))
	if err != nil {
		h.writeGitOpsAuditEvent(
			c,
			auditSvc.GitOpsAuditActionSourceVerify,
			"source:"+strconv.FormatUint(sourceID, 10),
			domain.AuditOutcomeDenied,
			map[string]any{
				"sourceId": sourceID,
				"error":    err.Error(),
			},
		)
		writeGitOpsError(c, err)
		return
	}
	h.writeGitOpsAuditEvent(
		c,
		auditSvc.GitOpsAuditActionSourceVerify,
		"source:"+strconv.FormatUint(sourceID, 10),
		domain.AuditOutcomeSuccess,
		map[string]any{
			"sourceId":    sourceID,
			"operationId": item.ID,
		},
	)
	c.JSON(http.StatusAccepted, toGitOpsOperationResponse(item))
}

func (h *GitOpsHandler) ListTargetGroups(c *gin.Context) {
	workspaceID, err := parseGitOpsRequiredQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseGitOpsOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListTargetGroups(c.Request.Context(), c.GetUint64(middleware.UserIDKey), workspaceID, projectID)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	res := make([]gin.H, 0, len(items))
	for i := range items {
		res = append(res, toGitOpsTargetGroupResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

func (h *GitOpsHandler) CreateTargetGroup(c *gin.Context) {
	var req createGitOpsTargetGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateTargetGroup(c.Request.Context(), c.GetUint64(middleware.UserIDKey), gitopsSvc.CreateTargetGroupInput{
		Name:            strings.TrimSpace(req.Name),
		WorkspaceID:     req.WorkspaceID,
		ProjectID:       req.ProjectID,
		ClusterRefs:     req.ClusterRefs,
		SelectorSummary: strings.TrimSpace(req.SelectorSummary),
		Description:     strings.TrimSpace(req.Description),
	})
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toGitOpsTargetGroupResponse(item))
}

func (h *GitOpsHandler) GetTargetGroup(c *gin.Context) {
	targetGroupID, err := parseGitOpsRequiredParamUint64(c, "targetGroupId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.GetTargetGroup(c.Request.Context(), c.GetUint64(middleware.UserIDKey), targetGroupID)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toGitOpsTargetGroupResponse(item))
}

func (h *GitOpsHandler) UpdateTargetGroup(c *gin.Context) {
	targetGroupID, err := parseGitOpsRequiredParamUint64(c, "targetGroupId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req updateGitOpsTargetGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.UpdateTargetGroup(c.Request.Context(), c.GetUint64(middleware.UserIDKey), targetGroupID, gitopsSvc.UpdateTargetGroupInput{
		Name:            req.Name,
		ClusterRefs:     req.ClusterRefs,
		SelectorSummary: req.SelectorSummary,
		Description:     req.Description,
		Disabled:        req.Disabled,
	})
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toGitOpsTargetGroupResponse(item))
}

func (h *GitOpsHandler) ListDeliveryUnits(c *gin.Context) {
	workspaceID, err := parseGitOpsRequiredQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseGitOpsOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListDeliveryUnits(c.Request.Context(), c.GetUint64(middleware.UserIDKey), workspaceID, projectID)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	res := make([]gin.H, 0, len(items))
	for i := range items {
		res = append(res, toGitOpsDeliveryUnitSummaryResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

func (h *GitOpsHandler) CreateDeliveryUnit(c *gin.Context) {
	var req createGitOpsDeliveryUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateDeliveryUnit(c.Request.Context(), c.GetUint64(middleware.UserIDKey), gitopsSvc.CreateDeliveryUnitInput{
		Name:                 strings.TrimSpace(req.Name),
		WorkspaceID:          req.WorkspaceID,
		ProjectID:            req.ProjectID,
		SourceID:             req.SourceID,
		SourcePath:           strings.TrimSpace(req.SourcePath),
		DefaultNamespace:     strings.TrimSpace(req.DefaultNamespace),
		SyncMode:             strings.TrimSpace(req.SyncMode),
		DesiredRevision:      strings.TrimSpace(req.DesiredRevision),
		DesiredAppVersion:    strings.TrimSpace(req.DesiredAppVersion),
		DesiredConfigVersion: strings.TrimSpace(req.DesiredConfigVersion),
		Environments:         toGitOpsSvcEnvironmentInputs(req.Environments),
		Overlays:             toGitOpsSvcOverlayInputs(req.Overlays),
	})
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusCreated, toGitOpsDeliveryUnitDetailResponse(item))
}

func (h *GitOpsHandler) GetDeliveryUnit(c *gin.Context) {
	unitID, err := parseGitOpsRequiredParamUint64(c, "unitId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.GetDeliveryUnit(c.Request.Context(), c.GetUint64(middleware.UserIDKey), unitID)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toGitOpsDeliveryUnitDetailResponse(item))
}

func (h *GitOpsHandler) UpdateDeliveryUnit(c *gin.Context) {
	unitID, err := parseGitOpsRequiredParamUint64(c, "unitId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req updateGitOpsDeliveryUnitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := gitopsSvc.UpdateDeliveryUnitInput{
		Name:                 req.Name,
		SourcePath:           req.SourcePath,
		DefaultNamespace:     req.DefaultNamespace,
		SyncMode:             req.SyncMode,
		DesiredRevision:      req.DesiredRevision,
		DesiredAppVersion:    req.DesiredAppVersion,
		DesiredConfigVersion: req.DesiredConfigVersion,
	}
	if req.Environments != nil {
		envs := toGitOpsSvcEnvironmentInputs(*req.Environments)
		input.Environments = &envs
	}
	if req.Overlays != nil {
		overlays := toGitOpsSvcOverlayInputs(*req.Overlays)
		input.Overlays = &overlays
	}
	item, err := h.svc.UpdateDeliveryUnit(c.Request.Context(), c.GetUint64(middleware.UserIDKey), unitID, input)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toGitOpsDeliveryUnitDetailResponse(item))
}

func (h *GitOpsHandler) GetDeliveryUnitStatus(c *gin.Context) {
	unitID, err := parseGitOpsRequiredParamUint64(c, "unitId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.GetDeliveryUnitStatus(
		c.Request.Context(),
		c.GetUint64(middleware.UserIDKey),
		unitID,
		strings.TrimSpace(c.Query("environment")),
	)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GitOpsHandler) GetDeliveryUnitDiff(c *gin.Context) {
	unitID, err := parseGitOpsRequiredParamUint64(c, "unitId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stageID, err := parseGitOpsOptionalQueryUint64(c, "stageId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	res, err := h.svc.GetDeliveryUnitDiff(c.Request.Context(), c.GetUint64(middleware.UserIDKey), unitID, stageID)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, res)
}

func (h *GitOpsHandler) ListReleaseRevisions(c *gin.Context) {
	unitID, err := parseGitOpsRequiredParamUint64(c, "unitId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListReleaseRevisions(c.Request.Context(), c.GetUint64(middleware.UserIDKey), unitID)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *GitOpsHandler) SubmitAction(c *gin.Context) {
	unitID, err := parseGitOpsRequiredParamUint64(c, "unitId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req submitGitOpsActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	payloadJSON, err := marshalPayload(req.Payload)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
		return
	}
	item, err := h.svc.SubmitAction(c.Request.Context(), c.GetUint64(middleware.UserIDKey), gitopsSvc.ActionRequest{
		RequestID:          strings.TrimSpace(req.RequestID),
		DeliveryUnitID:     unitID,
		EnvironmentStageID: req.EnvironmentStageID,
		ActionType:         domain.DeliveryActionType(strings.TrimSpace(req.ActionType)),
		TargetReleaseID:    req.TargetReleaseID,
		PayloadJSON:        payloadJSON,
	})
	if err != nil {
		h.writeGitOpsAuditEvent(
			c,
			auditActionForGitOpsOperation(strings.TrimSpace(req.ActionType)),
			"unit:"+strconv.FormatUint(unitID, 10),
			domain.AuditOutcomeDenied,
			map[string]any{
				"unitId":     unitID,
				"actionType": strings.TrimSpace(req.ActionType),
				"error":      err.Error(),
			},
		)
		writeGitOpsError(c, err)
		return
	}
	h.writeGitOpsAuditEvent(
		c,
		auditActionForGitOpsOperation(strings.TrimSpace(req.ActionType)),
		"unit:"+strconv.FormatUint(unitID, 10),
		domain.AuditOutcomeSuccess,
		map[string]any{
			"unitId":          unitID,
			"actionType":      strings.TrimSpace(req.ActionType),
			"operationId":     item.ID,
			"targetReleaseId": item.TargetReleaseID,
		},
	)
	c.JSON(http.StatusAccepted, toGitOpsOperationResponse(item))
}

func (h *GitOpsHandler) GetOperation(c *gin.Context) {
	operationID, err := parseGitOpsRequiredParamUint64(c, "operationId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.GetOperation(c.Request.Context(), c.GetUint64(middleware.UserIDKey), operationID)
	if err != nil {
		writeGitOpsError(c, err)
		return
	}
	c.JSON(http.StatusOK, toGitOpsOperationResponse(item))
}

func toGitOpsSourceResponse(item *domain.DeliverySource) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":               item.ID,
		"name":             item.Name,
		"sourceType":       item.SourceType,
		"endpoint":         item.Endpoint,
		"defaultRef":       item.DefaultRef,
		"credentialRef":    item.CredentialRef,
		"workspaceId":      item.WorkspaceID,
		"projectId":        item.ProjectID,
		"status":           item.Status,
		"lastVerifiedAt":   item.LastVerifiedAt,
		"lastErrorMessage": item.LastErrorMessage,
		"createdAt":        item.CreatedAt,
		"updatedAt":        item.UpdatedAt,
	}
}

func toGitOpsTargetGroupResponse(item *domain.ClusterTargetGroup) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":              item.ID,
		"name":            item.Name,
		"workspaceId":     item.WorkspaceID,
		"projectId":       item.ProjectID,
		"clusterRefs":     decodeGitOpsUint64Array(item.ClusterRefsJSON),
		"selectorSummary": item.ClusterSelectorSnapshot,
		"description":     item.Description,
		"status":          item.Status,
		"createdAt":       item.CreatedAt,
		"updatedAt":       item.UpdatedAt,
	}
}

func toGitOpsDeliveryUnitSummaryResponse(item *domain.ApplicationDeliveryUnit) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":                   item.ID,
		"name":                 item.Name,
		"workspaceId":          item.WorkspaceID,
		"projectId":            item.ProjectID,
		"sourceId":             item.SourceID,
		"sourcePath":           item.SourcePath,
		"defaultNamespace":     item.DefaultNamespace,
		"syncMode":             item.SyncMode,
		"desiredRevision":      item.DesiredRevision,
		"desiredAppVersion":    item.DesiredAppVersion,
		"desiredConfigVersion": item.DesiredConfigVersion,
		"paused":               item.Paused,
		"deliveryStatus":       item.DeliveryStatus,
		"lastSyncedAt":         item.LastSyncedAt,
		"lastReleaseId":        item.LastReleaseID,
		"createdAt":            item.CreatedAt,
		"updatedAt":            item.UpdatedAt,
	}
}

func toGitOpsDeliveryUnitDetailResponse(item *gitopsSvc.DeliveryUnitDetail) gin.H {
	if item == nil {
		return gin.H{}
	}
	envs := make([]gin.H, 0, len(item.Environments))
	for i := range item.Environments {
		envs = append(envs, gin.H{
			"id":            item.Environments[i].ID,
			"name":          item.Environments[i].Name,
			"orderIndex":    item.Environments[i].OrderIndex,
			"targetGroupId": item.Environments[i].TargetGroupID,
			"promotionMode": item.Environments[i].PromotionMode,
			"paused":        item.Environments[i].Paused,
			"status":        item.Environments[i].Status,
		})
	}
	overlays := make([]gin.H, 0, len(item.Overlays))
	for i := range item.Overlays {
		overlays = append(overlays, gin.H{
			"id":             item.Overlays[i].ID,
			"overlayType":    item.Overlays[i].OverlayType,
			"overlayRef":     item.Overlays[i].OverlayRef,
			"precedence":     item.Overlays[i].Precedence,
			"effectiveScope": item.Overlays[i].EffectiveScopeJSON,
		})
	}
	res := toGitOpsDeliveryUnitSummaryResponse(&item.Unit)
	res["environments"] = envs
	res["overlays"] = overlays
	return res
}

func toGitOpsSvcEnvironmentInputs(items []gitOpsDeliveryUnitEnvironmentRequest) []gitopsSvc.EnvironmentStageInput {
	res := make([]gitopsSvc.EnvironmentStageInput, 0, len(items))
	for i := range items {
		res = append(res, gitopsSvc.EnvironmentStageInput{
			Name:          strings.TrimSpace(items[i].Name),
			OrderIndex:    items[i].OrderIndex,
			TargetGroupID: items[i].TargetGroupID,
			PromotionMode: strings.TrimSpace(items[i].PromotionMode),
			Paused:        items[i].Paused,
		})
	}
	return res
}

func toGitOpsSvcOverlayInputs(items []gitOpsDeliveryUnitOverlayRequest) []gitopsSvc.ConfigurationOverlayInput {
	res := make([]gitopsSvc.ConfigurationOverlayInput, 0, len(items))
	for i := range items {
		res = append(res, gitopsSvc.ConfigurationOverlayInput{
			OverlayType:    strings.TrimSpace(items[i].OverlayType),
			OverlayRef:     strings.TrimSpace(items[i].OverlayRef),
			Precedence:     items[i].Precedence,
			EffectiveScope: strings.TrimSpace(items[i].EffectiveScope),
		})
	}
	return res
}

func decodeGitOpsUint64Array(raw string) []uint64 {
	if strings.TrimSpace(raw) == "" {
		return []uint64{}
	}
	items := make([]uint64, 0)
	if err := json.Unmarshal([]byte(raw), &items); err != nil {
		return []uint64{}
	}
	return items
}

func toGitOpsOperationResponse(item *domain.DeliveryOperation) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":                 item.ID,
		"requestId":          item.RequestID,
		"operatorId":         item.OperatorID,
		"deliveryUnitId":     item.DeliveryUnitID,
		"environmentStageId": item.EnvironmentStageID,
		"actionType":         item.ActionType,
		"targetReleaseId":    item.TargetReleaseID,
		"status":             item.Status,
		"progressPercent":    item.ProgressPercent,
		"resultSummary":      item.ResultSummary,
		"failureReason":      item.FailureReason,
		"startedAt":          item.StartedAt,
		"completedAt":        item.CompletedAt,
		"createdAt":          item.CreatedAt,
		"updatedAt":          item.UpdatedAt,
	}
}

func marshalPayload(payload map[string]any) (string, error) {
	if len(payload) == 0 {
		return "", nil
	}
	data, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func writeGitOpsError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	lower := strings.ToLower(err.Error())
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case strings.Contains(lower, "scope access denied"), strings.Contains(lower, "forbidden"):
		status = http.StatusForbidden
	case strings.Contains(lower, "required"), strings.Contains(lower, "invalid"):
		status = http.StatusBadRequest
	case strings.Contains(lower, "locked"):
		status = http.StatusConflict
	case strings.Contains(lower, "not configured"):
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, gin.H{"error": err.Error()})
}

func parseGitOpsRequiredParamUint64(c *gin.Context, key string) (uint64, error) {
	value := strings.TrimSpace(c.Param(key))
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		return 0, errors.New("invalid " + key)
	}
	return parsed, nil
}

func parseGitOpsRequiredQueryUint64(c *gin.Context, key string) (uint64, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return 0, errors.New(key + " is required")
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		return 0, errors.New("invalid " + key)
	}
	return parsed, nil
}

func parseGitOpsOptionalQueryUint64(c *gin.Context, key string) (uint64, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil {
		return 0, errors.New("invalid " + key)
	}
	return parsed, nil
}

func (h *GitOpsHandler) writeGitOpsAuditEvent(
	c *gin.Context,
	action string,
	resourceID string,
	outcome domain.AuditOutcome,
	details map[string]any,
) {
	if h == nil || h.auditWriter == nil {
		return
	}
	userID := c.GetUint64(middleware.UserIDKey)
	if userID == 0 {
		return
	}
	actorID := userID
	_ = h.auditWriter.WriteGitOpsEvent(
		c.Request.Context(),
		c.GetString(middleware.RequestIDKey),
		&actorID,
		action,
		resourceID,
		outcome,
		details,
	)
}

func auditActionForGitOpsOperation(actionType string) string {
	switch strings.ToLower(strings.TrimSpace(actionType)) {
	case "install":
		return auditSvc.GitOpsAuditActionInstallSubmit
	case "upgrade":
		return auditSvc.GitOpsAuditActionUpgradeSubmit
	case "resync":
		return auditSvc.GitOpsAuditActionResyncSubmit
	case "promote":
		return auditSvc.GitOpsAuditActionPromoteSubmit
	case "rollback":
		return auditSvc.GitOpsAuditActionRollbackSubmit
	case "pause":
		return auditSvc.GitOpsAuditActionPauseSubmit
	case "resume":
		return auditSvc.GitOpsAuditActionResumeSubmit
	case "uninstall":
		return auditSvc.GitOpsAuditActionUninstallSubmit
	default:
		return auditSvc.GitOpsAuditActionSyncSubmit
	}
}
