package handler

import (
	"net/http"
	"strconv"

	"kbmanage/backend/internal/service/cluster"

	"github.com/gin-gonic/gin"
)

type ClusterHandler struct {
	svc *cluster.Service
}

func NewClusterHandler(svc *cluster.Service) *ClusterHandler {
	return &ClusterHandler{svc: svc}
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

func (h *ClusterHandler) ValidateConnectivity(c *gin.Context) {
	clusterID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cluster id"})
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

	summary, err := h.svc.GetHealthSummary(c.Request.Context(), clusterID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}
