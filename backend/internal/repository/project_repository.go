package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ProjectRepository struct {
	db *gorm.DB
}

func NewProjectRepository(db *gorm.DB) *ProjectRepository {
	return &ProjectRepository{db: db}
}

func (r *ProjectRepository) Create(ctx context.Context, p *domain.Project) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *ProjectRepository) ListByWorkspace(ctx context.Context, workspaceID uint64) ([]domain.Project, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	var items []domain.Project
	err := r.db.WithContext(ctx).Where("workspace_id = ?", workspaceID).Find(&items).Error
	return items, err
}

func (r *ProjectRepository) GetByID(ctx context.Context, id uint64) (*domain.Project, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	var item domain.Project
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
