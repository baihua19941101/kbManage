package repository

import (
	"context"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ScaleEvidenceListFilter struct {
	WorkspaceIDs []uint64
	ProjectIDs   []uint64
	EvidenceType string
}

type ScaleEvidenceRepository struct{ db *gorm.DB }

func NewScaleEvidenceRepository(db *gorm.DB) *ScaleEvidenceRepository {
	return &ScaleEvidenceRepository{db: db}
}

func (r *ScaleEvidenceRepository) Create(ctx context.Context, item *domain.ScaleEvidence) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ScaleEvidenceRepository) List(ctx context.Context, filter ScaleEvidenceListFilter) ([]domain.ScaleEvidence, error) {
	query := r.db.WithContext(ctx).Model(&domain.ScaleEvidence{})
	if len(filter.WorkspaceIDs) > 0 {
		query = query.Where("workspace_id IN ?", filter.WorkspaceIDs)
	}
	if len(filter.ProjectIDs) > 0 {
		query = query.Where("project_id IN ?", filter.ProjectIDs)
	}
	if filter.EvidenceType != "" {
		query = query.Where("evidence_type = ?", filter.EvidenceType)
	}
	var items []domain.ScaleEvidence
	err := query.Order("id DESC").Find(&items).Error
	return items, err
}

func (r *ScaleEvidenceRepository) FindLatestByScope(ctx context.Context, workspaceID, projectID uint64) (*domain.ScaleEvidence, error) {
	var item domain.ScaleEvidence
	query := r.db.WithContext(ctx).Model(&domain.ScaleEvidence{}).Where("workspace_id = ?", workspaceID)
	if projectID != 0 {
		query = query.Where("project_id = ?", projectID)
	}
	if err := query.Order("captured_at DESC, id DESC").First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}
