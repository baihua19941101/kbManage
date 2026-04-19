package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	executorProvider "kbmanage/backend/internal/integration/backuprestore/executor"
	validatorProvider "kbmanage/backend/internal/integration/backuprestore/validator"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	backupRestoreSvc "kbmanage/backend/internal/service/backuprestore"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterBackupRestoreRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	policyRepo := repository.NewBackupPolicyRepository(db)
	restorePointRepo := repository.NewRestorePointRepository(db)
	restoreJobRepo := repository.NewRestoreJobRepository(db)
	migrationRepo := repository.NewMigrationPlanRepository(db)
	drillPlanRepo := repository.NewDRDrillPlanRepository(db)
	drillRecordRepo := repository.NewDRDrillRecordRepository(db)
	drillReportRepo := repository.NewDRDrillReportRepository(db)
	auditRepo := repository.NewBackupAuditRepository(db)
	bindingRepo := repository.NewScopeRoleBindingRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	svc := backupRestoreSvc.NewService(
		policyRepo,
		restorePointRepo,
		restoreJobRepo,
		migrationRepo,
		drillPlanRepo,
		drillRecordRepo,
		drillReportRepo,
		auditRepo,
		bindingRepo,
		projectRepo,
		backupRestoreSvc.NewProgressCache(rdb),
		backupRestoreSvc.NewPrecheckCache(rdb),
		backupRestoreSvc.NewOperationLock(rdb),
		executorProvider.NewStaticProvider(),
		validatorProvider.NewStaticProvider(),
		auditSvc.NewEventWriter(repository.NewAuditRepository(db)),
	)
	h := handler.NewBackupRestoreHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.BackupPolicy{},
			&domain.RestorePoint{},
			&domain.RestoreJob{},
			&domain.MigrationPlan{},
			&domain.DRDrillPlan{},
			&domain.DRDrillRecord{},
			&domain.DRDrillReport{},
			&domain.BackupAuditEvent{},
		)
	}

	backup := group.Group("/backup-restore")
	{
		backup.GET("/policies", h.ListPolicies)
		backup.POST("/policies", h.CreatePolicy)
		backup.POST("/policies/:policyId/run", h.RunPolicy)
		backup.GET("/restore-points", h.ListRestorePoints)
		backup.GET("/restore-points/:restorePointId", h.GetRestorePoint)
		backup.POST("/restore-jobs", h.CreateRestoreJob)
		backup.GET("/restore-jobs", h.ListRestoreJobs)
		backup.POST("/restore-jobs/:jobId/validate", h.ValidateRestoreJob)
		backup.POST("/migrations", h.CreateMigrationPlan)
		backup.GET("/drills/plans", h.ListDrillPlans)
		backup.POST("/drills/plans", h.CreateDrillPlan)
		backup.POST("/drills/plans/:planId/run", h.RunDrillPlan)
		backup.GET("/drills/records/:recordId", h.GetDrillRecord)
		backup.POST("/drills/records/:recordId/report", h.GenerateDrillReport)
	}
}
