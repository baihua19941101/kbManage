package handler

import (
	"net/http"
	"strconv"

	obsSvc "kbmanage/backend/internal/service/observability"

	"github.com/gin-gonic/gin"
)

type ObservabilityConfigHandler struct {
	service *obsSvc.Service
}

func NewObservabilityConfigHandler(service *obsSvc.Service) *ObservabilityConfigHandler {
	if service == nil {
		service = obsSvc.NewService(nil)
	}
	return &ObservabilityConfigHandler{service: service}
}

func (h *ObservabilityConfigHandler) GetClusterConfig(c *gin.Context) {
	clusterID, err := parseClusterIDParam(c)
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid clusterId")
		return
	}
	c.JSON(http.StatusOK, h.service.GetClusterConfig(clusterID))
}

func (h *ObservabilityConfigHandler) UpdateClusterConfig(c *gin.Context) {
	clusterID, err := parseClusterIDParam(c)
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", "invalid clusterId")
		return
	}

	var req obsSvc.UpdateObservabilityConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "invalid_parameter", err.Error())
		return
	}

	res, err := h.service.UpdateClusterConfig(clusterID, req)
	if err != nil {
		writeObservabilityError(c, http.StatusBadRequest, "update_config_failed", err.Error())
		return
	}
	c.JSON(http.StatusOK, res)
}

func parseClusterIDParam(c *gin.Context) (uint64, error) {
	id := c.Param("clusterId")
	if id == "" {
		id = c.Param("id")
	}
	return strconv.ParseUint(id, 10, 64)
}
