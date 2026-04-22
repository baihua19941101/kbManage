package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DeliveryReadinessBundleRepository struct{ db *gorm.DB }

func NewDeliveryReadinessBundleRepository(db *gorm.DB) *DeliveryReadinessBundleRepository {
	return &DeliveryReadinessBundleRepository{db: db}
}

func (r *DeliveryReadinessBundleRepository) Create(ctx context.Context, item *domain.DeliveryReadinessBundle) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DeliveryReadinessBundleRepository) List(ctx context.Context) ([]domain.DeliveryReadinessBundle, error) {
	var items []domain.DeliveryReadinessBundle
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}

func (r *DeliveryReadinessBundleRepository) GetByID(ctx context.Context, id uint64) (*domain.DeliveryReadinessBundle, error) {
	var item domain.DeliveryReadinessBundle
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
