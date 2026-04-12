package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
	authSvc "kbmanage/backend/internal/service/auth"
)

const (
	OperationAuditActionSubmit  = "operation.submit"
	OperationAuditActionStart   = "operation.start"
	OperationAuditActionSuccess = "operation.success"
	OperationAuditActionFailure = "operation.failure"

	OperationAuditResourceType = "operation"

	ObservabilityAuditResourceType = "observability"

	ObservabilityAuditActionAlertSync                 = "observability.alert.sync"
	ObservabilityAuditActionAlertAcknowledge          = "observability.alert.acknowledge"
	ObservabilityAuditActionAlertHandlingRecordCreate = "observability.alert.handling_record.create"
	ObservabilityAuditActionAlertRuleCreate           = "observability.alert_rule.create"
	ObservabilityAuditActionAlertRuleUpdate           = "observability.alert_rule.update"
	ObservabilityAuditActionAlertRuleDelete           = "observability.alert_rule.delete"
	ObservabilityAuditActionNotificationTargetCreate  = "observability.notification_target.create"
	ObservabilityAuditActionNotificationTargetUpdate  = "observability.notification_target.update"
	ObservabilityAuditActionNotificationTargetDelete  = "observability.notification_target.delete"
	ObservabilityAuditActionSilenceCreate             = "observability.silence.create"
	ObservabilityAuditActionSilenceCancel             = "observability.silence.cancel"
	ObservabilityAuditActionAccessRead                = "observability.access.read"
)

type EventWriter struct {
	repo *repository.AuditRepository
}

func NewEventWriter(repo *repository.AuditRepository) *EventWriter {
	return &EventWriter{repo: repo}
}

func (w *EventWriter) Write(
	ctx context.Context,
	requestID string,
	actorID *uint64,
	action string,
	resourceType string,
	resourceID string,
	outcome domain.AuditOutcome,
	details map[string]any,
) error {
	if w == nil || w.repo == nil {
		return nil
	}
	if details == nil {
		details = map[string]any{}
	}

	payload, err := json.Marshal(details)
	if err != nil {
		return err
	}

	scopeSnapshot, _ := buildScopeSnapshotJSON(resourceID, details)
	category, actionScope, tags := classifyAuditMetadata(action, details)

	event := &domain.AuditEvent{
		RequestID:     requestID,
		ActorID:       actorID,
		ClusterID:     resolveClusterID(resourceType, resourceID, details),
		WorkspaceID:   resolveWorkspaceID(resourceID, details),
		ProjectID:     resolveProjectID(resourceID, details),
		AuditCategory: category,
		ActionScope:   actionScope,
		ScopeSnapshot: scopeSnapshot,
		SearchTags:    strings.Join(tags, ","),
		Action:        action,
		ResourceType:  resourceType,
		ResourceID:    resourceID,
		Outcome:       outcome,
		Details:       payload,
	}
	return w.repo.Create(ctx, event)
}

func (w *EventWriter) WriteOperationEvent(
	ctx context.Context,
	requestID string,
	actorID *uint64,
	operationID uint64,
	action string,
	outcome domain.AuditOutcome,
	details map[string]any,
) error {
	if details == nil {
		details = map[string]any{}
	}
	details["operationId"] = operationID
	return w.Write(
		ctx,
		requestID,
		actorID,
		action,
		OperationAuditResourceType,
		strconv.FormatUint(operationID, 10),
		outcome,
		details,
	)
}

func (w *EventWriter) WriteObservabilityEvent(
	ctx context.Context,
	requestID string,
	actorID *uint64,
	action string,
	resourceID string,
	outcome domain.AuditOutcome,
	details map[string]any,
) error {
	if details == nil {
		details = map[string]any{}
	}
	return w.Write(
		ctx,
		requestID,
		actorID,
		action,
		ObservabilityAuditResourceType,
		strings.TrimSpace(resourceID),
		outcome,
		details,
	)
}

func resolveClusterID(resourceType, resourceID string, details map[string]any) *uint64 {
	if strings.EqualFold(strings.TrimSpace(resourceType), "cluster") {
		if id, err := strconv.ParseUint(strings.TrimSpace(resourceID), 10, 64); err == nil && id != 0 {
			return &id
		}
	}

	if id, ok := authSvc.ParseClusterIDFromReference(resourceID); ok && id != 0 {
		return &id
	}

	if details == nil {
		return nil
	}
	if id := resolveOptionalUint64(details, "clusterId", "clusterID"); id != nil {
		return id
	}
	targetRef := strings.TrimSpace(resolveOptionalString(details, "targetRef", "target_ref"))
	if targetRef == "" {
		return nil
	}
	if id, ok := authSvc.ParseClusterIDFromReference(targetRef); ok && id != 0 {
		return &id
	}
	return nil
}

func resolveWorkspaceID(resourceID string, details map[string]any) *uint64 {
	if id := resolveOptionalUint64(details, "workspaceId", "workspaceID"); id != nil {
		return id
	}
	if id := parseScopeIDFromReference(resourceID, "workspace"); id != nil {
		return id
	}
	targetRef := strings.TrimSpace(resolveOptionalString(details, "targetRef", "target_ref"))
	if targetRef == "" {
		return nil
	}
	return parseScopeIDFromReference(targetRef, "workspace")
}

func resolveProjectID(resourceID string, details map[string]any) *uint64 {
	if id := resolveOptionalUint64(details, "projectId", "projectID"); id != nil {
		return id
	}
	if id := parseScopeIDFromReference(resourceID, "project"); id != nil {
		return id
	}
	targetRef := strings.TrimSpace(resolveOptionalString(details, "targetRef", "target_ref"))
	if targetRef == "" {
		return nil
	}
	return parseScopeIDFromReference(targetRef, "project")
}

func parseScopeIDFromReference(rawRef, scope string) *uint64 {
	prefix := strings.ToLower(strings.TrimSpace(scope)) + ":"
	if prefix == ":" {
		return nil
	}
	for _, segment := range strings.Split(strings.TrimSpace(rawRef), "/") {
		part := strings.TrimSpace(segment)
		if !strings.HasPrefix(strings.ToLower(part), prefix) {
			continue
		}
		idText := strings.TrimSpace(part[len(prefix):])
		id, err := strconv.ParseUint(idText, 10, 64)
		if err != nil || id == 0 {
			return nil
		}
		return &id
	}
	return nil
}

func resolveOptionalUint64(data map[string]any, keys ...string) *uint64 {
	if data == nil {
		return nil
	}
	for _, key := range keys {
		v, ok := findValueByKey(data, key)
		if !ok {
			continue
		}
		if id, ok := convertToUint64(v); ok && id != 0 {
			return &id
		}
	}
	return nil
}

func resolveOptionalString(data map[string]any, keys ...string) string {
	if data == nil {
		return ""
	}
	for _, key := range keys {
		v, ok := findValueByKey(data, key)
		if !ok {
			continue
		}
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func findValueByKey(data map[string]any, key string) (any, bool) {
	for existing, value := range data {
		if strings.EqualFold(strings.TrimSpace(existing), strings.TrimSpace(key)) {
			return value, true
		}
	}
	return nil, false
}

func convertToUint64(value any) (uint64, bool) {
	switch v := value.(type) {
	case uint64:
		return v, true
	case uint:
		return uint64(v), true
	case int:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case int64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case float64:
		if v < 0 {
			return 0, false
		}
		return uint64(v), true
	case string:
		text := strings.TrimSpace(v)
		if text == "" {
			return 0, false
		}
		id, err := strconv.ParseUint(text, 10, 64)
		if err != nil {
			return 0, false
		}
		return id, true
	default:
		return 0, false
	}
}

func classifyAuditMetadata(action string, details map[string]any) (string, string, []string) {
	action = strings.TrimSpace(strings.ToLower(action))
	category := "general"
	actionScope := "unknown"
	if strings.HasPrefix(action, "observability.") {
		category = "observability"
		actionScope = "access"
		switch {
		case strings.Contains(action, "alert_rule."):
			actionScope = "alert_rule"
		case strings.Contains(action, "silence."):
			actionScope = "silence"
		case strings.Contains(action, "acknowledge"):
			actionScope = "acknowledge"
		case strings.Contains(action, "handling_record"):
			actionScope = "handling_record"
		case strings.Contains(action, "access.read"):
			actionScope = "access"
		}
	}

	tags := map[string]struct{}{
		fmt.Sprintf("category:%s", category): {},
		fmt.Sprintf("scope:%s", actionScope): {},
	}
	addTagFromDetails(tags, details, "operation", "operation")
	addTagFromDetails(tags, details, "resourceKind", "resourceKind")
	addTagFromDetails(tags, details, "resourceName", "resourceName")
	addTagFromDetails(tags, details, "subjectType", "subjectType")

	out := make([]string, 0, len(tags))
	for tag := range tags {
		out = append(out, tag)
	}
	sort.Strings(out)
	return category, actionScope, out
}

func addTagFromDetails(tags map[string]struct{}, details map[string]any, key, prefix string) {
	if details == nil {
		return
	}
	v, ok := findValueByKey(details, key)
	if !ok {
		return
	}
	text, ok := v.(string)
	if !ok {
		return
	}
	text = strings.TrimSpace(text)
	if text == "" {
		return
	}
	tags[fmt.Sprintf("%s:%s", prefix, text)] = struct{}{}
}

func buildScopeSnapshotJSON(resourceID string, details map[string]any) (json.RawMessage, error) {
	clusterIDs := make([]uint64, 0, 1)
	workspaceIDs := make([]uint64, 0, 1)
	projectIDs := make([]uint64, 0, 1)

	if id, ok := authSvc.ParseClusterIDFromReference(resourceID); ok && id != 0 {
		clusterIDs = append(clusterIDs, id)
	}
	if id := parseScopeIDFromReference(resourceID, "workspace"); id != nil {
		workspaceIDs = append(workspaceIDs, *id)
	}
	if id := parseScopeIDFromReference(resourceID, "project"); id != nil {
		projectIDs = append(projectIDs, *id)
	}

	if id := resolveOptionalUint64(details, "clusterId", "clusterID"); id != nil {
		clusterIDs = append(clusterIDs, *id)
	}
	if id := resolveOptionalUint64(details, "workspaceId", "workspaceID"); id != nil {
		workspaceIDs = append(workspaceIDs, *id)
	}
	if id := resolveOptionalUint64(details, "projectId", "projectID"); id != nil {
		projectIDs = append(projectIDs, *id)
	}

	payload := map[string]any{
		"clusterIds":   uniqueIDs(clusterIDs),
		"workspaceIds": uniqueIDs(workspaceIDs),
		"projectIds":   uniqueIDs(projectIDs),
	}
	encoded, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	return json.RawMessage(encoded), nil
}

func uniqueIDs(values []uint64) []uint64 {
	set := make(map[uint64]struct{}, len(values))
	for _, value := range values {
		if value == 0 {
			continue
		}
		set[value] = struct{}{}
	}
	out := make([]uint64, 0, len(set))
	for value := range set {
		out = append(out, value)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
