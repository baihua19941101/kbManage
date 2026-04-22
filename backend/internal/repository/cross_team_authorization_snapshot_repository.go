package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type CrossTeamAuthorizationSnapshotRepository struct{ db *gorm.DB }

func NewCrossTeamAuthorizationSnapshotRepository(db *gorm.DB) *CrossTeamAuthorizationSnapshotRepository {
	return &CrossTeamAuthorizationSnapshotRepository{db: db}
}

func (r *CrossTeamAuthorizationSnapshotRepository) Create(ctx context.Context, item *domain.CrossTeamAuthorizationSnapshot) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *CrossTeamAuthorizationSnapshotRepository) List(ctx context.Context) ([]domain.CrossTeamAuthorizationSnapshot, error) {
	var items []domain.CrossTeamAuthorizationSnapshot
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}
