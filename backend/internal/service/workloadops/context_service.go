package workloadops

import "context"

func (s *Service) buildContext(_ context.Context, target WorkloadReference) map[string]any {
	return map[string]any{
		"clusterId":           target.ClusterID,
		"workspaceId":         valueOrNil(target.WorkspaceID),
		"projectId":           valueOrNil(target.ProjectID),
		"namespace":           target.Namespace,
		"resourceKind":        target.ResourceKind,
		"resourceName":        target.ResourceName,
		"healthStatus":        "unknown",
		"rolloutStatus":       "unknown",
		"latestChangeSummary": "暂无最近变更记录",
		"latestActionSummary": "暂无最近动作记录",
		"availableActions":    []string{"scale", "restart", "redeploy", "replace-instance", "rollback"},
	}
}
