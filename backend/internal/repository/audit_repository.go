package repository

import (
	"context"
	"sort"
	"strings"
	"sync"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type AuditRepository struct {
	db *gorm.DB

	mu     sync.RWMutex
	nextID uint64
	events []domain.AuditEvent
}

type AuditQuery struct {
	StartAt *time.Time
	EndAt   *time.Time
	ActorID *uint64
	Action  string
	Outcome string
	Limit   int
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{
		db:     db,
		nextID: 1,
		events: make([]domain.AuditEvent, 0),
	}
}

func (r *AuditRepository) Create(ctx context.Context, event *domain.AuditEvent) error {
	if event == nil {
		return nil
	}
	if r.db == nil {
		r.mu.Lock()
		defer r.mu.Unlock()

		copyItem := *event
		if copyItem.ID == 0 {
			copyItem.ID = r.nextID
			r.nextID++
		}
		if copyItem.CreatedAt.IsZero() {
			copyItem.CreatedAt = time.Now()
		}
		r.events = append(r.events, copyItem)
		*event = copyItem
		return nil
	}
	return r.db.WithContext(ctx).Create(event).Error
}

func (r *AuditRepository) Query(ctx context.Context, q AuditQuery) ([]domain.AuditEvent, error) {
	limit := q.Limit
	if limit <= 0 {
		limit = 100
	}
	action := strings.TrimSpace(q.Action)
	outcome := strings.TrimSpace(strings.ToLower(q.Outcome))

	if r.db != nil {
		tx := r.db.WithContext(ctx).Model(&domain.AuditEvent{})
		if q.StartAt != nil {
			tx = tx.Where("created_at >= ?", *q.StartAt)
		}
		if q.EndAt != nil {
			tx = tx.Where("created_at <= ?", *q.EndAt)
		}
		if q.ActorID != nil {
			tx = tx.Where("actor_id = ?", *q.ActorID)
		}
		if action != "" {
			tx = tx.Where("action = ?", action)
		}
		if outcome != "" {
			tx = tx.Where("outcome = ?", outcome)
		}

		var events []domain.AuditEvent
		err := tx.Order("created_at DESC").Limit(limit).Find(&events).Error
		return events, err
	}

	r.mu.RLock()
	defer r.mu.RUnlock()

	var events []domain.AuditEvent
	for i := range r.events {
		event := r.events[i]

		if q.StartAt != nil && event.CreatedAt.Before(*q.StartAt) {
			continue
		}
		if q.EndAt != nil && event.CreatedAt.After(*q.EndAt) {
			continue
		}
		if q.ActorID != nil {
			if event.ActorID == nil || *event.ActorID != *q.ActorID {
				continue
			}
		}
		if action != "" && event.Action != action {
			continue
		}
		if outcome != "" && string(event.Outcome) != outcome {
			continue
		}

		events = append(events, event)
	}

	sort.Slice(events, func(i, j int) bool {
		return events[i].CreatedAt.After(events[j].CreatedAt)
	})
	if len(events) > limit {
		events = events[:limit]
	}
	return events, nil
}

func (r *AuditRepository) QueryByTimeRange(ctx context.Context, start, end time.Time, limit int) ([]domain.AuditEvent, error) {
	return r.Query(ctx, AuditQuery{
		StartAt: &start,
		EndAt:   &end,
		Limit:   limit,
	})
}
