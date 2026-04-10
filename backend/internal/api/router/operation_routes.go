package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	operationSvc "kbmanage/backend/internal/service/operation"
	"kbmanage/backend/internal/worker"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// RegisterOperationRoutes mounts US3 operation APIs.
func RegisterOperationRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	repo := repository.NewOperationRepository(db)
	idempotencySvc := operationSvc.NewIdempotencyService(rdb)
	queueSvc := operationSvc.NewQueueService(rdb)
	svc := operationSvc.NewService(repo, idempotencySvc, queueSvc)
	h := handler.NewOperationHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(&domain.OperationRequest{})
	}

	operationWorker := worker.NewOperationWorker(repo, queueSvc)
	operationWorker.Start(context.Background())

	group.POST("/operations", h.Create)
	group.GET("/operations/:operationId", h.GetByID)
}
