package securitypolicy

import (
	"context"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

// EffectivePolicyItem represents a policy in final hierarchy view.
type EffectivePolicyItem struct {
	PolicyID    uint64                       `json:"policyId"`
	PolicyName  string                       `json:"policyName"`
	ScopeLevel  domain.PolicyScopeLevel      `json:"scopeLevel"`
	Category    domain.PolicyCategory        `json:"category"`
	RiskLevel   domain.PolicyRiskLevel       `json:"riskLevel"`
	Enforcement domain.PolicyEnforcementMode `json:"enforcement"`
	SourceLevel string                       `json:"sourceLevel"`
}

// ResolvePolicyHierarchy returns policies in platform/workspace/project order,
// helping clients render final applicable policy set with source level.
func (s *Service) ResolvePolicyHierarchy(
	ctx context.Context,
	userID uint64,
	workspaceID uint64,
	projectID uint64,
) ([]EffectivePolicyItem, error) {
	if s == nil || s.policies == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	if err := s.validateScope(ctx, userID, workspaceID, projectID, PermissionSecurityPolicyRead); err != nil {
		return nil, err
	}

	build := func(scope domain.PolicyScopeLevel, wID *uint64, pID *uint64, source string) ([]EffectivePolicyItem, error) {
		items, err := s.policies.List(ctx, repository.SecurityPolicyListFilter{
			WorkspaceID: wID,
			ProjectID:   pID,
			ScopeLevel:  scope,
		})
		if err != nil {
			return nil, err
		}
		res := make([]EffectivePolicyItem, 0, len(items))
		for i := range items {
			res = append(res, EffectivePolicyItem{
				PolicyID:    items[i].ID,
				PolicyName:  items[i].Name,
				ScopeLevel:  items[i].ScopeLevel,
				Category:    items[i].Category,
				RiskLevel:   items[i].RiskLevel,
				Enforcement: items[i].DefaultEnforcementMode,
				SourceLevel: source,
			})
		}
		return res, nil
	}

	out := make([]EffectivePolicyItem, 0, 16)
	platformItems, err := build(domain.PolicyScopeLevelPlatform, nil, nil, "platform")
	if err != nil {
		return nil, err
	}
	out = append(out, platformItems...)

	if workspaceID > 0 {
		workspaceItems, err := build(domain.PolicyScopeLevelWorkspace, uint64PtrOrNil(workspaceID), nil, "workspace")
		if err != nil {
			return nil, err
		}
		out = append(out, workspaceItems...)
	}

	if projectID > 0 {
		projectItems, err := build(domain.PolicyScopeLevelProject, uint64PtrOrNil(workspaceID), uint64PtrOrNil(projectID), "project")
		if err != nil {
			return nil, err
		}
		out = append(out, projectItems...)
	}

	return out, nil
}
