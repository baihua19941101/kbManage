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
			"auditor": {
				"audit:read": {},
			},
			"operator": {
				"resource:read":     {},
				"operation:execute": {},
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
