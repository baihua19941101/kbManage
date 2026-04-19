package identitytenancy

import (
	"context"
	"strconv"

	"kbmanage/backend/internal/domain"
)

func (s *Service) ListMemberships(ctx context.Context, userID, unitID uint64) ([]domain.OrganizationMembership, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	items, err := s.memberships.ListByUnitID(ctx, unitID)
	if err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionMembershipQuery, ResourceTypeOrganizationUnit, strconv.FormatUint(unitID, 10), domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"count": len(items),
	})
	return items, nil
}
