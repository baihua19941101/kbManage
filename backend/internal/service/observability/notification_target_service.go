package observability

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type NotificationTargetService struct {
	repo *repository.NotificationTargetRepository
}

type UpsertNotificationTargetRequest struct {
	Name          string `json:"name"`
	TargetType    string `json:"targetType"`
	ConfigRef     string `json:"configRef"`
	ScopeSnapshot string `json:"scopeSnapshot"`
	Status        string `json:"status"`
}

func NewNotificationTargetService(repo *repository.NotificationTargetRepository) *NotificationTargetService {
	return &NotificationTargetService{repo: repo}
}

func (s *NotificationTargetService) List(ctx context.Context) ([]domain.NotificationTarget, error) {
	if s == nil || s.repo == nil {
		return []domain.NotificationTarget{}, nil
	}
	return s.repo.List(ctx)
}

func (s *NotificationTargetService) Get(ctx context.Context, id uint64) (*domain.NotificationTarget, error) {
	if s == nil || s.repo == nil {
		return nil, ErrObservabilityUnavailable
	}
	return s.repo.GetByID(ctx, id)
}

func (s *NotificationTargetService) Create(ctx context.Context, createdBy uint64, req UpsertNotificationTargetRequest) (*domain.NotificationTarget, error) {
	item := &domain.NotificationTarget{
		Name:          strings.TrimSpace(req.Name),
		TargetType:    strings.TrimSpace(req.TargetType),
		ConfigRef:     strings.TrimSpace(req.ConfigRef),
		ScopeSnapshot: strings.TrimSpace(req.ScopeSnapshot),
		Status:        normalizeTargetStatus(req.Status),
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

func (s *NotificationTargetService) Update(ctx context.Context, id uint64, req UpsertNotificationTargetRequest) (*domain.NotificationTarget, error) {
	if s == nil || s.repo == nil {
		return nil, ErrObservabilityUnavailable
	}
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	item.Name = strings.TrimSpace(req.Name)
	item.TargetType = strings.TrimSpace(req.TargetType)
	item.ConfigRef = strings.TrimSpace(req.ConfigRef)
	item.ScopeSnapshot = strings.TrimSpace(req.ScopeSnapshot)
	item.Status = normalizeTargetStatus(req.Status)
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *NotificationTargetService) Delete(ctx context.Context, id uint64) error {
	if s == nil || s.repo == nil {
		return nil
	}
	return s.repo.Delete(ctx, id)
}

func normalizeTargetStatus(in string) string {
	value := strings.TrimSpace(strings.ToLower(in))
	switch value {
	case "active", "disabled":
		return value
	default:
		return "active"
	}
}
