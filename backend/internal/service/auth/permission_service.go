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
				"resource:read":                 {},
				"operation:execute":             {},
				"observability:read":            {},
				"observability:write":           {},
				"workloadops:read":              {},
				"workloadops:execute":           {},
				"workloadops:terminal":          {},
				"workloadops:rollback":          {},
				"workloadops:batch":             {},
				"gitops:read":                   {},
				"gitops:manage-source":          {},
				"gitops:sync":                   {},
				"gitops:promote":                {},
				"gitops:rollback":               {},
				"gitops:override":               {},
				"securitypolicy:read":           {},
				"securitypolicy:manage":         {},
				"securitypolicy:enforce":        {},
				"compliance:read":               {},
				"compliance:manage-baseline":    {},
				"compliance:execute-scan":       {},
				"compliance:manage-remediation": {},
				"compliance:review-exception":   {},
				"compliance:export-archive":     {},
			},
			"auditor": {
				"audit:read":                {},
				"gitops:read":               {},
				"securitypolicy:read":       {},
				"compliance:read":           {},
				"compliance:export-archive": {},
			},
			"audit-reader": {
				"audit:read":                {},
				"observability:read":        {},
				"workloadops:read":          {},
				"gitops:read":               {},
				"securitypolicy:read":       {},
				"compliance:read":           {},
				"compliance:export-archive": {},
			},
			"operator": {
				"resource:read":                 {},
				"operation:execute":             {},
				"observability:read":            {},
				"observability:write":           {},
				"workloadops:read":              {},
				"workloadops:execute":           {},
				"workloadops:terminal":          {},
				"workloadops:rollback":          {},
				"workloadops:batch":             {},
				"gitops:read":                   {},
				"gitops:manage-source":          {},
				"gitops:sync":                   {},
				"gitops:promote":                {},
				"gitops:rollback":               {},
				"gitops:override":               {},
				"securitypolicy:read":           {},
				"securitypolicy:manage":         {},
				"securitypolicy:enforce":        {},
				"compliance:read":               {},
				"compliance:manage-baseline":    {},
				"compliance:execute-scan":       {},
				"compliance:manage-remediation": {},
				"compliance:review-exception":   {},
				"compliance:export-archive":     {},
			},
			"readonly": {
				"resource:read":             {},
				"audit:read":                {},
				"observability:read":        {},
				"workloadops:read":          {},
				"gitops:read":               {},
				"securitypolicy:read":       {},
				"compliance:read":           {},
				"compliance:export-archive": {},
			},
			"workspace-owner": {
				"access:workspace:read":         {},
				"access:workspace:write":        {},
				"access:project:read":           {},
				"access:project:write":          {},
				"access:binding:read":           {},
				"access:binding:write":          {},
				"observability:read":            {},
				"observability:write":           {},
				"workloadops:read":              {},
				"workloadops:execute":           {},
				"workloadops:terminal":          {},
				"workloadops:rollback":          {},
				"workloadops:batch":             {},
				"gitops:read":                   {},
				"gitops:manage-source":          {},
				"gitops:sync":                   {},
				"gitops:promote":                {},
				"gitops:rollback":               {},
				"gitops:override":               {},
				"securitypolicy:read":           {},
				"securitypolicy:manage":         {},
				"securitypolicy:enforce":        {},
				"compliance:read":               {},
				"compliance:manage-baseline":    {},
				"compliance:execute-scan":       {},
				"compliance:manage-remediation": {},
				"compliance:review-exception":   {},
				"compliance:export-archive":     {},
			},
			"workspace-viewer": {
				"access:workspace:read": {},
				"access:project:read":   {},
				"access:binding:read":   {},
				"observability:read":    {},
				"workloadops:read":      {},
				"gitops:read":           {},
				"securitypolicy:read":   {},
			},
			"project-owner": {
				"access:project:read":           {},
				"access:project:write":          {},
				"access:binding:read":           {},
				"access:binding:write":          {},
				"observability:read":            {},
				"observability:write":           {},
				"workloadops:read":              {},
				"workloadops:execute":           {},
				"workloadops:terminal":          {},
				"workloadops:rollback":          {},
				"workloadops:batch":             {},
				"gitops:read":                   {},
				"gitops:manage-source":          {},
				"gitops:sync":                   {},
				"gitops:promote":                {},
				"gitops:rollback":               {},
				"gitops:override":               {},
				"securitypolicy:read":           {},
				"securitypolicy:manage":         {},
				"securitypolicy:enforce":        {},
				"compliance:read":               {},
				"compliance:manage-baseline":    {},
				"compliance:execute-scan":       {},
				"compliance:manage-remediation": {},
				"compliance:review-exception":   {},
				"compliance:export-archive":     {},
			},
			"project-viewer": {
				"access:project:read":       {},
				"access:binding:read":       {},
				"observability:read":        {},
				"workloadops:read":          {},
				"gitops:read":               {},
				"securitypolicy:read":       {},
				"compliance:read":           {},
				"compliance:export-archive": {},
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
