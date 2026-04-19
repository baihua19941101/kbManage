package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type AlertBaselineRepository struct{ db *gorm.DB }

func NewAlertBaselineRepository(db *gorm.DB) *AlertBaselineRepository {
	return &AlertBaselineRepository{db: db}
}

func (r *AlertBaselineRepository) Create(ctx context.Context, item *domain.AlertBaseline) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *AlertBaselineRepository) ListByRunbookID(ctx context.Context, runbookID uint64) ([]domain.AlertBaseline, error) {
	var items []domain.AlertBaseline
	err := r.db.WithContext(ctx).Where("recommended_runbook_id = ?", runbookID).Order("id DESC").Find(&items).Error
	return items, err
}
