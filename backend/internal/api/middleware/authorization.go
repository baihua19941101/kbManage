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
	auditSvc "kbmanage/backend/internal/service/audit"
	"kbmanage/backend/internal/service/auth"
	obsSvc "kbmanage/backend/internal/service/observability"

	"github.com/gin-gonic/gin"
)

const (
	PermissionWorkspaceRead      = "access:workspace:read"
	PermissionProjectRead        = "access:project:read"
	PermissionProjectWrite       = "access:project:write"
	PermissionBindingRead        = "access:binding:read"
	PermissionBindingWrite       = "access:binding:write"
	PermissionObservabilityRead  = "observability:read"
	PermissionObservabilityWrite = "observability:write"
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

func RequireObservabilityScope(
	scopeAccess *auth.ScopeAccessService,
	scopeService *obsSvc.ScopeService,
	permission string,
	auditWriter ...*auditSvc.EventWriter,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		if scopeAccess == nil || scopeService == nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "observability authorization is not configured"})
			return
		}

		userID := c.GetUint64(UserIDKey)
		if userID == 0 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing authenticated user"})
			return
		}

		scopeFilter, err := parseObservabilityScopeFilter(c)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		workspaceIDs, err := scopeAccess.ListWorkspaceIDsByPermission(c.Request.Context(), userID, permission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if len(workspaceIDs) == 0 {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "observability scope access denied"})
			return
		}

		for _, workspaceID := range scopeFilter.WorkspaceIDs {
			allowed, err := scopeAccess.HasScopePermission(
				c.Request.Context(),
				userID,
				domain.ScopeTypeWorkspace,
				workspaceID,
				0,
				permission,
			)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "observability scope access denied"})
				return
			}
		}
		for _, projectID := range scopeFilter.ProjectIDs {
			allowed, err := scopeAccess.HasScopePermission(
				c.Request.Context(),
				userID,
				domain.ScopeTypeProject,
				0,
				projectID,
				permission,
			)
			if err != nil {
				c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			if !allowed {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "observability scope access denied"})
				return
			}
		}

		clusterIDs, constrained, err := scopeAccess.ListClusterIDsByPermission(c.Request.Context(), userID, permission)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		if constrained {
			allowedSet := make(map[uint64]struct{}, len(clusterIDs))
			for _, clusterID := range clusterIDs {
				allowedSet[clusterID] = struct{}{}
			}
			for _, clusterID := range scopeFilter.ClusterIDs {
				if _, ok := allowedSet[clusterID]; !ok {
					c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "observability scope access denied"})
					return
				}
			}
		}

		access := obsSvc.AccessContext{
			UserID:             userID,
			Permission:         permission,
			ClusterConstrained: constrained,
			WorkspaceIDs:       workspaceIDs,
			ProjectIDs:         scopeFilter.ProjectIDs,
			ClusterIDs:         clusterIDs,
		}
		if len(scopeFilter.ClusterIDs) == 0 && len(scopeFilter.WorkspaceIDs) == 0 && len(scopeFilter.ProjectIDs) == 0 {
			scopeFilter.WorkspaceIDs = append(scopeFilter.WorkspaceIDs, workspaceIDs...)
		}
		ctxWithAccess := obsSvc.WithAccessContext(c.Request.Context(), access)

		authorizedFilter, err := scopeService.FilterByScope(ctxWithAccess, userID, scopeFilter)
		if err != nil {
			if errors.Is(err, obsSvc.ErrObservabilityScopeDenied) || errors.Is(err, obsSvc.ErrInvalidObservabilityUser) {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "observability scope access denied"})
				return
			}
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		access.WorkspaceIDs = authorizedFilter.WorkspaceIDs
		access.ProjectIDs = authorizedFilter.ProjectIDs
		if constrained {
			access.ClusterIDs = authorizedFilter.ClusterIDs
		}
		c.Request = c.Request.WithContext(obsSvc.WithAccessContext(c.Request.Context(), access))

		if len(auditWriter) > 0 && auditWriter[0] != nil && permission == PermissionObservabilityRead {
			actorID := userID
			_ = auditWriter[0].WriteObservabilityEvent(
				c.Request.Context(),
				c.GetString(RequestIDKey),
				&actorID,
				auditSvc.ObservabilityAuditActionAccessRead,
				buildObservabilityResourceReference(authorizedFilter),
				domain.AuditOutcomeSuccess,
				map[string]any{
					"method":       c.Request.Method,
					"path":         c.FullPath(),
					"operation":    "read",
					"clusterIds":   authorizedFilter.ClusterIDs,
					"workspaceIds": authorizedFilter.WorkspaceIDs,
					"projectIds":   authorizedFilter.ProjectIDs,
				},
			)
		}
		c.Next()
	}
}

func buildObservabilityResourceReference(filter obsSvc.ScopeFilter) string {
	if len(filter.ClusterIDs) > 0 {
		return "cluster:" + strconv.FormatUint(filter.ClusterIDs[0], 10)
	}
	if len(filter.ProjectIDs) > 0 {
		return "project:" + strconv.FormatUint(filter.ProjectIDs[0], 10)
	}
	if len(filter.WorkspaceIDs) > 0 {
		return "workspace:" + strconv.FormatUint(filter.WorkspaceIDs[0], 10)
	}
	return "scope:unknown"
}

func parseObservabilityScopeFilter(c *gin.Context) (obsSvc.ScopeFilter, error) {
	filter := obsSvc.ScopeFilter{
		ClusterIDs:   parseUint64ListFromQuery(c, "clusterIds", "clusterId"),
		WorkspaceIDs: parseUint64ListFromQuery(c, "workspaceIds", "workspaceId"),
		ProjectIDs:   parseUint64ListFromQuery(c, "projectIds", "projectId"),
	}

	if c.Request.Method == http.MethodGet || c.Request.Method == http.MethodDelete {
		return filter, nil
	}

	body, err := c.GetRawData()
	if err != nil {
		return obsSvc.ScopeFilter{}, errors.New("invalid request body")
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
	if len(bytes.TrimSpace(body)) == 0 {
		return filter, nil
	}

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return obsSvc.ScopeFilter{}, errors.New("invalid request body")
	}

	filter.ClusterIDs = mergeUniqueIDs(filter.ClusterIDs, parseUint64ListFromMap(payload, "clusterIds", "clusterId"))
	filter.WorkspaceIDs = mergeUniqueIDs(filter.WorkspaceIDs, parseUint64ListFromMap(payload, "workspaceIds", "workspaceId"))
	filter.ProjectIDs = mergeUniqueIDs(filter.ProjectIDs, parseUint64ListFromMap(payload, "projectIds", "projectId"))

	if raw, ok := payload["scopeSnapshot"]; ok {
		snapshot := parseScopeSnapshot(raw)
		filter.ClusterIDs = mergeUniqueIDs(filter.ClusterIDs, snapshot.ClusterIDs)
		filter.WorkspaceIDs = mergeUniqueIDs(filter.WorkspaceIDs, snapshot.WorkspaceIDs)
		filter.ProjectIDs = mergeUniqueIDs(filter.ProjectIDs, snapshot.ProjectIDs)
	}
	if raw, ok := payload["scope"]; ok {
		snapshot := parseScopeSnapshot(raw)
		filter.ClusterIDs = mergeUniqueIDs(filter.ClusterIDs, snapshot.ClusterIDs)
		filter.WorkspaceIDs = mergeUniqueIDs(filter.WorkspaceIDs, snapshot.WorkspaceIDs)
		filter.ProjectIDs = mergeUniqueIDs(filter.ProjectIDs, snapshot.ProjectIDs)
	}

	return filter, nil
}

func parseScopeSnapshot(raw any) obsSvc.ScopeFilter {
	switch value := raw.(type) {
	case map[string]any:
		return obsSvc.ScopeFilter{
			ClusterIDs:   parseUint64ListFromMap(value, "clusterIds", "clusterId"),
			WorkspaceIDs: parseUint64ListFromMap(value, "workspaceIds", "workspaceId"),
			ProjectIDs:   parseUint64ListFromMap(value, "projectIds", "projectId"),
		}
	case string:
		text := strings.TrimSpace(value)
		if text == "" {
			return obsSvc.ScopeFilter{}
		}
		var payload map[string]any
		if err := json.Unmarshal([]byte(text), &payload); err != nil {
			return obsSvc.ScopeFilter{}
		}
		return obsSvc.ScopeFilter{
			ClusterIDs:   parseUint64ListFromMap(payload, "clusterIds", "clusterId"),
			WorkspaceIDs: parseUint64ListFromMap(payload, "workspaceIds", "workspaceId"),
			ProjectIDs:   parseUint64ListFromMap(payload, "projectIds", "projectId"),
		}
	default:
		return obsSvc.ScopeFilter{}
	}
}

func parseUint64ListFromQuery(c *gin.Context, listKey, singleKey string) []uint64 {
	items := make([]uint64, 0)
	for _, raw := range c.QueryArray(listKey) {
		items = append(items, parseCSVUint64(raw)...)
	}
	items = append(items, parseCSVUint64(c.Query(listKey))...)
	items = append(items, parseCSVUint64(c.Query(singleKey))...)
	return mergeUniqueIDs(nil, items)
}

func parseUint64ListFromMap(data map[string]any, listKey, singleKey string) []uint64 {
	items := make([]uint64, 0)
	if v, ok := findValueByKey(data, listKey); ok {
		items = append(items, parseUint64ListFromAny(v)...)
	}
	if v, ok := findValueByKey(data, singleKey); ok {
		items = append(items, parseUint64ListFromAny(v)...)
	}
	return mergeUniqueIDs(nil, items)
}

func parseUint64ListFromAny(v any) []uint64 {
	out := make([]uint64, 0)
	switch value := v.(type) {
	case []any:
		for _, item := range value {
			if id, err := parseUint64Any(item); err == nil && id != 0 {
				out = append(out, id)
			}
		}
	case string:
		out = append(out, parseCSVUint64(value)...)
	default:
		if id, err := parseUint64Any(value); err == nil && id != 0 {
			out = append(out, id)
		}
	}
	return mergeUniqueIDs(nil, out)
}

func parseCSVUint64(raw string) []uint64 {
	out := make([]uint64, 0)
	for _, token := range strings.Split(strings.TrimSpace(raw), ",") {
		text := strings.TrimSpace(token)
		if text == "" {
			continue
		}
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil || id == 0 {
			continue
		}
		out = append(out, id)
	}
	return out
}

func mergeUniqueIDs(left []uint64, right []uint64) []uint64 {
	set := make(map[uint64]struct{}, len(left)+len(right))
	for _, id := range left {
		if id == 0 {
			continue
		}
		set[id] = struct{}{}
	}
	for _, id := range right {
		if id == 0 {
			continue
		}
		set[id] = struct{}{}
	}
	out := make([]uint64, 0, len(set))
	for id := range set {
		out = append(out, id)
	}
	return out
}

func findValueByKey(data map[string]any, key string) (any, bool) {
	for existing, value := range data {
		if strings.EqualFold(strings.TrimSpace(existing), strings.TrimSpace(key)) {
			return value, true
		}
	}
	return nil, false
}
