package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type WorkspaceRepository struct {
	db *gorm.DB
}

func NewWorkspaceRepository(db *gorm.DB) *WorkspaceRepository {
	return &WorkspaceRepository{db: db}
}

func (r *WorkspaceRepository) Create(ctx context.Context, ws *domain.Workspace) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.WithContext(ctx).Create(ws).Error
}

func (r *WorkspaceRepository) List(ctx context.Context) ([]domain.Workspace, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	var items []domain.Workspace
	err := r.db.WithContext(ctx).Find(&items).Error
	return items, err
}

func (r *WorkspaceRepository) GetByID(ctx context.Context, id uint64) (*domain.Workspace, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	var item domain.Workspace
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
