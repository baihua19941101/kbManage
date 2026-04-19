package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type RunbookArticleListFilter struct {
	WorkspaceIDs []uint64
	ProjectIDs   []uint64
	Status       string
}

type RunbookArticleRepository struct{ db *gorm.DB }

func NewRunbookArticleRepository(db *gorm.DB) *RunbookArticleRepository {
	return &RunbookArticleRepository{db: db}
}

func (r *RunbookArticleRepository) Create(ctx context.Context, item *domain.RunbookArticle) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RunbookArticleRepository) List(ctx context.Context, filter RunbookArticleListFilter) ([]domain.RunbookArticle, error) {
	query := r.db.WithContext(ctx).Model(&domain.RunbookArticle{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if len(filter.ProjectIDs) > 0 {
		query = query.Where("project_id IN ?", filter.ProjectIDs)
	}
	if filter.Status != "" {
		query = query.Where("status = ?", filter.Status)
	}
	var items []domain.RunbookArticle
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *RunbookArticleRepository) GetByID(ctx context.Context, id uint64) (*domain.RunbookArticle, error) {
	var item domain.RunbookArticle
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
