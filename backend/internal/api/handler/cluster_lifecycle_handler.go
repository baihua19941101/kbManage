package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/service/clusterlifecycle"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ClusterLifecycleHandler struct {
	svc *clusterlifecycle.Service
}

func NewClusterLifecycleHandler(svc *clusterlifecycle.Service) *ClusterLifecycleHandler {
	return &ClusterLifecycleHandler{svc: svc}
}

func (h *ClusterLifecycleHandler) ListClusters(c *gin.Context) {
	items, err := h.svc.ListClusters(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterlifecycle.ClusterListFilter{
		Status:             strings.TrimSpace(c.Query("status")),
		InfrastructureType: strings.TrimSpace(c.Query("infrastructureType")),
		DriverKey:          strings.TrimSpace(c.Query("driverKey")),
		Keyword:            strings.TrimSpace(c.Query("keyword")),
	})
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ClusterLifecycleHandler) ImportCluster(c *gin.Context) {
	var req clusterlifecycle.ImportClusterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	op, cluster, err := h.svc.ImportCluster(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"operation": op, "cluster": cluster})
}

func (h *ClusterLifecycleHandler) RegisterCluster(c *gin.Context) {
	var req clusterlifecycle.RegisterClusterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	bundle, err := h.svc.RegisterCluster(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"clusterId":         bundle.ClusterID,
		"registrationToken": strconv.FormatUint(bundle.ClusterID, 10),
		"commandSnippet":    bundle.Command,
		"expiresAt":         "",
		"status":            bundle.Status,
		"instructions":      bundle.Instructions,
	})
}

func (h *ClusterLifecycleHandler) GetCluster(c *gin.Context) {
	clusterID, ok := parseUint64Param(c, "clusterId")
	if !ok {
		return
	}
	detail, err := h.svc.GetCluster(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterID)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusOK, detail)
}

func (h *ClusterLifecycleHandler) CreateCluster(c *gin.Context) {
	var req clusterlifecycle.CreateClusterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	op, cluster, err := h.svc.CreateCluster(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, gin.H{"operation": op, "cluster": cluster})
}

func (h *ClusterLifecycleHandler) ValidateClusterChange(c *gin.Context) {
	clusterID, ok := parseUint64Param(c, "clusterId")
	if !ok {
		return
	}
	var req clusterlifecycle.ValidationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ValidateChange(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterID, req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	blockers := make([]string, 0)
	warnings := make([]string, 0)
	passed := make([]string, 0)
	for _, check := range result.Checks {
		switch {
		case !check.Passed && string(check.Severity) == "blocker":
			blockers = append(blockers, check.Message)
		case !check.Passed:
			warnings = append(warnings, check.Message)
		default:
			passed = append(passed, check.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"overallStatus": result.Status,
		"blockers":      blockers,
		"warnings":      warnings,
		"passedChecks":  passed,
	})
}

func (h *ClusterLifecycleHandler) CreateUpgradePlan(c *gin.Context) {
	clusterID, ok := parseUint64Param(c, "clusterId")
	if !ok {
		return
	}
	var req clusterlifecycle.CreateUpgradePlanInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	plan, err := h.svc.CreateUpgradePlan(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterID, req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, plan)
}

func (h *ClusterLifecycleHandler) ExecuteUpgradePlan(c *gin.Context) {
	clusterID, ok := parseUint64Param(c, "clusterId")
	if !ok {
		return
	}
	planID, ok := parseUint64Param(c, "planId")
	if !ok {
		return
	}
	op, err := h.svc.ExecuteUpgradePlan(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterID, planID)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, op)
}

func (h *ClusterLifecycleHandler) ListNodePools(c *gin.Context) {
	clusterID, ok := parseUint64Param(c, "clusterId")
	if !ok {
		return
	}
	items, err := h.svc.ListNodePools(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterID)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ClusterLifecycleHandler) ScaleNodePool(c *gin.Context) {
	clusterID, ok := parseUint64Param(c, "clusterId")
	if !ok {
		return
	}
	nodePoolID, ok := parseUint64Param(c, "nodePoolId")
	if !ok {
		return
	}
	var req clusterlifecycle.ScaleNodePoolInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	op, err := h.svc.ScaleNodePool(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterID, nodePoolID, req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, op)
}

func (h *ClusterLifecycleHandler) DisableCluster(c *gin.Context) {
	clusterID, ok := parseUint64Param(c, "clusterId")
	if !ok {
		return
	}
	var req clusterlifecycle.DisableClusterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	op, err := h.svc.DisableCluster(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterID, req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, op)
}

func (h *ClusterLifecycleHandler) RetireCluster(c *gin.Context) {
	clusterID, ok := parseUint64Param(c, "clusterId")
	if !ok {
		return
	}
	var req clusterlifecycle.RetireClusterInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	op, err := h.svc.RetireCluster(c.Request.Context(), c.GetUint64(middleware.UserIDKey), clusterID, req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, op)
}

func (h *ClusterLifecycleHandler) ListDrivers(c *gin.Context) {
	items, err := h.svc.ListDrivers(c.Request.Context(), c.GetUint64(middleware.UserIDKey), strings.TrimSpace(c.Query("providerType")))
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ClusterLifecycleHandler) CreateDriver(c *gin.Context) {
	var req clusterlifecycle.CreateDriverInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.UpsertDriver(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ClusterLifecycleHandler) ListDriverCapabilities(c *gin.Context) {
	driverID, ok := parseUint64Param(c, "driverId")
	if !ok {
		return
	}
	items, err := h.svc.ListCapabilities(c.Request.Context(), c.GetUint64(middleware.UserIDKey), driverID)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ClusterLifecycleHandler) ListTemplates(c *gin.Context) {
	items, err := h.svc.ListTemplates(c.Request.Context(), c.GetUint64(middleware.UserIDKey), strings.TrimSpace(c.Query("driverKey")), strings.TrimSpace(c.Query("infrastructureType")))
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ClusterLifecycleHandler) CreateTemplate(c *gin.Context) {
	var req clusterlifecycle.CreateTemplateInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateTemplate(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ClusterLifecycleHandler) ValidateTemplate(c *gin.Context) {
	templateID, ok := parseUint64Param(c, "templateId")
	if !ok {
		return
	}
	var req clusterlifecycle.TemplateValidationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	result, err := h.svc.ValidateTemplate(c.Request.Context(), c.GetUint64(middleware.UserIDKey), templateID, req)
	if err != nil {
		writeClusterLifecycleError(c, err)
		return
	}
	blockers := make([]string, 0)
	warnings := make([]string, 0)
	passed := make([]string, 0)
	for _, check := range result.Checks {
		switch {
		case !check.Passed && string(check.Severity) == "blocker":
			blockers = append(blockers, check.Message)
		case !check.Passed:
			warnings = append(warnings, check.Message)
		default:
			passed = append(passed, check.Message)
		}
	}
	c.JSON(http.StatusOK, gin.H{
		"overallStatus": result.Status,
		"blockers":      blockers,
		"warnings":      warnings,
		"passedChecks":  passed,
	})
}

func parseUint64Param(c *gin.Context, key string) (uint64, bool) {
	value := strings.TrimSpace(c.Param(key))
	out, err := strconv.ParseUint(value, 10, 64)
	if err != nil || out == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
		return 0, false
	}
	return out, true
}

func writeClusterLifecycleError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, clusterlifecycle.ErrLifecycleScopeDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, clusterlifecycle.ErrLifecycleConflict), errors.Is(err, clusterlifecycle.ErrLifecycleBlocked):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, clusterlifecycle.ErrLifecycleInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
