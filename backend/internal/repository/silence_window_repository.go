package repository

import (
	"context"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type SilenceWindowRepository struct {
	db *gorm.DB
}

func NewSilenceWindowRepository(db *gorm.DB) *SilenceWindowRepository {
	return &SilenceWindowRepository{db: db}
}

func (r *SilenceWindowRepository) Create(ctx context.Context, item *domain.SilenceWindow) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *SilenceWindowRepository) GetByID(ctx context.Context, id uint64) (*domain.SilenceWindow, error) {
	var item domain.SilenceWindow
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *SilenceWindowRepository) ListActiveAt(ctx context.Context, at time.Time) ([]domain.SilenceWindow, error) {
	var items []domain.SilenceWindow
	err := r.db.WithContext(ctx).
		Where("starts_at <= ? AND ends_at >= ? AND status = ?", at, at, domain.SilenceWindowStatusActive).
		Order("starts_at DESC").
		Find(&items).Error
	return items, err
}

func (r *SilenceWindowRepository) List(ctx context.Context, status domain.SilenceWindowStatus) ([]domain.SilenceWindow, error) {
	tx := r.db.WithContext(ctx).Model(&domain.SilenceWindow{})
	if status != "" {
		tx = tx.Where("status = ?", status)
	}

	var items []domain.SilenceWindow
	err := tx.Order("starts_at DESC").Find(&items).Error
	return items, err
}

func (r *SilenceWindowRepository) Update(ctx context.Context, item *domain.SilenceWindow) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *SilenceWindowRepository) Delete(ctx context.Context, id uint64) error {
	return r.db.WithContext(ctx).Delete(&domain.SilenceWindow{}, id).Error
}
