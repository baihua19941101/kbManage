package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	securityPolicySvc "kbmanage/backend/internal/service/securitypolicy"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterSecurityPolicyRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	policyRepo := repository.NewSecurityPolicyRepository(db)
	assignmentRepo := repository.NewPolicyAssignmentRepository(db)
	hitRepo := repository.NewPolicyHitRepository(db)
	exceptionRepo := repository.NewPolicyExceptionRepository(db)
	scopeAccess := newScopeAccessService(db)
	scopeSvc := securityPolicySvc.NewScopeService(scopeAccess)
	distributionCache := securityPolicySvc.NewDistributionCache(rdb, 0)
	exceptionCache := securityPolicySvc.NewExceptionCache(rdb, 0)
	svc := securityPolicySvc.NewService(policyRepo, assignmentRepo, hitRepo, exceptionRepo, scopeSvc, distributionCache, exceptionCache)
	auditWriter := auditSvc.NewEventWriter(repository.NewAuditRepository(db))
	h := handler.NewSecurityPolicyHandler(svc, auditWriter)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.SecurityPolicy{},
			&domain.PolicyAssignment{},
			&domain.PolicyDistributionTask{},
			&domain.PolicyHitRecord{},
			&domain.PolicyExceptionRequest{},
		)
	}

	policies := group.Group("/security-policies")
	{
		policies.GET("", middleware.RequireSecurityPolicyScopeFromRequest(scopeAccess, middleware.PermissionSecurityPolicyRead), h.ListPolicies)
		policies.POST("", middleware.RequireSecurityPolicyScopeFromRequest(scopeAccess, middleware.PermissionSecurityPolicyManage), h.CreatePolicy)

		policies.GET("/:policyId", middleware.RequireSecurityPolicyEntityScope(scopeAccess, policyRepo, middleware.PermissionSecurityPolicyRead), h.GetPolicy)
		policies.PATCH("/:policyId", middleware.RequireSecurityPolicyEntityScope(scopeAccess, policyRepo, middleware.PermissionSecurityPolicyManage), h.UpdatePolicy)
		policies.PUT("/:policyId", middleware.RequireSecurityPolicyEntityScope(scopeAccess, policyRepo, middleware.PermissionSecurityPolicyManage), h.UpdatePolicy)

		policies.GET("/:policyId/assignments", middleware.RequireSecurityPolicyEntityScope(scopeAccess, policyRepo, middleware.PermissionSecurityPolicyRead), h.ListAssignments)
		policies.POST("/:policyId/assignments", middleware.RequireSecurityPolicyEntityScope(scopeAccess, policyRepo, middleware.PermissionSecurityPolicyEnforce), h.CreateAssignment)

		policies.POST("/:policyId/mode-switch", middleware.RequireSecurityPolicyEntityScope(scopeAccess, policyRepo, middleware.PermissionSecurityPolicyEnforce), h.SwitchPolicyMode)
		policies.GET("/hits", middleware.RequireSecurityPolicyScopeFromRequest(scopeAccess, middleware.PermissionSecurityPolicyRead), h.ListHits)
		policies.POST("/hits/:hitId/exceptions", middleware.RequireSecurityPolicyScopeFromRequest(scopeAccess, middleware.PermissionSecurityPolicyEnforce), h.CreateException)
		policies.PATCH("/hits/:hitId/remediation", middleware.RequireSecurityPolicyScopeFromRequest(scopeAccess, middleware.PermissionSecurityPolicyManage), h.UpdateRemediation)
		policies.GET("/exceptions", middleware.RequireSecurityPolicyScopeFromRequest(scopeAccess, middleware.PermissionSecurityPolicyRead), h.ListExceptions)
		policies.POST("/exceptions/:exceptionId/review", middleware.RequireSecurityPolicyScopeFromRequest(scopeAccess, middleware.PermissionSecurityPolicyManage), h.ReviewException)
	}
}
