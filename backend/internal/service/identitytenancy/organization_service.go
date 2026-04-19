package identitytenancy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListOrganizationUnits(ctx context.Context, userID uint64, filter OrganizationUnitListFilter) ([]domain.OrganizationUnit, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	return s.orgUnits.List(ctx, repository.OrganizationUnitListFilter{
		UnitType:     filter.UnitType,
		ParentUnitID: filter.ParentUnitID,
	})
}

func (s *Service) CreateOrganizationUnit(ctx context.Context, userID uint64, input CreateOrganizationUnitInput) (*domain.OrganizationUnit, error) {
	scopeType, scopeRef, err := s.scope.DefaultScope(ctx, userID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.EnsureManageOrg(ctx, userID, scopeType, scopeRef); err != nil {
		return nil, err
	}
	if normalizeName(input.Name) == "" || normalizeName(input.UnitType) == "" {
		return nil, ErrIdentityTenancyInvalid
	}
	parentID := uint64PtrIf(input.ParentUnitID)
	if _, err := s.orgUnits.FindByParentName(ctx, parentID, input.Name); err == nil {
		return nil, ErrIdentityTenancyConflict
	} else if !notFoundOrNil(err) {
		return nil, err
	}
	if input.ParentUnitID != 0 {
		parent, err := s.orgUnits.GetByID(ctx, input.ParentUnitID)
		if err != nil {
			return nil, err
		}
		if parent.ID == input.ParentUnitID && parent.ParentUnitID != nil && *parent.ParentUnitID == input.ParentUnitID {
			return nil, ErrIdentityTenancyBlocked
		}
	}
	status := input.Status
	if strings.TrimSpace(status) == "" {
		status = string(domain.OrganizationUnitStatusActive)
	}
	item := &domain.OrganizationUnit{
		UnitType:         domain.OrganizationUnitType(strings.ToLower(strings.TrimSpace(input.UnitType))),
		Name:             normalizeName(input.Name),
		Description:      strings.TrimSpace(input.Description),
		ParentUnitID:     uint64PtrIf(input.ParentUnitID),
		IdentitySourceID: uint64PtrIf(input.IdentitySourceID),
		OwnerUserID:      userID,
		Status:           domain.OrganizationUnitStatus(status),
	}
	if err := s.orgUnits.Create(ctx, item); err != nil {
		return nil, err
	}
	if err := s.memberships.Create(ctx, &domain.OrganizationMembership{
		UnitID:         item.ID,
		MemberType:     "user",
		MemberRef:      strconv.FormatUint(userID, 10),
		MembershipRole: domain.OrganizationMembershipRoleOwner,
		Status:         domain.OrganizationMembershipStatusActive,
		JoinedAt:       time.Now(),
	}); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionOrganizationCreate, ResourceTypeOrganizationUnit, strconv.FormatUint(item.ID, 10), domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"name":     item.Name,
		"unitType": item.UnitType,
	})
	return item, nil
}
