package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	marketplaceSvc "kbmanage/backend/internal/service/marketplace"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type MarketplaceHandler struct {
	svc *marketplaceSvc.Service
}

func NewMarketplaceHandler(svc *marketplaceSvc.Service) *MarketplaceHandler {
	return &MarketplaceHandler{svc: svc}
}

func (h *MarketplaceHandler) ListCatalogSources(c *gin.Context) {
	items, err := h.svc.ListCatalogSources(c.Request.Context(), c.GetUint64(middleware.UserIDKey), marketplaceSvc.CatalogSourceListFilter{
		SourceType: strings.TrimSpace(c.Query("sourceType")),
		Status:     strings.TrimSpace(c.Query("status")),
		Keyword:    strings.TrimSpace(c.Query("keyword")),
	})
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *MarketplaceHandler) CreateCatalogSource(c *gin.Context) {
	var req marketplaceSvc.CreateCatalogSourceInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateCatalogSource(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *MarketplaceHandler) SyncCatalogSource(c *gin.Context) {
	sourceID, ok := parseMarketplaceUint64(c, "sourceId")
	if !ok {
		return
	}
	item, err := h.svc.SyncCatalogSource(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sourceID)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, item)
}

func (h *MarketplaceHandler) ListTemplates(c *gin.Context) {
	sourceID, _ := parseOptionalMarketplaceUint64(firstNonEmptyMarketplaceQuery(c, "catalogSourceId", "sourceId"))
	items, err := h.svc.ListTemplates(c.Request.Context(), c.GetUint64(middleware.UserIDKey), marketplaceSvc.TemplateListFilter{
		CatalogSourceID: sourceID,
		Category:        strings.TrimSpace(c.Query("category")),
		Status:          strings.TrimSpace(c.Query("status")),
		Keyword:         strings.TrimSpace(c.Query("keyword")),
	})
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *MarketplaceHandler) GetTemplateDetail(c *gin.Context) {
	templateID, ok := parseMarketplaceUint64(c, "templateId")
	if !ok {
		return
	}
	item, err := h.svc.GetTemplateDetail(c.Request.Context(), c.GetUint64(middleware.UserIDKey), templateID)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *MarketplaceHandler) ListTemplateReleases(c *gin.Context) {
	templateID, ok := parseMarketplaceUint64(c, "templateId")
	if !ok {
		return
	}
	items, err := h.svc.ListTemplateReleases(c.Request.Context(), c.GetUint64(middleware.UserIDKey), templateID)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *MarketplaceHandler) CreateTemplateRelease(c *gin.Context) {
	templateID, ok := parseMarketplaceUint64(c, "templateId")
	if !ok {
		return
	}
	var req marketplaceSvc.CreateTemplateReleaseInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateTemplateRelease(c.Request.Context(), c.GetUint64(middleware.UserIDKey), templateID, req)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *MarketplaceHandler) ListInstallations(c *gin.Context) {
	scopeID, _ := parseOptionalMarketplaceUint64(c.Query("scopeId"))
	items, err := h.svc.ListInstallationRecords(c.Request.Context(), c.GetUint64(middleware.UserIDKey), marketplaceSvc.InstallationListFilter{
		ScopeType: strings.TrimSpace(c.Query("scopeType")),
		ScopeID:   scopeID,
		Status:    strings.TrimSpace(c.Query("status")),
	})
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *MarketplaceHandler) ListExtensions(c *gin.Context) {
	items, err := h.svc.ListExtensions(c.Request.Context(), c.GetUint64(middleware.UserIDKey), marketplaceSvc.ExtensionListFilter{
		Type:    strings.TrimSpace(c.Query("extensionType")),
		Status:  strings.TrimSpace(c.Query("status")),
		Keyword: strings.TrimSpace(c.Query("keyword")),
	})
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *MarketplaceHandler) RegisterExtension(c *gin.Context) {
	var req marketplaceSvc.CreateExtensionPackageInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.RegisterExtension(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *MarketplaceHandler) EnableExtension(c *gin.Context) {
	extensionID, ok := parseMarketplaceUint64(c, "extensionId")
	if !ok {
		return
	}
	var req marketplaceSvc.ExtensionLifecycleInput
	_ = c.ShouldBindJSON(&req)
	item, err := h.svc.EnableExtension(c.Request.Context(), c.GetUint64(middleware.UserIDKey), extensionID, req)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, item)
}

func (h *MarketplaceHandler) DisableExtension(c *gin.Context) {
	extensionID, ok := parseMarketplaceUint64(c, "extensionId")
	if !ok {
		return
	}
	var req marketplaceSvc.ExtensionLifecycleInput
	_ = c.ShouldBindJSON(&req)
	item, err := h.svc.DisableExtension(c.Request.Context(), c.GetUint64(middleware.UserIDKey), extensionID, req)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusAccepted, item)
}

func (h *MarketplaceHandler) GetExtensionCompatibility(c *gin.Context) {
	extensionID, ok := parseMarketplaceUint64(c, "extensionId")
	if !ok {
		return
	}
	item, err := h.svc.GetExtensionCompatibility(c.Request.Context(), c.GetUint64(middleware.UserIDKey), extensionID)
	if err != nil {
		writeMarketplaceError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func writeMarketplaceError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, marketplaceSvc.ErrMarketplaceScopeDenied):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, marketplaceSvc.ErrMarketplaceConflict), errors.Is(err, marketplaceSvc.ErrMarketplaceBlocked):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, marketplaceSvc.ErrMarketplaceInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func parseMarketplaceUint64(c *gin.Context, name string) (uint64, bool) {
	value := strings.TrimSpace(c.Param(name))
	id, err := strconv.ParseUint(value, 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + name})
		return 0, false
	}
	return id, true
}

func parseOptionalMarketplaceUint64(v string) (uint64, bool) {
	if strings.TrimSpace(v) == "" {
		return 0, false
	}
	out, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, false
	}
	return out, true
}

func firstNonEmptyMarketplaceQuery(c *gin.Context, keys ...string) string {
	for _, key := range keys {
		if value := strings.TrimSpace(c.Query(key)); value != "" {
			return value
		}
	}
	return ""
}
