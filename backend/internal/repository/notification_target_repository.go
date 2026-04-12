package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type NotificationTargetRepository struct {
	db *gorm.DB
}

func NewNotificationTargetRepository(db *gorm.DB) *NotificationTargetRepository {
	return &NotificationTargetRepository{db: db}
}

func (r *NotificationTargetRepository) Create(ctx context.Context, item *domain.NotificationTarget) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *NotificationTargetRepository) GetByID(ctx context.Context, id uint64) (*domain.NotificationTarget, error) {
	var item domain.NotificationTarget
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *NotificationTargetRepository) List(ctx context.Context) ([]domain.NotificationTarget, error) {
	var items []domain.NotificationTarget
	err := r.db.WithContext(ctx).Order("id DESC").Find(&items).Error
	return items, err
}

func (r *NotificationTargetRepository) Update(ctx context.Context, item *domain.NotificationTarget) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *NotificationTargetRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.NotificationTarget{}, id).Error
}

