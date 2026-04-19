package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type DRDrillPlanRepository struct {
	db *gorm.DB
}

func NewDRDrillPlanRepository(db *gorm.DB) *DRDrillPlanRepository {
	return &DRDrillPlanRepository{db: db}
}

func (r *DRDrillPlanRepository) Create(ctx context.Context, item *domain.DRDrillPlan) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *DRDrillPlanRepository) Update(ctx context.Context, item *domain.DRDrillPlan) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *DRDrillPlanRepository) GetByID(ctx context.Context, id uint64) (*domain.DRDrillPlan, error) {
	var item domain.DRDrillPlan
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *DRDrillPlanRepository) List(ctx context.Context, workspaceIDs []uint64) ([]domain.DRDrillPlan, error) {
	query := r.db.WithContext(ctx).Model(&domain.DRDrillPlan{})
	if len(workspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", workspaceIDs)
	}
	var items []domain.DRDrillPlan
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}
