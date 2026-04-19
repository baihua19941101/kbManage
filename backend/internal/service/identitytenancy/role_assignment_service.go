package identitytenancy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListRoleAssignments(ctx context.Context, userID uint64, filter RoleAssignmentListFilter) ([]domain.RoleAssignment, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.reconcileExpiredAssignments(ctx); err != nil {
		return nil, err
	}
	return s.assignments.List(ctx, repository.RoleAssignmentListFilter{
		SubjectRef: filter.SubjectRef,
		ScopeType:  filter.ScopeType,
		Status:     filter.Status,
	})
}

func (s *Service) CreateRoleAssignment(ctx context.Context, userID uint64, input CreateRoleAssignmentInput) (*domain.RoleAssignment, error) {
	if err := s.scope.EnsureManageRole(ctx, userID, input.ScopeType, input.ScopeRef); err != nil {
		return nil, err
	}
	if normalizeName(input.SubjectType) == "" || normalizeName(input.SubjectRef) == "" || input.RoleDefinitionID == 0 || normalizeName(input.ScopeType) == "" || normalizeName(input.ScopeRef) == "" {
		return nil, ErrIdentityTenancyInvalid
	}
	role, err := s.roles.GetByID(ctx, input.RoleDefinitionID)
	if err != nil {
		return nil, err
	}
	if input.DelegationGrantID != 0 {
		grant, err := s.delegations.GetByID(ctx, input.DelegationGrantID)
		if err != nil {
			return nil, err
		}
		if grant.Status != domain.DelegationGrantStatusActive || time.Now().After(grant.ValidUntil) {
			return nil, ErrIdentityTenancyBlocked
		}
		if !strings.Contains(strings.ToLower(grant.AllowedRoleLevels), strings.ToLower(string(role.RoleLevel))) {
			return nil, ErrIdentityTenancyBlocked
		}
	}
	sourceType := input.SourceType
	if strings.TrimSpace(sourceType) == "" {
		sourceType = string(domain.RoleAssignmentSourceDirect)
	}
	now := time.Now()
	item := &domain.RoleAssignment{
		SubjectType:       strings.ToLower(strings.TrimSpace(input.SubjectType)),
		SubjectRef:        strings.TrimSpace(input.SubjectRef),
		RoleDefinitionID:  input.RoleDefinitionID,
		ScopeType:         strings.ToLower(strings.TrimSpace(input.ScopeType)),
		ScopeRef:          strings.TrimSpace(input.ScopeRef),
		SourceType:        domain.RoleAssignmentSourceType(strings.ToLower(strings.TrimSpace(sourceType))),
		DelegationGrantID: uint64PtrIf(input.DelegationGrantID),
		ValidFrom:         now,
		ValidUntil:        input.ValidUntil,
		Status:            roleAssignmentStatusAt(now, input.ValidUntil),
		GrantedBy:         userID,
	}
	if item.SourceType == domain.RoleAssignmentSourceTemporary && item.ValidUntil == nil {
		return nil, ErrIdentityTenancyInvalid
	}
	if err := s.assignments.Create(ctx, item); err != nil {
		return nil, err
	}
	if err := s.recordRiskForAssignment(ctx, role, item); err != nil {
		return nil, err
	}
	if item.SubjectType == "user" {
		if parsedUserID, err := strconv.ParseUint(item.SubjectRef, 10, 64); err == nil {
			_ = s.permCache.Mark(ctx, parsedUserID, strconv.FormatInt(now.Unix(), 10))
		}
	}
	s.writeAudit(ctx, userID, ActionRoleAssignmentCreate, ResourceTypeRoleAssignment, strconv.FormatUint(item.ID, 10), domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"subjectRef": item.SubjectRef,
		"scopeType":  item.ScopeType,
		"scopeRef":   item.ScopeRef,
	})
	return item, nil
}
