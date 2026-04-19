package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/service/identitytenancy"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type IdentityTenancyHandler struct {
	svc *identitytenancy.Service
}

func NewIdentityTenancyHandler(svc *identitytenancy.Service) *IdentityTenancyHandler {
	return &IdentityTenancyHandler{svc: svc}
}

func (h *IdentityTenancyHandler) ListIdentitySources(c *gin.Context) {
	items, err := h.svc.ListIdentitySources(c.Request.Context(), c.GetUint64(middleware.UserIDKey), identitytenancy.IdentitySourceListFilter{
		SourceType: strings.TrimSpace(c.Query("sourceType")),
		Status:     strings.TrimSpace(c.Query("status")),
	})
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	loginModes, err := h.svc.ListAvailableLoginModes(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items, "loginModes": loginModes})
}

func (h *IdentityTenancyHandler) CreateIdentitySource(c *gin.Context) {
	var req identitytenancy.CreateIdentitySourceInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateIdentitySource(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *IdentityTenancyHandler) GetIdentitySource(c *gin.Context) {
	id, ok := parseIdentityUint64Param(c, "sourceId")
	if !ok {
		return
	}
	item, err := h.svc.GetIdentitySource(c.Request.Context(), c.GetUint64(middleware.UserIDKey), id)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *IdentityTenancyHandler) UpdatePreferredLoginMode(c *gin.Context) {
	var req struct {
		LoginMode string `json:"loginMode"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	loginMode, err := h.svc.UpdatePreferredLoginMode(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req.LoginMode)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"loginMode": loginMode})
}

func (h *IdentityTenancyHandler) ListSessions(c *gin.Context) {
	items, err := h.svc.ListSessions(c.Request.Context(), c.GetUint64(middleware.UserIDKey), identitytenancy.SessionListFilter{
		Status:    strings.TrimSpace(c.Query("status")),
		RiskLevel: strings.TrimSpace(c.Query("riskLevel")),
	})
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *IdentityTenancyHandler) RevokeSession(c *gin.Context) {
	sessionID, ok := parseIdentityUint64Param(c, "sessionId")
	if !ok {
		return
	}
	item, err := h.svc.RevokeSession(c.Request.Context(), c.GetUint64(middleware.UserIDKey), sessionID)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, item)
}

func (h *IdentityTenancyHandler) ListOrganizationUnits(c *gin.Context) {
	parentID, _ := parseOptionalIdentityUint64(c.Query("parentUnitId"))
	items, err := h.svc.ListOrganizationUnits(c.Request.Context(), c.GetUint64(middleware.UserIDKey), identitytenancy.OrganizationUnitListFilter{
		UnitType:     strings.TrimSpace(c.Query("unitType")),
		ParentUnitID: parentID,
	})
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *IdentityTenancyHandler) CreateOrganizationUnit(c *gin.Context) {
	var req identitytenancy.CreateOrganizationUnitInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateOrganizationUnit(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *IdentityTenancyHandler) ListMemberships(c *gin.Context) {
	unitID, ok := parseIdentityUint64Param(c, "unitId")
	if !ok {
		return
	}
	items, err := h.svc.ListMemberships(c.Request.Context(), c.GetUint64(middleware.UserIDKey), unitID)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *IdentityTenancyHandler) ListTenantScopeMappings(c *gin.Context) {
	unitID, ok := parseIdentityUint64Param(c, "unitId")
	if !ok {
		return
	}
	items, err := h.svc.ListTenantScopeMappings(c.Request.Context(), c.GetUint64(middleware.UserIDKey), unitID)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *IdentityTenancyHandler) CreateTenantScopeMapping(c *gin.Context) {
	unitID, ok := parseIdentityUint64Param(c, "unitId")
	if !ok {
		return
	}
	var req identitytenancy.CreateTenantScopeMappingInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateTenantScopeMapping(c.Request.Context(), c.GetUint64(middleware.UserIDKey), unitID, req)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *IdentityTenancyHandler) ListRoleDefinitions(c *gin.Context) {
	items, err := h.svc.ListRoleDefinitions(c.Request.Context(), c.GetUint64(middleware.UserIDKey), identitytenancy.RoleDefinitionListFilter{
		RoleLevel: strings.TrimSpace(c.Query("roleLevel")),
		Status:    strings.TrimSpace(c.Query("status")),
	})
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *IdentityTenancyHandler) CreateRoleDefinition(c *gin.Context) {
	var req identitytenancy.CreateRoleDefinitionInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateRoleDefinition(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *IdentityTenancyHandler) ListRoleAssignments(c *gin.Context) {
	items, err := h.svc.ListRoleAssignments(c.Request.Context(), c.GetUint64(middleware.UserIDKey), identitytenancy.RoleAssignmentListFilter{
		SubjectRef: strings.TrimSpace(c.Query("subjectRef")),
		ScopeType:  strings.TrimSpace(c.Query("scopeType")),
		Status:     strings.TrimSpace(c.Query("status")),
	})
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *IdentityTenancyHandler) CreateRoleAssignment(c *gin.Context) {
	var req struct {
		SubjectType       string `json:"subjectType"`
		SubjectRef        string `json:"subjectRef"`
		RoleDefinitionID  uint64 `json:"roleDefinitionId"`
		ScopeType         string `json:"scopeType"`
		ScopeRef          string `json:"scopeRef"`
		SourceType        string `json:"sourceType"`
		DelegationGrantID uint64 `json:"delegationGrantId"`
		ValidUntil        string `json:"validUntil"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var validUntil *time.Time
	if strings.TrimSpace(req.ValidUntil) != "" {
		parsed, err := time.Parse(time.RFC3339, strings.TrimSpace(req.ValidUntil))
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "validUntil must be RFC3339"})
			return
		}
		validUntil = &parsed
	}
	item, err := h.svc.CreateRoleAssignment(c.Request.Context(), c.GetUint64(middleware.UserIDKey), identitytenancy.CreateRoleAssignmentInput{
		SubjectType:       req.SubjectType,
		SubjectRef:        req.SubjectRef,
		RoleDefinitionID:  req.RoleDefinitionID,
		ScopeType:         req.ScopeType,
		ScopeRef:          req.ScopeRef,
		SourceType:        req.SourceType,
		DelegationGrantID: req.DelegationGrantID,
		ValidUntil:        validUntil,
	})
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *IdentityTenancyHandler) CreateDelegationGrant(c *gin.Context) {
	var req identitytenancy.CreateDelegationGrantInput
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreateDelegationGrant(c.Request.Context(), c.GetUint64(middleware.UserIDKey), req)
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusCreated, item)
}

func (h *IdentityTenancyHandler) ListDelegationGrants(c *gin.Context) {
	items, err := h.svc.ListDelegationGrants(c.Request.Context(), c.GetUint64(middleware.UserIDKey))
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (h *IdentityTenancyHandler) ListAccessRisks(c *gin.Context) {
	items, err := h.svc.ListAccessRisks(c.Request.Context(), c.GetUint64(middleware.UserIDKey), identitytenancy.AccessRiskListFilter{
		SubjectType: strings.TrimSpace(c.Query("subjectType")),
		Severity:    strings.TrimSpace(c.Query("severity")),
	})
	if err != nil {
		writeIdentityTenancyError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func writeIdentityTenancyError(c *gin.Context, err error) {
	switch {
	case errors.Is(err, identitytenancy.ErrIdentityTenancyForbidden):
		c.JSON(http.StatusForbidden, gin.H{"error": err.Error()})
	case errors.Is(err, identitytenancy.ErrIdentityTenancyConflict), errors.Is(err, identitytenancy.ErrIdentityTenancyBlocked):
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
	case errors.Is(err, identitytenancy.ErrIdentityTenancyInvalid):
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
	case errors.Is(err, gorm.ErrRecordNotFound):
		c.JSON(http.StatusNotFound, gin.H{"error": "resource not found"})
	default:
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
}

func parseIdentityUint64Param(c *gin.Context, key string) (uint64, bool) {
	id, err := strconv.ParseUint(strings.TrimSpace(c.Param(key)), 10, 64)
	if err != nil || id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid " + key})
		return 0, false
	}
	return id, true
}

func parseOptionalIdentityUint64(v string) (uint64, bool) {
	if strings.TrimSpace(v) == "" {
		return 0, false
	}
	parsed, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
	if err != nil {
		return 0, false
	}
	return parsed, true
}
