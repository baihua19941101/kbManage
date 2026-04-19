package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type IdentityAuditListFilter struct {
	Action     string
	Outcome    string
	TargetType string
}

type IdentityAuditRepository struct {
	db *gorm.DB
}

func NewIdentityAuditRepository(db *gorm.DB) *IdentityAuditRepository {
	return &IdentityAuditRepository{db: db}
}

func (r *IdentityAuditRepository) Create(ctx context.Context, item *domain.IdentityGovernanceAuditEvent) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *IdentityAuditRepository) List(ctx context.Context, filter IdentityAuditListFilter) ([]domain.IdentityGovernanceAuditEvent, error) {
	query := r.db.WithContext(ctx).Model(&domain.IdentityGovernanceAuditEvent{})
	if v := strings.TrimSpace(filter.Action); v != "" {
		query = query.Where("action = ?", v)
	}
	if v := strings.TrimSpace(filter.Outcome); v != "" {
		query = query.Where("outcome = ?", v)
	}
	if v := strings.TrimSpace(filter.TargetType); v != "" {
		query = query.Where("target_type = ?", v)
	}
	var items []domain.IdentityGovernanceAuditEvent
	err := query.Order("occurred_at DESC, id DESC").Find(&items).Error
	return items, err
}
