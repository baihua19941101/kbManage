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

	WorkloadOpsAuditResourceType   = "workloadops"
	WorkloadOpsAuditActionSubmit   = "workloadops.action.submit"
	WorkloadOpsAuditActionStart    = "workloadops.action.start"
	WorkloadOpsAuditActionSuccess  = "workloadops.action.success"
	WorkloadOpsAuditActionFailure  = "workloadops.action.failure"
	WorkloadOpsAuditBatchSubmit    = "workloadops.batch.submit"
	WorkloadOpsAuditBatchFinish    = "workloadops.batch.finish"
	WorkloadOpsAuditRollbackSubmit = "workloadops.rollback.submit"
	WorkloadOpsAuditRollbackFinish = "workloadops.rollback.finish"
	WorkloadOpsAuditTerminalOpen   = "workloadops.terminal.open"
	WorkloadOpsAuditTerminalClose  = "workloadops.terminal.close"

	GitOpsAuditResourceType          = "gitops"
	GitOpsAuditActionSourceVerify    = "gitops.source.verify"
	GitOpsAuditActionSyncSubmit      = "gitops.sync.submit"
	GitOpsAuditActionResyncSubmit    = "gitops.resync.submit"
	GitOpsAuditActionInstallSubmit   = "gitops.install.submit"
	GitOpsAuditActionUpgradeSubmit   = "gitops.upgrade.submit"
	GitOpsAuditActionPromoteSubmit   = "gitops.promote.submit"
	GitOpsAuditActionRollbackSubmit  = "gitops.rollback.submit"
	GitOpsAuditActionPauseSubmit     = "gitops.pause.submit"
	GitOpsAuditActionResumeSubmit    = "gitops.resume.submit"
	GitOpsAuditActionUninstallSubmit = "gitops.uninstall.submit"

	SecurityPolicyAuditResourceType = "securitypolicy"

	SecurityPolicyAuditActionPolicyCreate         = "securitypolicy.policy.create"
	SecurityPolicyAuditActionPolicyUpdate         = "securitypolicy.policy.update"
	SecurityPolicyAuditActionAssignmentCreate     = "securitypolicy.assignment.create"
	SecurityPolicyAuditActionAssignmentUpdate     = "securitypolicy.assignment.update"
	SecurityPolicyAuditActionModeSwitch           = "securitypolicy.mode_switch.submit"
	SecurityPolicyAuditActionExceptionCreate      = "securitypolicy.exception.create"
	SecurityPolicyAuditActionExceptionReview      = "securitypolicy.exception.review"
	SecurityPolicyAuditActionHitQuery             = "securitypolicy.hit.query"
	SecurityPolicyAuditActionHitRemediationUpdate = "securitypolicy.hit.remediation.update"

	ComplianceAuditResourceType            = "compliance"
	ComplianceAuditActionBaselineCreate    = "compliance.baseline.create"
	ComplianceAuditActionBaselineUpdate    = "compliance.baseline.update"
	ComplianceAuditActionScanExecute       = "compliance.scan.execute"
	ComplianceAuditActionRemediationCreate = "compliance.remediation.create"
	ComplianceAuditActionRemediationUpdate = "compliance.remediation.update"
	ComplianceAuditActionExceptionRequest  = "compliance.exception.request"
	ComplianceAuditActionExceptionReview   = "compliance.exception.review"
	ComplianceAuditActionRecheckCreate     = "compliance.recheck.create"
	ComplianceAuditActionRecheckComplete   = "compliance.recheck.complete"
	ComplianceAuditActionArchiveExport     = "compliance.archive_export.create"

	ClusterLifecycleAuditResourceType         = "clusterlifecycle"
	ClusterLifecycleAuditActionImportSubmit   = "clusterlifecycle.import.submit"
	ClusterLifecycleAuditActionRegisterSubmit = "clusterlifecycle.register.submit"
	ClusterLifecycleAuditActionCreateSubmit   = "clusterlifecycle.create.submit"
	ClusterLifecycleAuditActionValidateSubmit = "clusterlifecycle.validate.submit"
	ClusterLifecycleAuditActionUpgradeSubmit  = "clusterlifecycle.upgrade.submit"
	ClusterLifecycleAuditActionNodePoolScale  = "clusterlifecycle.nodepool.scale.submit"
	ClusterLifecycleAuditActionDisableSubmit  = "clusterlifecycle.disable.submit"
	ClusterLifecycleAuditActionRetireSubmit   = "clusterlifecycle.retire.submit"
	ClusterLifecycleAuditActionDriverUpsert   = "clusterlifecycle.driver.upsert"
	ClusterLifecycleAuditActionTemplateUpsert = "clusterlifecycle.template.upsert"
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

func (w *EventWriter) WriteWorkloadOpsEvent(
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
		WorkloadOpsAuditResourceType,
		strings.TrimSpace(resourceID),
		outcome,
		details,
	)
}

func (w *EventWriter) WriteGitOpsEvent(
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
		GitOpsAuditResourceType,
		strings.TrimSpace(resourceID),
		outcome,
		details,
	)
}

func (w *EventWriter) WriteComplianceEvent(
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
	return w.Write(ctx, requestID, actorID, action, ComplianceAuditResourceType, strings.TrimSpace(resourceID), outcome, details)
}

func (w *EventWriter) WriteClusterLifecycleEvent(
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
	return w.Write(ctx, requestID, actorID, action, ClusterLifecycleAuditResourceType, strings.TrimSpace(resourceID), outcome, details)
}

func (w *EventWriter) WriteSecurityPolicyEvent(
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
		SecurityPolicyAuditResourceType,
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
	case *uint64:
		if v == nil {
			return 0, false
		}
		return *v, true
	case *uint:
		if v == nil {
			return 0, false
		}
		return uint64(*v), true
	case *int:
		if v == nil || *v < 0 {
			return 0, false
		}
		return uint64(*v), true
	case *int64:
		if v == nil || *v < 0 {
			return 0, false
		}
		return uint64(*v), true
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
	if strings.HasPrefix(action, "workloadops.") {
		category = "workloadops"
		actionScope = "action"
		switch {
		case strings.Contains(action, ".terminal."):
			actionScope = "terminal"
		case strings.Contains(action, ".batch."):
			actionScope = "batch"
		case strings.Contains(action, ".rollback."):
			actionScope = "rollback"
		}
	}
	if strings.HasPrefix(action, "gitops.") {
		category = "gitops"
		actionScope = "action"
		switch {
		case strings.Contains(action, ".source."):
			actionScope = "source"
		case strings.Contains(action, ".promote."):
			actionScope = "promote"
		case strings.Contains(action, ".rollback."):
			actionScope = "rollback"
		case strings.Contains(action, ".pause.") || strings.Contains(action, ".resume."):
			actionScope = "sync_control"
		case strings.Contains(action, ".sync.") || strings.Contains(action, ".resync."):
			actionScope = "sync"
		}
	}
	if strings.HasPrefix(action, "securitypolicy.") {
		category = "securitypolicy"
		actionScope = "policy"
		switch {
		case strings.Contains(action, ".assignment."):
			actionScope = "assignment"
		case strings.Contains(action, ".mode_switch."):
			actionScope = "mode_switch"
		case strings.Contains(action, ".exception."):
			actionScope = "exception"
		case strings.Contains(action, ".hit.") || strings.Contains(action, ".remediation."):
			actionScope = "hit"
		}
	}
	if strings.HasPrefix(action, "clusterlifecycle.") {
		category = "clusterlifecycle"
		actionScope = "cluster"
		switch {
		case strings.Contains(action, ".driver."):
			actionScope = "driver"
		case strings.Contains(action, ".template."):
			actionScope = "template"
		case strings.Contains(action, ".nodepool."):
			actionScope = "nodepool"
		case strings.Contains(action, ".upgrade."):
			actionScope = "upgrade"
		case strings.Contains(action, ".retire.") || strings.Contains(action, ".disable."):
			actionScope = "retirement"
		case strings.Contains(action, ".validate."):
			actionScope = "validation"
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
