package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type UpgradePlanRepository struct {
	db *gorm.DB
}

func NewUpgradePlanRepository(db *gorm.DB) *UpgradePlanRepository {
	return &UpgradePlanRepository{db: db}
}

func (r *UpgradePlanRepository) Create(ctx context.Context, item *domain.UpgradePlan) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *UpgradePlanRepository) Update(ctx context.Context, item *domain.UpgradePlan) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *UpgradePlanRepository) GetByID(ctx context.Context, id uint64) (*domain.UpgradePlan, error) {
	var item domain.UpgradePlan
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *UpgradePlanRepository) ListByClusterID(ctx context.Context, clusterID uint64) ([]domain.UpgradePlan, error) {
	var items []domain.UpgradePlan
	err := r.db.WithContext(ctx).Where("cluster_id = ?", clusterID).Order("id DESC").Find(&items).Error
	return items, err
}
