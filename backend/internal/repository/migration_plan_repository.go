package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type MigrationPlanRepository struct {
	db *gorm.DB
}

func NewMigrationPlanRepository(db *gorm.DB) *MigrationPlanRepository {
	return &MigrationPlanRepository{db: db}
}

func (r *MigrationPlanRepository) Create(ctx context.Context, item *domain.MigrationPlan) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *MigrationPlanRepository) Update(ctx context.Context, item *domain.MigrationPlan) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *MigrationPlanRepository) GetByID(ctx context.Context, id uint64) (*domain.MigrationPlan, error) {
	var item domain.MigrationPlan
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
