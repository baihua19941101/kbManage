package identitytenancy

import (
	"context"
	"strconv"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

func (s *Service) ListSessions(ctx context.Context, userID uint64, filter SessionListFilter) ([]domain.SessionRecord, error) {
	if err := s.scope.EnsureSessionGovern(ctx, userID); err != nil {
		return nil, err
	}
	if err := s.reconcileExpiredAssignments(ctx); err != nil {
		return nil, err
	}
	items, err := s.sessions.List(ctx, repository.SessionRecordListFilter{
		Status:    filter.Status,
		RiskLevel: filter.RiskLevel,
	})
	if err != nil {
		return nil, err
	}
	if len(items) == 0 {
		source, err := s.sources.FindLocal(ctx)
		if err == nil {
			_ = s.bootstrapSessionForUser(ctx, userID, source)
			items, err = s.sessions.List(ctx, repository.SessionRecordListFilter{
				Status:    filter.Status,
				RiskLevel: filter.RiskLevel,
			})
			if err != nil {
				return nil, err
			}
		}
	}
	_ = s.sessionCache.Store(ctx, userID, items)
	s.writeAudit(ctx, userID, ActionSessionGovernanceRead, ResourceTypeSessionRecord, "list", domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"count": len(items),
	})
	return items, nil
}

func (s *Service) RevokeSession(ctx context.Context, actorUserID, sessionID uint64) (*domain.SessionRecord, error) {
	if err := s.scope.EnsureSessionGovern(ctx, actorUserID); err != nil {
		return nil, err
	}
	item, err := s.sessions.GetByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	item.Status = domain.IdentitySessionStatusRevoked
	item.RevokedAt = &now
	item.LastSeenAt = &now
	if err := s.sessions.Update(ctx, item); err != nil {
		return nil, err
	}
	_ = s.revocations.Mark(ctx, item.UserID, "manual-session-revoke")
	s.writeAudit(ctx, actorUserID, ActionSessionRevoke, ResourceTypeSessionRecord, strconv.FormatUint(item.ID, 10), domain.IdentityAuditOutcomeSucceeded, map[string]any{
		"userId": item.UserID,
	})
	return item, nil
}

func (s *Service) reconcileExpiredAssignments(ctx context.Context) error {
	items, err := s.assignments.List(ctx, repository.RoleAssignmentListFilter{})
	if err != nil {
		return err
	}
	now := time.Now()
	for i := range items {
		if items[i].Status == domain.RoleAssignmentStatusActive && items[i].ValidUntil != nil && !items[i].ValidUntil.After(now) {
			items[i].Status = domain.RoleAssignmentStatusExpired
			if err := s.assignments.Update(ctx, &items[i]); err != nil {
				return err
			}
			if items[i].SubjectType == "user" {
				if parsedUserID, err := strconv.ParseUint(items[i].SubjectRef, 10, 64); err == nil {
					_ = s.revocations.Mark(ctx, parsedUserID, "assignment-expired")
				}
			}
		}
	}
	return nil
}
