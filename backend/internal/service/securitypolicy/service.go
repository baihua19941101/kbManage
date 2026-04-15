package securitypolicy

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"

	"gorm.io/gorm"
)

const (
	PermissionSecurityPolicyRead    = "securitypolicy:read"
	PermissionSecurityPolicyManage  = "securitypolicy:manage"
	PermissionSecurityPolicyEnforce = "securitypolicy:enforce"
)

var (
	ErrSecurityPolicyNotConfigured = errors.New("security policy service is not configured")
	ErrSecurityPolicyScopeDenied   = errors.New("security policy scope access denied")
)

type PolicyListFilter struct {
	ScopeLevel string
	Status     string
	Category   string
}

type CreatePolicyInput struct {
	Name                   string
	WorkspaceID            uint64
	ProjectID              uint64
	ScopeLevel             string
	Category               string
	RuleTemplate           map[string]any
	DefaultEnforcementMode string
	RiskLevel              string
}

type UpdatePolicyInput struct {
	Name                   *string
	RuleTemplate           map[string]any
	DefaultEnforcementMode *string
	Status                 *string
}

type CreateAssignmentInput struct {
	WorkspaceID     uint64
	ProjectID       uint64
	ClusterRefs     []string
	NamespaceRefs   []string
	ResourceKinds   []string
	EnforcementMode string
	RolloutStage    string
	EffectiveFrom   *time.Time
	EffectiveTo     *time.Time
}

type SwitchPolicyModeInput struct {
	TargetMode    string
	AssignmentIDs []uint64
	Reason        string
}

type ListHitsInput struct {
	PolicyID          uint64
	WorkspaceID       uint64
	ProjectID         uint64
	ClusterID         uint64
	Namespace         string
	EnforcementMode   string
	RiskLevel         string
	RemediationStatus string
	From              *time.Time
	To                *time.Time
	Limit             int
}

type UpdateRemediationInput struct {
	Status  string
	Comment string
}

type CreateExceptionInput struct {
	Reason    string
	StartsAt  *time.Time
	ExpiresAt *time.Time
}

type ListExceptionsInput struct {
	WorkspaceID uint64
	ProjectID   uint64
	PolicyID    uint64
	Status      string
}

type ReviewExceptionInput struct {
	Decision string
	Comment  string
}

type Service struct {
	policies       *repository.SecurityPolicyRepository
	assignments    *repository.PolicyAssignmentRepository
	hits           *repository.PolicyHitRepository
	exceptions     *repository.PolicyExceptionRepository
	scope          *ScopeService
	distribution   *DistributionCache
	exceptionCache *ExceptionCache
}

func NewService(
	policies *repository.SecurityPolicyRepository,
	assignments *repository.PolicyAssignmentRepository,
	hits *repository.PolicyHitRepository,
	exceptions *repository.PolicyExceptionRepository,
	scope *ScopeService,
	distribution *DistributionCache,
	exceptionCache *ExceptionCache,
) *Service {
	return &Service{
		policies:       policies,
		assignments:    assignments,
		hits:           hits,
		exceptions:     exceptions,
		scope:          scope,
		distribution:   distribution,
		exceptionCache: exceptionCache,
	}
}

func (s *Service) ListPolicies(
	ctx context.Context,
	userID uint64,
	workspaceID uint64,
	projectID uint64,
	filter PolicyListFilter,
) ([]domain.SecurityPolicy, error) {
	if s == nil || s.policies == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	if err := s.validateScope(ctx, userID, workspaceID, projectID, PermissionSecurityPolicyRead); err != nil {
		return nil, err
	}
	scopeLevel, err := normalizeScopeLevel(filter.ScopeLevel, true)
	if err != nil {
		return nil, err
	}
	category, err := normalizeCategory(filter.Category, true)
	if err != nil {
		return nil, err
	}
	status, err := normalizePolicyStatus(filter.Status, true)
	if err != nil {
		return nil, err
	}
	return s.policies.List(ctx, repository.SecurityPolicyListFilter{
		WorkspaceID: uint64PtrOrNil(workspaceID),
		ProjectID:   uint64PtrOrNil(projectID),
		ScopeLevel:  scopeLevel,
		Category:    category,
		Status:      status,
	})
}

func (s *Service) CreatePolicy(ctx context.Context, userID uint64, input CreatePolicyInput) (*domain.SecurityPolicy, error) {
	if s == nil || s.policies == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	if err := s.validateScope(ctx, userID, input.WorkspaceID, input.ProjectID, PermissionSecurityPolicyManage); err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, errors.New("name is required")
	}
	scopeLevel, err := normalizeScopeLevel(input.ScopeLevel, false)
	if err != nil {
		return nil, err
	}
	category, err := normalizeCategory(input.Category, false)
	if err != nil {
		return nil, err
	}
	mode, err := normalizeEnforcementMode(input.DefaultEnforcementMode, false)
	if err != nil {
		return nil, err
	}
	riskLevel, err := normalizeRiskLevel(input.RiskLevel, true)
	if err != nil {
		return nil, err
	}
	if riskLevel == "" {
		riskLevel = domain.PolicyRiskLevelMedium
	}
	templateJSON, err := marshalJSONMap(input.RuleTemplate)
	if err != nil {
		return nil, err
	}
	createdBy := userID
	item := &domain.SecurityPolicy{
		Name:                   name,
		WorkspaceID:            uint64PtrOrNil(input.WorkspaceID),
		ProjectID:              uint64PtrOrNil(input.ProjectID),
		ScopeLevel:             scopeLevel,
		Category:               category,
		RuleTemplateJSON:       templateJSON,
		DefaultEnforcementMode: mode,
		RiskLevel:              riskLevel,
		Status:                 domain.PolicyStatusDraft,
		CreatedBy:              &createdBy,
		UpdatedBy:              &createdBy,
	}
	if err := s.policies.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) GetPolicy(ctx context.Context, userID uint64, policyID uint64) (*domain.SecurityPolicy, error) {
	if s == nil || s.policies == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	item, err := s.policies.GetByID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if err := s.validatePolicyScope(ctx, userID, item, PermissionSecurityPolicyRead); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *Service) UpdatePolicy(
	ctx context.Context,
	userID uint64,
	policyID uint64,
	input UpdatePolicyInput,
) (*domain.SecurityPolicy, error) {
	if s == nil || s.policies == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	item, err := s.policies.GetByID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if err := s.validatePolicyScope(ctx, userID, item, PermissionSecurityPolicyManage); err != nil {
		return nil, err
	}

	updates := make(map[string]any)
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, errors.New("name is required")
		}
		updates["name"] = name
	}
	if input.RuleTemplate != nil {
		templateJSON, marshalErr := marshalJSONMap(input.RuleTemplate)
		if marshalErr != nil {
			return nil, marshalErr
		}
		updates["rule_template_json"] = templateJSON
	}
	if input.DefaultEnforcementMode != nil {
		mode, normalizeErr := normalizeEnforcementMode(*input.DefaultEnforcementMode, false)
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		updates["default_enforcement_mode"] = mode
	}
	if input.Status != nil {
		status, normalizeErr := normalizePolicyStatus(*input.Status, false)
		if normalizeErr != nil {
			return nil, normalizeErr
		}
		updates["status"] = status
	}
	if len(updates) == 0 {
		return item, nil
	}
	updatedBy := userID
	updates["updated_by"] = &updatedBy
	if err := s.policies.UpdateFields(ctx, policyID, updates); err != nil {
		return nil, err
	}
	return s.policies.GetByID(ctx, policyID)
}

func (s *Service) ListAssignments(ctx context.Context, userID uint64, policyID uint64) ([]domain.PolicyAssignment, error) {
	if s == nil || s.policies == nil || s.assignments == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	policy, err := s.policies.GetByID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if err := s.validatePolicyScope(ctx, userID, policy, PermissionSecurityPolicyRead); err != nil {
		return nil, err
	}
	return s.assignments.ListByPolicyID(ctx, policyID)
}

func (s *Service) CreateAssignment(
	ctx context.Context,
	userID uint64,
	policyID uint64,
	input CreateAssignmentInput,
) (*domain.PolicyAssignment, *domain.PolicyDistributionTask, error) {
	if s == nil || s.policies == nil || s.assignments == nil {
		return nil, nil, ErrSecurityPolicyNotConfigured
	}
	policy, err := s.policies.GetByID(ctx, policyID)
	if err != nil {
		return nil, nil, err
	}
	if err := s.validatePolicyScope(ctx, userID, policy, PermissionSecurityPolicyEnforce); err != nil {
		return nil, nil, err
	}
	mode, err := normalizeEnforcementMode(input.EnforcementMode, false)
	if err != nil {
		return nil, nil, err
	}
	rolloutStage, err := normalizeRolloutStage(input.RolloutStage, false)
	if err != nil {
		return nil, nil, err
	}
	clusterRefsJSON, err := marshalStringArray(input.ClusterRefs)
	if err != nil {
		return nil, nil, err
	}
	namespaceRefsJSON, err := marshalStringArray(input.NamespaceRefs)
	if err != nil {
		return nil, nil, err
	}
	resourceKindsJSON, err := marshalStringArray(input.ResourceKinds)
	if err != nil {
		return nil, nil, err
	}
	createdBy := userID
	now := time.Now()
	task := &domain.PolicyDistributionTask{
		PolicyID:       policyID,
		Operation:      domain.PolicyDistributionOperationAssign,
		Status:         domain.PolicyDistributionStatusSucceeded,
		TargetCount:    len(input.ClusterRefs),
		SucceededCount: len(input.ClusterRefs),
		FailedCount:    0,
		ResultSummary:  "assignment distributed",
		CreatedBy:      &createdBy,
		StartedAt:      &now,
		CompletedAt:    &now,
	}
	if err := s.assignments.CreateDistributionTask(ctx, task); err != nil {
		return nil, nil, err
	}

	assignment := &domain.PolicyAssignment{
		PolicyID:          policyID,
		WorkspaceID:       uint64PtrOrNil(input.WorkspaceID),
		ProjectID:         uint64PtrOrNil(input.ProjectID),
		ClusterRefsJSON:   clusterRefsJSON,
		NamespaceRefsJSON: namespaceRefsJSON,
		ResourceKindsJSON: resourceKindsJSON,
		EnforcementMode:   mode,
		RolloutStage:      rolloutStage,
		Status:            domain.PolicyAssignmentStatusActive,
		EffectiveFrom:     input.EffectiveFrom,
		EffectiveTo:       input.EffectiveTo,
		LastTaskID:        &task.ID,
		CreatedBy:         &createdBy,
		UpdatedBy:         &createdBy,
	}
	if err := s.assignments.Create(ctx, assignment); err != nil {
		return nil, nil, err
	}
	if s.distribution != nil {
		_ = s.distribution.SetTaskSnapshot(ctx, task.ID, DistributionTaskSnapshot{
			Status:         string(task.Status),
			TargetCount:    task.TargetCount,
			SucceededCount: task.SucceededCount,
			FailedCount:    task.FailedCount,
			UpdatedAt:      now,
		})
	}
	return assignment, task, nil
}

func (s *Service) SwitchPolicyMode(
	ctx context.Context,
	userID uint64,
	policyID uint64,
	input SwitchPolicyModeInput,
) (*domain.PolicyDistributionTask, error) {
	if s == nil || s.policies == nil || s.assignments == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	policy, err := s.policies.GetByID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if err := s.validatePolicyScope(ctx, userID, policy, PermissionSecurityPolicyEnforce); err != nil {
		return nil, err
	}
	mode, err := normalizeEnforcementMode(input.TargetMode, false)
	if err != nil {
		return nil, err
	}
	assignmentIDs := uniqueUint64s(input.AssignmentIDs)
	assignments, err := s.assignments.ListByPolicyID(ctx, policyID)
	if err != nil {
		return nil, err
	}
	if len(assignmentIDs) > 0 {
		assignments, err = s.assignments.ListByPolicyAndIDs(ctx, policyID, assignmentIDs)
		if err != nil {
			return nil, err
		}
		if len(assignments) == 0 {
			return nil, errors.New("assignmentIds not found for policy")
		}
	}

	now := time.Now()
	createdBy := userID
	task := &domain.PolicyDistributionTask{
		PolicyID:       policyID,
		Operation:      domain.PolicyDistributionOperationModeSwitch,
		Status:         domain.PolicyDistributionStatusSucceeded,
		TargetCount:    len(assignments),
		SucceededCount: len(assignments),
		FailedCount:    0,
		ResultSummary:  "mode switch distributed",
		CreatedBy:      &createdBy,
		StartedAt:      &now,
		CompletedAt:    &now,
	}
	if trimmedReason := strings.TrimSpace(input.Reason); trimmedReason != "" {
		task.ResultSummary = "mode switch distributed: " + trimmedReason
	}
	if err := s.assignments.CreateDistributionTask(ctx, task); err != nil {
		return nil, err
	}

	if len(assignments) > 0 {
		ids := make([]uint64, 0, len(assignments))
		for i := range assignments {
			ids = append(ids, assignments[i].ID)
		}
		if err := s.assignments.BulkUpdateEnforcementMode(ctx, ids, mode, "", &createdBy); err != nil {
			return nil, err
		}
	}

	if err := s.policies.UpdateFields(ctx, policyID, map[string]any{
		"default_enforcement_mode": mode,
		"updated_by":               &createdBy,
	}); err != nil {
		return nil, err
	}

	if s.distribution != nil {
		_ = s.distribution.SetTaskSnapshot(ctx, task.ID, DistributionTaskSnapshot{
			Status:         string(task.Status),
			TargetCount:    task.TargetCount,
			SucceededCount: task.SucceededCount,
			FailedCount:    task.FailedCount,
			UpdatedAt:      now,
		})
	}
	return task, nil
}

func (s *Service) ListHits(ctx context.Context, userID uint64, input ListHitsInput) ([]domain.PolicyHitRecord, error) {
	if s == nil || s.hits == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	if input.PolicyID != 0 {
		policy, err := s.policies.GetByID(ctx, input.PolicyID)
		if err != nil {
			return nil, err
		}
		if err := s.validatePolicyScope(ctx, userID, policy, PermissionSecurityPolicyRead); err != nil {
			return nil, err
		}
	} else if err := s.validateScope(ctx, userID, input.WorkspaceID, input.ProjectID, PermissionSecurityPolicyRead); err != nil {
		return nil, err
	}

	riskLevel, err := normalizeRiskLevel(input.RiskLevel, true)
	if err != nil {
		return nil, err
	}
	remediation, err := normalizeRemediationStatus(input.RemediationStatus, true)
	if err != nil {
		return nil, err
	}
	enforcementMode, err := normalizeEnforcementMode(input.EnforcementMode, true)
	if err != nil {
		return nil, err
	}
	if input.To != nil && input.From != nil && input.To.Before(*input.From) {
		return nil, errors.New("to must be after from")
	}
	return s.hits.List(ctx, repository.PolicyHitQuery{
		PolicyID:          uint64PtrOrNil(input.PolicyID),
		WorkspaceID:       uint64PtrOrNil(input.WorkspaceID),
		ProjectID:         uint64PtrOrNil(input.ProjectID),
		ClusterID:         uint64PtrOrNil(input.ClusterID),
		Namespace:         strings.TrimSpace(input.Namespace),
		EnforcementMode:   enforcementMode,
		RiskLevel:         riskLevel,
		RemediationStatus: remediation,
		DetectedFrom:      input.From,
		DetectedTo:        input.To,
		Limit:             input.Limit,
	})
}

func (s *Service) UpdateRemediation(
	ctx context.Context,
	userID uint64,
	hitID uint64,
	input UpdateRemediationInput,
) (*domain.PolicyHitRecord, error) {
	if s == nil || s.hits == nil || s.policies == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	hit, err := s.hits.GetByID(ctx, hitID)
	if err != nil {
		return nil, err
	}
	policy, err := s.policies.GetByID(ctx, hit.PolicyID)
	if err != nil {
		return nil, err
	}
	if err := s.validatePolicyScope(ctx, userID, policy, PermissionSecurityPolicyManage); err != nil {
		return nil, err
	}

	nextStatus, err := normalizeRemediationStatus(input.Status, false)
	if err != nil {
		return nil, err
	}
	if err := validateRemediationTransition(hit.RemediationStatus, nextStatus); err != nil {
		return nil, err
	}

	updates := map[string]any{
		"remediation_status": nextStatus,
	}
	switch nextStatus {
	case domain.PolicyRemediationMitigated, domain.PolicyRemediationClosed:
		resolvedAt := time.Now().UTC()
		updates["resolved_at"] = &resolvedAt
	default:
		updates["resolved_at"] = nil
	}

	trimmedComment := strings.TrimSpace(input.Comment)
	if trimmedComment != "" {
		updates["message"] = appendRemediationComment(hit.Message, trimmedComment, time.Now().UTC())
	}
	if err := s.hits.UpdateFields(ctx, hitID, updates); err != nil {
		return nil, err
	}
	return s.hits.GetByID(ctx, hitID)
}

func (s *Service) CreateException(
	ctx context.Context,
	userID uint64,
	hitID uint64,
	input CreateExceptionInput,
) (*domain.PolicyExceptionRequest, error) {
	if s == nil || s.hits == nil || s.exceptions == nil || s.policies == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	hit, err := s.hits.GetByID(ctx, hitID)
	if err != nil {
		return nil, err
	}
	policy, err := s.policies.GetByID(ctx, hit.PolicyID)
	if err != nil {
		return nil, err
	}
	if err := s.validatePolicyScope(ctx, userID, policy, PermissionSecurityPolicyEnforce); err != nil {
		return nil, err
	}
	reason := strings.TrimSpace(input.Reason)
	if reason == "" {
		return nil, errors.New("reason is required")
	}
	now := time.Now()
	startsAt := now
	if input.StartsAt != nil {
		startsAt = *input.StartsAt
	}
	if input.ExpiresAt == nil {
		return nil, errors.New("expiresAt is required")
	}
	expiresAt := *input.ExpiresAt
	if !expiresAt.After(startsAt) {
		return nil, errors.New("expiresAt must be after startsAt")
	}

	requestedBy := userID
	item := &domain.PolicyExceptionRequest{
		HitID:       &hitID,
		PolicyID:    policy.ID,
		WorkspaceID: policy.WorkspaceID,
		ProjectID:   policy.ProjectID,
		Status:      domain.PolicyExceptionPending,
		Reason:      reason,
		ExpiresAt:   &expiresAt,
		RequestedBy: &requestedBy,
		CreatedAt:   startsAt,
	}
	if err := s.exceptions.Create(ctx, item); err != nil {
		return nil, err
	}
	if s.exceptionCache != nil {
		_ = s.exceptionCache.SetExceptionStatus(ctx, item.ID, string(item.Status))
	}
	return item, nil
}

func (s *Service) ListExceptions(
	ctx context.Context,
	userID uint64,
	input ListExceptionsInput,
) ([]domain.PolicyExceptionRequest, error) {
	if s == nil || s.exceptions == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	if input.PolicyID != 0 {
		policy, err := s.policies.GetByID(ctx, input.PolicyID)
		if err != nil {
			return nil, err
		}
		if err := s.validatePolicyScope(ctx, userID, policy, PermissionSecurityPolicyRead); err != nil {
			return nil, err
		}
	} else if err := s.validateScope(ctx, userID, input.WorkspaceID, input.ProjectID, PermissionSecurityPolicyRead); err != nil {
		return nil, err
	}
	status, err := normalizeExceptionStatus(input.Status, true)
	if err != nil {
		return nil, err
	}
	return s.exceptions.List(ctx, repository.PolicyExceptionListFilter{
		WorkspaceID: uint64PtrOrNil(input.WorkspaceID),
		ProjectID:   uint64PtrOrNil(input.ProjectID),
		PolicyID:    uint64PtrOrNil(input.PolicyID),
		Status:      status,
	})
}

func (s *Service) ReviewException(
	ctx context.Context,
	userID uint64,
	exceptionID uint64,
	input ReviewExceptionInput,
) (*domain.PolicyExceptionRequest, error) {
	if s == nil || s.exceptions == nil || s.policies == nil {
		return nil, ErrSecurityPolicyNotConfigured
	}
	item, err := s.exceptions.GetByID(ctx, exceptionID)
	if err != nil {
		return nil, err
	}
	policy, err := s.policies.GetByID(ctx, item.PolicyID)
	if err != nil {
		return nil, err
	}
	if err := s.validatePolicyScope(ctx, userID, policy, PermissionSecurityPolicyManage); err != nil {
		return nil, err
	}
	decision := strings.ToLower(strings.TrimSpace(input.Decision))
	if decision == "" {
		return nil, errors.New("decision is required")
	}
	nextStatus, err := resolveReviewedExceptionStatus(decision, item, time.Now())
	if err != nil {
		return nil, err
	}
	reviewedBy := userID
	reviewedAt := time.Now()
	updates := map[string]any{
		"status":         nextStatus,
		"reviewed_by":    &reviewedBy,
		"reviewed_at":    &reviewedAt,
		"review_comment": strings.TrimSpace(input.Comment),
	}
	if err := s.exceptions.UpdateFields(ctx, exceptionID, updates); err != nil {
		return nil, err
	}
	updated, err := s.exceptions.GetByID(ctx, exceptionID)
	if err != nil {
		return nil, err
	}
	if s.exceptionCache != nil {
		_ = s.exceptionCache.SetExceptionStatus(ctx, updated.ID, string(updated.Status))
	}
	return updated, nil
}

func (s *Service) validateScope(ctx context.Context, userID uint64, workspaceID uint64, projectID uint64, permission string) error {
	if s == nil || s.scope == nil {
		return nil
	}
	return s.scope.ValidateScope(ctx, userID, workspaceID, projectID, permission)
}

func (s *Service) validatePolicyScope(ctx context.Context, userID uint64, policy *domain.SecurityPolicy, permission string) error {
	if s == nil || s.scope == nil {
		return nil
	}
	return s.scope.ValidatePolicyScope(ctx, userID, policy, permission)
}

func normalizeScopeLevel(raw string, optional bool) (domain.PolicyScopeLevel, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if optional {
			return "", nil
		}
		return "", errors.New("scopeLevel is required")
	}
	scope := domain.PolicyScopeLevel(value)
	switch scope {
	case domain.PolicyScopeLevelPlatform, domain.PolicyScopeLevelWorkspace, domain.PolicyScopeLevelProject:
		return scope, nil
	default:
		return "", errors.New("invalid scopeLevel")
	}
}

func normalizeCategory(raw string, optional bool) (domain.PolicyCategory, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if optional {
			return "", nil
		}
		return "", errors.New("category is required")
	}
	category := domain.PolicyCategory(value)
	switch category {
	case domain.PolicyCategoryPodSecurity, domain.PolicyCategoryImage, domain.PolicyCategoryResource,
		domain.PolicyCategoryLabel, domain.PolicyCategoryNetwork, domain.PolicyCategoryAdmission:
		return category, nil
	default:
		return "", errors.New("invalid category")
	}
}

func normalizeEnforcementMode(raw string, optional bool) (domain.PolicyEnforcementMode, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if optional {
			return "", nil
		}
		return "", errors.New("enforcementMode is required")
	}
	mode := domain.PolicyEnforcementMode(value)
	switch mode {
	case domain.PolicyEnforcementModeAudit, domain.PolicyEnforcementModeAlert, domain.PolicyEnforcementModeWarn, domain.PolicyEnforcementModeEnforce:
		return mode, nil
	default:
		return "", errors.New("invalid enforcementMode")
	}
}

func normalizeRiskLevel(raw string, optional bool) (domain.PolicyRiskLevel, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if optional {
			return "", nil
		}
		return "", errors.New("riskLevel is required")
	}
	risk := domain.PolicyRiskLevel(value)
	switch risk {
	case domain.PolicyRiskLevelLow, domain.PolicyRiskLevelMedium, domain.PolicyRiskLevelHigh, domain.PolicyRiskLevelCritical:
		return risk, nil
	default:
		return "", errors.New("invalid riskLevel")
	}
}

func normalizePolicyStatus(raw string, optional bool) (domain.PolicyStatus, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if optional {
			return "", nil
		}
		return "", errors.New("status is required")
	}
	status := domain.PolicyStatus(value)
	switch status {
	case domain.PolicyStatusDraft, domain.PolicyStatusActive, domain.PolicyStatusDisabled, domain.PolicyStatusArchived:
		return status, nil
	default:
		return "", errors.New("invalid status")
	}
}

func normalizeRolloutStage(raw string, optional bool) (domain.PolicyRolloutStage, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if optional {
			return "", nil
		}
		return "", errors.New("rolloutStage is required")
	}
	stage := domain.PolicyRolloutStage(value)
	switch stage {
	case domain.PolicyRolloutStagePilot, domain.PolicyRolloutStageCanary, domain.PolicyRolloutStageBroad, domain.PolicyRolloutStageFull:
		return stage, nil
	default:
		return "", errors.New("invalid rolloutStage")
	}
}

func normalizeRemediationStatus(raw string, optional bool) (domain.PolicyRemediationStatus, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if optional {
			return "", nil
		}
		return "", errors.New("remediationStatus is required")
	}
	status := domain.PolicyRemediationStatus(value)
	switch status {
	case domain.PolicyRemediationOpen, domain.PolicyRemediationInProgress, domain.PolicyRemediationMitigated, domain.PolicyRemediationClosed:
		return status, nil
	default:
		return "", errors.New("invalid remediationStatus")
	}
}

func normalizeExceptionStatus(raw string, optional bool) (domain.PolicyExceptionStatus, error) {
	value := strings.ToLower(strings.TrimSpace(raw))
	if value == "" {
		if optional {
			return "", nil
		}
		return "", errors.New("status is required")
	}
	status := domain.PolicyExceptionStatus(value)
	switch status {
	case domain.PolicyExceptionPending,
		domain.PolicyExceptionApproved,
		domain.PolicyExceptionRejected,
		domain.PolicyExceptionActive,
		domain.PolicyExceptionExpired,
		domain.PolicyExceptionRevoked:
		return status, nil
	default:
		return "", errors.New("invalid status")
	}
}

func resolveReviewedExceptionStatus(
	decision string,
	item *domain.PolicyExceptionRequest,
	now time.Time,
) (domain.PolicyExceptionStatus, error) {
	if item == nil {
		return "", errors.New("exception not found")
	}
	switch decision {
	case "approve":
		if item.ExpiresAt != nil && !item.ExpiresAt.After(now) {
			return domain.PolicyExceptionExpired, nil
		}
		if item.CreatedAt.After(now) {
			return domain.PolicyExceptionApproved, nil
		}
		return domain.PolicyExceptionActive, nil
	case "reject":
		return domain.PolicyExceptionRejected, nil
	case "revoke":
		return domain.PolicyExceptionRevoked, nil
	default:
		return "", errors.New("invalid decision")
	}
}

func validateRemediationTransition(
	current domain.PolicyRemediationStatus,
	next domain.PolicyRemediationStatus,
) error {
	if current == "" {
		current = domain.PolicyRemediationOpen
	}
	if current == next {
		return nil
	}
	allowed := map[domain.PolicyRemediationStatus]map[domain.PolicyRemediationStatus]struct{}{
		domain.PolicyRemediationOpen: {
			domain.PolicyRemediationInProgress: {},
			domain.PolicyRemediationMitigated:  {},
			domain.PolicyRemediationClosed:     {},
		},
		domain.PolicyRemediationInProgress: {
			domain.PolicyRemediationMitigated: {},
			domain.PolicyRemediationClosed:    {},
			domain.PolicyRemediationOpen:      {},
		},
		domain.PolicyRemediationMitigated: {
			domain.PolicyRemediationClosed:     {},
			domain.PolicyRemediationOpen:       {},
			domain.PolicyRemediationInProgress: {},
		},
		domain.PolicyRemediationClosed: {
			domain.PolicyRemediationOpen: {},
		},
	}
	if nextSet, ok := allowed[current]; ok {
		if _, ok := nextSet[next]; ok {
			return nil
		}
	}
	return errors.New("invalid remediation transition")
}

func appendRemediationComment(base, comment string, at time.Time) string {
	trimmedComment := strings.TrimSpace(comment)
	if trimmedComment == "" {
		return base
	}
	line := "[remediation] " + at.Format(time.RFC3339) + " " + trimmedComment
	trimmedBase := strings.TrimSpace(base)
	if trimmedBase == "" {
		return line
	}
	return trimmedBase + "\n" + line
}

func marshalJSONMap(value map[string]any) (string, error) {
	if value == nil {
		return "{}", nil
	}
	payload, err := json.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func marshalStringArray(items []string) (string, error) {
	cleaned := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed == "" {
			continue
		}
		cleaned = append(cleaned, trimmed)
	}
	payload, err := json.Marshal(cleaned)
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func uint64PtrOrNil(value uint64) *uint64 {
	if value == 0 {
		return nil
	}
	v := value
	return &v
}

func uniqueUint64s(items []uint64) []uint64 {
	set := make(map[uint64]struct{}, len(items))
	out := make([]uint64, 0, len(items))
	for _, item := range items {
		if item == 0 {
			continue
		}
		if _, exists := set[item]; exists {
			continue
		}
		set[item] = struct{}{}
		out = append(out, item)
	}
	return out
}

func IsScopeDenied(err error) bool {
	return errors.Is(err, ErrSecurityPolicyScopeDenied)
}

func IsNotConfigured(err error) bool {
	return errors.Is(err, ErrSecurityPolicyNotConfigured)
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}
