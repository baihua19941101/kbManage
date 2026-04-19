package router

import (
	"context"

	"kbmanage/backend/internal/api/handler"
	"kbmanage/backend/internal/domain"
	identityint "kbmanage/backend/internal/integration/identity"
	"kbmanage/backend/internal/repository"
	auditSvc "kbmanage/backend/internal/service/audit"
	identitySvc "kbmanage/backend/internal/service/identitytenancy"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func RegisterIdentityTenancyRoutes(group *gin.RouterGroup, db *gorm.DB, rdb *redis.Client) {
	sourceRepo := repository.NewIdentitySourceRepository(db)
	accountRepo := repository.NewIdentityAccountRepository(db)
	orgRepo := repository.NewOrganizationUnitRepository(db)
	membershipRepo := repository.NewOrganizationMembershipRepository(db)
	mappingRepo := repository.NewTenantScopeMappingRepository(db)
	roleRepo := repository.NewRoleDefinitionRepository(db)
	assignmentRepo := repository.NewRoleAssignmentRepository(db)
	delegationRepo := repository.NewDelegationGrantRepository(db)
	sessionRepo := repository.NewSessionRecordRepository(db)
	riskRepo := repository.NewAccessRiskRepository(db)
	auditRepo := repository.NewIdentityAuditRepository(db)
	bindingRepo := repository.NewScopeRoleBindingRepository(db)
	projectRepo := repository.NewProjectRepository(db)

	svc := identitySvc.NewService(
		sourceRepo,
		accountRepo,
		orgRepo,
		membershipRepo,
		mappingRepo,
		roleRepo,
		assignmentRepo,
		delegationRepo,
		sessionRepo,
		riskRepo,
		auditRepo,
		bindingRepo,
		projectRepo,
		identitySvc.NewSessionCache(rdb),
		identitySvc.NewPermissionCache(rdb),
		identitySvc.NewRevocationCoordinator(rdb),
		identityint.NewStaticProvider(),
		identityint.NewStaticSyncProvider(),
		auditSvc.NewEventWriter(repository.NewAuditRepository(db)),
	)
	h := handler.NewIdentityTenancyHandler(svc)

	if db != nil {
		_ = db.WithContext(context.Background()).AutoMigrate(
			&domain.IdentitySource{},
			&domain.IdentityAccount{},
			&domain.OrganizationUnit{},
			&domain.OrganizationMembership{},
			&domain.TenantScopeMapping{},
			&domain.RoleDefinition{},
			&domain.RoleAssignment{},
			&domain.DelegationGrant{},
			&domain.SessionRecord{},
			&domain.AccessRiskSnapshot{},
			&domain.IdentityGovernanceAuditEvent{},
		)
	}

	identity := group.Group("/identity")
	{
		identity.GET("/sources", h.ListIdentitySources)
		identity.POST("/sources", h.CreateIdentitySource)
		identity.GET("/sources/:sourceId", h.GetIdentitySource)
		identity.POST("/login-mode", h.UpdatePreferredLoginMode)
		identity.GET("/sessions", h.ListSessions)
		identity.POST("/sessions/:sessionId/revoke", h.RevokeSession)
		identity.GET("/organizations", h.ListOrganizationUnits)
		identity.POST("/organizations", h.CreateOrganizationUnit)
		identity.GET("/organizations/:unitId/memberships", h.ListMemberships)
		identity.GET("/organizations/:unitId/mappings", h.ListTenantScopeMappings)
		identity.POST("/organizations/:unitId/mappings", h.CreateTenantScopeMapping)
		identity.GET("/roles", h.ListRoleDefinitions)
		identity.POST("/roles", h.CreateRoleDefinition)
		identity.GET("/assignments", h.ListRoleAssignments)
		identity.POST("/assignments", h.CreateRoleAssignment)
		identity.GET("/delegations", h.ListDelegationGrants)
		identity.POST("/delegations", h.CreateDelegationGrant)
		identity.GET("/access-risks", h.ListAccessRisks)
	}
}
