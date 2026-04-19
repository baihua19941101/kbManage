package identitytenancy

import (
	"context"
	"fmt"
	"time"

	"kbmanage/backend/internal/domain"
	identityint "kbmanage/backend/internal/integration/identity"
)

func (s *Service) ensureIdentityAccount(ctx context.Context, userID, sourceID uint64, externalRef string) error {
	if s.accounts == nil || userID == 0 || sourceID == 0 {
		return nil
	}
	_, err := s.accounts.FindBySourceExternalRef(ctx, sourceID, externalRef)
	if err == nil {
		return nil
	}
	if !notFoundOrNil(err) {
		return err
	}
	now := time.Now()
	return s.accounts.Create(ctx, &domain.IdentityAccount{
		UserID:           userID,
		IdentitySourceID: sourceID,
		ExternalRef:      externalRef,
		PrincipalType:    domain.IdentityPrincipalTypeUser,
		Status:           domain.IdentityAccountStatusActive,
		LastLoginAt:      &now,
	})
}

func (s *Service) syncSourcePrincipals(ctx context.Context, actorID uint64, source *domain.IdentitySource) error {
	if s.syncProvider == nil || source == nil {
		return nil
	}
	result, err := s.syncProvider.SyncDirectory(ctx, identityint.SyncRequest{
		SourceID:   source.ID,
		SourceType: string(source.SourceType),
	})
	if err != nil {
		return err
	}
	source.SyncState = domain.IdentitySyncState(result.State)
	if result.CompletedAt != nil {
		source.LastCheckedAt = result.CompletedAt
	}
	if err := s.sources.Update(ctx, source); err != nil {
		return err
	}
	if actorID != 0 {
		externalRef := fmt.Sprintf("%s:%d", source.SourceType, actorID)
		if len(result.Principals) > 0 && result.Principals[0].ExternalRef != "" {
			externalRef = result.Principals[0].ExternalRef
		}
		return s.ensureIdentityAccount(ctx, actorID, source.ID, externalRef)
	}
	return nil
}
