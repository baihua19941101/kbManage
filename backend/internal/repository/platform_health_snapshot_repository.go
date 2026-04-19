package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type PlatformHealthSnapshotRepository struct{ db *gorm.DB }

func NewPlatformHealthSnapshotRepository(db *gorm.DB) *PlatformHealthSnapshotRepository {
	return &PlatformHealthSnapshotRepository{db: db}
}

func (r *PlatformHealthSnapshotRepository) Create(ctx context.Context, item *domain.PlatformHealthSnapshot) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PlatformHealthSnapshotRepository) FindLatestByScope(ctx context.Context, workspaceID, projectID uint64) (*domain.PlatformHealthSnapshot, error) {
	var item domain.PlatformHealthSnapshot
	query := r.db.WithContext(ctx).Model(&domain.PlatformHealthSnapshot{}).Where("workspace_id = ?", workspaceID)
	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Order("snapshot_at DESC, id DESC").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
