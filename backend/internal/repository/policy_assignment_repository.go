package repository

import (
	"context"
	"errors"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type PolicyAssignmentRepository struct {
	db *gorm.DB
}

func NewPolicyAssignmentRepository(db *gorm.DB) *PolicyAssignmentRepository {
	return &PolicyAssignmentRepository{db: db}
}

func (r *PolicyAssignmentRepository) Create(ctx context.Context, item *domain.PolicyAssignment) error {
	if item == nil {
		return errors.New("policy assignment is required")
	}
	if r == nil || r.db == nil {
		return errors.New("policy assignment repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *PolicyAssignmentRepository) ListByPolicyID(ctx context.Context, policyID uint64) ([]domain.PolicyAssignment, error) {
	if policyID == 0 {
		return nil, errors.New("policy id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("policy assignment repository is not configured")
	}
	items := make([]domain.PolicyAssignment, 0)
	if err := r.db.WithContext(ctx).
		Where("policy_id = ?", policyID).
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PolicyAssignmentRepository) ListByPolicyAndIDs(ctx context.Context, policyID uint64, assignmentIDs []uint64) ([]domain.PolicyAssignment, error) {
	if policyID == 0 {
		return nil, errors.New("policy id is required")
	}
	if r == nil || r.db == nil {
		return nil, errors.New("policy assignment repository is not configured")
	}
	if len(assignmentIDs) == 0 {
		return []domain.PolicyAssignment{}, nil
	}
	items := make([]domain.PolicyAssignment, 0, len(assignmentIDs))
	if err := r.db.WithContext(ctx).
		Where("policy_id = ? AND id IN ?", policyID, assignmentIDs).
		Order("id DESC").
		Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *PolicyAssignmentRepository) BulkUpdateEnforcementMode(
	ctx context.Context,
	assignmentIDs []uint64,
	mode domain.PolicyEnforcementMode,
	rolloutStage domain.PolicyRolloutStage,
	updatedBy *uint64,
) error {
	if r == nil || r.db == nil {
		return errors.New("policy assignment repository is not configured")
	}
	if len(assignmentIDs) == 0 {
		return nil
	}
	updates := map[string]any{
		"enforcement_mode": mode,
		"updated_at":       time.Now(),
	}
	if rolloutStage != "" {
		updates["rollout_stage"] = rolloutStage
	}
	if updatedBy != nil {
		updates["updated_by"] = updatedBy
	}
	return r.db.WithContext(ctx).
		Model(&domain.PolicyAssignment{}).
		Where("id IN ?", assignmentIDs).
		Updates(updates).Error
}

func (r *PolicyAssignmentRepository) UpdateFields(ctx context.Context, assignmentID uint64, updates map[string]any) error {
	if assignmentID == 0 {
		return errors.New("assignment id is required")
	}
	if len(updates) == 0 {
		return nil
	}
	if r == nil || r.db == nil {
		return errors.New("policy assignment repository is not configured")
	}
	updates["updated_at"] = time.Now()
	res := r.db.WithContext(ctx).Model(&domain.PolicyAssignment{}).Where("id = ?", assignmentID).Updates(updates)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *PolicyAssignmentRepository) CreateDistributionTask(ctx context.Context, item *domain.PolicyDistributionTask) error {
	if item == nil {
		return errors.New("policy distribution task is required")
	}
	if r == nil || r.db == nil {
		return errors.New("policy assignment repository is not configured")
	}
	return r.db.WithContext(ctx).Create(item).Error
}
