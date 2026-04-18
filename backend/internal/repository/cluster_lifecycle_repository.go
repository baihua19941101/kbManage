package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ClusterLifecycleListFilter struct {
	WorkspaceIDs       []uint64
	ProjectIDs         []uint64
	Status             string
	InfrastructureType string
	DriverRef          string
	Keyword            string
}

type ClusterLifecycleRepository struct {
	db *gorm.DB
}

func NewClusterLifecycleRepository(db *gorm.DB) *ClusterLifecycleRepository {
	return &ClusterLifecycleRepository{db: db}
}

func (r *ClusterLifecycleRepository) Create(ctx context.Context, item *domain.ClusterLifecycleRecord) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ClusterLifecycleRepository) Update(ctx context.Context, item *domain.ClusterLifecycleRecord) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ClusterLifecycleRepository) GetByID(ctx context.Context, id uint64) (*domain.ClusterLifecycleRecord, error) {
	var item domain.ClusterLifecycleRecord
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ClusterLifecycleRepository) FindByNameInWorkspace(ctx context.Context, workspaceID uint64, name string) (*domain.ClusterLifecycleRecord, error) {
	var item domain.ClusterLifecycleRecord
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND lower(name) = ?", workspaceID, strings.ToLower(strings.TrimSpace(name))).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ClusterLifecycleRepository) List(ctx context.Context, filter ClusterLifecycleListFilter) ([]domain.ClusterLifecycleRecord, error) {
	query := r.db.WithContext(ctx).Model(&domain.ClusterLifecycleRecord{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if len(filter.ProjectIDs) > 0 {
		query = query.Where("project_id IN ?", filter.ProjectIDs)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	if v := strings.TrimSpace(filter.InfrastructureType); v != "" {
		query = query.Where("infrastructure_type = ?", v)
	}
	if v := strings.TrimSpace(filter.DriverRef); v != "" {
		query = query.Where("driver_ref = ?", v)
	}
	if v := strings.TrimSpace(filter.Keyword); v != "" {
		like := "%" + strings.ToLower(v) + "%"
		query = query.Where("lower(name) LIKE ? OR lower(display_name) LIKE ?", like, like)
	}

	var items []domain.ClusterLifecycleRecord
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

type ClusterLifecycleOperationRepository struct {
	db *gorm.DB
}

func NewClusterLifecycleOperationRepository(db *gorm.DB) *ClusterLifecycleOperationRepository {
	return &ClusterLifecycleOperationRepository{db: db}
}

func (r *ClusterLifecycleOperationRepository) Create(ctx context.Context, item *domain.LifecycleOperation) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ClusterLifecycleOperationRepository) Update(ctx context.Context, item *domain.LifecycleOperation) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ClusterLifecycleOperationRepository) GetByID(ctx context.Context, id uint64) (*domain.LifecycleOperation, error) {
	var item domain.LifecycleOperation
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ClusterLifecycleOperationRepository) FindRunningByClusterID(ctx context.Context, clusterID uint64) (*domain.LifecycleOperation, error) {
	var item domain.LifecycleOperation
	err := r.db.WithContext(ctx).
		Where("cluster_id = ? AND status IN ?", clusterID, []domain.LifecycleOperationStatus{domain.LifecycleOperationPending, domain.LifecycleOperationRunning}).
		Order("id DESC").
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}
