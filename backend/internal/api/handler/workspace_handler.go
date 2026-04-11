package handler

import (
	"errors"
	"net/http"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"
	workspaceSvc "kbmanage/backend/internal/service/workspace"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type WorkspaceHandler struct {
	svc         *workspaceSvc.Service
	roleRepo    *repository.ScopeRoleRepository
	bindingRepo *repository.ScopeRoleBindingRepository
	scopeAccess *authSvc.ScopeAccessService
}

func NewWorkspaceHandler(
	svc *workspaceSvc.Service,
	roleRepo *repository.ScopeRoleRepository,
	bindingRepo *repository.ScopeRoleBindingRepository,
	scopeAccess *authSvc.ScopeAccessService,
) *WorkspaceHandler {
	return &WorkspaceHandler{
		svc:         svc,
		roleRepo:    roleRepo,
		bindingRepo: bindingRepo,
		scopeAccess: scopeAccess,
	}
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

	userID := c.GetUint64(middleware.UserIDKey)
	if userID != 0 && h.roleRepo != nil && h.bindingRepo != nil {
		role, err := h.roleRepo.GetByScopeAndRoleKey(c.Request.Context(), "workspace", "workspace-owner")
		if err != nil {
			writeAccessError(c, err)
			return
		}

		binding := &repository.ScopeRoleBinding{
			SubjectType: "user",
			SubjectID:   userID,
			ScopeType:   "workspace",
			ScopeID:     item.ID,
			ScopeRoleID: role.ID,
			GrantedBy:   userID,
		}
		if err := h.bindingRepo.Create(c.Request.Context(), binding); err != nil && !errors.Is(err, gorm.ErrDuplicatedKey) {
			writeAccessError(c, err)
			return
		}
	}

	c.JSON(http.StatusCreated, item)
}

func (h *WorkspaceHandler) List(c *gin.Context) {
	items, err := h.svc.List(c.Request.Context())
	if err != nil {
		writeAccessError(c, err)
		return
	}

	if h.scopeAccess != nil {
		userID := c.GetUint64(middleware.UserIDKey)
		allowedIDs, err := h.scopeAccess.ListWorkspaceIDsByPermission(c.Request.Context(), userID, middleware.PermissionWorkspaceRead)
		if err != nil {
			writeAccessError(c, err)
			return
		}

		allowedSet := make(map[uint64]struct{}, len(allowedIDs))
		for _, id := range allowedIDs {
			allowedSet[id] = struct{}{}
		}

		filtered := make([]domain.Workspace, 0, len(items))
		for _, item := range items {
			if _, ok := allowedSet[item.ID]; ok {
				filtered = append(filtered, item)
			}
		}
		c.JSON(http.StatusOK, gin.H{"items": filtered})
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
