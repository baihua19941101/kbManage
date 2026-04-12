package testutil

import (
	"fmt"
	"testing"
	"time"

	"kbmanage/backend/internal/api/router"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type App struct {
	Router *gin.Engine
	DB     *gorm.DB
	Config repository.Config
}

func NewApp(t *testing.T) *App {
	t.Helper()

	gin.SetMode(gin.TestMode)

	cfg := repository.Config{
		JWTSecret:         "test-secret",
		AccessTokenTTL:    15 * time.Minute,
		RefreshTokenTTL:   24 * time.Hour,
		AdminSeedEnabled:  false,
		CORSAllowOrigins:  []string{"*"},
		CORSAllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		CORSAllowHeaders:  []string{"Authorization", "Content-Type", "X-Request-Id"},
		CORSExposeHeaders: []string{"X-Request-Id"},
	}

	dsn := fmt.Sprintf("file:kbmanage-test-%d?mode=memory&cache=shared&_busy_timeout=5000", time.Now().UnixNano())
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("open sqlite db failed: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		t.Fatalf("open sql db handle failed: %v", err)
	}
	t.Cleanup(func() {
		_ = sqlDB.Close()
	})

	if err := db.AutoMigrate(
		&domain.User{},
		&domain.Session{},
		&domain.Workspace{},
		&domain.Project{},
		&repository.WorkspaceClusterBinding{},
		&repository.ProjectClusterBinding{},
		&repository.ScopeRole{},
		&repository.ScopeRoleBinding{},
		&domain.OperationRequest{},
		&domain.AuditEvent{},
		&domain.Cluster{},
		&repository.ClusterCredential{},
		&repository.ResourceInventory{},
		&domain.ObservabilityDataSource{},
		&domain.AlertRule{},
		&domain.NotificationTarget{},
		&domain.SilenceWindow{},
		&domain.AlertIncidentSnapshot{},
		&domain.AlertHandlingRecord{},
	); err != nil {
		t.Fatalf("auto-migrate test schema failed: %v", err)
	}

	return &App{
		Router: router.NewRouter(db, nil, cfg),
		DB:     db,
		Config: cfg,
	}
}
