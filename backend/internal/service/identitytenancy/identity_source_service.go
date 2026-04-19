package identitytenancy

import (
	"context"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	identityint "kbmanage/backend/internal/integration/identity"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListIdentitySources(ctx context.Context, userID uint64, filter IdentitySourceListFilter) ([]domain.IdentitySource, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	items, err := s.sources.List(ctx, repository.IdentitySourceListFilter{
		SourceType: filter.SourceType,
		Status:     filter.Status,
	})
	if err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionIdentitySourceRead, ResourceTypeIdentitySource, "list", domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"count": len(items),
	})
	return items, nil
}

func (s *Service) CreateIdentitySource(ctx context.Context, userID uint64, input CreateIdentitySourceInput) (*domain.IdentitySource, error) {
	if err := s.scope.EnsureManageSource(ctx, userID); err != nil {
		return nil, err
	}
	if normalizeName(input.Name) == "" || normalizeName(input.SourceType) == "" || normalizeName(input.LoginMode) == "" || normalizeName(input.ScopeMode) == "" {
		return nil, ErrIdentityTenancyInvalid
	}
	if _, err := s.sources.FindByName(ctx, input.Name); err == nil {
		return nil, ErrIdentityTenancyConflict
	} else if !notFoundOrNil(err) {
		return nil, err
	}
	if strings.ToLower(strings.TrimSpace(input.SourceType)) != string(domain.IdentitySourceTypeLocal) {
		if err := s.ensureLocalFallback(ctx, userID); err != nil {
			return nil, err
		}
	}
	status := input.Status
	if strings.TrimSpace(status) == "" {
		status = string(domain.IdentitySourceStatusActive)
	}
	health, err := s.provider.CheckHealth(ctx, identityint.HealthCheckRequest{
		Name:       input.Name,
		SourceType: input.SourceType,
	})
	if err != nil {
		return nil, err
	}
	item := &domain.IdentitySource{
		Name:          normalizeName(input.Name),
		SourceType:    domain.IdentitySourceType(strings.ToLower(strings.TrimSpace(input.SourceType))),
		Status:        domain.IdentitySourceStatus(status),
		LoginMode:     domain.IdentityLoginMode(strings.ToLower(strings.TrimSpace(input.LoginMode))),
		ScopeMode:     domain.IdentityScopeMode(strings.ToLower(strings.TrimSpace(input.ScopeMode))),
		SyncState:     domain.IdentitySyncStateIdle,
		ConfigSummary: marshalJSON(map[string]any{"available": health.Available, "message": health.Message}),
		OwnerUserID:   userID,
		LastCheckedAt: health.LastCheckedAt,
	}
	if !health.Available {
		item.Status = domain.IdentitySourceStatusUnavailable
		item.LastError = health.Message
	}
	if err := s.sources.Create(ctx, item); err != nil {
		return nil, err
	}
	if err := s.syncSourcePrincipals(ctx, userID, item); err != nil {
		return nil, err
	}
	if err := s.bootstrapSessionForUser(ctx, userID, item); err != nil {
		return nil, err
	}
	s.writeAudit(ctx, userID, ActionIdentitySourceCreate, ResourceTypeIdentitySource, strconv.FormatUint(item.ID, 10), domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"name":       item.Name,
		"sourceType": item.SourceType,
		"loginMode":  item.LoginMode,
	})
	return item, nil
}

func (s *Service) GetIdentitySource(ctx context.Context, userID, sourceID uint64) (*domain.IdentitySource, error) {
	if err := s.scope.EnsureReadAny(ctx, userID); err != nil {
		return nil, err
	}
	return s.sources.GetByID(ctx, sourceID)
}

func (s *Service) bootstrapSessionForUser(ctx context.Context, userID uint64, source *domain.IdentitySource) error {
	if s.sessions == nil || source == nil || userID == 0 {
		return nil
	}
	existing, err := s.sessions.FindByUserSource(ctx, userID, source.ID)
	if err == nil && existing != nil {
		return nil
	}
	if err != nil && !notFoundOrNil(err) {
		return err
	}
	now := time.Now()
	session := &domain.SessionRecord{
		UserID:            userID,
		IdentitySourceID:  source.ID,
		LoginMethod:       string(source.LoginMode),
		Status:            domain.IdentitySessionStatusActive,
		RiskLevel:         domain.IdentityRiskLevelLow,
		PermissionVersion: strconv.FormatInt(now.Unix(), 10),
		LastSeenAt:        &now,
	}
	return s.sessions.Create(ctx, session)
}
