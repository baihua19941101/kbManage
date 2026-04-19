package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type AccessRiskListFilter struct {
	SubjectType string
	Severity    string
	Status      string
}

type AccessRiskRepository struct {
	db *gorm.DB
}

func NewAccessRiskRepository(db *gorm.DB) *AccessRiskRepository {
	return &AccessRiskRepository{db: db}
}

func (r *AccessRiskRepository) Create(ctx context.Context, item *domain.AccessRiskSnapshot) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *AccessRiskRepository) Update(ctx context.Context, item *domain.AccessRiskSnapshot) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *AccessRiskRepository) List(ctx context.Context, filter AccessRiskListFilter) ([]domain.AccessRiskSnapshot, error) {
	query := r.db.WithContext(ctx).Model(&domain.AccessRiskSnapshot{})
	if v := strings.TrimSpace(filter.SubjectType); v != "" {
		query = query.Where("subject_type = ?", v)
	}
	if v := strings.TrimSpace(filter.Severity); v != "" {
		query = query.Where("severity = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	var items []domain.AccessRiskSnapshot
	err := query.Order("generated_at DESC, id DESC").Find(&items).Error
	return items, err
}
