package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type HAPolicyListFilter struct {
	WorkspaceIDs []uint64
	ProjectIDs   []uint64
	Status       string
	Keyword      string
}

type HAPolicyRepository struct{ db *gorm.DB }

func NewHAPolicyRepository(db *gorm.DB) *HAPolicyRepository { return &HAPolicyRepository{db: db} }

func (r *HAPolicyRepository) Create(ctx context.Context, item *domain.HAPolicy) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *HAPolicyRepository) Update(ctx context.Context, item *domain.HAPolicy) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *HAPolicyRepository) GetByID(ctx context.Context, id uint64) (*domain.HAPolicy, error) {
	var item domain.HAPolicy
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *HAPolicyRepository) FindByName(ctx context.Context, workspaceID uint64, name string) (*domain.HAPolicy, error) {
	var item domain.HAPolicy
	err := r.db.WithContext(ctx).Where("workspace_id = ? AND lower(name) = ?", workspaceID, strings.ToLower(strings.TrimSpace(name))).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *HAPolicyRepository) List(ctx context.Context, filter HAPolicyListFilter) ([]domain.HAPolicy, error) {
	query := r.db.WithContext(ctx).Model(&domain.HAPolicy{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if len(filter.ProjectIDs) > 0 {
		query = query.Where("project_id IN ?", filter.ProjectIDs)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	if v := strings.TrimSpace(filter.Keyword); v != "" {
		like := "%" + strings.ToLower(v) + "%"
		query = query.Where("lower(name) LIKE ? OR lower(deployment_mode) LIKE ?", like, like)
	}
	var items []domain.HAPolicy
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *HAPolicyRepository) FindLatestByScope(ctx context.Context, workspaceID, projectID uint64) (*domain.HAPolicy, error) {
	var item domain.HAPolicy
	query := r.db.WithContext(ctx).Model(&domain.HAPolicy{}).Where("workspace_id = ?", workspaceID)
	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Order("id DESC").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
