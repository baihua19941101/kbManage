package repository

import (
	"context"
	"errors"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type PolicyHitQuery struct {
	PolicyID          *uint64
	WorkspaceID       *uint64
	ProjectID         *uint64
	ClusterID         *uint64
	Namespace         string
	EnforcementMode   domain.PolicyEnforcementMode
	RiskLevel         domain.PolicyRiskLevel
	RemediationStatus domain.PolicyRemediationStatus
	DetectedFrom      *time.Time
	DetectedTo        *time.Time
	Limit             int
}

type PolicyHitRepository struct {
	db *gorm.DB
}

func NewPolicyHitRepository(db *gorm.DB) *PolicyHitRepository {
	return &PolicyHitRepository{db: db}
}

func (r *PolicyHitRepository) Create(ctx context.Context, item *domain.PolicyHitRecord) error {
	if item == nil {
		return errors.New("policy hit record is required")
	}
	if r == nil || r.db == nil {
		return errors.New("policy hit repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PolicyHitRepository) GetByID(ctx context.Context, hitID uint64) (*domain.PolicyHitRecord, error) {
	if hitID == 0 {
		return nil, errors.New("hit id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("policy hit repository is not configured")
	}
	var item domain.PolicyHitRecord
	if err := r.db.WithContext(ctx).First(&item, hitID).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *PolicyHitRepository) List(ctx context.Context, query PolicyHitQuery) ([]domain.PolicyHitRecord, error) {
	if r == nil || r.db == nil {
		return nil, errors.New("policy hit repository is not configured")
	}
	dbQuery := r.db.WithContext(ctx).
		Model(&domain.PolicyHitRecord{}).
		Joins("LEFT JOIN security_policies ON security_policies.id = policy_hit_records.policy_id").
		Joins("LEFT JOIN policy_assignments ON policy_assignments.id = policy_hit_records.assignment_id")
	if query.WorkspaceID != nil {
		dbQuery = dbQuery.Where("security_policies.workspace_id = ?", *query.WorkspaceID)
	}
	if query.ProjectID != nil {
		dbQuery = dbQuery.Where("security_policies.project_id = ?", *query.ProjectID)
	}
	if query.PolicyID != nil {
		dbQuery = dbQuery.Where("policy_hit_records.policy_id = ?", *query.PolicyID)
	}
	if query.ClusterID != nil {
		dbQuery = dbQuery.Where("policy_hit_records.cluster_id = ?", *query.ClusterID)
	}
	if query.Namespace != "" {
		dbQuery = dbQuery.Where("policy_hit_records.namespace = ?", query.Namespace)
	}
	if query.EnforcementMode != "" {
		dbQuery = dbQuery.Where("policy_assignments.enforcement_mode = ?", query.EnforcementMode)
	}
	if query.RiskLevel != "" {
		dbQuery = dbQuery.Where("policy_hit_records.risk_level = ?", query.RiskLevel)
	}
	if query.RemediationStatus != "" {
		dbQuery = dbQuery.Where("policy_hit_records.remediation_status = ?", query.RemediationStatus)
	}
	if query.DetectedFrom != nil {
		dbQuery = dbQuery.Where("policy_hit_records.detected_at >= ?", *query.DetectedFrom)
	}
	if query.DetectedTo != nil {
		dbQuery = dbQuery.Where("policy_hit_records.detected_at <= ?", *query.DetectedTo)
	}
	if query.Limit <= 0 || query.Limit > 500 {
		query.Limit = 100
	}
	items := make([]domain.PolicyHitRecord, 0)
	if err := dbQuery.Order("policy_hit_records.detected_at DESC, policy_hit_records.id DESC").Limit(query.Limit).Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PolicyHitRepository) UpdateFields(ctx context.Context, hitID uint64, updates map[string]any) error {
	if hitID == 0 {
		return errors.New("hit id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("policy hit repository is not configured")
	}
	return r.db.WithContext(ctx).Model(&domain.PolicyHitRecord{}).Where("id = ?", hitID).Updates(updates).Error
}
