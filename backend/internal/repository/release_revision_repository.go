package repository

import (
	"context"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ReleaseRevisionStore interface {
	Create(ctx context.Context, item *domain.ReleaseRevision) error
	GetByID(ctx context.Context, id uint64) (*domain.ReleaseRevision, error)
	ListByDeliveryUnit(ctx context.Context, deliveryUnitID uint64) ([]domain.ReleaseRevision, error)
	GetLatestByDeliveryUnit(ctx context.Context, deliveryUnitID uint64) (*domain.ReleaseRevision, error)
	GetLatestActiveByDeliveryUnit(ctx context.Context, deliveryUnitID uint64) (*domain.ReleaseRevision, error)
	Update(ctx context.Context, item *domain.ReleaseRevision) error
	UpdateStatus(ctx context.Context, id uint64, status domain.ReleaseRevisionStatus, rollbackAvailable bool) error
	MarkOthersHistorical(ctx context.Context, deliveryUnitID uint64, keepID uint64) error
}

type ReleaseRevisionRepository struct {
	db *gorm.DB
}

func NewReleaseRevisionRepository(db *gorm.DB) *ReleaseRevisionRepository {
	return &ReleaseRevisionRepository{db: db}
}

func (r *ReleaseRevisionRepository) Create(ctx context.Context, item *domain.ReleaseRevision) error {
	if item == nil {
		return errors.New("release revision is required")
	}
	if r == nil || r.db == nil {
		return errors.New("release revision repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ReleaseRevisionRepository) GetByID(ctx context.Context, id uint64) (*domain.ReleaseRevision, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("release revision repository is not configured")
	}
	var item domain.ReleaseRevision
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ReleaseRevisionRepository) ListByDeliveryUnit(ctx context.Context, deliveryUnitID uint64) ([]domain.ReleaseRevision, error) {
	if deliveryUnitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("release revision repository is not configured")
	}
	items := make([]domain.ReleaseRevision, 0)
	if err := r.db.WithContext(ctx).
		Where("delivery_unit_id = ?", deliveryUnitID).
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ReleaseRevisionRepository) GetLatestByDeliveryUnit(ctx context.Context, deliveryUnitID uint64) (*domain.ReleaseRevision, error) {
	if deliveryUnitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("release revision repository is not configured")
	}
	var item domain.ReleaseRevision
	if err := r.db.WithContext(ctx).
		Where("delivery_unit_id = ?", deliveryUnitID).
		Order("id DESC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ReleaseRevisionRepository) GetLatestActiveByDeliveryUnit(ctx context.Context, deliveryUnitID uint64) (*domain.ReleaseRevision, error) {
	if deliveryUnitID == 0 {
		return nil, errors.New("delivery unit id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("release revision repository is not configured")
	}
	var item domain.ReleaseRevision
	if err := r.db.WithContext(ctx).
		Where("delivery_unit_id = ? AND status = ?", deliveryUnitID, domain.ReleaseRevisionStatusActive).
		Order("id DESC").
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ReleaseRevisionRepository) Update(ctx context.Context, item *domain.ReleaseRevision) error {
	if item == nil || item.ID == 0 {
		return errors.New("release revision id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("release revision repository is not configured")
	}
	res := r.db.WithContext(ctx).Model(&domain.ReleaseRevision{}).Where("id = ?", item.ID).Updates(item)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ReleaseRevisionRepository) UpdateStatus(
	ctx context.Context,
	id uint64,
	status domain.ReleaseRevisionStatus,
	rollbackAvailable bool,
) error {
	if id == 0 {
		return errors.New("release revision id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("release revision repository is not configured")
	}
	res := r.db.WithContext(ctx).
		Model(&domain.ReleaseRevision{}).
		Where("id = ?", id).
		Updates(map[string]any{
			"status":              status,
			"rollback_available":  rollbackAvailable,
			"release_notes_summary": gorm.Expr("COALESCE(release_notes_summary, '')"),
		})
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ReleaseRevisionRepository) MarkOthersHistorical(ctx context.Context, deliveryUnitID uint64, keepID uint64) error {
	if deliveryUnitID == 0 {
		return errors.New("delivery unit id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("release revision repository is not configured")
	}
	query := r.db.WithContext(ctx).
		Model(&domain.ReleaseRevision{}).
		Where("delivery_unit_id = ?", deliveryUnitID)
	if keepID > 0 {
		query = query.Where("id <> ?", keepID)
	}
	return query.Updates(map[string]any{
		"status":             domain.ReleaseRevisionStatusHistorical,
		"rollback_available": true,
		"release_notes_summary": strings.TrimSpace(""),
	}).Error
}
