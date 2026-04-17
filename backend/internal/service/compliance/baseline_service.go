package compliance

import (
	"context"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"
	baselineProvider "kbmanage/backend/internal/integration/compliance/baseline"
	"kbmanage/backend/internal/repository"
)

type BaselineListFilter struct {
	StandardType string
	Status       string
}

type CreateBaselineInput struct {
	Name         string
	StandardType string
	Version      string
	Description  string
	TargetLevels []string
	Rules        map[string]any
	Status       string
}

type UpdateBaselineInput struct {
	Name         *string
	Version      *string
	Description  *string
	TargetLevels []string
	Rules        map[string]any
	Status       *string
}

type BaselineService struct {
	repo             *repository.ComplianceBaselineRepository
	scope            *ScopeService
	snapshotProvider baselineProvider.Provider
}

func NewBaselineService(repo *repository.ComplianceBaselineRepository, scope *ScopeService, snapshotProvider baselineProvider.Provider) *BaselineService {
	if snapshotProvider == nil {
		snapshotProvider = baselineProvider.NewStaticProvider()
	}
	return &BaselineService{repo: repo, scope: scope, snapshotProvider: snapshotProvider}
}

func (s *BaselineService) List(ctx context.Context, userID uint64, filter BaselineListFilter) ([]domain.ComplianceBaseline, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	if err := s.scope.ValidateScope(ctx, userID, 0, 0, PermissionComplianceRead); err != nil {
		return nil, err
	}
	return s.repo.List(ctx, repository.ComplianceBaselineListFilter{
		StandardType: domain.ComplianceStandardType(strings.TrimSpace(filter.StandardType)),
		Status:       domain.ComplianceBaselineStatus(strings.TrimSpace(filter.Status)),
	})
}

func (s *BaselineService) Create(ctx context.Context, userID uint64, input CreateBaselineInput) (*domain.ComplianceBaseline, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	if err := s.scope.ValidateScope(ctx, userID, 0, 0, PermissionComplianceManageBaseline); err != nil {
		return nil, err
	}
	standardType, err := normalizeBaselineStandardType(input.StandardType)
	if err != nil {
		return nil, err
	}
	status, err := normalizeBaselineStatus(input.Status, true)
	if err != nil {
		return nil, err
	}
	if status == "" {
		status = domain.ComplianceBaselineStatusDraft
	}
	rulesJSON, err := marshalJSON(cloneMap(input.Rules))
	if err != nil {
		return nil, err
	}
	targetLevels := sortAndCompactStrings(input.TargetLevels)
	targetLevelsJSON, err := marshalJSON(targetLevels)
	if err != nil {
		return nil, err
	}
	name := strings.TrimSpace(input.Name)
	version := strings.TrimSpace(input.Version)
	if name == "" || version == "" {
		return nil, errors.New("name and version are required")
	}
	createdBy := userID
	item := &domain.ComplianceBaseline{
		Name:             name,
		StandardType:     standardType,
		Version:          version,
		Description:      strings.TrimSpace(input.Description),
		TargetLevelsJSON: targetLevelsJSON,
		RulesJSON:        rulesJSON,
		RuleCount:        ruleCountFromRules(input.Rules),
		Status:           status,
		CreatedBy:        &createdBy,
		UpdatedBy:        &createdBy,
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *BaselineService) Get(ctx context.Context, userID, baselineID uint64) (*domain.ComplianceBaseline, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	item, err := s.repo.GetByID(ctx, baselineID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.ValidateBaselineScope(ctx, userID, item, PermissionComplianceRead); err != nil {
		return nil, err
	}
	return item, nil
}

func (s *BaselineService) Update(ctx context.Context, userID, baselineID uint64, input UpdateBaselineInput) (*domain.ComplianceBaseline, error) {
	if s == nil || s.repo == nil {
		return nil, ErrComplianceNotConfigured
	}
	item, err := s.repo.GetByID(ctx, baselineID)
	if err != nil {
		return nil, err
	}
	if err := s.scope.ValidateBaselineScope(ctx, userID, item, PermissionComplianceManageBaseline); err != nil {
		return nil, err
	}
	updates := map[string]any{}
	if input.Name != nil {
		name := strings.TrimSpace(*input.Name)
		if name == "" {
			return nil, errors.New("name is required")
		}
		updates["name"] = name
	}
	if input.Version != nil {
		version := strings.TrimSpace(*input.Version)
		if version == "" {
			return nil, errors.New("version is required")
		}
		updates["version"] = version
	}
	if input.Description != nil {
		updates["description"] = strings.TrimSpace(*input.Description)
	}
	if input.TargetLevels != nil {
		payload, err := marshalJSON(sortAndCompactStrings(input.TargetLevels))
		if err != nil {
			return nil, err
		}
		updates["target_levels_json"] = payload
	}
	if input.Rules != nil {
		payload, err := marshalJSON(cloneMap(input.Rules))
		if err != nil {
			return nil, err
		}
		updates["rules_json"] = payload
		updates["rule_count"] = ruleCountFromRules(input.Rules)
	}
	if input.Status != nil {
		status, err := normalizeBaselineStatus(*input.Status, false)
		if err != nil {
			return nil, err
		}
		updates["status"] = status
	}
	if len(updates) == 0 {
		return item, nil
	}
	updates["updated_by"] = uint64Ptr(userID)
	if err := s.repo.UpdateFields(ctx, baselineID, updates); err != nil {
		return nil, err
	}
	return s.repo.GetByID(ctx, baselineID)
}

func (s *BaselineService) GetSnapshot(ctx context.Context, userID, baselineID uint64) (domain.ComplianceBaselineSnapshot, error) {
	item, err := s.Get(ctx, userID, baselineID)
	if err != nil {
		return domain.ComplianceBaselineSnapshot{}, err
	}
	return s.snapshotProvider.BuildSnapshot(ctx, item)
}

func normalizeBaselineStandardType(value string) (domain.ComplianceStandardType, error) {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case string(domain.ComplianceStandardTypeCIS):
		return domain.ComplianceStandardTypeCIS, nil
	case string(domain.ComplianceStandardTypeSTIG):
		return domain.ComplianceStandardTypeSTIG, nil
	case string(domain.ComplianceStandardTypePlatformBaseline):
		return domain.ComplianceStandardTypePlatformBaseline, nil
	default:
		return "", errors.New("standardType must be cis, stig or platform-baseline")
	}
}

func normalizeBaselineStatus(value string, allowEmpty bool) (domain.ComplianceBaselineStatus, error) {
	switch strings.TrimSpace(strings.ToLower(value)) {
	case "":
		if allowEmpty {
			return "", nil
		}
		return "", errors.New("status is required")
	case string(domain.ComplianceBaselineStatusDraft):
		return domain.ComplianceBaselineStatusDraft, nil
	case string(domain.ComplianceBaselineStatusActive):
		return domain.ComplianceBaselineStatusActive, nil
	case string(domain.ComplianceBaselineStatusDisabled):
		return domain.ComplianceBaselineStatusDisabled, nil
	case string(domain.ComplianceBaselineStatusArchived):
		return domain.ComplianceBaselineStatusArchived, nil
	default:
		return "", errors.New("invalid baseline status")
	}
}

func ruleCountFromRules(rules map[string]any) int {
	if len(rules) == 0 {
		return 0
	}
	if controls, ok := rules["controls"].([]any); ok {
		return len(controls)
	}
	return len(rules)
}
