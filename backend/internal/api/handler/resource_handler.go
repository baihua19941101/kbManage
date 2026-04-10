package handler

import (
	"net/http"
	"strconv"

	"kbmanage/backend/internal/repository"
	"kbmanage/backend/internal/service/cluster"

	"github.com/gin-gonic/gin"
)

type ResourceHandler struct {
	svc *cluster.Service
}

func NewResourceHandler(svc *cluster.Service) *ResourceHandler {
	return &ResourceHandler{svc: svc}
}

func (h *ResourceHandler) List(c *gin.Context) {
	var filter repository.ResourceListFilter

	if clusterIDText := c.Query("clusterId"); clusterIDText != "" {
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
