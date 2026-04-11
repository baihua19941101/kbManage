package middleware

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/service/auth"

	"github.com/gin-gonic/gin"
)

const (
	PermissionWorkspaceRead = "access:workspace:read"
	PermissionProjectRead   = "access:project:read"
	PermissionProjectWrite  = "access:project:write"
	PermissionBindingRead   = "access:binding:read"
	PermissionBindingWrite  = "access:binding:write"
)

func RequireWorkspaceScope(scopeAccess *auth.ScopeAccessService, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if scopeAccess == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "scope authorization is not configured"})
			return
		}

		userID := c.GetUint64(UserIDKey)
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
			return
		}

		workspaceID, err := strconv.ParseUint(c.Param("workspaceId"), 10, 64)
		if err != nil || workspaceID == 0 {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid workspace id"})
			return
		}

		allowed, err := scopeAccess.HasScopePermission(c.Request.Context(), userID, domain.ScopeTypeWorkspace, workspaceID, 0, permission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

func RequireRoleBindingScope(scopeAccess *auth.ScopeAccessService, permission string) gin.HandlerFunc {
	return func(c *gin.Context) {
		if scopeAccess == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "scope authorization is not configured"})
			return
		}

		userID := c.GetUint64(UserIDKey)
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
			return
		}

		targetType, workspaceID, projectID, err := parseScopeTarget(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		allowed, err := scopeAccess.HasScopePermission(c.Request.Context(), userID, targetType, workspaceID, projectID, permission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if !allowed {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "forbidden"})
			return
		}

		c.Next()
	}
}

func parseScopeTarget(c *gin.Context) (domain.ScopeType, uint64, uint64, error) {
	scopeType := strings.TrimSpace(c.Query("scopeType"))
	var scopeID uint64

	if c.Request.Method == http.MethodPost {
		body, err := c.GetRawData()
		if err != nil {
			return "", 0, 0, errors.New("invalid request body")
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

		var payload struct {
			ScopeType string `json:"scopeType"`
			ScopeID   any    `json:"scopeId"`
		}
		if err := json.Unmarshal(body, &payload); err != nil {
			return "", 0, 0, errors.New("invalid request body")
		}

		if scopeType == "" {
			scopeType = strings.TrimSpace(payload.ScopeType)
		}
		scopeID, err = parseUint64Any(payload.ScopeID)
		if err != nil {
			return "", 0, 0, errors.New("invalid scopeId")
		}
	} else {
		scopeIDText := strings.TrimSpace(c.Query("scopeId"))
		if scopeType == "" {
			return "", 0, 0, errors.New("scopeType is required")
		}
		if scopeIDText == "" {
			return "", 0, 0, errors.New("scopeId is required")
		}
		parsed, err := strconv.ParseUint(scopeIDText, 10, 64)
		if err != nil || parsed == 0 {
			return "", 0, 0, errors.New("invalid scopeId")
		}
		scopeID = parsed
	}

	if scopeID == 0 {
		return "", 0, 0, errors.New("invalid scopeId")
	}

	switch strings.ToLower(scopeType) {
	case string(domain.ScopeTypeWorkspace):
		return domain.ScopeTypeWorkspace, scopeID, 0, nil
	case string(domain.ScopeTypeProject):
		return domain.ScopeTypeProject, 0, scopeID, nil
	default:
		return "", 0, 0, errors.New("scopeType must be workspace or project")
	}
}

func parseUint64Any(v any) (uint64, error) {
	switch value := v.(type) {
	case nil:
		return 0, nil
	case uint64:
		return value, nil
	case uint:
		return uint64(value), nil
	case int:
		if value < 0 {
			return 0, errors.New("negative")
		}
		return uint64(value), nil
	case int64:
		if value < 0 {
			return 0, errors.New("negative")
		}
		return uint64(value), nil
	case float64:
		if value < 0 {
			return 0, errors.New("negative")
		}
		if value != float64(uint64(value)) {
			return 0, errors.New("non-integer")
		}
		return uint64(value), nil
	case string:
		trimmed := strings.TrimSpace(value)
		if trimmed == "" {
			return 0, nil
		}
		return strconv.ParseUint(trimmed, 10, 64)
	default:
		return 0, errors.New("unsupported")
	}
}
