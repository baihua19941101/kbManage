package workloadops

import (
	"context"
	"time"

	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
)

const defaultTerminalActiveTTL = 30 * time.Minute

func (s *Service) expireSessionIfNeeded(ctx context.Context, item *domain.TerminalSession) (*domain.TerminalSession, error) {
	if item == nil {
		return nil, nil
	}
	if item.Status != domain.TerminalSessionStatusActive || item.StartedAt == nil {
		return item, nil
	}
	if time.Since(*item.StartedAt) <= defaultTerminalActiveTTL {
		return item, nil
	}
	if s.sessions == nil {
		return item, nil
	}
	if err := s.sessions.UpdateStatus(ctx, item.ID, domain.TerminalSessionStatusExpired, "session timeout"); err != nil {
		return nil, err
	}
	latest, err := s.sessions.GetByID(ctx, item.ID)
	if err != nil {
		return nil, err
	}
	s.writeAudit(
		ctx,
		normalizeRequestID("", latest.OperatorID),
		latest.OperatorID,
		auditSvc.WorkloadOpsAuditTerminalClose,
		domain.AuditOutcomeSuccess,
		nil,
		withTerminalAuditBoundary(latest),
		nil,
	)
	return latest, nil
}
