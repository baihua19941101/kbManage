package repository

import (
	"context"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"

	"gorm.io/gorm"
)

type PlatformRoleRepository struct {
	db *gorm.DB
}

func NewPlatformRoleRepository(db *gorm.DB) *PlatformRoleRepository {
	return &PlatformRoleRepository{db: db}
}

func (r *PlatformRoleRepository) ListByUserID(ctx context.Context, userID uint64) ([]domain.PlatformRole, error) {
	var roles []domain.PlatformRole
	err := r.db.WithContext(ctx).
		Table("platform_roles pr").
		Select("pr.*").
		Joins("JOIN user_platform_roles upr ON upr.role_id = pr.id").
		Where("upr.user_id = ?", userID).
		Scan(&roles).Error
	return roles, err
}

func (r *PlatformRoleRepository) BindRoleToUser(ctx context.Context, userID, roleID uint64) error {
	binding := domain.UserPlatformRole{UserID: userID, RoleID: roleID}
	return r.db.WithContext(ctx).Create(&binding).Error
}

func (r *PlatformRoleRepository) EnsureDefaults(ctx context.Context) error {
	if r == nil || r.db == nil {
		return gorm.ErrInvalidDB
	}

	defaults := []domain.PlatformRole{
		{Name: "platform-admin", Description: "Platform administrator"},
		{Name: "ops-operator", Description: "Operations operator"},
		{Name: "audit-reader", Description: "Audit reader"},
		{Name: "readonly", Description: "Read only user"},
	}

	for i := range defaults {
		item := defaults[i]
		if strings.TrimSpace(item.Name) == "" {
			continue
		}
		if err := r.db.WithContext(ctx).
			Where("name = ?", item.Name).
			FirstOrCreate(&item).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *PlatformRoleRepository) GetByName(ctx context.Context, roleName string) (*domain.PlatformRole, error) {
	if r == nil || r.db == nil {
		return nil, gorm.ErrInvalidDB
	}

	name := strings.TrimSpace(roleName)
	if name == "" {
		return nil, gorm.ErrRecordNotFound
	}

	var role domain.PlatformRole
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *PlatformRoleRepository) EnsureUserRoleByName(ctx context.Context, userID uint64, roleName string) error {
	if r == nil || r.db == nil {
		return gorm.ErrInvalidDB
	}
	if userID == 0 {
		return gorm.ErrInvalidData
	}

	role, err := r.GetByName(ctx, roleName)
	if err != nil {
		return err
	}

	var binding domain.UserPlatformRole
	err = r.db.WithContext(ctx).
		Where("user_id = ? AND role_id = ?", userID, role.ID).
		First(&binding).Error
	switch {
	case err == nil:
		return nil
	case errors.Is(err, gorm.ErrRecordNotFound):
		return r.BindRoleToUser(ctx, userID, role.ID)
	default:
		return err
	}
}
