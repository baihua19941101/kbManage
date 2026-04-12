package observability

import (
	"context"
	"time"

	"kbmanage/backend/internal/domain"
	alertProvider "kbmanage/backend/internal/integration/observability/alerts"
	"kbmanage/backend/internal/repository"
)

type AlertCenterService struct {
	alerts         alertProvider.Provider
	incidentRepo   *repository.AlertIncidentRepository
	handlingRecord *HandlingRecordService
}

func NewAlertCenterService(
	alerts alertProvider.Provider,
	incidentRepo *repository.AlertIncidentRepository,
	handlingRecord *HandlingRecordService,
) *AlertCenterService {
	if alerts == nil {
		alerts = alertProvider.NewMockProvider()
	}
	return &AlertCenterService{
		alerts:         alerts,
		incidentRepo:   incidentRepo,
		handlingRecord: handlingRecord,
	}
}

func (s *AlertCenterService) List(ctx context.Context, status domain.AlertIncidentStatus, limit int) ([]domain.AlertIncidentSnapshot, error) {
	if s == nil {
		return []domain.AlertIncidentSnapshot{}, nil
	}
	if s.incidentRepo != nil {
		items, err := s.incidentRepo.List(ctx, status, limit)
		if err != nil {
			return nil, err
		}
		if len(items) > 0 {
			return items, nil
		}
	}

	providerItems, err := s.alerts.ListAlerts(ctx, alertProvider.AlertQuery{Status: string(status)})
	if err != nil {
		return nil, err
	}
	out := make([]domain.AlertIncidentSnapshot, 0, len(providerItems))
	for _, item := range providerItems {
		out = append(out, domain.AlertIncidentSnapshot{
			SourceIncidentKey: item.IncidentKey,
			Severity:          normalizeAlertSeverity(domain.AlertSeverity(item.Severity)),
			Status:            normalizeIncidentStatus(item.Status),
			Summary:           item.Summary,
		})
	}
	return out, nil
}

func (s *AlertCenterService) Get(ctx context.Context, id uint64) (*domain.AlertIncidentSnapshot, error) {
	if s == nil || s.incidentRepo == nil {
		return nil, ErrObservabilityUnavailable
	}
	return s.incidentRepo.GetByID(ctx, id)
}

func (s *AlertCenterService) Acknowledge(ctx context.Context, id uint64, actedBy uint64, note string) (*domain.AlertIncidentSnapshot, error) {
	if s == nil || s.incidentRepo == nil {
		return nil, ErrObservabilityUnavailable
	}
	item, err := s.incidentRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now().UTC()
	item.Status = domain.AlertIncidentStatusAcknowledged
	item.AcknowledgedAt = &now
	if err := s.incidentRepo.Update(ctx, item); err != nil {
		return nil, err
	}
	if s.handlingRecord != nil {
		_, _ = s.handlingRecord.Create(ctx, id, actedBy, "acknowledge", note)
	}
	return item, nil
}
