package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type RoleDefinitionListFilter struct {
	RoleLevel string
	Status    string
}

type RoleDefinitionRepository struct {
	db *gorm.DB
}

func NewRoleDefinitionRepository(db *gorm.DB) *RoleDefinitionRepository {
	return &RoleDefinitionRepository{db: db}
}

func (r *RoleDefinitionRepository) Create(ctx context.Context, item *domain.RoleDefinition) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RoleDefinitionRepository) Update(ctx context.Context, item *domain.RoleDefinition) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *RoleDefinitionRepository) GetByID(ctx context.Context, id uint64) (*domain.RoleDefinition, error) {
	var item domain.RoleDefinition
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RoleDefinitionRepository) FindByLevelName(ctx context.Context, roleLevel, name string) (*domain.RoleDefinition, error) {
	var item domain.RoleDefinition
	err := r.db.WithContext(ctx).
		Where("role_level = ? AND lower(name) = ?", roleLevel, strings.ToLower(strings.TrimSpace(name))).
		First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RoleDefinitionRepository) List(ctx context.Context, filter RoleDefinitionListFilter) ([]domain.RoleDefinition, error) {
	query := r.db.WithContext(ctx).Model(&domain.RoleDefinition{})
	if v := strings.TrimSpace(filter.RoleLevel); v != "" {
		query = query.Where("role_level = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	var items []domain.RoleDefinition
	err := query.Order("id ASC").Find(&items).Error
	return items, err
}
