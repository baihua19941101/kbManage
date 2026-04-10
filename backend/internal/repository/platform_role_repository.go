package repository

import (
	"context"

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
