package identitytenancy

import (
	"context"
	"strconv"
	"strings"

	"kbmanage/backend/internal/domain"
)

func (s *Service) CreateTenantScopeMapping(ctx context.Context, userID, unitID uint64, input CreateTenantScopeMappingInput) (*domain.TenantScopeMapping, error) {
	if _, err := s.orgUnits.GetByID(ctx, unitID); err != nil {
		return nil, err
	}
	if err := s.scope.EnsureManageOrg(ctx, userID, input.ScopeType, input.ScopeRef); err != nil {
		return nil, err
	}
	if normalizeName(input.ScopeType) == "" || normalizeName(input.ScopeRef) == "" || normalizeName(input.InheritanceMode) == "" {
		return nil, ErrIdentityTenancyInvalid
	}
	if _, err := s.mappings.FindByUnitScope(ctx, unitID, input.ScopeType, input.ScopeRef); err == nil {
		return nil, ErrIdentityTenancyConflict
	} else if !notFoundOrNil(err) {
		return nil, err
	}
	status := input.Status
	if strings.TrimSpace(status) == "" {
		status = string(domain.TenantScopeMappingStatusActive)
	}
	item := &domain.TenantScopeMapping{
		UnitID:          unitID,
		ScopeType:       domain.TenantScopeType(strings.ToLower(strings.TrimSpace(input.ScopeType))),
		ScopeRef:        strings.TrimSpace(input.ScopeRef),
		InheritanceMode: domain.TenantInheritanceMode(strings.ToLower(strings.TrimSpace(input.InheritanceMode))),
		Status:          domain.TenantScopeMappingStatus(status),
		CreatedBy:       userID,
	}
	if err := s.mappings.Create(ctx, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionTenantMappingCreate, ResourceTypeTenantMapping, strconv.FormatUint(item.ID, 10), domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"unitId":    unitID,
		"scopeType": item.ScopeType,
		"scopeRef":  item.ScopeRef,
	})
	return item, nil
}

func (s *Service) ListTenantScopeMappings(ctx context.Context, userID, unitID uint64) ([]domain.TenantScopeMapping, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	return s.mappings.ListByUnitID(ctx, unitID)
}
