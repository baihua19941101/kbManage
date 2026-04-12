package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// WorkspaceClusterBinding maps a workspace to a managed cluster.
type WorkspaceClusterBinding struct {
	ID                uint64    `json:"id" gorm:"primaryKey"`
	WorkspaceID       uint64    `json:"workspaceId" gorm:"not null;uniqueIndex:uk_workspace_cluster"`
	ClusterID         uint64    `json:"clusterId" gorm:"not null;uniqueIndex:uk_workspace_cluster"`
	DefaultNamespaces string    `json:"defaultNamespaces,omitempty" gorm:"size:1024"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (WorkspaceClusterBinding) TableName() string { return "workspace_cluster_bindings" }

// ProjectClusterBinding maps a project to a managed cluster.
type ProjectClusterBinding struct {
	ID                uint64    `json:"id" gorm:"primaryKey"`
	ProjectID         uint64    `json:"projectId" gorm:"not null;uniqueIndex:uk_project_cluster"`
	ClusterID         uint64    `json:"clusterId" gorm:"not null;uniqueIndex:uk_project_cluster"`
	DefaultNamespaces string    `json:"defaultNamespaces,omitempty" gorm:"size:1024"`
	CreatedAt         time.Time `json:"createdAt"`
}

func (ProjectClusterBinding) TableName() string { return "project_cluster_bindings" }

type WorkspaceClusterRepository struct {
	db *gorm.DB
}

func NewWorkspaceClusterRepository(db *gorm.DB) *WorkspaceClusterRepository {
	return &WorkspaceClusterRepository{db: db}
}

func (r *WorkspaceClusterRepository) Create(ctx context.Context, item *WorkspaceClusterBinding) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *WorkspaceClusterRepository) ListByWorkspace(ctx context.Context, workspaceID uint64) ([]WorkspaceClusterBinding, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	var items []WorkspaceClusterBinding
	err := r.db.WithContext(ctx).
		Where("workspace_id = ?", workspaceID).
		Order("id DESC").
		Find(&items).Error
	return items, err
}

func (r *WorkspaceClusterRepository) ListByCluster(ctx context.Context, clusterID uint64) ([]WorkspaceClusterBinding, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	var items []WorkspaceClusterBinding
	err := r.db.WithContext(ctx).
		Where("cluster_id = ?", clusterID).
		Order("id DESC").
		Find(&items).Error
	return items, err
}

func (r *WorkspaceClusterRepository) ListByWorkspaces(ctx context.Context, workspaceIDs []uint64) ([]WorkspaceClusterBinding, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}
	if len(workspaceIDs) == 0 {
		return []WorkspaceClusterBinding{}, nil
	}

	var items []WorkspaceClusterBinding
	err := r.db.WithContext(ctx).
		Where("workspace_id IN ?", workspaceIDs).
		Order("id DESC").
		Find(&items).Error
	return items, err
}

func (r *WorkspaceClusterRepository) CreateProjectBinding(ctx context.Context, item *ProjectClusterBinding) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *WorkspaceClusterRepository) ListByProject(ctx context.Context, projectID uint64) ([]ProjectClusterBinding, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	var items []ProjectClusterBinding
	err := r.db.WithContext(ctx).
		Where("project_id = ?", projectID).
		Order("id DESC").
		Find(&items).Error
	return items, err
}

func (r *WorkspaceClusterRepository) ListProjectBindingsByCluster(ctx context.Context, clusterID uint64) ([]ProjectClusterBinding, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	var items []ProjectClusterBinding
	err := r.db.WithContext(ctx).
		Where("cluster_id = ?", clusterID).
		Order("id DESC").
		Find(&items).Error
	return items, err
}
