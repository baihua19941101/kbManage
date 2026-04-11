package handler

import (
	"errors"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RoleBindingHandler struct {
	roleRepo    *repository.ScopeRoleRepository
	bindingRepo *repository.ScopeRoleBindingRepository
}

func NewRoleBindingHandler(roleRepo *repository.ScopeRoleRepository, bindingRepo *repository.ScopeRoleBindingRepository) *RoleBindingHandler {
	return &RoleBindingHandler{roleRepo: roleRepo, bindingRepo: bindingRepo}
}

type createRoleBindingRequest struct {
	SubjectType string `json:"subjectType"`
	SubjectID   any    `json:"subjectId"`
	ScopeType   string `json:"scopeType"`
	ScopeID     any    `json:"scopeId"`
	RoleKey     string `json:"roleKey"`
	GrantedBy   any    `json:"grantedBy"`
}

func (h *RoleBindingHandler) Create(c *gin.Context) {
	var req createRoleBindingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subjectType := strings.TrimSpace(req.SubjectType)
	if subjectType != "user" && subjectType != "group" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "subjectType must be user or group"})
		return
	}
	scopeType := strings.TrimSpace(req.ScopeType)
	if scopeType != "workspace" && scopeType != "project" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "scopeType must be workspace or project"})
		return
	}
	roleKey := strings.TrimSpace(req.RoleKey)
	if roleKey == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "roleKey is required"})
		return
	}

	subjectID, err := toUint64(req.SubjectID)
	if err != nil || subjectID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subjectId"})
		return
	}
	scopeID, err := toUint64(req.ScopeID)
	if err != nil || scopeID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scopeId"})
		return
	}

	grantedBy, err := toUint64(req.GrantedBy)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid grantedBy"})
		return
	}
	if grantedBy == 0 {
		grantedBy = c.GetUint64(middleware.UserIDKey)
	}

	role, err := h.roleRepo.GetByScopeAndRoleKey(c.Request.Context(), scopeType, roleKey)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("roleKey %q is not defined for scopeType %q", roleKey, scopeType)})
			return
		}
		writeAccessError(c, err)
		return
	}

	item := &repository.ScopeRoleBinding{
		SubjectType: subjectType,
		SubjectID:   subjectID,
		ScopeType:   scopeType,
		ScopeID:     scopeID,
		ScopeRoleID: role.ID,
		GrantedBy:   grantedBy,
	}
	if err := h.bindingRepo.Create(c.Request.Context(), item); err != nil {
		writeAccessError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"id":          item.ID,
		"subjectType": item.SubjectType,
		"subjectId":   item.SubjectID,
		"scopeType":   item.ScopeType,
		"scopeId":     item.ScopeID,
		"scopeRoleId": item.ScopeRoleID,
		"roleKey":     role.RoleKey,
		"grantedBy":   item.GrantedBy,
		"createdAt":   item.CreatedAt,
	})
}

func (h *RoleBindingHandler) List(c *gin.Context) {
	filter := repository.ScopeRoleBindingFilter{
		SubjectType: strings.TrimSpace(c.Query("subjectType")),
		ScopeType:   strings.TrimSpace(c.Query("scopeType")),
	}

	if text := c.Query("subjectId"); text != "" {
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subjectId"})
			return
		}
		filter.SubjectID = id
	}
	if text := c.Query("scopeId"); text != "" {
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid scopeId"})
			return
		}
		filter.ScopeID = id
	}
	if text := c.Query("limit"); text != "" {
		v, err := strconv.Atoi(text)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		filter.Limit = v
	}
	if text := c.Query("offset"); text != "" {
		v, err := strconv.Atoi(text)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid offset"})
			return
		}
		filter.Offset = v
	}

	items, err := h.bindingRepo.List(c.Request.Context(), filter)
	if err != nil {
		writeAccessError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func toUint64(v any) (uint64, error) {
	switch val := v.(type) {
	case nil:
		return 0, nil
	case uint64:
		return val, nil
	case uint:
		return uint64(val), nil
	case int:
		if val < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint64(val), nil
	case int64:
		if val < 0 {
			return 0, fmt.Errorf("negative value")
		}
		return uint64(val), nil
	case float64:
		if val < 0 || val > math.MaxUint64 || math.Trunc(val) != val {
			return 0, fmt.Errorf("invalid number")
		}
		return uint64(val), nil
	case string:
		if strings.TrimSpace(val) == "" {
			return 0, nil
		}
		n, err := strconv.ParseUint(strings.TrimSpace(val), 10, 64)
		if err != nil {
			return 0, err
		}
		return n, nil
	default:
		return 0, fmt.Errorf("unsupported id type")
	}
}
