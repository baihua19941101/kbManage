package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	authSvc "kbmanage/backend/internal/service/auth"
	"kbmanage/backend/internal/service/cluster"

	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	svc         *cluster.Service
	scopeAccess *authSvc.ScopeAccessService
}

var errMissingAuthenticatedUser = errors.New("missing authenticated user")

func NewClusterHandler(svc *cluster.Service, scopeAccess *authSvc.ScopeAccessService) *ClusterHandler {
	return &ClusterHandler{svc: svc, scopeAccess: scopeAccess}
}

func (h *ClusterHandler) Register(c *gin.Context) {
	var req cluster.RegisterClusterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.RegisterCluster(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ClusterHandler) List(c *gin.Context) {
	allowedClusterIDs, constrained, err := h.listVisibleClusterIDs(c)
	if err != nil {
		if errors.Is(err, errMissingAuthenticatedUser) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if constrained {
		clusterItems, err := h.svc.ListClustersByIDs(c.Request.Context(), allowedClusterIDs)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"items": clusterItems})
		return
	}

	items, err := h.svc.ListClusters(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *ClusterHandler) ValidateConnectivity(c *gin.Context) {
	clusterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
		return
	}
	if ok := h.ensureClusterVisible(c, clusterID); !ok {
		return
	}

	result, err := h.svc.ValidateConnectivity(c.Request.Context(), clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, result)
}

func (h *ClusterHandler) HealthSummary(c *gin.Context) {
	clusterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
		return
	}
	if ok := h.ensureClusterVisible(c, clusterID); !ok {
		return
	}

	summary, err := h.svc.GetHealthSummary(c.Request.Context(), clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

func (h *ClusterHandler) SyncResources(c *gin.Context) {
	clusterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || clusterID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
		return
	}
	if ok := h.ensureClusterVisible(c, clusterID); !ok {
		return
	}

	result, err := h.svc.TriggerResourceSync(c.Request.Context(), clusterID)
	if err != nil {
		status := http.StatusInternalServerError
		lower := strings.ToLower(err.Error())
		if strings.Contains(lower, "required") || strings.Contains(lower, "invalid") {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusAccepted, result)
}

func (h *ClusterHandler) ensureClusterVisible(c *gin.Context, clusterID uint64) bool {
	if h.scopeAccess == nil {
		return true
	}
	userID := c.GetUint64(middleware.UserIDKey)
	if userID == 0 {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
		return false
	}

	allowed, err := h.scopeAccess.CanAccessClusterByPermission(c.Request.Context(), userID, clusterID, middleware.PermissionProjectRead)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return false
	}
	if !allowed {
		c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
		return false
	}
	return true
}

func (h *ClusterHandler) listVisibleClusterIDs(c *gin.Context) ([]uint64, bool, error) {
	if h.scopeAccess == nil {
		return []uint64{}, false, nil
	}
	userID := c.GetUint64(middleware.UserIDKey)
	if userID == 0 {
		return nil, false, errMissingAuthenticatedUser
	}
	return h.scopeAccess.ListClusterIDsByPermission(c.Request.Context(), userID, middleware.PermissionProjectRead)
}
