package repository

import (
	"context"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type ComplianceFindingListFilter struct {
	ScanExecutionID *uint64
	WorkspaceID     *uint64
	ProjectID       *uint64
	Result          domain.ComplianceFindingResult
	RiskLevel       domain.ComplianceRiskLevel
}

type ComplianceFindingRepository struct{ db *gorm.DB }

type ComplianceEvidenceRepository struct{ db *gorm.DB }

func NewComplianceFindingRepository(db *gorm.DB) *ComplianceFindingRepository {
	return &ComplianceFindingRepository{db: db}
}
func NewComplianceEvidenceRepository(db *gorm.DB) *ComplianceEvidenceRepository {
	return &ComplianceEvidenceRepository{db: db}
}

func (r *ComplianceFindingRepository) ReplaceByScanExecution(ctx context.Context, scanID uint64, items []domain.ComplianceFinding) error {
	if scanID == 0 {
		return errors.New("scan id is required")
	}
	if r == nil || r.db == nil {
		return errors.New("finding repository is not configured")
	}
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		var findingIDs []uint64
		if err := tx.Model(&domain.ComplianceFinding{}).Where("scan_execution_id = ?", scanID).Pluck("id", &findingIDs).Error; err != nil {
			return err
		}
		if len(findingIDs) > 0 {
			if err := tx.Where("finding_id IN ?", findingIDs).Delete(&domain.EvidenceRecord{}).Error; err != nil {
				return err
			}
		}
		if err := tx.Where("scan_execution_id = ?", scanID).Delete(&domain.ComplianceFinding{}).Error; err != nil {
			return err
		}
		if len(items) == 0 {
			return nil
		}
		return tx.Create(&items).Error
	})
}

func (r *ComplianceFindingRepository) List(ctx context.Context, filter ComplianceFindingListFilter) ([]domain.ComplianceFinding, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("finding repository is not configured")
	}
	query := r.db.WithContext(ctx).Table(domain.ComplianceFinding{}.TableName() + " AS f").Select("f.*")
	if filter.WorkspaceID != nil || filter.ProjectID != nil {
		query = query.Joins("JOIN compliance_scan_executions se ON se.id = f.scan_execution_id")
		if filter.WorkspaceID != nil {
			query = query.Where("se.workspace_id = ?", *filter.WorkspaceID)
		}
		if filter.ProjectID != nil {
			query = query.Where("se.project_id = ?", *filter.ProjectID)
		}
	}
	if filter.ScanExecutionID != nil {
		query = query.Where("f.scan_execution_id = ?", *filter.ScanExecutionID)
	}
	if strings.TrimSpace(string(filter.Result)) != "" {
		query = query.Where("f.result = ?", filter.Result)
	}
	if strings.TrimSpace(string(filter.RiskLevel)) != "" {
		query = query.Where("f.risk_level = ?", filter.RiskLevel)
	}
	items := make([]domain.ComplianceFinding, 0)
	if err := query.Order("f.id DESC").Scan(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *ComplianceFindingRepository) GetByID(ctx context.Context, id uint64) (*domain.ComplianceFinding, error) {
	if id == 0 {
		return nil, errors.New("finding id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("finding repository is not configured")
	}
	var item domain.ComplianceFinding
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ComplianceEvidenceRepository) CreateBatch(ctx context.Context, items []domain.EvidenceRecord) error {
	if r == nil || r.db == nil {
		return errors.New("evidence repository is not configured")
	}
	if len(items) == 0 {
		return nil
	}
	return r.db.WithContext(ctx).Create(&items).Error
}

func (r *ComplianceEvidenceRepository) ListByFindingID(ctx context.Context, findingID uint64) ([]domain.EvidenceRecord, error) {
	if findingID == 0 {
		return nil, errors.New("finding id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("evidence repository is not configured")
	}
	items := make([]domain.EvidenceRecord, 0)
	if err := r.db.WithContext(ctx).Where("finding_id = ?", findingID).Order("id ASC").Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}
