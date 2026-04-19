package identitytenancy

import (
	"context"
	"strconv"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListRoleDefinitions(ctx context.Context, userID uint64, filter RoleDefinitionListFilter) ([]domain.RoleDefinition, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	return s.roles.List(ctx, repository.RoleDefinitionListFilter{
		RoleLevel: filter.RoleLevel,
		Status:    filter.Status,
	})
}

func (s *Service) CreateRoleDefinition(ctx context.Context, userID uint64, input CreateRoleDefinitionInput) (*domain.RoleDefinition, error) {
	if err := s.scope.EnsureManageRole(ctx, userID, string(domain.ScopeTypePlatform), "platform"); err != nil {
		return nil, err
	}
	if normalizeName(input.Name) == "" || normalizeName(input.RoleLevel) == "" || normalizeName(input.PermissionSummary) == "" || normalizeName(input.InheritancePolicy) == "" {
		return nil, ErrIdentityTenancyInvalid
	}
	if _, err := s.roles.FindByLevelName(ctx, input.RoleLevel, input.Name); err == nil {
		return nil, ErrIdentityTenancyConflict
	} else if !notFoundOrNil(err) {
		return nil, err
	}
	status := input.Status
	if strings.TrimSpace(status) == "" {
		status = string(domain.RoleDefinitionStatusActive)
	}
	item := &domain.RoleDefinition{
		Name:              normalizeName(input.Name),
		RoleLevel:         domain.RoleLevel(strings.ToLower(strings.TrimSpace(input.RoleLevel))),
		Description:       strings.TrimSpace(input.Description),
		PermissionSummary: strings.TrimSpace(input.PermissionSummary),
		InheritancePolicy: domain.RoleInheritancePolicy(strings.ToLower(strings.TrimSpace(input.InheritancePolicy))),
		Delegable:         input.Delegable,
		Status:            domain.RoleDefinitionStatus(status),
		CreatedBy:         userID,
	}
	if err := s.roles.Create(ctx, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionRoleDefinitionCreate, ResourceTypeRoleDefinition, strconv.FormatUint(item.ID, 10), domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"name":      item.Name,
		"roleLevel": item.RoleLevel,
	})
	return item, nil
}
