package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type RollbackValidationRepository struct{ db *gorm.DB }

func NewRollbackValidationRepository(db *gorm.DB) *RollbackValidationRepository {
	return &RollbackValidationRepository{db: db}
}

func (r *RollbackValidationRepository) Create(ctx context.Context, item *domain.RollbackValidation) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RollbackValidationRepository) ListByUpgradePlanID(ctx context.Context, upgradePlanID uint64) ([]domain.RollbackValidation, error) {
	var items []domain.RollbackValidation
	err := r.db.WithContext(ctx).Where("upgrade_plan_id = ?", upgradePlanID).Order("id DESC").Find(&items).Error
	return items, err
}
