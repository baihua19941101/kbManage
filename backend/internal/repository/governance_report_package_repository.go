package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type GovernanceReportPackageRepository struct{ db *gorm.DB }

func NewGovernanceReportPackageRepository(db *gorm.DB) *GovernanceReportPackageRepository {
	return &GovernanceReportPackageRepository{db: db}
}

func (r *GovernanceReportPackageRepository) Create(ctx context.Context, item *domain.GovernanceReportPackage) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *GovernanceReportPackageRepository) Update(ctx context.Context, item *domain.GovernanceReportPackage) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *GovernanceReportPackageRepository) List(ctx context.Context) ([]domain.GovernanceReportPackage, error) {
	var items []domain.GovernanceReportPackage
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}

func (r *GovernanceReportPackageRepository) GetByID(ctx context.Context, id uint64) (*domain.GovernanceReportPackage, error) {
	var item domain.GovernanceReportPackage
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
