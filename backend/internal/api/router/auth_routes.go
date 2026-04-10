package router

import (
	"context"
	"log"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	"kbmanage/backend/internal/service/auth"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterAuthRoutes mounts login/refresh endpoints under /api/v1/auth.
func RegisterAuthRoutes(group *gin.RouterGroup, db *gorm.DB, tokenSvc *auth.TokenService, cfg repository.Config) {
	userRepo := repository.NewUserRepository(db)
	sessionRepo := repository.NewSessionRepository(db)
	loginSvc := auth.NewLoginService(
		userRepo,
		sessionRepo,
		auth.NewPasswordService(0),
		tokenSvc,
		auth.DefaultAdminSeed{
			Enabled:     cfg.AdminSeedEnabled,
			Username:    cfg.AdminSeedUsername,
			Password:    cfg.AdminSeedPassword,
			DisplayName: cfg.AdminSeedDisplayName,
			Email:       cfg.AdminSeedEmail,
		},
	)
	h := handler.NewAuthHandler(loginSvc)

	if db != nil {
		if err := db.WithContext(context.Background()).AutoMigrate(&domain.User{}, &domain.Session{}); err != nil {
			log.Printf("auth auto-migrate failed: %v", err)
		}
	}
	if err := loginSvc.EnsureDefaultAdmin(context.Background()); err != nil {
		log.Printf("default admin seed failed: %v", err)
	}

	group.POST("/auth/login", h.Login)
	group.POST("/auth/refresh", h.Refresh)
}
