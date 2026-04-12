package auth

import "strings"

type PermissionService struct {
	rolePermissions map[string]map[string]struct{}
}

func NewPermissionService() *PermissionService {
	return &PermissionService{
		rolePermissions: map[string]map[string]struct{}{
			"platform-admin": {
				"*": {},
			},
			"ops-operator": {
				"resource:read":       {},
				"operation:execute":   {},
				"observability:read":  {},
				"observability:write": {},
			},
			"auditor": {
				"audit:read": {},
			},
			"audit-reader": {
				"audit:read":         {},
				"observability:read": {},
			},
			"operator": {
				"resource:read":       {},
				"operation:execute":   {},
				"observability:read":  {},
				"observability:write": {},
			},
			"readonly": {
				"resource:read":      {},
				"audit:read":         {},
				"observability:read": {},
			},
			"workspace-owner": {
				"access:workspace:read":  {},
				"access:workspace:write": {},
				"access:project:read":    {},
				"access:project:write":   {},
				"access:binding:read":    {},
				"access:binding:write":   {},
				"observability:read":     {},
				"observability:write":    {},
			},
			"workspace-viewer": {
				"access:workspace:read": {},
				"access:project:read":   {},
				"access:binding:read":   {},
				"observability:read":    {},
			},
			"project-owner": {
				"access:project:read":  {},
				"access:project:write": {},
				"access:binding:read":  {},
				"access:binding:write": {},
				"observability:read":   {},
				"observability:write":  {},
			},
			"project-viewer": {
				"access:project:read": {},
				"access:binding:read": {},
				"observability:read":  {},
			},
		},
	}
}

func (s *PermissionService) HasPermission(roleNames []string, permission string) bool {
	for _, role := range roleNames {
		grants, ok := s.rolePermissions[strings.ToLower(role)]
		if !ok {
			continue
		}
		if _, ok = grants["*"]; ok {
			return true
		}
		if _, ok = grants[permission]; ok {
			return true
		}
	}
	return false
}
