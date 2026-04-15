package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type SecurityPolicyListFilter struct {
	WorkspaceID *uint64
	ProjectID   *uint64
	ScopeLevel  domain.PolicyScopeLevel
	Category    domain.PolicyCategory
	Status      domain.PolicyStatus
}

type SecurityPolicyRepository struct {
	db *gorm.DB
}

func NewSecurityPolicyRepository(db *gorm.DB) *SecurityPolicyRepository {
	return &SecurityPolicyRepository{db: db}
}

func (r *SecurityPolicyRepository) Create(ctx context.Context, item *domain.SecurityPolicy) error {
	if item == nil {
		return errors.New("security policy is required")
	}
	if r == nil || r.db == nil {
		return errors.New("security policy repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *SecurityPolicyRepository) GetByID(ctx context.Context, id uint64) (*domain.SecurityPolicy, error) {
	if id == 0 {
		return nil, errors.New("policy id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("security policy repository is not configured")
	}
	var item domain.SecurityPolicy
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *SecurityPolicyRepository) List(ctx context.Context, filter SecurityPolicyListFilter) ([]domain.SecurityPolicy, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("security policy repository is not configured")
	}
	query := r.db.WithContext(ctx).Model(&domain.SecurityPolicy{})
	if filter.WorkspaceID != nil {
		query = query.Where("workspace_id = ?", *filter.WorkspaceID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if strings.TrimSpace(string(filter.ScopeLevel)) != "" {
		query = query.Where("scope_level = ?", filter.ScopeLevel)
	}
	if strings.TrimSpace(string(filter.Category)) != "" {
		query = query.Where("category = ?", filter.Category)
	}
	if strings.TrimSpace(string(filter.Status)) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	items := make([]domain.SecurityPolicy, 0)
	if err := query.Order("updated_at DESC, id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *SecurityPolicyRepository) UpdateFields(ctx context.Context, policyID uint64, updates map[string]any) error {
	if policyID == 0 {
		return errors.New("policy id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("security policy repository is not configured")
	}
	updates["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&domain.SecurityPolicy{}).Where("id = ?", policyID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
