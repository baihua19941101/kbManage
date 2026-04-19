package marketplace

import (
	"context"
	"time"

	"kbmanage/backend/internal/domain"
)

func (s *Service) appendLifecycleRecord(ctx context.Context, extensionID uint64, action domain.ExtensionLifecycleAction, scopeType, scopeRef, outcome, reason string, userID uint64) error {
	return s.lifecycle.Create(ctx, &domain.ExtensionLifecycleRecord{
		ExtensionPackageID: extensionID,
		Action:             action,
		ScopeType:          scopeType,
		ScopeRef:           scopeRef,
		Outcome:            outcome,
		Reason:             reason,
		ExecutedBy:         userID,
		ExecutedAt:         time.Now(),
	})
}
