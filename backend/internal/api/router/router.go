package router

import (
	"net/http"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/repository"
	"kbmanage/backend/internal/service/auth"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, rdb *redis.Client, cfg repository.Config) *gin.Engine {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.CORS(cfg))
	engine.Use(middleware.RequestID())
	engine.Use(middleware.ErrorHandler())

	tokenSvc := auth.NewTokenService(cfg.JWTSecret, cfg.AccessTokenTTL, cfg.RefreshTokenTTL)

	engine.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := engine.Group("/api/v1")
	{
		v1.GET("/ping", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{"message": "pong"})
		})
		RegisterAuthRoutes(v1, db, tokenSvc, cfg)

		authed := v1.Group("/")
		authed.Use(middleware.AuthRequired(tokenSvc))
		authed.GET("/me", func(c *gin.Context) {
			c.JSON(http.StatusOK, gin.H{
				"user_id":    c.GetUint64(middleware.UserIDKey),
				"request_id": c.GetString(middleware.RequestIDKey),
			})
		})
		RegisterAccessRoutes(authed, db)
		RegisterOperationRoutes(authed, db, rdb)
		RegisterAuditRoutes(authed, db)
		RegisterObservabilityRoutes(authed, db, nil)
		RegisterWorkloadOpsRoutes(authed, db, rdb)
		RegisterGitOpsRoutes(authed, db, rdb)
		RegisterSecurityPolicyRoutes(authed, db, rdb)
		RegisterComplianceRoutes(authed, db, rdb)
		RegisterClusterLifecycleRoutes(authed, db, rdb)
		RegisterBackupRestoreRoutes(authed, db, rdb)
		RegisterIdentityTenancyRoutes(authed, db, rdb)
		RegisterMarketplaceRoutes(authed, db, rdb)

		if db != nil {
			RegisterClusterRoutes(authed, db)
		}
	}

	return engine
}
