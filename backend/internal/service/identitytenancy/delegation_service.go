package identitytenancy

import (
	"context"
	"strconv"
	"strings"

	"kbmanage/backend/internal/domain"
)

func (s *Service) ListDelegationGrants(ctx context.Context, userID uint64) ([]domain.DelegationGrant, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	items, err := s.delegations.List(ctx, "")
	if err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionDelegationGrantRead, ResourceTypeDelegationGrant, "list", domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"count": len(items),
	})
	return items, nil
}

func (s *Service) CreateDelegationGrant(ctx context.Context, userID uint64, input CreateDelegationGrantInput) (*domain.DelegationGrant, error) {
	if err := s.scope.EnsureDelegate(ctx, userID, string(domain.ScopeTypePlatform), "platform"); err != nil {
		return nil, err
	}
	if normalizeName(input.GrantorRef) == "" || normalizeName(input.DelegateRef) == "" || len(input.AllowedRoleLevels) == 0 {
		return nil, ErrIdentityTenancyInvalid
	}
	if !input.ValidUntil.After(input.ValidFrom) {
		return nil, ErrIdentityTenancyInvalid
	}
	item := &domain.DelegationGrant{
		GrantorRef:           strings.TrimSpace(input.GrantorRef),
		DelegateRef:          strings.TrimSpace(input.DelegateRef),
		AllowedRoleLevels:    strings.ToLower(strings.Join(input.AllowedRoleLevels, ",")),
		AllowedScopeSnapshot: marshalJSON(map[string]any{"scope": "bounded-by-grantor"}),
		Status:               domain.DelegationGrantStatusActive,
		ValidFrom:            input.ValidFrom,
		ValidUntil:           input.ValidUntil,
		Reason:               strings.TrimSpace(input.Reason),
		CreatedBy:            userID,
	}
	if err := s.delegations.Create(ctx, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionDelegationGrantCreate, ResourceTypeDelegationGrant, strconv.FormatUint(item.ID, 10), domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"grantorRef":  item.GrantorRef,
		"delegateRef": item.DelegateRef,
	})
	return item, nil
}
