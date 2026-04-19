package handler

import (
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	sreint "kbmanage/backend/internal/integration/sre"
	sreSvc "kbmanage/backend/internal/service/sre"

	"github.com/gin-gonic/gin"
)

type SREHandler struct {
	svc *sreSvc.Service
}

func NewSREHandler(svc *sreSvc.Service) *SREHandler { return &SREHandler{svc: svc} }

func (h *SREHandler) ListHAPolicies(c *gin.Context) {
	items, err := h.svc.ListHAPolicies(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sreSvc.HAPolicyListFilter{
		Status:  strings.TrimSpace(c.Query("status")),
		Keyword: strings.TrimSpace(c.Query("keyword")),
	})
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *SREHandler) UpsertHAPolicy(c *gin.Context) {
	var req sreSvc.HAPolicyInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.UpsertHAPolicy(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *SREHandler) GetHealthOverview(c *gin.Context) {
	workspaceID, _ := parseOptionalSREUint64(c.Query("workspaceId"))
	projectID, hasProject := parseOptionalSREUint64(c.Query("projectId"))
	var pid *uint64
	if hasProject {
		pid = &projectID
	}
	item, err := h.svc.GetHealthOverview(c.Request.Context(), c.GetUint64(middleware.UserIDKey), workspaceID, pid)
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *SREHandler) ListMaintenanceWindows(c *gin.Context) {
	items, err := h.svc.ListMaintenanceWindows(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *SREHandler) UpsertMaintenanceWindow(c *gin.Context) {
	var req sreSvc.MaintenanceWindowInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.UpsertMaintenanceWindow(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *SREHandler) RunUpgradePrecheck(c *gin.Context) {
	var req sreSvc.UpgradePrecheckInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.RunUpgradePrecheck(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *SREHandler) ListUpgradePlans(c *gin.Context) {
	items, err := h.svc.ListUpgradePlans(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *SREHandler) CreateUpgradePlan(c *gin.Context) {
	var req sreSvc.SREUpgradePlanInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateUpgradePlan(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *SREHandler) CreateRollbackValidation(c *gin.Context) {
	upgradeID, ok := parseSREUint64Param(c, "upgradeId")
	if !ok {
		return
	}
	var req sreSvc.RollbackValidationInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateRollbackValidation(c.Request.Context(), c.GetUint64(middleware.UserIDKey), upgradeID, req)
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *SREHandler) ListCapacityBaselines(c *gin.Context) {
	items, err := h.svc.ListCapacityBaselines(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sreSvc.CapacityBaselineListFilter{
		Status: strings.TrimSpace(c.Query("status")),
	})
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *SREHandler) ListScaleEvidence(c *gin.Context) {
	items, err := h.svc.ListScaleEvidence(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sreSvc.ScaleEvidenceListFilter{
		EvidenceType: strings.TrimSpace(c.Query("evidenceType")),
	})
	if err != nil {
		writeSREError(c, err)
		return
	}
	for i := range items {
		analysis := sreint.NewStaticScaleAnalyzer().Analyze(c.Request.Context(), sreint.ScaleEvidenceInput{
			EvidenceType:    items[i].EvidenceType,
			Summary:         items[i].Summary,
			ForecastSummary: items[i].ForecastSummary,
			ConfidenceLevel: items[i].ConfidenceLevel,
		})
		if strings.TrimSpace(items[i].BottleneckSummary) == "" {
			items[i].BottleneckSummary = analysis.Bottleneck
		}
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *SREHandler) ListRunbooks(c *gin.Context) {
	items, err := h.svc.ListRunbooks(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeSREError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func parseSREUint64Param(c *gin.Context, name string) (uint64, bool) {
	value := strings.TrimSpace(c.Param(name))
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return id, true
}

func parseOptionalSREUint64(v string) (uint64, bool) {
	if strings.TrimSpace(v) == "" {
		return 0, false
	}
	out, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, false
	}
	return out, true
}

func writeSREError(c *gin.Context, err error) {
	switch err {
	case sreSvc.ErrSREScopeDenied:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case sreSvc.ErrSREInvalid:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
