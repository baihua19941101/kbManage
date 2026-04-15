package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/domain"
	auditSvc "kbmanage/backend/internal/service/audit"
	securityPolicySvc "kbmanage/backend/internal/service/securitypolicy"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SecurityPolicyHandler struct {
	svc         *securityPolicySvc.Service
	auditWriter *auditSvc.EventWriter
}

func NewSecurityPolicyHandler(svc *securityPolicySvc.Service, auditWriter ...*auditSvc.EventWriter) *SecurityPolicyHandler {
	if svc == nil {
		svc = securityPolicySvc.NewService(nil, nil, nil, nil, nil, nil, nil)
	}
	var writer *auditSvc.EventWriter
	if len(auditWriter) > 0 {
		writer = auditWriter[0]
	}
	return &SecurityPolicyHandler{svc: svc, auditWriter: writer}
}

type createSecurityPolicyRequest struct {
	Name                   string         `json:"name"`
	WorkspaceID            uint64         `json:"workspaceId"`
	ProjectID              uint64         `json:"projectId"`
	ScopeLevel             string         `json:"scopeLevel"`
	Category               string         `json:"category"`
	RuleTemplate           map[string]any `json:"ruleTemplate"`
	DefaultEnforcementMode string         `json:"defaultEnforcementMode"`
	RiskLevel              string         `json:"riskLevel"`
}

type updateSecurityPolicyRequest struct {
	Name                   *string        `json:"name"`
	RuleTemplate           map[string]any `json:"ruleTemplate"`
	DefaultEnforcementMode *string        `json:"defaultEnforcementMode"`
	Status                 *string        `json:"status"`
}

type createPolicyAssignmentRequest struct {
	WorkspaceID     uint64    `json:"workspaceId"`
	ProjectID       uint64    `json:"projectId"`
	ClusterRefs     []string  `json:"clusterRefs"`
	NamespaceRefs   []string  `json:"namespaceRefs"`
	ResourceKinds   []string  `json:"resourceKinds"`
	EnforcementMode string    `json:"enforcementMode"`
	RolloutStage    string    `json:"rolloutStage"`
	EffectiveFrom   time.Time `json:"effectiveFrom"`
	EffectiveTo     time.Time `json:"effectiveTo"`
}

type switchPolicyModeRequest struct {
	TargetMode    string   `json:"targetMode"`
	AssignmentIDs []uint64 `json:"assignmentIds"`
	Reason        string   `json:"reason"`
}

type createPolicyExceptionRequest struct {
	Reason    string    `json:"reason"`
	StartsAt  time.Time `json:"startsAt"`
	ExpiresAt time.Time `json:"expiresAt"`
}

type reviewPolicyExceptionRequest struct {
	Decision string `json:"decision"`
	Comment  string `json:"comment"`
}

type updatePolicyRemediationRequest struct {
	Status  string `json:"status"`
	Comment string `json:"comment"`
}

func (h *SecurityPolicyHandler) ListPolicies(c *gin.Context) {
	workspaceID, err := parseSecurityPolicyOptionalQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseSecurityPolicyOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListPolicies(
		c.Request.Context(),
		c.GetUint64(middleware.UserIDKey),
		workspaceID,
		projectID,
		securityPolicySvc.PolicyListFilter{
			ScopeLevel: strings.TrimSpace(c.Query("scopeLevel")),
			Status:     strings.TrimSpace(c.Query("status")),
			Category:   strings.TrimSpace(c.Query("category")),
		},
	)
	if err != nil {
		writeSecurityPolicyError(c, err)
		return
	}
	res := make([]gin.H, 0, len(items))
	for i := range items {
		res = append(res, toSecurityPolicyResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

func (h *SecurityPolicyHandler) CreatePolicy(c *gin.Context) {
	var req createSecurityPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.CreatePolicy(c.Request.Context(), c.GetUint64(middleware.UserIDKey), securityPolicySvc.CreatePolicyInput{
		Name:                   strings.TrimSpace(req.Name),
		WorkspaceID:            req.WorkspaceID,
		ProjectID:              req.ProjectID,
		ScopeLevel:             strings.TrimSpace(req.ScopeLevel),
		Category:               strings.TrimSpace(req.Category),
		RuleTemplate:           req.RuleTemplate,
		DefaultEnforcementMode: strings.TrimSpace(req.DefaultEnforcementMode),
		RiskLevel:              strings.TrimSpace(req.RiskLevel),
	})
	if err != nil {
		h.writePolicyAuditEvent(c, auditSvc.SecurityPolicyAuditActionPolicyCreate, "", domain.AuditOutcomeDenied, map[string]any{"error": err.Error()})
		writeSecurityPolicyError(c, err)
		return
	}
	h.writePolicyAuditEvent(
		c,
		auditSvc.SecurityPolicyAuditActionPolicyCreate,
		"policy:"+strconv.FormatUint(item.ID, 10),
		domain.AuditOutcomeSuccess,
		map[string]any{
			"policyId":    item.ID,
			"workspaceId": item.WorkspaceID,
			"projectId":   item.ProjectID,
			"category":    item.Category,
			"scopeLevel":  item.ScopeLevel,
			"policyName":  item.Name,
			"riskLevel":   item.RiskLevel,
			"enforcement": item.DefaultEnforcementMode,
		},
	)
	c.JSON(http.StatusCreated, toSecurityPolicyResponse(item))
}

func (h *SecurityPolicyHandler) GetPolicy(c *gin.Context) {
	policyID, err := parseSecurityPolicyRequiredParamUint64(c, "policyId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.GetPolicy(c.Request.Context(), c.GetUint64(middleware.UserIDKey), policyID)
	if err != nil {
		writeSecurityPolicyError(c, err)
		return
	}
	c.JSON(http.StatusOK, toSecurityPolicyResponse(item))
}

func (h *SecurityPolicyHandler) UpdatePolicy(c *gin.Context) {
	policyID, err := parseSecurityPolicyRequiredParamUint64(c, "policyId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req updateSecurityPolicyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.UpdatePolicy(c.Request.Context(), c.GetUint64(middleware.UserIDKey), policyID, securityPolicySvc.UpdatePolicyInput{
		Name:                   req.Name,
		RuleTemplate:           req.RuleTemplate,
		DefaultEnforcementMode: req.DefaultEnforcementMode,
		Status:                 req.Status,
	})
	if err != nil {
		h.writePolicyAuditEvent(
			c,
			auditSvc.SecurityPolicyAuditActionPolicyUpdate,
			"policy:"+strconv.FormatUint(policyID, 10),
			domain.AuditOutcomeDenied,
			map[string]any{"policyId": policyID, "error": err.Error()},
		)
		writeSecurityPolicyError(c, err)
		return
	}
	h.writePolicyAuditEvent(
		c,
		auditSvc.SecurityPolicyAuditActionPolicyUpdate,
		"policy:"+strconv.FormatUint(policyID, 10),
		domain.AuditOutcomeSuccess,
		map[string]any{"policyId": policyID},
	)
	c.JSON(http.StatusOK, toSecurityPolicyResponse(item))
}

func (h *SecurityPolicyHandler) ListAssignments(c *gin.Context) {
	policyID, err := parseSecurityPolicyRequiredParamUint64(c, "policyId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListAssignments(c.Request.Context(), c.GetUint64(middleware.UserIDKey), policyID)
	if err != nil {
		writeSecurityPolicyError(c, err)
		return
	}
	res := make([]gin.H, 0, len(items))
	for i := range items {
		res = append(res, toPolicyAssignmentResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

func (h *SecurityPolicyHandler) CreateAssignment(c *gin.Context) {
	policyID, err := parseSecurityPolicyRequiredParamUint64(c, "policyId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req createPolicyAssignmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := securityPolicySvc.CreateAssignmentInput{
		WorkspaceID:     req.WorkspaceID,
		ProjectID:       req.ProjectID,
		ClusterRefs:     req.ClusterRefs,
		NamespaceRefs:   req.NamespaceRefs,
		ResourceKinds:   req.ResourceKinds,
		EnforcementMode: strings.TrimSpace(req.EnforcementMode),
		RolloutStage:    strings.TrimSpace(req.RolloutStage),
	}
	if !req.EffectiveFrom.IsZero() {
		from := req.EffectiveFrom
		input.EffectiveFrom = &from
	}
	if !req.EffectiveTo.IsZero() {
		to := req.EffectiveTo
		input.EffectiveTo = &to
	}
	assignment, task, err := h.svc.CreateAssignment(c.Request.Context(), c.GetUint64(middleware.UserIDKey), policyID, input)
	if err != nil {
		h.writePolicyAuditEvent(
			c,
			auditSvc.SecurityPolicyAuditActionAssignmentCreate,
			"policy:"+strconv.FormatUint(policyID, 10),
			domain.AuditOutcomeDenied,
			map[string]any{"policyId": policyID, "error": err.Error()},
		)
		writeSecurityPolicyError(c, err)
		return
	}
	h.writePolicyAuditEvent(
		c,
		auditSvc.SecurityPolicyAuditActionAssignmentCreate,
		"policy:"+strconv.FormatUint(policyID, 10),
		domain.AuditOutcomeSuccess,
		map[string]any{
			"policyId":     policyID,
			"assignmentId": assignment.ID,
			"taskId":       task.ID,
		},
	)
	c.JSON(http.StatusAccepted, toPolicyDistributionTaskResponse(task))
}

func (h *SecurityPolicyHandler) SwitchPolicyMode(c *gin.Context) {
	policyID, err := parseSecurityPolicyRequiredParamUint64(c, "policyId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req switchPolicyModeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	task, err := h.svc.SwitchPolicyMode(c.Request.Context(), c.GetUint64(middleware.UserIDKey), policyID, securityPolicySvc.SwitchPolicyModeInput{
		TargetMode:    strings.TrimSpace(req.TargetMode),
		AssignmentIDs: req.AssignmentIDs,
		Reason:        strings.TrimSpace(req.Reason),
	})
	if err != nil {
		h.writePolicyAuditEvent(
			c,
			auditSvc.SecurityPolicyAuditActionModeSwitch,
			"policy:"+strconv.FormatUint(policyID, 10),
			domain.AuditOutcomeDenied,
			map[string]any{"policyId": policyID, "error": err.Error()},
		)
		writeSecurityPolicyError(c, err)
		return
	}
	h.writePolicyAuditEvent(
		c,
		auditSvc.SecurityPolicyAuditActionModeSwitch,
		"policy:"+strconv.FormatUint(policyID, 10),
		domain.AuditOutcomeSuccess,
		map[string]any{
			"policyId":      policyID,
			"targetMode":    req.TargetMode,
			"assignmentIds": req.AssignmentIDs,
			"taskId":        task.ID,
		},
	)
	c.JSON(http.StatusAccepted, toPolicyDistributionTaskResponse(task))
}

func (h *SecurityPolicyHandler) ListHits(c *gin.Context) {
	policyID, err := parseSecurityPolicyOptionalQueryUint64(c, "policyId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	workspaceID, err := parseSecurityPolicyOptionalQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseSecurityPolicyOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clusterID, err := parseSecurityPolicyOptionalQueryUint64(c, "clusterId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	limit := 100
	if rawLimit := strings.TrimSpace(c.Query("limit")); rawLimit != "" {
		parsed, parseErr := strconv.Atoi(rawLimit)
		if parseErr != nil || parsed <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit"})
			return
		}
		limit = parsed
	}
	from, err := parseSecurityPolicyOptionalQueryTime(c, "from")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	to, err := parseSecurityPolicyOptionalQueryTime(c, "to")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	items, err := h.svc.ListHits(c.Request.Context(), c.GetUint64(middleware.UserIDKey), securityPolicySvc.ListHitsInput{
		PolicyID:          policyID,
		WorkspaceID:       workspaceID,
		ProjectID:         projectID,
		ClusterID:         clusterID,
		Namespace:         strings.TrimSpace(c.Query("namespace")),
		EnforcementMode:   strings.TrimSpace(c.Query("enforcementMode")),
		RiskLevel:         strings.TrimSpace(c.Query("riskLevel")),
		RemediationStatus: strings.TrimSpace(c.Query("remediationStatus")),
		From:              from,
		To:                to,
		Limit:             limit,
	})
	if err != nil {
		h.writePolicyAuditEvent(
			c,
			auditSvc.SecurityPolicyAuditActionHitQuery,
			"",
			domain.AuditOutcomeDenied,
			map[string]any{"policyId": policyID, "error": err.Error()},
		)
		writeSecurityPolicyError(c, err)
		return
	}
	res := make([]gin.H, 0, len(items))
	for i := range items {
		res = append(res, toPolicyHitResponse(&items[i]))
	}
	h.writePolicyAuditEvent(
		c,
		auditSvc.SecurityPolicyAuditActionHitQuery,
		"",
		domain.AuditOutcomeSuccess,
		map[string]any{
			"policyId":          policyID,
			"workspaceId":       workspaceID,
			"projectId":         projectID,
			"clusterId":         clusterID,
			"namespace":         strings.TrimSpace(c.Query("namespace")),
			"enforcementMode":   strings.TrimSpace(c.Query("enforcementMode")),
			"riskLevel":         strings.TrimSpace(c.Query("riskLevel")),
			"remediationStatus": strings.TrimSpace(c.Query("remediationStatus")),
			"resultCount":       len(items),
		},
	)
	c.JSON(http.StatusOK, gin.H{"items": res})
}

func (h *SecurityPolicyHandler) CreateException(c *gin.Context) {
	hitID, err := parseSecurityPolicyRequiredParamUint64(c, "hitId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req createPolicyExceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	input := securityPolicySvc.CreateExceptionInput{
		Reason: strings.TrimSpace(req.Reason),
	}
	if !req.StartsAt.IsZero() {
		startsAt := req.StartsAt
		input.StartsAt = &startsAt
	}
	if !req.ExpiresAt.IsZero() {
		expiresAt := req.ExpiresAt
		input.ExpiresAt = &expiresAt
	}
	item, err := h.svc.CreateException(c.Request.Context(), c.GetUint64(middleware.UserIDKey), hitID, input)
	if err != nil {
		h.writePolicyAuditEvent(
			c,
			auditSvc.SecurityPolicyAuditActionExceptionCreate,
			"hit:"+strconv.FormatUint(hitID, 10),
			domain.AuditOutcomeDenied,
			map[string]any{"hitId": hitID, "error": err.Error()},
		)
		writeSecurityPolicyError(c, err)
		return
	}
	h.writePolicyAuditEvent(
		c,
		auditSvc.SecurityPolicyAuditActionExceptionCreate,
		"exception:"+strconv.FormatUint(item.ID, 10),
		domain.AuditOutcomeSuccess,
		map[string]any{
			"exceptionId": item.ID,
			"hitId":       item.HitID,
			"policyId":    item.PolicyID,
			"status":      item.Status,
		},
	)
	c.JSON(http.StatusCreated, toPolicyExceptionResponse(item))
}

func (h *SecurityPolicyHandler) ListExceptions(c *gin.Context) {
	workspaceID, err := parseSecurityPolicyOptionalQueryUint64(c, "workspaceId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	projectID, err := parseSecurityPolicyOptionalQueryUint64(c, "projectId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	policyID, err := parseSecurityPolicyOptionalQueryUint64(c, "policyId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	items, err := h.svc.ListExceptions(c.Request.Context(), c.GetUint64(middleware.UserIDKey), securityPolicySvc.ListExceptionsInput{
		WorkspaceID: workspaceID,
		ProjectID:   projectID,
		PolicyID:    policyID,
		Status:      strings.TrimSpace(c.Query("status")),
	})
	if err != nil {
		writeSecurityPolicyError(c, err)
		return
	}
	res := make([]gin.H, 0, len(items))
	for i := range items {
		res = append(res, toPolicyExceptionResponse(&items[i]))
	}
	c.JSON(http.StatusOK, gin.H{"items": res})
}

func (h *SecurityPolicyHandler) ReviewException(c *gin.Context) {
	exceptionID, err := parseSecurityPolicyRequiredParamUint64(c, "exceptionId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req reviewPolicyExceptionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.ReviewException(c.Request.Context(), c.GetUint64(middleware.UserIDKey), exceptionID, securityPolicySvc.ReviewExceptionInput{
		Decision: strings.TrimSpace(req.Decision),
		Comment:  strings.TrimSpace(req.Comment),
	})
	if err != nil {
		h.writePolicyAuditEvent(
			c,
			auditSvc.SecurityPolicyAuditActionExceptionReview,
			"exception:"+strconv.FormatUint(exceptionID, 10),
			domain.AuditOutcomeDenied,
			map[string]any{"exceptionId": exceptionID, "error": err.Error()},
		)
		writeSecurityPolicyError(c, err)
		return
	}
	h.writePolicyAuditEvent(
		c,
		auditSvc.SecurityPolicyAuditActionExceptionReview,
		"exception:"+strconv.FormatUint(exceptionID, 10),
		domain.AuditOutcomeSuccess,
		map[string]any{"exceptionId": exceptionID, "decision": req.Decision, "status": item.Status},
	)
	c.JSON(http.StatusOK, toPolicyExceptionResponse(item))
}

func (h *SecurityPolicyHandler) UpdateRemediation(c *gin.Context) {
	hitID, err := parseSecurityPolicyRequiredParamUint64(c, "hitId")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var req updatePolicyRemediationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	item, err := h.svc.UpdateRemediation(c.Request.Context(), c.GetUint64(middleware.UserIDKey), hitID, securityPolicySvc.UpdateRemediationInput{
		Status:  strings.TrimSpace(req.Status),
		Comment: strings.TrimSpace(req.Comment),
	})
	if err != nil {
		h.writePolicyAuditEvent(
			c,
			auditSvc.SecurityPolicyAuditActionHitRemediationUpdate,
			"hit:"+strconv.FormatUint(hitID, 10),
			domain.AuditOutcomeDenied,
			map[string]any{"hitId": hitID, "error": err.Error()},
		)
		writeSecurityPolicyError(c, err)
		return
	}
	h.writePolicyAuditEvent(
		c,
		auditSvc.SecurityPolicyAuditActionHitRemediationUpdate,
		"hit:"+strconv.FormatUint(hitID, 10),
		domain.AuditOutcomeSuccess,
		func() map[string]any {
			details := map[string]any{
				"hitId":             item.ID,
				"policyId":          item.PolicyID,
				"clusterId":         item.ClusterID,
				"resourceKind":      item.ResourceKind,
				"resourceName":      item.ResourceName,
				"remediationStatus": item.RemediationStatus,
			}
			if policy, getErr := h.svc.GetPolicy(c.Request.Context(), c.GetUint64(middleware.UserIDKey), item.PolicyID); getErr == nil && policy != nil {
				details["workspaceId"] = policy.WorkspaceID
				details["projectId"] = policy.ProjectID
			}
			return details
		}(),
	)
	c.JSON(http.StatusOK, toPolicyHitResponse(item))
}

func toSecurityPolicyResponse(item *domain.SecurityPolicy) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":                     item.ID,
		"name":                   item.Name,
		"workspaceId":            item.WorkspaceID,
		"projectId":              item.ProjectID,
		"scopeLevel":             item.ScopeLevel,
		"category":               item.Category,
		"ruleTemplate":           decodePolicyJSONMap(item.RuleTemplateJSON),
		"defaultEnforcementMode": item.DefaultEnforcementMode,
		"riskLevel":              item.RiskLevel,
		"status":                 item.Status,
		"updatedAt":              item.UpdatedAt,
	}
}

func toPolicyAssignmentResponse(item *domain.PolicyAssignment) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":              item.ID,
		"policyId":        item.PolicyID,
		"workspaceId":     item.WorkspaceID,
		"projectId":       item.ProjectID,
		"clusterRefs":     decodePolicyStringArray(item.ClusterRefsJSON),
		"namespaceRefs":   decodePolicyStringArray(item.NamespaceRefsJSON),
		"resourceKinds":   decodePolicyStringArray(item.ResourceKindsJSON),
		"enforcementMode": item.EnforcementMode,
		"rolloutStage":    item.RolloutStage,
		"status":          item.Status,
		"effectiveFrom":   item.EffectiveFrom,
		"effectiveTo":     item.EffectiveTo,
	}
}

func toPolicyDistributionTaskResponse(item *domain.PolicyDistributionTask) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":             item.ID,
		"policyId":       item.PolicyID,
		"operation":      item.Operation,
		"status":         item.Status,
		"targetCount":    item.TargetCount,
		"succeededCount": item.SucceededCount,
		"failedCount":    item.FailedCount,
		"resultSummary":  item.ResultSummary,
	}
}

func toPolicyHitResponse(item *domain.PolicyHitRecord) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":                item.ID,
		"policyId":          item.PolicyID,
		"assignmentId":      item.AssignmentID,
		"clusterId":         item.ClusterID,
		"namespace":         item.Namespace,
		"resourceKind":      item.ResourceKind,
		"resourceName":      item.ResourceName,
		"hitResult":         item.HitResult,
		"riskLevel":         item.RiskLevel,
		"message":           item.Message,
		"remediationStatus": item.RemediationStatus,
		"detectedAt":        item.DetectedAt,
	}
}

func toPolicyExceptionResponse(item *domain.PolicyExceptionRequest) gin.H {
	if item == nil {
		return gin.H{}
	}
	return gin.H{
		"id":            item.ID,
		"policyId":      item.PolicyID,
		"hitRecordId":   item.HitID,
		"workspaceId":   item.WorkspaceID,
		"projectId":     item.ProjectID,
		"reason":        item.Reason,
		"status":        item.Status,
		"startsAt":      item.CreatedAt,
		"expiresAt":     item.ExpiresAt,
		"reviewComment": item.ReviewComment,
	}
}

func decodePolicyJSONMap(raw string) map[string]any {
	if strings.TrimSpace(raw) == "" {
		return map[string]any{}
	}
	var res map[string]any
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return map[string]any{}
	}
	if res == nil {
		return map[string]any{}
	}
	return res
}

func decodePolicyStringArray(raw string) []string {
	if strings.TrimSpace(raw) == "" {
		return []string{}
	}
	var res []string
	if err := json.Unmarshal([]byte(raw), &res); err != nil {
		return []string{}
	}
	if res == nil {
		return []string{}
	}
	return res
}

func writeSecurityPolicyError(c *gin.Context, err error) {
	status := http.StatusInternalServerError
	lower := strings.ToLower(err.Error())
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		status = http.StatusNotFound
	case strings.Contains(lower, "scope access denied"), strings.Contains(lower, "forbidden"):
		status = http.StatusForbidden
	case strings.Contains(lower, "required"), strings.Contains(lower, "invalid"):
		status = http.StatusBadRequest
	case strings.Contains(lower, "not configured"):
		status = http.StatusServiceUnavailable
	}
	c.JSON(status, gin.H{"error": err.Error()})
}

func parseSecurityPolicyRequiredParamUint64(c *gin.Context, key string) (uint64, error) {
	value := strings.TrimSpace(c.Param(key))
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		return 0, errors.New("invalid " + key)
	}
	return parsed, nil
}

func parseSecurityPolicyOptionalQueryUint64(c *gin.Context, key string) (uint64, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.ParseUint(value, 10, 64)
	if err != nil || parsed == 0 {
		return 0, errors.New("invalid " + key)
	}
	return parsed, nil
}

func parseSecurityPolicyOptionalQueryTime(c *gin.Context, key string) (*time.Time, error) {
	value := strings.TrimSpace(c.Query(key))
	if value == "" {
		return nil, nil
	}
	parsed, err := time.Parse(time.RFC3339, value)
	if err != nil {
		return nil, errors.New("invalid " + key)
	}
	return &parsed, nil
}

func (h *SecurityPolicyHandler) writePolicyAuditEvent(
	c *gin.Context,
	action string,
	resourceID string,
	outcome domain.AuditOutcome,
	details map[string]any,
) {
	if h == nil || h.auditWriter == nil {
		return
	}
	userID := c.GetUint64(middleware.UserIDKey)
	if userID == 0 {
		return
	}
	actorID := userID
	_ = h.auditWriter.WriteSecurityPolicyEvent(
		c.Request.Context(),
		c.GetString(middleware.RequestIDKey),
		&actorID,
		action,
		resourceID,
		outcome,
		details,
	)
}
