package audit

import (
	"bytes"
	"context"
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"

	"gorm.io/gorm"
)

var (
	errInvalidTimeRange      = errors.New("startAt must be earlier than endAt")
	errOperatorIDRequired    = errors.New("operator id is required")
	errAuditExportIDRequired = errors.New("task id is required")
	errAuditExportNotReady   = errors.New("export task is not ready")
	emailPattern             = regexp.MustCompile(`^[^@\s]+@[^@\s]+\.[^@\s]+$`)
)

type QueryEventsRequest struct {
	StartAt     *time.Time
	EndAt       *time.Time
	ActorID     *uint64
	ClusterID   *uint64
	WorkspaceID *uint64
	ProjectID   *uint64
	Action      string
	Outcome     string
	Result      string
	Resource    string
	Limit       int
	ViewerID    uint64
}

type SubmitExportRequest struct {
	StartAt     *time.Time
	EndAt       *time.Time
	ActorID     *uint64
	ClusterID   *uint64
	WorkspaceID *uint64
	ProjectID   *uint64
	Action      string
	Outcome     string
	Result      string
	Resource    string
}

type ExportDownload struct {
	FileName    string
	ContentType string
	Data        []byte
}

type Service struct {
	auditRepo       *repository.AuditRepository
	auditExportRepo *repository.AuditExportRepository
	scopeAccess     *authSvc.ScopeAccessService
}

func NewService(auditRepo *repository.AuditRepository, auditExportRepo *repository.AuditExportRepository, scopeAccess ...*authSvc.ScopeAccessService) *Service {
	var scopedAccess *authSvc.ScopeAccessService
	if len(scopeAccess) > 0 {
		scopedAccess = scopeAccess[0]
	}
	return &Service{
		auditRepo:       auditRepo,
		auditExportRepo: auditExportRepo,
		scopeAccess:     scopedAccess,
	}
}

func (s *Service) QueryEvents(ctx context.Context, req QueryEventsRequest) ([]domain.AuditEvent, error) {
	if req.StartAt != nil && req.EndAt != nil && req.StartAt.After(*req.EndAt) {
		return nil, errInvalidTimeRange
	}
	if s.auditRepo == nil {
		return []domain.AuditEvent{}, nil
	}

	outcome, err := normalizeOutcome(firstNonEmpty(req.Result, req.Outcome))
	if err != nil {
		return nil, err
	}

	items, err := s.auditRepo.Query(ctx, repository.AuditQuery{
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		ActorID:     req.ActorID,
		ClusterID:   req.ClusterID,
		WorkspaceID: req.WorkspaceID,
		ProjectID:   req.ProjectID,
		Action:      strings.TrimSpace(req.Action),
		Outcome:     outcome,
		Result:      outcome,
		Resource:    strings.TrimSpace(req.Resource),
		Limit:       req.Limit,
	})
	if err != nil {
		return nil, err
	}

	return s.filterVisibleEvents(ctx, req.ViewerID, req.ClusterID, items)
}

func (s *Service) SubmitExport(ctx context.Context, operatorID uint64, req SubmitExportRequest) (*repository.AuditExportTask, error) {
	if operatorID == 0 {
		return nil, errOperatorIDRequired
	}
	if req.StartAt != nil && req.EndAt != nil && req.StartAt.After(*req.EndAt) {
		return nil, errInvalidTimeRange
	}
	if s.auditExportRepo == nil {
		return nil, errors.New("audit export repository is not initialized")
	}

	outcome, err := normalizeOutcome(firstNonEmpty(req.Result, req.Outcome))
	if err != nil {
		return nil, err
	}

	task, err := s.auditExportRepo.Create(ctx, operatorID, repository.AuditQuery{
		StartAt:     req.StartAt,
		EndAt:       req.EndAt,
		ActorID:     req.ActorID,
		ClusterID:   req.ClusterID,
		WorkspaceID: req.WorkspaceID,
		ProjectID:   req.ProjectID,
		Action:      strings.TrimSpace(req.Action),
		Outcome:     outcome,
		Result:      outcome,
		Resource:    strings.TrimSpace(req.Resource),
		Limit:       10000, // export uses a higher upper bound than list query.
	})
	if err != nil {
		return nil, err
	}

	if err := s.auditExportRepo.Enqueue(ctx, task.ID); err != nil {
		return nil, err
	}
	return task, nil
}

func (s *Service) GetExportTask(ctx context.Context, taskID string) (*repository.AuditExportTask, error) {
	if strings.TrimSpace(taskID) == "" {
		return nil, errAuditExportIDRequired
	}
	if s.auditExportRepo == nil {
		return nil, errors.New("audit export repository is not initialized")
	}
	return s.auditExportRepo.Get(ctx, taskID)
}

func (s *Service) GetExportTaskForViewer(ctx context.Context, taskID string, viewerID uint64) (*repository.AuditExportTask, error) {
	task, err := s.GetExportTask(ctx, taskID)
	if err != nil {
		return nil, err
	}
	if viewerID == 0 || task.OperatorID != viewerID {
		return nil, gorm.ErrRecordNotFound
	}
	return task, nil
}

func (s *Service) GetExportDownloadForViewer(ctx context.Context, taskID string, viewerID uint64) (*ExportDownload, error) {
	task, err := s.GetExportTaskForViewer(ctx, taskID, viewerID)
	if err != nil {
		return nil, err
	}
	if task.Status != repository.AuditExportStatusSucceeded {
		return nil, errAuditExportNotReady
	}
	if s.auditExportRepo == nil {
		return nil, errors.New("audit export repository is not initialized")
	}

	artifact, err := s.auditExportRepo.GetArtifact(ctx, taskID)
	if err != nil {
		return nil, err
	}
	return &ExportDownload{
		FileName:    artifact.FileName,
		ContentType: artifact.ContentType,
		Data:        append([]byte(nil), artifact.Data...),
	}, nil
}

func (s *Service) ProcessExportTask(ctx context.Context, taskID string) error {
	if strings.TrimSpace(taskID) == "" {
		return errAuditExportIDRequired
	}
	if s.auditExportRepo == nil {
		return errors.New("audit export repository is not initialized")
	}

	task, err := s.auditExportRepo.Get(ctx, taskID)
	if err != nil {
		return err
	}
	if err := s.auditExportRepo.MarkRunning(ctx, taskID); err != nil {
		return err
	}

	result, err := s.QueryEvents(ctx, QueryEventsRequest{
		StartAt:     task.Filters.StartAt,
		EndAt:       task.Filters.EndAt,
		ActorID:     task.Filters.ActorID,
		ClusterID:   task.Filters.ClusterID,
		WorkspaceID: task.Filters.WorkspaceID,
		ProjectID:   task.Filters.ProjectID,
		Action:      task.Filters.Action,
		Outcome:     task.Filters.Outcome,
		Result:      task.Filters.Result,
		Resource:    task.Filters.Resource,
		Limit:       task.Filters.Limit,
		ViewerID:    task.OperatorID,
	})
	if err != nil {
		_ = s.auditExportRepo.MarkFailed(ctx, taskID, err.Error())
		return err
	}

	payload, err := buildExportCSV(result)
	if err != nil {
		_ = s.auditExportRepo.MarkFailed(ctx, taskID, err.Error())
		return err
	}
	fileName := fmt.Sprintf("audit-export-%s.csv", taskID)
	if err := s.auditExportRepo.SaveArtifact(ctx, taskID, fileName, "text/csv; charset=utf-8", payload); err != nil {
		_ = s.auditExportRepo.MarkFailed(ctx, taskID, err.Error())
		return err
	}

	downloadURL := fmt.Sprintf("/api/v1/audits/exports/%s/download", taskID)
	return s.auditExportRepo.MarkSucceeded(ctx, taskID, len(result), downloadURL)
}

func normalizeOutcome(raw string) (string, error) {
	outcome := strings.TrimSpace(strings.ToLower(raw))
	if outcome == "" {
		return "", nil
	}

	switch domain.AuditOutcome(outcome) {
	case domain.AuditOutcomeSuccess, domain.AuditOutcomeDenied, domain.AuditOutcomeFailed:
		return outcome, nil
	default:
		return "", fmt.Errorf("invalid outcome: %s", raw)
	}
}

func (s *Service) filterVisibleEvents(
	ctx context.Context,
	viewerID uint64,
	clusterID *uint64,
	items []domain.AuditEvent,
) ([]domain.AuditEvent, error) {
	if len(items) == 0 {
		return []domain.AuditEvent{}, nil
	}

	constrained := false
	allowedClusterSet := make(map[uint64]struct{})
	if viewerID != 0 && s.scopeAccess != nil {
		constrained = true
		permissions := []string{"access:project:read", "gitops:read"}
		for _, permission := range permissions {
			allowedClusterIDs, hasScopeConstraint, err := s.scopeAccess.ListClusterIDsByPermission(ctx, viewerID, permission)
			if err != nil {
				return nil, err
			}
			if !hasScopeConstraint {
				constrained = false
			}
			for _, id := range allowedClusterIDs {
				allowedClusterSet[id] = struct{}{}
			}
		}
	}

	filtered := make([]domain.AuditEvent, 0, len(items))
	for _, event := range items {
		eventClusterID, hasEventCluster := resolveEventClusterID(event)

		if clusterID != nil {
			if !hasEventCluster || eventClusterID != *clusterID {
				continue
			}
		}

		if viewerID == 0 {
			filtered = append(filtered, event)
			continue
		}

		isOwner := event.ActorID != nil && *event.ActorID == viewerID
		if isOwner {
			filtered = append(filtered, event)
			continue
		}

		if hasEventCluster {
			if !constrained {
				filtered = append(filtered, event)
				continue
			}
			if _, ok := allowedClusterSet[eventClusterID]; ok {
				filtered = append(filtered, event)
			}
			continue
		}

		if !constrained {
			filtered = append(filtered, event)
		}
	}
	return filtered, nil
}

func resolveEventClusterID(event domain.AuditEvent) (uint64, bool) {
	if event.ClusterID != nil && *event.ClusterID != 0 {
		return *event.ClusterID, true
	}

	resourceType := strings.ToLower(strings.TrimSpace(event.ResourceType))
	resourceID := strings.TrimSpace(event.ResourceID)

	if resourceType == "cluster" {
		if id, err := strconv.ParseUint(resourceID, 10, 64); err == nil && id != 0 {
			return id, true
		}
	}
	if id, ok := authSvc.ParseClusterIDFromReference(resourceID); ok {
		return id, true
	}

	if len(event.Details) == 0 {
		return 0, false
	}
	var details map[string]any
	if err := json.Unmarshal(event.Details, &details); err != nil {
		return 0, false
	}
	if id, ok := extractUint64(details, "clusterId", "clusterID"); ok && id != 0 {
		return id, true
	}
	targetRef := strings.TrimSpace(extractString(details, "targetRef", "target_ref"))
	if targetRef == "" {
		return 0, false
	}
	return authSvc.ParseClusterIDFromReference(targetRef)
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		trimmed := strings.TrimSpace(value)
		if trimmed != "" {
			return trimmed
		}
	}
	return ""
}

func buildExportCSV(items []domain.AuditEvent) ([]byte, error) {
	var buf bytes.Buffer
	writer := csv.NewWriter(&buf)
	header := []string{
		"id", "requestId", "actorId", "clusterId", "workspaceId", "projectId",
		"auditCategory", "actionScope", "searchTags", "scopeSnapshot",
		"action", "resourceType", "resourceId", "result", "details", "createdAt",
	}
	if err := writer.Write(header); err != nil {
		return nil, err
	}

	for _, item := range items {
		detailsJSON, err := marshalMaskedDetails(item.Details)
		if err != nil {
			return nil, err
		}
		row := []string{
			strconv.FormatUint(item.ID, 10),
			item.RequestID,
			formatOptionalUint64(item.ActorID),
			formatOptionalUint64(item.ClusterID),
			formatOptionalUint64(item.WorkspaceID),
			formatOptionalUint64(item.ProjectID),
			item.AuditCategory,
			item.ActionScope,
			item.SearchTags,
			string(item.ScopeSnapshot),
			item.Action,
			item.ResourceType,
			item.ResourceID,
			string(item.Outcome),
			detailsJSON,
			item.CreatedAt.UTC().Format(time.RFC3339Nano),
		}
		if err := writer.Write(row); err != nil {
			return nil, err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func marshalMaskedDetails(raw json.RawMessage) (string, error) {
	trimmed := strings.TrimSpace(string(raw))
	if trimmed == "" {
		return "{}", nil
	}

	var payload any
	if err := json.Unmarshal(raw, &payload); err != nil {
		// Keep malformed legacy payload readable in export.
		masked := maskSensitiveValue(trimmed, "details")
		if text, ok := masked.(string); ok {
			return text, nil
		}
		return fmt.Sprintf("%v", masked), nil
	}

	masked := maskSensitiveValue(payload, "details")
	encoded, err := json.Marshal(masked)
	if err != nil {
		return "", err
	}
	return string(encoded), nil
}

func formatOptionalUint64(value *uint64) string {
	if value == nil || *value == 0 {
		return ""
	}
	return strconv.FormatUint(*value, 10)
}

func extractUint64(data map[string]any, keys ...string) (uint64, bool) {
	for _, key := range keys {
		value, ok := extractByKey(data, key)
		if !ok {
			continue
		}
		switch v := value.(type) {
		case uint64:
			if v != 0 {
				return v, true
			}
		case uint:
			if v != 0 {
				return uint64(v), true
			}
		case int:
			if v > 0 {
				return uint64(v), true
			}
		case int64:
			if v > 0 {
				return uint64(v), true
			}
		case float64:
			if v > 0 {
				return uint64(v), true
			}
		case string:
			id, err := strconv.ParseUint(strings.TrimSpace(v), 10, 64)
			if err == nil && id != 0 {
				return id, true
			}
		}
	}
	return 0, false
}

func extractString(data map[string]any, keys ...string) string {
	for _, key := range keys {
		value, ok := extractByKey(data, key)
		if !ok {
			continue
		}
		if text, ok := value.(string); ok {
			return text
		}
	}
	return ""
}

func extractByKey(data map[string]any, key string) (any, bool) {
	for existing, value := range data {
		if strings.EqualFold(strings.TrimSpace(existing), strings.TrimSpace(key)) {
			return value, true
		}
	}
	return nil, false
}

func maskSensitiveValue(value any, fieldName string) any {
	switch v := value.(type) {
	case map[string]any:
		masked := make(map[string]any, len(v))
		for key, item := range v {
			masked[key] = maskSensitiveValue(item, key)
		}
		return masked
	case []any:
		masked := make([]any, 0, len(v))
		for _, item := range v {
			masked = append(masked, maskSensitiveValue(item, fieldName))
		}
		return masked
	case string:
		return maskStringValue(v, fieldName)
	default:
		if isSensitiveFieldName(fieldName) {
			return "[REDACTED]"
		}
		return value
	}
}

func maskStringValue(value, fieldName string) string {
	if isSensitiveFieldName(fieldName) {
		return "[REDACTED]"
	}
	if looksLikeEmail(value) {
		return maskEmail(value)
	}
	if looksLikePhone(value) {
		return maskPhone(value)
	}
	return value
}

func isSensitiveFieldName(fieldName string) bool {
	normalized := strings.ToLower(strings.TrimSpace(fieldName))
	if normalized == "" {
		return false
	}
	if strings.Contains(normalized, "token") ||
		strings.Contains(normalized, "secret") ||
		strings.Contains(normalized, "password") ||
		strings.Contains(normalized, "credential") ||
		strings.Contains(normalized, "email") ||
		strings.Contains(normalized, "phone") {
		return true
	}
	return normalized == "key" ||
		strings.Contains(normalized, "api_key") ||
		strings.Contains(normalized, "apikey") ||
		strings.HasSuffix(normalized, "key") ||
		strings.HasSuffix(normalized, "_key")
}

func looksLikeEmail(value string) bool {
	return emailPattern.MatchString(strings.TrimSpace(value))
}

func maskEmail(value string) string {
	parts := strings.Split(strings.TrimSpace(value), "@")
	if len(parts) != 2 {
		return "[REDACTED]"
	}
	local := strings.TrimSpace(parts[0])
	domain := strings.TrimSpace(parts[1])
	if local == "" || domain == "" {
		return "[REDACTED]"
	}
	return string(local[0]) + "***@" + domain
}

func looksLikePhone(value string) bool {
	digits := onlyDigits(value)
	return len(digits) >= 7 && len(digits) <= 15
}

func maskPhone(value string) string {
	digits := onlyDigits(value)
	if len(digits) < 4 {
		return "[REDACTED]"
	}
	last4 := digits[len(digits)-4:]
	return "***-***-" + last4
}

func onlyDigits(value string) string {
	var b strings.Builder
	for _, r := range value {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
