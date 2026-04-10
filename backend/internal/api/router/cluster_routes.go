package router

import (
	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/repository"
	clusterSvc "kbmanage/backend/internal/service/cluster"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterClusterRoutes mounts US1 cluster/resource APIs under the provided group.
func RegisterClusterRoutes(group *gin.RouterGroup, db *gorm.DB) {
	clusterRepo := repository.NewClusterRepository(db)
	credRepo := repository.NewClusterCredentialRepository(db)
	resourceRepo := repository.NewResourceInventoryRepository(db)

	svc := clusterSvc.NewService(clusterRepo, credRepo, resourceRepo, clusterSvc.NewCredentialCipher())
	clusterHandler := handler.NewClusterHandler(svc)
	resourceHandler := handler.NewResourceHandler(svc)

	group.POST("/clusters", clusterHandler.Register)
	group.POST("/clusters/:id/connectivity", clusterHandler.ValidateConnectivity)
	group.GET("/clusters/:id/health-summary", clusterHandler.HealthSummary)

	group.GET("/resources", resourceHandler.List)
}
