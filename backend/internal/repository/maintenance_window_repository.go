package repository

import (
	"context"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type MaintenanceWindowListFilter struct {
	WorkspaceIDs []uint64
	ProjectIDs   []uint64
	Status       string
}

type MaintenanceWindowRepository struct{ db *gorm.DB }

func NewMaintenanceWindowRepository(db *gorm.DB) *MaintenanceWindowRepository {
	return &MaintenanceWindowRepository{db: db}
}

func (r *MaintenanceWindowRepository) Create(ctx context.Context, item *domain.MaintenanceWindow) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *MaintenanceWindowRepository) Update(ctx context.Context, item *domain.MaintenanceWindow) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *MaintenanceWindowRepository) List(ctx context.Context, filter MaintenanceWindowListFilter) ([]domain.MaintenanceWindow, error) {
	query := r.db.WithContext(ctx).Model(&domain.MaintenanceWindow{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if len(filter.ProjectIDs) > 0 {
		query = query.Where("project_id IN ?", filter.ProjectIDs)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	var items []domain.MaintenanceWindow
	err := query.Order("start_at DESC").Find(&items).Error
	return items, err
}

func (r *MaintenanceWindowRepository) FindActiveByScope(ctx context.Context, workspaceID, projectID uint64, now time.Time) (*domain.MaintenanceWindow, error) {
	var item domain.MaintenanceWindow
	query := r.db.WithContext(ctx).Model(&domain.MaintenanceWindow{}).
		Where("workspace_id = ? AND start_at <= ? AND end_at >= ?", workspaceID, now, now)
	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Order("id DESC").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
