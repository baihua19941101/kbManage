package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type GovernanceActionItemRepository struct{ db *gorm.DB }

func NewGovernanceActionItemRepository(db *gorm.DB) *GovernanceActionItemRepository {
	return &GovernanceActionItemRepository{db: db}
}

func (r *GovernanceActionItemRepository) Create(ctx context.Context, item *domain.GovernanceActionItem) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GovernanceActionItemRepository) List(ctx context.Context) ([]domain.GovernanceActionItem, error) {
	var items []domain.GovernanceActionItem
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}
