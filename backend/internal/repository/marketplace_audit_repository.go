package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type MarketplaceAuditListFilter struct {
	Action     string
	TargetType string
	Outcome    string
}

type MarketplaceAuditRepository struct {
	db *gorm.DB
}

func NewMarketplaceAuditRepository(db *gorm.DB) *MarketplaceAuditRepository {
	return &MarketplaceAuditRepository{db: db}
}

func (r *MarketplaceAuditRepository) Create(ctx context.Context, item *domain.MarketplaceAuditEvent) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *MarketplaceAuditRepository) List(ctx context.Context, filter MarketplaceAuditListFilter) ([]domain.MarketplaceAuditEvent, error) {
	query := r.db.WithContext(ctx).Model(&domain.MarketplaceAuditEvent{})
	if v := strings.TrimSpace(filter.Action); v != "" {
		query = query.Where("action = ?", v)
	}
	if v := strings.TrimSpace(filter.TargetType); v != "" {
		query = query.Where("target_type = ?", v)
	}
	if v := strings.TrimSpace(filter.Outcome); v != "" {
		query = query.Where("outcome = ?", v)
	}
	var items []domain.MarketplaceAuditEvent
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}
