package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type KeyOperationTraceRepository struct{ db *gorm.DB }

func NewKeyOperationTraceRepository(db *gorm.DB) *KeyOperationTraceRepository {
	return &KeyOperationTraceRepository{db: db}
}

func (r *KeyOperationTraceRepository) Create(ctx context.Context, item *domain.KeyOperationTrace) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *KeyOperationTraceRepository) List(ctx context.Context) ([]domain.KeyOperationTrace, error) {
	var items []domain.KeyOperationTrace
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}
