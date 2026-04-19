package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type SREUpgradePlanListFilter struct {
	WorkspaceIDs []uint64
	ProjectIDs   []uint64
	Status       string
}

type SREUpgradePlanRepository struct{ db *gorm.DB }

func NewSREUpgradePlanRepository(db *gorm.DB) *SREUpgradePlanRepository {
	return &SREUpgradePlanRepository{db: db}
}

func (r *SREUpgradePlanRepository) Create(ctx context.Context, item *domain.SREUpgradePlan) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *SREUpgradePlanRepository) Update(ctx context.Context, item *domain.SREUpgradePlan) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *SREUpgradePlanRepository) GetByID(ctx context.Context, id uint64) (*domain.SREUpgradePlan, error) {
	var item domain.SREUpgradePlan
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *SREUpgradePlanRepository) List(ctx context.Context, filter SREUpgradePlanListFilter) ([]domain.SREUpgradePlan, error) {
	query := r.db.WithContext(ctx).Model(&domain.SREUpgradePlan{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if len(filter.ProjectIDs) > 0 {
		query = query.Where("project_id IN ?", filter.ProjectIDs)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	var items []domain.SREUpgradePlan
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *SREUpgradePlanRepository) FindLatestByScope(ctx context.Context, workspaceID, projectID uint64) (*domain.SREUpgradePlan, error) {
	var item domain.SREUpgradePlan
	query := r.db.WithContext(ctx).Model(&domain.SREUpgradePlan{}).Where("workspace_id = ?", workspaceID)
	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Order("id DESC").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
