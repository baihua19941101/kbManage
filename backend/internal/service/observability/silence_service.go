package observability

import (
	"context"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type SilenceService struct {
	repo *repository.SilenceWindowRepository
}

type CreateSilenceRequest struct {
	Name          string `json:"name"`
	ScopeSnapshot string `json:"scopeSnapshot"`
	Reason        string `json:"reason"`
	StartsAt      string `json:"startsAt"`
	EndsAt        string `json:"endsAt"`
}

func NewSilenceService(repo *repository.SilenceWindowRepository) *SilenceService {
	return &SilenceService{repo: repo}
}

func (s *SilenceService) List(ctx context.Context, status domain.SilenceWindowStatus) ([]domain.SilenceWindow, error) {
	if s == nil || s.repo == nil {
		return []domain.SilenceWindow{}, nil
	}
	items, err := s.repo.List(ctx, status)
	if err != nil {
		return nil, err
	}
	return applySilenceRuntimeStatus(items), nil
}

func (s *SilenceService) Create(ctx context.Context, createdBy uint64, req CreateSilenceRequest) (*domain.SilenceWindow, error) {
	startsAt := parseTimeOrDefault(req.StartsAt, time.Now().UTC())
	endsAt := parseTimeOrDefault(req.EndsAt, startsAt.Add(30*time.Minute))
	item := &domain.SilenceWindow{
		Name:          strings.TrimSpace(req.Name),
		ScopeSnapshot: strings.TrimSpace(req.ScopeSnapshot),
		Reason:        strings.TrimSpace(req.Reason),
		StartsAt:      startsAt,
		EndsAt:        endsAt,
		Status:        calculateSilenceStatus(startsAt, endsAt, false),
		CreatedBy:     createdBy,
	}
	if s == nil || s.repo == nil {
		return item, nil
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *SilenceService) Cancel(ctx context.Context, id uint64, canceledBy uint64) (*domain.SilenceWindow, error) {
	if s == nil || s.repo == nil {
		return nil, ErrObservabilityUnavailable
	}
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Status = domain.SilenceWindowStatusCanceled
	if canceledBy != 0 {
		item.CanceledBy = &canceledBy
	}
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func applySilenceRuntimeStatus(items []domain.SilenceWindow) []domain.SilenceWindow {
	out := make([]domain.SilenceWindow, 0, len(items))
	now := time.Now().UTC()
	for _, item := range items {
		item.Status = calculateSilenceStatus(item.StartsAt, item.EndsAt, item.Status == domain.SilenceWindowStatusCanceled)
		out = append(out, item)
	}
	_ = now
	return out
}

func calculateSilenceStatus(startsAt, endsAt time.Time, canceled bool) domain.SilenceWindowStatus {
	if canceled {
		return domain.SilenceWindowStatusCanceled
	}
	now := time.Now().UTC()
	if now.Before(startsAt) {
		return domain.SilenceWindowStatusScheduled
	}
	if now.After(endsAt) {
		return domain.SilenceWindowStatusExpired
	}
	return domain.SilenceWindowStatusActive
}

func parseTimeOrDefault(raw string, fallback time.Time) time.Time {
	value := strings.TrimSpace(raw)
	if value == "" {
		return fallback
	}
	t, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return fallback
	}
	return t.UTC()
}
