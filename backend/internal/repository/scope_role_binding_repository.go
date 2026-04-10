package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// ScopeRoleBinding grants a scope role to a user/group.
type ScopeRoleBinding struct {
	ID          uint64    `json:"id" gorm:"primaryKey"`
	SubjectType string    `json:"subjectType" gorm:"size:32;not null;uniqueIndex:uk_scope_binding"`
	SubjectID   uint64    `json:"subjectId" gorm:"not null;uniqueIndex:uk_scope_binding"`
	ScopeType   string    `json:"scopeType" gorm:"size:32;not null;uniqueIndex:uk_scope_binding"`
	ScopeID     uint64    `json:"scopeId" gorm:"not null;uniqueIndex:uk_scope_binding"`
	ScopeRoleID uint64    `json:"scopeRoleId" gorm:"not null;uniqueIndex:uk_scope_binding"`
	GrantedBy   uint64    `json:"grantedBy" gorm:"not null;default:0"`
	CreatedAt   time.Time `json:"createdAt"`
}

func (ScopeRoleBinding) TableName() string { return "scope_role_bindings" }

type ScopeRoleBindingFilter struct {
	SubjectType string
	SubjectID   uint64
	ScopeType   string
	ScopeID     uint64
	Limit       int
	Offset      int
}

type ScopeRoleBindingWithRole struct {
	ID          uint64    `json:"id"`
	SubjectType string    `json:"subjectType"`
	SubjectID   uint64    `json:"subjectId"`
	ScopeType   string    `json:"scopeType"`
	ScopeID     uint64    `json:"scopeId"`
	ScopeRoleID uint64    `json:"scopeRoleId"`
	RoleKey     string    `json:"roleKey"`
	GrantedBy   uint64    `json:"grantedBy"`
	CreatedAt   time.Time `json:"createdAt"`
}

type ScopeRoleBindingRepository struct {
	db *gorm.DB
}

func NewScopeRoleBindingRepository(db *gorm.DB) *ScopeRoleBindingRepository {
	return &ScopeRoleBindingRepository{db: db}
}

func (r *ScopeRoleBindingRepository) Create(ctx context.Context, item *ScopeRoleBinding) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.WithContext(ctx).Create(item).Error
}

func (r *ScopeRoleBindingRepository) List(ctx context.Context, filter ScopeRoleBindingFilter) ([]ScopeRoleBindingWithRole, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	if filter.Limit <= 0 || filter.Limit > 200 {
		filter.Limit = 50
	}
	if filter.Offset < 0 {
		filter.Offset = 0
	}

	query := r.db.WithContext(ctx).
		Table("scope_role_bindings srb").
		Select("srb.id, srb.subject_type, srb.subject_id, srb.scope_type, srb.scope_id, srb.scope_role_id, sr.role_key, srb.granted_by, srb.created_at").
		Joins("LEFT JOIN scope_roles sr ON sr.id = srb.scope_role_id")

	if filter.SubjectType != "" {
		query = query.Where("srb.subject_type = ?", filter.SubjectType)
	}
	if filter.SubjectID != 0 {
		query = query.Where("srb.subject_id = ?", filter.SubjectID)
	}
	if filter.ScopeType != "" {
		query = query.Where("srb.scope_type = ?", filter.ScopeType)
	}
	if filter.ScopeID != 0 {
		query = query.Where("srb.scope_id = ?", filter.ScopeID)
	}

	var items []ScopeRoleBindingWithRole
	err := query.Order("srb.id DESC").Limit(filter.Limit).Offset(filter.Offset).Scan(&items).Error
	return items, err
}
