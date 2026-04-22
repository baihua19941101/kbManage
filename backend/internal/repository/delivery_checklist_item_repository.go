package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DeliveryChecklistItemRepository struct{ db *gorm.DB }

func NewDeliveryChecklistItemRepository(db *gorm.DB) *DeliveryChecklistItemRepository {
	return &DeliveryChecklistItemRepository{db: db}
}

func (r *DeliveryChecklistItemRepository) Create(ctx context.Context, item *domain.DeliveryChecklistItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DeliveryChecklistItemRepository) ListByBundleID(ctx context.Context, bundleID uint64) ([]domain.DeliveryChecklistItem, error) {
	var items []domain.DeliveryChecklistItem
	err := r.db.WithContext(ctx).Where("bundle_id = ?", bundleID).Order("id ASC").Find(&items).Error
	return items, err
}
