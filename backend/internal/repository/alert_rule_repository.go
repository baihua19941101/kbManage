package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type AlertRuleRepository struct {
	db *gorm.DB
}

func NewAlertRuleRepository(db *gorm.DB) *AlertRuleRepository {
	return &AlertRuleRepository{db: db}
}

func (r *AlertRuleRepository) Create(ctx context.Context, item *domain.AlertRule) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *AlertRuleRepository) GetByID(ctx context.Context, id uint64) (*domain.AlertRule, error) {
	var item domain.AlertRule
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *AlertRuleRepository) List(ctx context.Context, status domain.AlertRuleStatus) ([]domain.AlertRule, error) {
	tx := r.db.WithContext(ctx).Model(&domain.AlertRule{})
	if status != "" {
		tx = tx.Where("status = ?", status)
	}

	var items []domain.AlertRule
	err := tx.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *AlertRuleRepository) Update(ctx context.Context, item *domain.AlertRule) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *AlertRuleRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.AlertRule{}, id).Error
}

