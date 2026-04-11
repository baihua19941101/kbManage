package repository

import (
	"context"
	"time"

	"gorm.io/gorm"
)

// ScopeRole defines a role in workspace/project scope.
type ScopeRole struct {
	ID           uint64    `json:"id" gorm:"primaryKey"`
	ScopeType    string    `json:"scopeType" gorm:"size:32;not null;uniqueIndex:uk_scope_role"`
	RoleKey      string    `json:"roleKey" gorm:"size:128;not null;uniqueIndex:uk_scope_role"`
	Name         string    `json:"name" gorm:"size:128;not null"`
	Description  string    `json:"description" gorm:"size:512"`
	MetadataJSON string    `json:"metadataJson,omitempty" gorm:"type:json"`
	IsSystem     bool      `json:"isSystem" gorm:"not null;default:false"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

func (ScopeRole) TableName() string { return "scope_roles" }

type ScopeRoleRepository struct {
	db *gorm.DB
}

func NewScopeRoleRepository(db *gorm.DB) *ScopeRoleRepository {
	return &ScopeRoleRepository{db: db}
}

func (r *ScopeRoleRepository) Create(ctx context.Context, role *ScopeRole) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}
	return r.db.WithContext(ctx).Create(role).Error
}

func (r *ScopeRoleRepository) List(ctx context.Context, scopeType string) ([]ScopeRole, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	var items []ScopeRole
	tx := r.db.WithContext(ctx).Order("id ASC")
	if scopeType != "" {
		tx = tx.Where("scope_type = ?", scopeType)
	}
	err := tx.Find(&items).Error
	return items, err
}

func (r *ScopeRoleRepository) GetByScopeAndRoleKey(ctx context.Context, scopeType, roleKey string) (*ScopeRole, error) {
	if r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	var item ScopeRole
	if err := r.db.WithContext(ctx).
		Where("scope_type = ? AND role_key = ?", scopeType, roleKey).
		First(&item).Error; err != nil {
		return nil, err
	}
	return &item, nil
}

func (r *ScopeRoleRepository) EnsureDefaults(ctx context.Context) error {
	if r.db == nil {
		return gorm.ErrInvalidDB
	}

	defaults := []ScopeRole{
		{ScopeType: "workspace", RoleKey: "platform-admin", Name: "Platform Admin", Description: "Workspace scoped super access", MetadataJSON: `{"matrix":"v1","tier":"core"}`, IsSystem: true},
		{ScopeType: "workspace", RoleKey: "ops-operator", Name: "Ops Operator", Description: "Workspace scoped operations access", MetadataJSON: `{"matrix":"v1","tier":"core"}`, IsSystem: true},
		{ScopeType: "workspace", RoleKey: "audit-reader", Name: "Audit Reader", Description: "Workspace scoped audit read access", MetadataJSON: `{"matrix":"v1","tier":"core"}`, IsSystem: true},
		{ScopeType: "workspace", RoleKey: "readonly", Name: "Read Only", Description: "Workspace scoped read-only access", MetadataJSON: `{"matrix":"v1","tier":"core"}`, IsSystem: true},
		{ScopeType: "project", RoleKey: "platform-admin", Name: "Platform Admin", Description: "Project scoped super access", MetadataJSON: `{"matrix":"v1","tier":"core"}`, IsSystem: true},
		{ScopeType: "project", RoleKey: "ops-operator", Name: "Ops Operator", Description: "Project scoped operations access", MetadataJSON: `{"matrix":"v1","tier":"core"}`, IsSystem: true},
		{ScopeType: "project", RoleKey: "audit-reader", Name: "Audit Reader", Description: "Project scoped audit read access", MetadataJSON: `{"matrix":"v1","tier":"core"}`, IsSystem: true},
		{ScopeType: "project", RoleKey: "readonly", Name: "Read Only", Description: "Project scoped read-only access", MetadataJSON: `{"matrix":"v1","tier":"core"}`, IsSystem: true},
		// Backward-compatible aliases kept for current handlers and contract tests.
		{ScopeType: "workspace", RoleKey: "workspace-owner", Name: "Workspace Owner", Description: "Workspace full access", IsSystem: true},
		{ScopeType: "workspace", RoleKey: "workspace-viewer", Name: "Workspace Viewer", Description: "Workspace read-only", IsSystem: true},
		{ScopeType: "project", RoleKey: "project-owner", Name: "Project Owner", Description: "Project full access", IsSystem: true},
		{ScopeType: "project", RoleKey: "project-viewer", Name: "Project Viewer", Description: "Project read-only", IsSystem: true},
	}

	for _, item := range defaults {
		copy := item
		if err := r.db.WithContext(ctx).
			Where("scope_type = ? AND role_key = ?", copy.ScopeType, copy.RoleKey).
			FirstOrCreate(&copy).Error; err != nil {
			return err
		}
	}
	return nil
}
