package handler

import (
	"errors"
	"net/http"
	"strings"

	workspaceSvc "kbmanage/backend/internal/service/workspace"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkspaceHandler struct {
	svc *workspaceSvc.Service
}

func NewWorkspaceHandler(svc *workspaceSvc.Service) *WorkspaceHandler {
	return &WorkspaceHandler{svc: svc}
}

func (h *WorkspaceHandler) Create(c *gin.Context) {
	var req workspaceSvc.CreateWorkspaceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	item, err := h.svc.Create(c.Request.Context(), req)
	if err != nil {
		writeAccessError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *WorkspaceHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		writeAccessError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func writeAccessError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case errors.Is(err, gorm.ErrDuplicatedKey):
		status = http.StatusConflict
	case errors.Is(err, gorm.ErrInvalidDB):
		status = http.StatusNotImplemented
	case strings.Contains(strings.ToLower(err.Error()), "required"):
		status = http.StatusBadRequest
	}
	c.JSON(status, gin.H{"error": err.Error()})
}
