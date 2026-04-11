package router

import (
	"context"
	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/kube/adapter"
	kubeclient "kbmanage/backend/internal/kube/client"
	"kbmanage/backend/internal/repository"
	clusterSvc "kbmanage/backend/internal/service/cluster"
	"kbmanage/backend/internal/worker"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// RegisterClusterRoutes mounts US1 cluster/resource APIs under the provided group.
func RegisterClusterRoutes(group *gin.RouterGroup, db *gorm.DB) {
	clusterRepo := repository.NewClusterRepository(db)
	credRepo := repository.NewClusterCredentialRepository(db)
	resourceRepo := repository.NewResourceInventoryRepository(db)
	clientManager := kubeclient.NewManager()
	cipher := clusterSvc.NewCredentialCipher()
	resourceIndexer := adapter.NewResourceIndexer(clusterRepo, credRepo, resourceRepo, cipher, clientManager)
	syncWorker := worker.NewClusterSyncWorker(resourceIndexer, 128, 20*time.Second)
	syncWorker.Start(context.Background())

	svc := clusterSvc.NewService(clusterRepo, credRepo, resourceRepo, cipher, resourceIndexer, syncWorker, clientManager)
	scopeAccess := newScopeAccessService(db)
	clusterHandler := handler.NewClusterHandler(svc, scopeAccess)
	resourceHandler := handler.NewResourceHandler(svc, scopeAccess)

	group.GET("/clusters", clusterHandler.List)
	group.POST("/clusters", clusterHandler.Register)
	group.POST("/clusters/:id/connectivity", clusterHandler.ValidateConnectivity)
	group.POST("/clusters/:id/sync", clusterHandler.SyncResources)
	group.GET("/clusters/:id/health-summary", clusterHandler.HealthSummary)
	group.GET("/clusters/:id/resources", resourceHandler.List)
	group.GET("/clusters/:id/resources/detail", resourceHandler.Detail)

	group.GET("/resources", resourceHandler.List)
}
