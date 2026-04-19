package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type CapacityBaselineListFilter struct {
	WorkspaceIDs []uint64
	ProjectIDs   []uint64
	Status       string
}

type CapacityBaselineRepository struct{ db *gorm.DB }

func NewCapacityBaselineRepository(db *gorm.DB) *CapacityBaselineRepository {
	return &CapacityBaselineRepository{db: db}
}

func (r *CapacityBaselineRepository) Create(ctx context.Context, item *domain.CapacityBaseline) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *CapacityBaselineRepository) Update(ctx context.Context, item *domain.CapacityBaseline) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *CapacityBaselineRepository) List(ctx context.Context, filter CapacityBaselineListFilter) ([]domain.CapacityBaseline, error) {
	query := r.db.WithContext(ctx).Model(&domain.CapacityBaseline{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if len(filter.ProjectIDs) > 0 {
		query = query.Where("project_id IN ?", filter.ProjectIDs)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	var items []domain.CapacityBaseline
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *CapacityBaselineRepository) FindLatestByScope(ctx context.Context, workspaceID, projectID uint64) (*domain.CapacityBaseline, error) {
	var item domain.CapacityBaseline
	query := r.db.WithContext(ctx).Model(&domain.CapacityBaseline{}).Where("workspace_id = ?", workspaceID)
	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Order("id DESC").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
