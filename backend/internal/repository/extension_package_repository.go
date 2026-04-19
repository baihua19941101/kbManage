package repository

import (
	"context"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ExtensionPackageListFilter struct {
	Type    string
	Status  string
	Keyword string
}

type ExtensionPackageRepository struct{ db *gorm.DB }

func NewExtensionPackageRepository(db *gorm.DB) *ExtensionPackageRepository {
	return &ExtensionPackageRepository{db: db}
}

func (r *ExtensionPackageRepository) Create(ctx context.Context, item *domain.ExtensionPackage) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ExtensionPackageRepository) Update(ctx context.Context, item *domain.ExtensionPackage) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *ExtensionPackageRepository) GetByID(ctx context.Context, id uint64) (*domain.ExtensionPackage, error) {
	var item domain.ExtensionPackage
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ExtensionPackageRepository) FindByName(ctx context.Context, name string) (*domain.ExtensionPackage, error) {
	var item domain.ExtensionPackage
	if err := r.db.WithContext(ctx).Where("name = ?", strings.TrimSpace(name)).First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ExtensionPackageRepository) List(ctx context.Context, filter ExtensionPackageListFilter) ([]domain.ExtensionPackage, error) {
	query := r.db.WithContext(ctx).Model(&domain.ExtensionPackage{})
	if v := strings.TrimSpace(filter.Type); v != "" {
		query = query.Where("extension_type = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	if v := strings.TrimSpace(filter.Keyword); v != "" {
		query = query.Where("name LIKE ?", "%"+v+"%")
	}
	var items []domain.ExtensionPackage
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}
