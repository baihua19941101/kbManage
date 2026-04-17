package compliance

import (
	"context"
	"errors"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type ScanProfileListFilter struct {
	WorkspaceID  uint64
	ProjectID    uint64
	ScopeType    string
	ScheduleMode string
	Status       string
}

type CreateScanProfileInput struct {
	Name           string
	BaselineID     uint64
	WorkspaceID    uint64
	ProjectID      uint64
	ScopeType      string
	ClusterRefs    []uint64
	NodeSelectors  []map[string]string
	NamespaceRefs  []string
	ResourceKinds  []string
	ScheduleMode   string
	CronExpression string
	Status         string
}

type UpdateScanProfileInput struct {
	Name           *string
	BaselineID     *uint64
	ScopeType      *string
	ClusterRefs    []uint64
	NodeSelectors  []map[string]string
	NamespaceRefs  []string
	ResourceKinds  []string
	ScheduleMode   *string
	CronExpression *string
	Status         *string
}

type ScanProfileService struct {
	repo          *repository.ComplianceScanProfileRepository
	baselineRepo  *repository.ComplianceBaselineRepository
	scope         *ScopeService
	scheduleCache *ScheduleCache
}

func NewScanProfileService(repo *repository.ComplianceScanProfileRepository, baselineRepo *repository.ComplianceBaselineRepository, scope *ScopeService, scheduleCache *ScheduleCache) *ScanProfileService {
	return &ScanProfileService{repo: repo, baselineRepo: baselineRepo, scope: scope, scheduleCache: scheduleCache}
}

func (s *ScanProfileService) List(ctx context.Context, userID uint64, filter ScanProfileListFilter) ([]domain.ScanProfile, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	if err := s.scope.ValidateScope(ctx, userID, filter.WorkspaceID, filter.ProjectID, PermissionComplianceRead); err != nil {
		return nil, err
	}
	items, err := s.repo.List(ctx, repository.ComplianceScanProfileListFilter{
		WorkspaceID:  uint64Ptr(filter.WorkspaceID),
		ProjectID:    uint64Ptr(filter.ProjectID),
		ScopeType:    domain.ComplianceScopeType(strings.TrimSpace(filter.ScopeType)),
		ScheduleMode: domain.ComplianceScheduleMode(strings.TrimSpace(filter.ScheduleMode)),
		Status:       domain.ComplianceScanProfileStatus(strings.TrimSpace(filter.Status)),
	})
	if err != nil {
		return nil, err
	}
	return s.scope.FilterScanProfiles(ctx, userID, items, PermissionComplianceRead), nil
}

func (s *ScanProfileService) Create(ctx context.Context, userID uint64, input CreateScanProfileInput) (*domain.ScanProfile, error) {
	if s == nil || s.repo == nil || s.baselineRepo == nil {
		return nil, ErrComplianceNotConfigured
	}
	if err := s.scope.ValidateScope(ctx, userID, input.WorkspaceID, input.ProjectID, PermissionComplianceExecuteScan); err != nil {
		return nil, err
	}
	baseline, err := s.baselineRepo.GetByID(ctx, input.BaselineID)
	if err != nil {
		return nil, err
	}
	if baseline.Status == domain.ComplianceBaselineStatusArchived {
		return nil, errors.New("archived baseline cannot be referenced")
	}
	item, err := s.buildProfile(userID, input)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	if err := s.persistSchedule(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ScanProfileService) Get(ctx context.Context, userID, profileID uint64) (*domain.ScanProfile, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	item, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.ValidateProfileScope(ctx, userID, item, PermissionComplianceRead); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ScanProfileService) Update(ctx context.Context, userID, profileID uint64, input UpdateScanProfileInput) (*domain.ScanProfile, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	item, err := s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.ValidateProfileScope(ctx, userID, item, PermissionComplianceExecuteScan); err != nil {
		return nil, err
	}
	updates := map[string]any{"updated_by": uint64Ptr(userID)}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, errors.New("name is required")
		}
		updates["name"] = name
	}
	if input.BaselineID != nil {
		if *input.BaselineID == 0 {
			return nil, errors.New("baselineId is required")
		}
		updates["baseline_id"] = *input.BaselineID
	}
	if input.ScopeType != nil {
		scopeType, err := normalizeProfileScopeType(*input.ScopeType)
		if err != nil {
			return nil, err
		}
		updates["scope_type"] = scopeType
	}
	if input.ClusterRefs != nil {
		payload, err := marshalJSON(input.ClusterRefs)
		if err != nil {
			return nil, err
		}
		updates["cluster_refs_json"] = payload
	}
	if input.NodeSelectors != nil {
		payload, err := marshalJSON(input.NodeSelectors)
		if err != nil {
			return nil, err
		}
		updates["node_selectors_json"] = payload
	}
	if input.NamespaceRefs != nil {
		payload, err := marshalJSON(sortAndCompactStrings(input.NamespaceRefs))
		if err != nil {
			return nil, err
		}
		updates["namespace_refs_json"] = payload
	}
	if input.ResourceKinds != nil {
		payload, err := marshalJSON(sortAndCompactStrings(input.ResourceKinds))
		if err != nil {
			return nil, err
		}
		updates["resource_kinds_json"] = payload
	}
	if input.ScheduleMode != nil {
		mode, err := normalizeScheduleMode(*input.ScheduleMode)
		if err != nil {
			return nil, err
		}
		updates["schedule_mode"] = mode
	}
	if input.CronExpression != nil {
		updates["cron_expression"] = strings.TrimSpace(*input.CronExpression)
	}
	if input.Status != nil {
		status, err := normalizeProfileStatus(*input.Status, false)
		if err != nil {
			return nil, err
		}
		updates["status"] = status
	}
	candidate := *item
	applyProfileUpdates(&candidate, updates)
	if err := validateProfile(candidate); err != nil {
		return nil, err
	}
	if err := s.repo.UpdateFields(ctx, profileID, updates); err != nil {
		return nil, err
	}
	item, err = s.repo.GetByID(ctx, profileID)
	if err != nil {
		return nil, err
	}
	if err := s.persistSchedule(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ScanProfileService) buildProfile(userID uint64, input CreateScanProfileInput) (*domain.ScanProfile, error) {
	scopeType, err := normalizeProfileScopeType(input.ScopeType)
	if err != nil {
		return nil, err
	}
	scheduleMode, err := normalizeScheduleMode(input.ScheduleMode)
	if err != nil {
		return nil, err
	}
	status, err := normalizeProfileStatus(input.Status, true)
	if err != nil {
		return nil, err
	}
	if status == "" {
		status = domain.ComplianceScanProfileStatusDraft
	}
	clusterRefsJSON, err := marshalJSON(input.ClusterRefs)
	if err != nil {
		return nil, err
	}
	nodeSelectorsJSON, err := marshalJSON(input.NodeSelectors)
	if err != nil {
		return nil, err
	}
	namespaceRefsJSON, err := marshalJSON(sortAndCompactStrings(input.NamespaceRefs))
	if err != nil {
		return nil, err
	}
	resourceKindsJSON, err := marshalJSON(sortAndCompactStrings(input.ResourceKinds))
	if err != nil {
		return nil, err
	}
	createdBy := userID
	item := &domain.ScanProfile{
		Name:              strings.TrimSpace(input.Name),
		BaselineID:        input.BaselineID,
		WorkspaceID:       uint64Ptr(input.WorkspaceID),
		ProjectID:         uint64Ptr(input.ProjectID),
		ScopeType:         scopeType,
		ClusterRefsJSON:   clusterRefsJSON,
		NodeSelectorsJSON: nodeSelectorsJSON,
		NamespaceRefsJSON: namespaceRefsJSON,
		ResourceKindsJSON: resourceKindsJSON,
		ScheduleMode:      scheduleMode,
		CronExpression:    strings.TrimSpace(input.CronExpression),
		Status:            status,
		CreatedBy:         &createdBy,
		UpdatedBy:         &createdBy,
	}
	if err := validateProfile(*item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *ScanProfileService) persistSchedule(ctx context.Context, item *domain.ScanProfile) error {
	if s == nil || s.scheduleCache == nil || item == nil {
		return nil
	}
	return s.scheduleCache.SetProfileSchedule(ctx, item.ID, ProfileScheduleSnapshot{CronExpression: item.CronExpression, ScheduleMode: string(item.ScheduleMode), UpdatedAt: time.Now().UTC()})
}

func applyProfileUpdates(item *domain.ScanProfile, updates map[string]any) {
	if item == nil {
		return
	}
	if v, ok := updates["name"].(string); ok {
		item.Name = v
	}
	if v, ok := updates["baseline_id"].(uint64); ok {
		item.BaselineID = v
	}
	if v, ok := updates["scope_type"].(domain.ComplianceScopeType); ok {
		item.ScopeType = v
	}
	if v, ok := updates["cluster_refs_json"].(string); ok {
		item.ClusterRefsJSON = v
	}
	if v, ok := updates["node_selectors_json"].(string); ok {
		item.NodeSelectorsJSON = v
	}
	if v, ok := updates["namespace_refs_json"].(string); ok {
		item.NamespaceRefsJSON = v
	}
	if v, ok := updates["resource_kinds_json"].(string); ok {
		item.ResourceKindsJSON = v
	}
	if v, ok := updates["schedule_mode"].(domain.ComplianceScheduleMode); ok {
		item.ScheduleMode = v
	}
	if v, ok := updates["cron_expression"].(string); ok {
		item.CronExpression = v
	}
	if v, ok := updates["status"].(domain.ComplianceScanProfileStatus); ok {
		item.Status = v
	}
}

func validateProfile(item domain.ScanProfile) error {
	if strings.TrimSpace(item.Name) == "" {
		return errors.New("name is required")
	}
	if item.BaselineID == 0 {
		return errors.New("baselineId is required")
	}
	if item.ScopeType == "" {
		return errors.New("scopeType is required")
	}
	if len(unmarshalUint64Slice(item.ClusterRefsJSON)) == 0 {
		return errors.New("clusterRefs is required")
	}
	if item.ScopeType == domain.ComplianceScopeTypeNode && len(unmarshalNodeSelectors(item.NodeSelectorsJSON)) == 0 {
		return errors.New("nodeSelectors is required for node scope")
	}
	if item.ScopeType == domain.ComplianceScopeTypeNamespace && len(unmarshalStringSlice(item.NamespaceRefsJSON)) == 0 {
		return errors.New("namespaceRefs is required for namespace scope")
	}
	if item.ScopeType == domain.ComplianceScopeTypeResourceSet && len(unmarshalStringSlice(item.ResourceKindsJSON)) == 0 {
		return errors.New("resourceKinds is required for resource-set scope")
	}
	if item.ScheduleMode == domain.ComplianceScheduleModeScheduled && strings.TrimSpace(item.CronExpression) == "" {
		return errors.New("cronExpression is required for scheduled profiles")
	}
	return nil
}

func normalizeProfileScopeType(value string) (domain.ComplianceScopeType, error) {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case string(domain.ComplianceScopeTypeCluster):
		return domain.ComplianceScopeTypeCluster, nil
	case string(domain.ComplianceScopeTypeNode):
		return domain.ComplianceScopeTypeNode, nil
	case string(domain.ComplianceScopeTypeNamespace):
		return domain.ComplianceScopeTypeNamespace, nil
	case string(domain.ComplianceScopeTypeResourceSet):
		return domain.ComplianceScopeTypeResourceSet, nil
	default:
		return "", errors.New("scopeType must be cluster, node, namespace or resource-set")
	}
}

func normalizeScheduleMode(value string) (domain.ComplianceScheduleMode, error) {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "", string(domain.ComplianceScheduleModeManual):
		return domain.ComplianceScheduleModeManual, nil
	case string(domain.ComplianceScheduleModeScheduled):
		return domain.ComplianceScheduleModeScheduled, nil
	default:
		return "", errors.New("scheduleMode must be manual or scheduled")
	}
}

func normalizeProfileStatus(value string, allowEmpty bool) (domain.ComplianceScanProfileStatus, error) {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "":
		if allowEmpty {
			return "", nil
		}
		return "", errors.New("status is required")
	case string(domain.ComplianceScanProfileStatusDraft):
		return domain.ComplianceScanProfileStatusDraft, nil
	case string(domain.ComplianceScanProfileStatusActive):
		return domain.ComplianceScanProfileStatusActive, nil
	case string(domain.ComplianceScanProfileStatusPaused):
		return domain.ComplianceScanProfileStatusPaused, nil
	case string(domain.ComplianceScanProfileStatusArchived):
		return domain.ComplianceScanProfileStatusArchived, nil
	default:
		return "", errors.New("invalid scan profile status")
	}
}
