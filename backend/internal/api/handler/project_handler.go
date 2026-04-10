package handler

import (
	"net/http"
	"strconv"

	projectSvc "kbmanage/backend/internal/service/project"

	"github.com/gin-gonic/gin"
)

type ProjectHandler struct {
	svc *projectSvc.Service
}

func NewProjectHandler(svc *projectSvc.Service) *ProjectHandler {
	return &ProjectHandler{svc: svc}
}

func (h *ProjectHandler) Create(c *gin.Context) {
	workspaceID, err := strconv.ParseUint(c.Param("workspaceId"), 10, 64)
	if err != nil || workspaceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}

	var req projectSvc.CreateProjectRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.Create(c.Request.Context(), workspaceID, req)
	if err != nil {
		writeAccessError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *ProjectHandler) ListByWorkspace(c *gin.Context) {
	workspaceID, err := strconv.ParseUint(c.Param("workspaceId"), 10, 64)
	if err != nil || workspaceID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
		return
	}

	items, err := h.svc.ListByWorkspace(c.Request.Context(), workspaceID)
	if err != nil {
		writeAccessError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}
