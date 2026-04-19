package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type BackupPolicyListFilter struct {
	WorkspaceIDs []uint64
	Status       string
	ScopeType    string
}

type BackupPolicyRepository struct {
	db *gorm.DB
}

func NewBackupPolicyRepository(db *gorm.DB) *BackupPolicyRepository {
	return &BackupPolicyRepository{db: db}
}

func (r *BackupPolicyRepository) Create(ctx context.Context, item *domain.BackupPolicy) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *BackupPolicyRepository) Update(ctx context.Context, item *domain.BackupPolicy) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *BackupPolicyRepository) GetByID(ctx context.Context, id uint64) (*domain.BackupPolicy, error) {
	var item domain.BackupPolicy
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *BackupPolicyRepository) FindByScopeName(ctx context.Context, workspaceID uint64, scopeType, scopeRef, name string) (*domain.BackupPolicy, error) {
	var item domain.BackupPolicy
	err := r.db.WithContext(ctx).
		Where("workspace_id = ? AND scope_type = ? AND scope_ref = ? AND lower(name) = ?", workspaceID, strings.TrimSpace(scopeType), strings.TrimSpace(scopeRef), strings.ToLower(strings.TrimSpace(name))).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *BackupPolicyRepository) List(ctx context.Context, filter BackupPolicyListFilter) ([]domain.BackupPolicy, error) {
	query := r.db.WithContext(ctx).Model(&domain.BackupPolicy{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	if v := strings.TrimSpace(filter.ScopeType); v != "" {
		query = query.Where("scope_type = ?", v)
	}
	var items []domain.BackupPolicy
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}
