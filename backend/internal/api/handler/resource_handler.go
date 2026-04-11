package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"
	"kbmanage/backend/internal/service/cluster"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ResourceHandler struct {
	svc         *cluster.Service
	scopeAccess *authSvc.ScopeAccessService
}

func NewResourceHandler(svc *cluster.Service, scopeAccess *authSvc.ScopeAccessService) *ResourceHandler {
	return &ResourceHandler{svc: svc, scopeAccess: scopeAccess}
}

func (h *ResourceHandler) List(c *gin.Context) {
	var filter repository.ResourceListFilter

	if clusterIDText := c.Param("id"); clusterIDText != "" {
		clusterID, err := strconv.ParseUint(clusterIDText, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
			return
		}
		filter.ClusterID = clusterID
	} else if clusterIDText := c.Query("clusterId"); clusterIDText != "" {
		clusterID, err := strconv.ParseUint(clusterIDText, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clusterId"})
			return
		}
		filter.ClusterID = clusterID
	}
	filter.Namespace = c.Query("namespace")
	filter.Kind = c.Query("kind")
	filter.Health = c.Query("health")
	filter.Keyword = c.Query("keyword")

	if limitText := c.Query("limit"); limitText != "" {
		limit, err := strconv.Atoi(limitText)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		filter.Limit = limit
	}
	if offsetText := c.Query("offset"); offsetText != "" {
		offset, err := strconv.Atoi(offsetText)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
		filter.Offset = offset
	}

	if h.scopeAccess != nil {
		userID := c.GetUint64(middleware.UserIDKey)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
			return
		}
		allowedClusterIDs, constrained, err := h.scopeAccess.ListClusterIDsByPermission(c.Request.Context(), userID, middleware.PermissionProjectRead)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if constrained {
			if filter.ClusterID > 0 {
				if !containsUint64(allowedClusterIDs, filter.ClusterID) {
					c.JSON(http.StatusOK, gin.H{
						"items":  []repository.ResourceInventory{},
						"filter": filter,
					})
					return
				}
			} else {
				filter.ClusterIDs = allowedClusterIDs
			}
		}
	}

	items, err := h.svc.ListResources(c.Request.Context(), filter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"items":  items,
		"filter": filter,
	})
}

func (h *ResourceHandler) Detail(c *gin.Context) {
	clusterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil || clusterID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
		return
	}

	filter := repository.ResourceDetailFilter{
		ClusterID: clusterID,
		Namespace: c.Query("namespace"),
		Kind:      c.Query("kind"),
		Name:      c.Query("name"),
	}

	if h.scopeAccess != nil {
		userID := c.GetUint64(middleware.UserIDKey)
		if userID == 0 {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
			return
		}
		allowed, err := h.scopeAccess.CanAccessClusterByPermission(c.Request.Context(), userID, clusterID, middleware.PermissionProjectRead)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !allowed {
			c.JSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}
	}

	item, err := h.svc.GetResourceDetail(c.Request.Context(), filter)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
			return
		}
		if strings.Contains(strings.ToLower(err.Error()), "required") {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, item)
}

func containsUint64(values []uint64, target uint64) bool {
	for _, item := range values {
		if item == target {
			return true
		}
	}
	return false
}
