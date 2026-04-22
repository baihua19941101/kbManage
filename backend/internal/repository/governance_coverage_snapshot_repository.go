package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type GovernanceCoverageSnapshotRepository struct{ db *gorm.DB }

func NewGovernanceCoverageSnapshotRepository(db *gorm.DB) *GovernanceCoverageSnapshotRepository {
	return &GovernanceCoverageSnapshotRepository{db: db}
}

func (r *GovernanceCoverageSnapshotRepository) Create(ctx context.Context, item *domain.GovernanceCoverageSnapshot) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GovernanceCoverageSnapshotRepository) List(ctx context.Context) ([]domain.GovernanceCoverageSnapshot, error) {
	var items []domain.GovernanceCoverageSnapshot
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}
