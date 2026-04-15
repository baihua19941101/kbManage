package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type PolicyExceptionListFilter struct {
	WorkspaceID *uint64
	ProjectID   *uint64
	PolicyID    *uint64
	HitID       *uint64
	Status      domain.PolicyExceptionStatus
}

type PolicyExceptionRepository struct {
	db *gorm.DB
}

func NewPolicyExceptionRepository(db *gorm.DB) *PolicyExceptionRepository {
	return &PolicyExceptionRepository{db: db}
}

func (r *PolicyExceptionRepository) Create(ctx context.Context, item *domain.PolicyExceptionRequest) error {
	if item == nil {
		return errors.New("policy exception request is required")
	}
	if r == nil || r.db == nil {
		return errors.New("policy exception repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PolicyExceptionRepository) GetByID(ctx context.Context, exceptionID uint64) (*domain.PolicyExceptionRequest, error) {
	if exceptionID == 0 {
		return nil, errors.New("exception id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("policy exception repository is not configured")
	}
	var item domain.PolicyExceptionRequest
	if err := r.db.WithContext(ctx).First(&item, exceptionID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PolicyExceptionRepository) List(ctx context.Context, filter PolicyExceptionListFilter) ([]domain.PolicyExceptionRequest, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("policy exception repository is not configured")
	}
	query := r.db.WithContext(ctx).Model(&domain.PolicyExceptionRequest{})
	if filter.WorkspaceID != nil {
		query = query.Where("workspace_id = ?", *filter.WorkspaceID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if filter.PolicyID != nil {
		query = query.Where("policy_id = ?", *filter.PolicyID)
	}
	if filter.HitID != nil {
		query = query.Where("hit_id = ?", *filter.HitID)
	}
	if strings.TrimSpace(string(filter.Status)) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	items := make([]domain.PolicyExceptionRequest, 0)
	if err := query.Order("created_at DESC, id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PolicyExceptionRepository) ExpireActiveBefore(ctx context.Context, cutoff time.Time) ([]uint64, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("policy exception repository is not configured")
	}
	var items []domain.PolicyExceptionRequest
	if err := r.db.WithContext(ctx).
		Model(&domain.PolicyExceptionRequest{}).
		Where("status = ? AND expires_at IS NOT NULL AND expires_at <= ?", domain.PolicyExceptionActive, cutoff).
		Find(&items).Error; err != nil {
		return nil, err
	}
	if len(items) == 0 {
		return []uint64{}, nil
	}
	ids := make([]uint64, 0, len(items))
	for i := range items {
		ids = append(ids, items[i].ID)
	}
	if err := r.db.WithContext(ctx).
		Model(&domain.PolicyExceptionRequest{}).
		Where("id IN ?", ids).
		Updates(map[string]any{
			"status":     domain.PolicyExceptionExpired,
			"updated_at": time.Now(),
		}).Error; err != nil {
		return nil, err
	}
	return ids, nil
}

func (r *PolicyExceptionRepository) UpdateFields(ctx context.Context, exceptionID uint64, updates map[string]any) error {
	if exceptionID == 0 {
		return errors.New("exception id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("policy exception repository is not configured")
	}
	updates["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&domain.PolicyExceptionRequest{}).Where("id = ?", exceptionID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
