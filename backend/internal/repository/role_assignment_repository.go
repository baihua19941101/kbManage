package repository

import (
	"context"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type RoleAssignmentListFilter struct {
	SubjectRef string
	ScopeType  string
	Status     string
}

type RoleAssignmentRepository struct {
	db *gorm.DB
}

func NewRoleAssignmentRepository(db *gorm.DB) *RoleAssignmentRepository {
	return &RoleAssignmentRepository{db: db}
}

func (r *RoleAssignmentRepository) Create(ctx context.Context, item *domain.RoleAssignment) error {
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *RoleAssignmentRepository) Update(ctx context.Context, item *domain.RoleAssignment) error {
	return r.db.WithContext(ctx).Save(item).Error
}

func (r *RoleAssignmentRepository) GetByID(ctx context.Context, id uint64) (*domain.RoleAssignment, error) {
	var item domain.RoleAssignment
	if err := r.db.WithContext(ctx).First(&item, id).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *RoleAssignmentRepository) List(ctx context.Context, filter RoleAssignmentListFilter) ([]domain.RoleAssignment, error) {
	query := r.db.WithContext(ctx).Model(&domain.RoleAssignment{})
	if v := strings.TrimSpace(filter.SubjectRef); v != "" {
		query = query.Where("subject_ref = ?", v)
	}
	if v := strings.TrimSpace(filter.ScopeType); v != "" {
		query = query.Where("scope_type = ?", v)
	}
	if v := strings.TrimSpace(filter.Status); v != "" {
		query = query.Where("status = ?", v)
	}
	var items []domain.RoleAssignment
	err := query.Order("id ASC").Find(&items).Error
	return items, err
}

func (r *RoleAssignmentRepository) ListActiveBySubject(ctx context.Context, subjectType, subjectRef string, now time.Time) ([]domain.RoleAssignment, error) {
	var items []domain.RoleAssignment
	err := r.db.WithContext(ctx).
		Where("subject_type = ? AND subject_ref = ? AND status = ?", subjectType, subjectRef, domain.RoleAssignmentStatusActive).
		Where("(valid_until IS NULL OR valid_until > ?)", now).
		Order("id ASC").
		Find(&items).Error
	return items, err
}
