package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type GovernanceRiskEventRepository struct{ db *gorm.DB }

func NewGovernanceRiskEventRepository(db *gorm.DB) *GovernanceRiskEventRepository {
	return &GovernanceRiskEventRepository{db: db}
}

func (r *GovernanceRiskEventRepository) Create(ctx context.Context, item *domain.GovernanceRiskEvent) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GovernanceRiskEventRepository) List(ctx context.Context) ([]domain.GovernanceRiskEvent, error) {
	var items []domain.GovernanceRiskEvent
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}
