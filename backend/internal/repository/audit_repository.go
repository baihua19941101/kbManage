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
	StartAt     *time.Time
	EndAt       *time.Time
	ActorID     *uint64
	ClusterID   *uint64
	WorkspaceID *uint64
	ProjectID   *uint64
	Action      string
	Outcome     string
	Result      string
	Resource    string
	Limit       int
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
	outcome := strings.TrimSpace(strings.ToLower(q.Result))
	if outcome == "" {
		outcome = strings.TrimSpace(strings.ToLower(q.Outcome))
	}
	resource := strings.TrimSpace(strings.ToLower(q.Resource))

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
		if q.ClusterID != nil {
			tx = tx.Where("cluster_id = ?", *q.ClusterID)
		}
		if q.WorkspaceID != nil {
			tx = tx.Where("workspace_id = ?", *q.WorkspaceID)
		}
		if q.ProjectID != nil {
			tx = tx.Where("project_id = ?", *q.ProjectID)
		}
		if action != "" {
			tx = tx.Where("action = ?", action)
		}
		if outcome != "" {
			tx = tx.Where("outcome = ?", outcome)
		}
		if resource != "" {
			like := "%" + resource + "%"
			detailsExpr := "LOWER(CAST(details AS TEXT))"
			if r.db.Dialector != nil && strings.EqualFold(r.db.Dialector.Name(), "mysql") {
				detailsExpr = "LOWER(CAST(details AS CHAR))"
			}
			tx = tx.Where("(LOWER(resource_type) LIKE ? OR LOWER(resource_id) LIKE ? OR "+detailsExpr+" LIKE ?)", like, like, like)
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
		if q.ClusterID != nil {
			if event.ClusterID == nil || *event.ClusterID != *q.ClusterID {
				continue
			}
		}
		if q.WorkspaceID != nil {
			if event.WorkspaceID == nil || *event.WorkspaceID != *q.WorkspaceID {
				continue
			}
		}
		if q.ProjectID != nil {
			if event.ProjectID == nil || *event.ProjectID != *q.ProjectID {
				continue
			}
		}
		if action != "" && event.Action != action {
			continue
		}
		if outcome != "" && string(event.Outcome) != outcome {
			continue
		}
		if resource != "" {
			if !strings.Contains(strings.ToLower(strings.TrimSpace(event.ResourceType)), resource) &&
				!strings.Contains(strings.ToLower(strings.TrimSpace(event.ResourceID)), resource) &&
				!strings.Contains(strings.ToLower(strings.TrimSpace(string(event.Details))), resource) {
				continue
			}
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

func (r *AuditRepository) DeleteOlderThan(ctx context.Context, cutoff time.Time) (int64, error) {
	if r.db != nil {
		res := r.db.WithContext(ctx).
			Where("created_at < ?", cutoff).
			Delete(&domain.AuditEvent{})
		return res.RowsAffected, res.Error
	}

	r.mu.Lock()
	defer r.mu.Unlock()

	filtered := r.events[:0]
	var deleted int64
	for _, event := range r.events {
		if event.CreatedAt.Before(cutoff) {
			deleted++
			continue
		}
		filtered = append(filtered, event)
	}
	r.events = filtered
	return deleted, nil
}
