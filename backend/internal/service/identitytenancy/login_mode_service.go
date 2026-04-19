package identitytenancy

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListAvailableLoginModes(ctx context.Context, userID uint64) ([]map[string]any, error) {
	items, err := s.ListIdentitySources(ctx, userID, IdentitySourceListFilter{})
	if err != nil {
		return nil, err
	}
	out := make([]map[string]any, 0, len(items))
	for _, item := range items {
		out = append(out, map[string]any{
			"sourceId":    item.ID,
			"sourceType":  item.SourceType,
			"status":      item.Status,
			"loginMode":   item.LoginMode,
			"isFallback":  item.SourceType == domain.IdentitySourceTypeLocal || item.LoginMode == domain.IdentityLoginModeFallback,
			"userVisible": item.Status == domain.IdentitySourceStatusActive,
		})
	}
	return out, nil
}

func (s *Service) UpdatePreferredLoginMode(ctx context.Context, userID uint64, loginMode string) (string, error) {
	if err := s.scope.EnsureManageSource(ctx, userID); err != nil {
		return "", err
	}
	normalized := strings.ToLower(strings.TrimSpace(loginMode))
	var target domain.IdentityLoginMode
	switch normalized {
	case "local":
		target = domain.IdentityLoginModeFallback
	case "external":
		target = domain.IdentityLoginModeExclusive
	case "mixed":
		target = domain.IdentityLoginModeOptional
	default:
		return "", ErrIdentityTenancyInvalid
	}
	items, err := s.sources.List(ctx, repository.IdentitySourceListFilter{})
	if err != nil {
		return "", err
	}
	for i := range items {
		switch target {
		case domain.IdentityLoginModeFallback:
			if items[i].SourceType == domain.IdentitySourceTypeLocal {
				items[i].LoginMode = target
			} else {
				items[i].LoginMode = domain.IdentityLoginModeOptional
			}
		case domain.IdentityLoginModeExclusive:
			if items[i].SourceType == domain.IdentitySourceTypeLocal {
				items[i].LoginMode = domain.IdentityLoginModeFallback
			} else {
				items[i].LoginMode = target
			}
		default:
			items[i].LoginMode = domain.IdentityLoginModeOptional
		}
		if err := s.sources.Update(ctx, &items[i]); err != nil {
			return "", err
		}
	}
	s.writeAudit(ctx, userID, ActionLoginModeUpdate, ResourceTypeIdentitySource, "preferred", domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"loginMode": normalized,
		"count":     len(items),
	})
	return normalized, nil
}
