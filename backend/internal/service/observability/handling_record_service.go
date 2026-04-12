package observability

import (
	"context"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type HandlingRecordService struct {
	repo *repository.AlertIncidentRepository
}

func NewHandlingRecordService(repo *repository.AlertIncidentRepository) *HandlingRecordService {
	return &HandlingRecordService{repo: repo}
}

func (s *HandlingRecordService) Create(
	ctx context.Context,
	incidentID uint64,
	actedBy uint64,
	actionType string,
	content string,
) (*domain.AlertHandlingRecord, error) {
	item := &domain.AlertHandlingRecord{
		IncidentID: incidentID,
		ActionType: strings.TrimSpace(actionType),
		Content:    strings.TrimSpace(content),
		ActedBy:    actedBy,
		ActedAt:    time.Now().UTC(),
	}
	if s == nil || s.repo == nil {
		return item, nil
	}
	if err := s.repo.CreateHandlingRecord(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *HandlingRecordService) List(ctx context.Context, incidentID uint64) ([]domain.AlertHandlingRecord, error) {
	if s == nil || s.repo == nil {
		return []domain.AlertHandlingRecord{}, nil
	}
	return s.repo.ListHandlingRecords(ctx, incidentID)
}
