package repository

import (
	"context"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceScanProfileListFilter struct {
	WorkspaceID  *uint64
	ProjectID    *uint64
	ScopeType    domain.ComplianceScopeType
	ScheduleMode domain.ComplianceScheduleMode
	Status       domain.ComplianceScanProfileStatus
}

type ComplianceScanProfileRepository struct{ db *gorm.DB }

func NewComplianceScanProfileRepository(db *gorm.DB) *ComplianceScanProfileRepository {
	return &ComplianceScanProfileRepository{db: db}
}

func (r *ComplianceScanProfileRepository) Create(ctx context.Context, item *domain.ScanProfile) error {
	if item == nil {
		return errors.New("scan profile is required")
	}
	if r == nil || r.db == nil {
		return errors.New("scan profile repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ComplianceScanProfileRepository) GetByID(ctx context.Context, id uint64) (*domain.ScanProfile, error) {
	if id == 0 {
		return nil, errors.New("profile id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("scan profile repository is not configured")
	}
	var item domain.ScanProfile
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ComplianceScanProfileRepository) List(ctx context.Context, filter ComplianceScanProfileListFilter) ([]domain.ScanProfile, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("scan profile repository is not configured")
	}
	query := r.db.WithContext(ctx).Model(&domain.ScanProfile{})
	if filter.WorkspaceID != nil {
		query = query.Where("workspace_id = ?", *filter.WorkspaceID)
	}
	if filter.ProjectID != nil {
		query = query.Where("project_id = ?", *filter.ProjectID)
	}
	if strings.TrimSpace(string(filter.ScopeType)) != "" {
		query = query.Where("scope_type = ?", filter.ScopeType)
	}
	if strings.TrimSpace(string(filter.ScheduleMode)) != "" {
		query = query.Where("schedule_mode = ?", filter.ScheduleMode)
	}
	if strings.TrimSpace(string(filter.Status)) != "" {
		query = query.Where("status = ?", filter.Status)
	}
	items := make([]domain.ScanProfile, 0)
	if err := query.Order("updated_at DESC, id DESC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ComplianceScanProfileRepository) UpdateFields(ctx context.Context, id uint64, updates map[string]any) error {
	if id == 0 {
		return errors.New("profile id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("scan profile repository is not configured")
	}
	updates["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&domain.ScanProfile{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
