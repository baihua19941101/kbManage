package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type RestoreJobListFilter struct {
	WorkspaceIDs []uint64
	JobType      string
	Status       string
}

type RestoreJobRepository struct {
	db *gorm.DB
}

func NewRestoreJobRepository(db *gorm.DB) *RestoreJobRepository {
	return &RestoreJobRepository{db: db}
}

func (r *RestoreJobRepository) Create(ctx context.Context, item *domain.RestoreJob) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RestoreJobRepository) Update(ctx context.Context, item *domain.RestoreJob) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *RestoreJobRepository) GetByID(ctx context.Context, id uint64) (*domain.RestoreJob, error) {
	var item domain.RestoreJob
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RestoreJobRepository) List(ctx context.Context, filter RestoreJobListFilter) ([]domain.RestoreJob, error) {
	query := r.db.WithContext(ctx).Model(&domain.RestoreJob{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if v := strings.TrimSpace(filter.JobType); v != "" {
		query = query.Where("job_type = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	var items []domain.RestoreJob
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}
