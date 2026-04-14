package repository

import (
	"context"
	"errors"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ClusterTargetGroupStore interface {
	Create(ctx context.Context, item *domain.ClusterTargetGroup) error
	GetByID(ctx context.Context, id uint64) (*domain.ClusterTargetGroup, error)
	ListByScope(ctx context.Context, workspaceID uint64, projectID *uint64) ([]domain.ClusterTargetGroup, error)
	ListByWorkspace(ctx context.Context, workspaceID uint64) ([]domain.ClusterTargetGroup, error)
	Update(ctx context.Context, item *domain.ClusterTargetGroup) error
}

type ClusterTargetGroupRepository struct {
	db *gorm.DB
}

func NewClusterTargetGroupRepository(db *gorm.DB) *ClusterTargetGroupRepository {
	return &ClusterTargetGroupRepository{db: db}
}

func (r *ClusterTargetGroupRepository) Create(ctx context.Context, item *domain.ClusterTargetGroup) error {
	if item == nil {
		return errors.New("cluster target group is required")
	}
	if r == nil || r.db == nil {
		return errors.New("cluster target group repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ClusterTargetGroupRepository) GetByID(ctx context.Context, id uint64) (*domain.ClusterTargetGroup, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("cluster target group repository is not configured")
	}
	var item domain.ClusterTargetGroup
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ClusterTargetGroupRepository) ListByWorkspace(ctx context.Context, workspaceID uint64) ([]domain.ClusterTargetGroup, error) {
	return r.ListByScope(ctx, workspaceID, nil)
}

func (r *ClusterTargetGroupRepository) ListByScope(ctx context.Context, workspaceID uint64, projectID *uint64) ([]domain.ClusterTargetGroup, error) {
	if workspaceID == 0 {
		return nil, errors.New("workspace id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("cluster target group repository is not configured")
	}
	query := r.db.WithContext(ctx).Model(&domain.ClusterTargetGroup{}).
		Where("workspace_id = ?", workspaceID)
	if projectID != nil {
		query = query.Where("project_id = ?", *projectID)
	}
	items := make([]domain.ClusterTargetGroup, 0)
	if err := query.Order("id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ClusterTargetGroupRepository) Update(ctx context.Context, item *domain.ClusterTargetGroup) error {
	if item == nil || item.ID == 0 {
		return errors.New("cluster target group id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("cluster target group repository is not configured")
	}
	return r.db.WithContext(ctx).Model(&domain.ClusterTargetGroup{}).Where("id = ?", item.ID).Updates(item).Error
}
