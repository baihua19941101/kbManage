package observability

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type AlertRuleService struct {
	repo  *repository.AlertRuleRepository
	scope *ScopeService
}

type UpsertAlertRuleRequest struct {
	Name                 string                 `json:"name"`
	Description          string                 `json:"description"`
	Severity             domain.AlertSeverity   `json:"severity"`
	ScopeSnapshotJSON    string                 `json:"scopeSnapshot"`
	ConditionExpression  string                 `json:"conditionExpression"`
	EvaluationWindow     string                 `json:"evaluationWindow"`
	NotificationStrategy string                 `json:"notificationStrategy"`
	Status               domain.AlertRuleStatus `json:"status"`
}

func NewAlertRuleService(repo *repository.AlertRuleRepository, scope *ScopeService) *AlertRuleService {
	return &AlertRuleService{repo: repo, scope: scope}
}

func (s *AlertRuleService) List(ctx context.Context, status domain.AlertRuleStatus) ([]domain.AlertRule, error) {
	if s == nil || s.repo == nil {
		return []domain.AlertRule{}, nil
	}
	items, err := s.repo.List(ctx, status)
	if err != nil {
		return nil, normalizeObservabilityProviderError(err)
	}
	access, ok := AccessContextFromContext(ctx)
	if !ok || s.scope == nil {
		return items, nil
	}

	filtered := make([]domain.AlertRule, 0, len(items))
	for _, item := range items {
		if s.canAccessRuleScope(ctx, access.UserID, item.ScopeSnapshotJSON) {
			filtered = append(filtered, item)
		}
	}
	return filtered, nil
}

func (s *AlertRuleService) Get(ctx context.Context, id uint64) (*domain.AlertRule, error) {
	if s == nil || s.repo == nil {
		return nil, ErrObservabilityUnavailable
	}
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, normalizeObservabilityProviderError(err)
	}
	access, ok := AccessContextFromContext(ctx)
	if ok && s.scope != nil && !s.canAccessRuleScope(ctx, access.UserID, item.ScopeSnapshotJSON) {
		return nil, ErrObservabilityScopeDenied
	}
	return item, nil
}

func (s *AlertRuleService) Create(ctx context.Context, createdBy uint64, req UpsertAlertRuleRequest) (*domain.AlertRule, error) {
	scopeSnapshot, err := s.normalizeAndAuthorizeScopeSnapshot(ctx, req.ScopeSnapshotJSON)
	if err != nil {
		return nil, err
	}
	item := &domain.AlertRule{
		Name:                 strings.TrimSpace(req.Name),
		Description:          strings.TrimSpace(req.Description),
		Severity:             normalizeAlertSeverity(req.Severity),
		ScopeSnapshotJSON:    scopeSnapshot,
		ConditionExpression:  strings.TrimSpace(req.ConditionExpression),
		EvaluationWindow:     strings.TrimSpace(req.EvaluationWindow),
		NotificationStrategy: strings.TrimSpace(req.NotificationStrategy),
		Status:               normalizeAlertRuleStatus(req.Status),
		CreatedBy:            createdBy,
	}
	if s == nil || s.repo == nil {
		return item, nil
	}
	if err := s.repo.Create(ctx, item); err != nil {
		return nil, normalizeObservabilityProviderError(err)
	}
	return item, nil
}

func (s *AlertRuleService) Update(ctx context.Context, id uint64, req UpsertAlertRuleRequest) (*domain.AlertRule, error) {
	if s == nil || s.repo == nil {
		return nil, ErrObservabilityUnavailable
	}
	scopeSnapshot, err := s.normalizeAndAuthorizeScopeSnapshot(ctx, req.ScopeSnapshotJSON)
	if err != nil {
		return nil, err
	}
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, normalizeObservabilityProviderError(err)
	}
	access, ok := AccessContextFromContext(ctx)
	if ok && s.scope != nil && !s.canAccessRuleScope(ctx, access.UserID, item.ScopeSnapshotJSON) {
		return nil, ErrObservabilityScopeDenied
	}
	item.Name = strings.TrimSpace(req.Name)
	item.Description = strings.TrimSpace(req.Description)
	item.Severity = normalizeAlertSeverity(req.Severity)
	item.ScopeSnapshotJSON = scopeSnapshot
	item.ConditionExpression = strings.TrimSpace(req.ConditionExpression)
	item.EvaluationWindow = strings.TrimSpace(req.EvaluationWindow)
	item.NotificationStrategy = strings.TrimSpace(req.NotificationStrategy)
	item.Status = normalizeAlertRuleStatus(req.Status)
	if err := s.repo.Update(ctx, item); err != nil {
		return nil, normalizeObservabilityProviderError(err)
	}
	return item, nil
}

func (s *AlertRuleService) Delete(ctx context.Context, id uint64) error {
	if s == nil || s.repo == nil {
		return nil
	}
	item, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return normalizeObservabilityProviderError(err)
	}
	access, ok := AccessContextFromContext(ctx)
	if ok && s.scope != nil && !s.canAccessRuleScope(ctx, access.UserID, item.ScopeSnapshotJSON) {
		return ErrObservabilityScopeDenied
	}
	return s.repo.Delete(ctx, id)
}

func (s *AlertRuleService) canAccessRuleScope(ctx context.Context, userID uint64, scopeSnapshot string) bool {
	if s == nil || s.scope == nil {
		return true
	}
	scopeFilter := decodeScopeSnapshot(scopeSnapshot)
	if len(scopeFilter.ClusterIDs) == 0 && len(scopeFilter.WorkspaceIDs) == 0 && len(scopeFilter.ProjectIDs) == 0 {
		return false
	}
	_, err := s.scope.FilterByScope(ctx, userID, scopeFilter)
	return err == nil
}

func (s *AlertRuleService) normalizeAndAuthorizeScopeSnapshot(ctx context.Context, raw string) (string, error) {
	access, ok := AccessContextFromContext(ctx)
	if !ok || s == nil || s.scope == nil {
		return strings.TrimSpace(raw), nil
	}

	scopeFilter := decodeScopeSnapshot(raw)
	if len(scopeFilter.ClusterIDs) == 0 && len(scopeFilter.WorkspaceIDs) == 0 && len(scopeFilter.ProjectIDs) == 0 {
		scopeFilter = ScopeFilter{
			ClusterIDs:   append([]uint64(nil), access.ClusterIDs...),
			WorkspaceIDs: append([]uint64(nil), access.WorkspaceIDs...),
			ProjectIDs:   append([]uint64(nil), access.ProjectIDs...),
		}
	}
	filtered, err := s.scope.FilterByScope(ctx, access.UserID, scopeFilter)
	if err != nil {
		return "", ErrObservabilityScopeDenied
	}
	if len(filtered.ClusterIDs) == 0 && len(filtered.WorkspaceIDs) == 0 && len(filtered.ProjectIDs) == 0 {
		return "", ErrObservabilityScopeDenied
	}
	normalized, err := json.Marshal(map[string]any{
		"clusterIds":   filtered.ClusterIDs,
		"workspaceIds": filtered.WorkspaceIDs,
		"projectIds":   filtered.ProjectIDs,
	})
	if err != nil {
		return "", errors.New("invalid observability scope snapshot")
	}
	return string(normalized), nil
}

func decodeScopeSnapshot(raw string) ScopeFilter {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return ScopeFilter{}
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(trimmed), &payload); err != nil {
		return ScopeFilter{}
	}
	return ScopeFilter{
		ClusterIDs:   parseScopeIDs(payload, "clusterIds", "clusterId"),
		WorkspaceIDs: parseScopeIDs(payload, "workspaceIds", "workspaceId"),
		ProjectIDs:   parseScopeIDs(payload, "projectIds", "projectId"),
	}
}

func parseScopeIDs(payload map[string]any, keys ...string) []uint64 {
	out := make([]uint64, 0)
	for _, key := range keys {
		v, ok := payload[key]
		if !ok {
			continue
		}
		switch typed := v.(type) {
		case float64:
			if typed > 0 {
				out = append(out, uint64(typed))
			}
		case string:
			parts := strings.Split(typed, ",")
			for _, part := range parts {
				idText := strings.TrimSpace(part)
				if idText == "" {
					continue
				}
				if id, err := jsonNumberToUint64(idText); err == nil {
					out = append(out, id)
				}
			}
		case []any:
			for _, item := range typed {
				switch n := item.(type) {
				case float64:
					if n > 0 {
						out = append(out, uint64(n))
					}
				case string:
					if id, err := jsonNumberToUint64(strings.TrimSpace(n)); err == nil {
						out = append(out, id)
					}
				}
			}
		}
	}
	return dedupeUint64(out)
}

func jsonNumberToUint64(raw string) (uint64, error) {
	if strings.TrimSpace(raw) == "" {
		return 0, errors.New("empty number")
	}
	var id uint64
	if err := json.Unmarshal([]byte(raw), &id); err != nil {
		return 0, err
	}
	if id == 0 {
		return 0, errors.New("zero")
	}
	return id, nil
}

func normalizeAlertSeverity(in domain.AlertSeverity) domain.AlertSeverity {
	switch in {
	case domain.AlertSeverityInfo, domain.AlertSeverityWarning, domain.AlertSeverityCritical:
		return in
	default:
		return domain.AlertSeverityWarning
	}
}

func normalizeAlertRuleStatus(in domain.AlertRuleStatus) domain.AlertRuleStatus {
	switch in {
	case domain.AlertRuleStatusEnabled, domain.AlertRuleStatusDisabled:
		return in
	default:
		return domain.AlertRuleStatusEnabled
	}
}
