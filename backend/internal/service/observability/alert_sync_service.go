package observability

import (
	"context"
	"time"

	"kbmanage/backend/internal/domain"
	alertProvider "kbmanage/backend/internal/integration/observability/alerts"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
)

type AlertSyncService struct {
	alerts       alertProvider.Provider
	incidentRepo *repository.AlertIncidentRepository
	auditWriter  *auditSvc.EventWriter
}

func NewAlertSyncService(alerts alertProvider.Provider, incidentRepo *repository.AlertIncidentRepository) *AlertSyncService {
	return &AlertSyncService{
		alerts:       alerts,
		incidentRepo: incidentRepo,
	}
}

func (s *AlertSyncService) SetAuditWriter(w *auditSvc.EventWriter) {
	if s == nil {
		return
	}
	s.auditWriter = w
}

func (s *AlertSyncService) Sync(ctx context.Context) error {
	if s == nil || s.alerts == nil || s.incidentRepo == nil {
		return nil
	}
	items, err := s.alerts.ListAlerts(ctx, alertProvider.AlertQuery{})
	if err != nil {
		if s.auditWriter != nil {
			_ = s.auditWriter.WriteObservabilityEvent(
				ctx,
				"observability-sync-worker",
				nil,
				auditSvc.ObservabilityAuditActionAlertSync,
				"alerts",
				domain.AuditOutcomeFailed,
				map[string]any{"error": err.Error()},
			)
		}
		return err
	}

	now := time.Now().UTC()
	for _, item := range items {
		snapshot := &domain.AlertIncidentSnapshot{
			SourceIncidentKey: item.IncidentKey,
			Severity:          normalizeAlertSeverity(domain.AlertSeverity(item.Severity)),
			Status:            normalizeIncidentStatus(item.Status),
			Summary:           item.Summary,
			LastSyncedAt:      &now,
		}
		if err := s.incidentRepo.UpsertBySourceKey(ctx, snapshot); err != nil {
			if s.auditWriter != nil {
				_ = s.auditWriter.WriteObservabilityEvent(
					ctx,
					"observability-sync-worker",
					nil,
					auditSvc.ObservabilityAuditActionAlertSync,
					"alerts",
					domain.AuditOutcomeFailed,
					map[string]any{"error": err.Error(), "sourceIncidentKey": item.IncidentKey},
				)
			}
			return err
		}
	}
	if s.auditWriter != nil {
		_ = s.auditWriter.WriteObservabilityEvent(
			ctx,
			"observability-sync-worker",
			nil,
			auditSvc.ObservabilityAuditActionAlertSync,
			"alerts",
			domain.AuditOutcomeSuccess,
			map[string]any{"syncedCount": len(items)},
		)
	}
	return nil
}

func normalizeIncidentStatus(in string) domain.AlertIncidentStatus {
	switch domain.AlertIncidentStatus(in) {
	case domain.AlertIncidentStatusFiring, domain.AlertIncidentStatusAcknowledged, domain.AlertIncidentStatusSilenced, domain.AlertIncidentStatusResolved:
		return domain.AlertIncidentStatus(in)
	default:
		return domain.AlertIncidentStatusFiring
	}
}
