package handler

import (
	"net/http"
	"strconv"

	"kbmanage/backend/internal/api/middleware"
	entSvc "kbmanage/backend/internal/service/enterprise"

	"github.com/gin-gonic/gin"
)

type EnterpriseHandler struct{ svc *entSvc.Service }

func NewEnterpriseHandler(svc *entSvc.Service) *EnterpriseHandler {
	return &EnterpriseHandler{svc: svc}
}

func (h *EnterpriseHandler) ListPermissionTrails(c *gin.Context) {
	items, err := h.svc.ListPermissionTrails(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *EnterpriseHandler) ListKeyOperations(c *gin.Context) {
	items, err := h.svc.ListKeyOperations(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *EnterpriseHandler) ListCoverage(c *gin.Context) {
	items, err := h.svc.ListCoverageSnapshots(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *EnterpriseHandler) ListActionItems(c *gin.Context) {
	items, err := h.svc.ListActionItems(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *EnterpriseHandler) ListReports(c *gin.Context) {
	items, err := h.svc.ListReports(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *EnterpriseHandler) CreateReport(c *gin.Context) {
	var req entSvc.GovernanceReportInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateReport(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *EnterpriseHandler) CreateExportRecord(c *gin.Context) {
	reportID, ok := parseEnterpriseUint64Param(c, "reportId")
	if !ok {
		return
	}
	var req entSvc.ExportRecordInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateExportRecord(c.Request.Context(), c.GetUint64(middleware.UserIDKey), reportID, req)
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *EnterpriseHandler) ListDeliveryArtifacts(c *gin.Context) {
	items, err := h.svc.ListDeliveryArtifacts(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *EnterpriseHandler) ListDeliveryBundles(c *gin.Context) {
	items, err := h.svc.ListDeliveryBundles(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *EnterpriseHandler) ListDeliveryChecklist(c *gin.Context) {
	bundleID, ok := parseEnterpriseUint64Param(c, "bundleId")
	if !ok {
		return
	}
	items, err := h.svc.ListDeliveryChecklist(c.Request.Context(), c.GetUint64(middleware.UserIDKey), bundleID)
	if err != nil {
		writeEnterpriseError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func parseEnterpriseUint64Param(c *gin.Context, name string) (uint64, bool) {
	id, err := strconv.ParseUint(c.Param(name), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return id, true
}

func writeEnterpriseError(c *gin.Context, err error) {
	switch err {
	case entSvc.ErrEnterpriseScopeDenied:
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case entSvc.ErrEnterpriseInvalid:
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}
