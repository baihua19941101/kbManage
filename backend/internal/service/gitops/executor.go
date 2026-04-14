package gitops

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"kbmanage/backend/internal/domain"
	"kbmanage/backend/internal/repository"
)

type ExecuteInput struct {
	Operation *domain.DeliveryOperation
	Unit      *domain.ApplicationDeliveryUnit
	Detail    *repository.DeliveryUnitDetail
	Payload   map[string]any
}

type ExecuteResult struct {
	Status         domain.DeliveryOperationStatus
	Progress       int
	ResultSummary  string
	FailureReason  string
	UnitUpdates    map[string]any
	StageResults   []StageExecutionResult
	OperationPatch map[string]any
}

type Executor struct {
	revisions  *RevisionService
	promotions *PromotionService
}

func NewExecutor(revisions *RevisionService, promotions *PromotionService) *Executor {
	return &Executor{
		revisions:  revisions,
		promotions: promotions,
	}
}

func (e *Executor) Execute(ctx context.Context, in ExecuteInput) (ExecuteResult, error) {
	if in.Operation == nil || in.Unit == nil {
		return ExecuteResult{}, errors.New("operation and unit are required")
	}
	action := normalizeDeliveryActionType(in.Operation.ActionType)
	if action == "" {
		action = domain.DeliveryActionTypeSync
	}
	payload := clonePayload(in.Payload)
	if payload == nil {
		payload = map[string]any{}
	}

	switch action {
	case domain.DeliveryActionTypeInstall:
		return e.executeInstall(ctx, in, action, payload)
	case domain.DeliveryActionTypeSync, domain.DeliveryActionTypeResync:
		return e.executeSync(in, action, payload), nil
	case domain.DeliveryActionTypeUpgrade:
		return e.executeUpgrade(ctx, in, action, payload)
	case domain.DeliveryActionTypePause:
		return e.executePause(in, payload), nil
	case domain.DeliveryActionTypeResume:
		return e.executeResume(in, payload), nil
	case domain.DeliveryActionTypeUninstall:
		return e.executeUninstall(in, payload), nil
	case domain.DeliveryActionTypePromote:
		return e.executePromote(ctx, in, payload)
	case domain.DeliveryActionTypeRollback:
		return e.executeRollback(ctx, in, payload)
	default:
		return ExecuteResult{}, fmt.Errorf("unsupported action type: %s", in.Operation.ActionType)
	}
}

func (e *Executor) executeInstall(
	ctx context.Context,
	in ExecuteInput,
	action domain.DeliveryActionType,
	payload map[string]any,
) (ExecuteResult, error) {
	if e == nil || e.revisions == nil {
		return ExecuteResult{}, ErrGitOpsNotConfigured
	}
	revision, err := e.revisions.RecordRelease(ctx, ReleaseRecordInput{
		DeliveryUnitID: in.Unit.ID,
		OperatorID:     in.Operation.OperatorID,
		SourceRevision: firstNonEmpty(readString(payload, "targetRevision"), in.Unit.DesiredRevision),
		AppVersion:     firstNonEmpty(readString(payload, "targetAppVersion"), in.Unit.DesiredAppVersion),
		ConfigVersion:  firstNonEmpty(readString(payload, "targetConfigVersion"), in.Unit.DesiredConfigVersion),
		Environment:    readString(payload, "environment"),
		Notes:          firstNonEmpty(readString(payload, "reason"), string(action)+" executed"),
	})
	if err != nil {
		return ExecuteResult{}, err
	}
	payload["revisionId"] = revision.ID
	return ExecuteResult{
		Status:        domain.DeliveryOperationStatusSucceeded,
		Progress:      100,
		ResultSummary: "install completed",
		UnitUpdates: map[string]any{
			"desired_revision":       revision.SourceRevision,
			"desired_app_version":    revision.AppVersion,
			"desired_config_version": revision.ConfigVersion,
			"paused":                false,
			"delivery_status":       domain.DeliveryUnitStatusReady,
			"last_synced_at":        time.Now(),
			"last_release_id":       revision.ID,
		},
		OperationPatch: payload,
	}, nil
}

func (e *Executor) executeSync(in ExecuteInput, action domain.DeliveryActionType, payload map[string]any) ExecuteResult {
	payload["syncAction"] = string(action)
	return ExecuteResult{
		Status:        domain.DeliveryOperationStatusSucceeded,
		Progress:      100,
		ResultSummary: fmt.Sprintf("%s completed", action),
		UnitUpdates: map[string]any{
			"delivery_status": domain.DeliveryUnitStatusReady,
			"last_synced_at":  time.Now(),
		},
		OperationPatch: payload,
	}
}

func (e *Executor) executeUpgrade(
	ctx context.Context,
	in ExecuteInput,
	action domain.DeliveryActionType,
	payload map[string]any,
) (ExecuteResult, error) {
	if e == nil || e.revisions == nil {
		return ExecuteResult{}, ErrGitOpsNotConfigured
	}
	revision, err := e.revisions.RecordRelease(ctx, ReleaseRecordInput{
		DeliveryUnitID: in.Unit.ID,
		OperatorID:     in.Operation.OperatorID,
		SourceRevision: firstNonEmpty(readString(payload, "targetRevision"), in.Unit.DesiredRevision),
		AppVersion:     firstNonEmpty(readString(payload, "targetAppVersion"), in.Unit.DesiredAppVersion),
		ConfigVersion:  firstNonEmpty(readString(payload, "targetConfigVersion"), in.Unit.DesiredConfigVersion),
		Environment:    readString(payload, "environment"),
		Notes:          firstNonEmpty(readString(payload, "reason"), string(action)+" executed"),
	})
	if err != nil {
		return ExecuteResult{}, err
	}
	payload["revisionId"] = revision.ID
	return ExecuteResult{
		Status:        domain.DeliveryOperationStatusSucceeded,
		Progress:      100,
		ResultSummary: "upgrade completed",
		UnitUpdates: map[string]any{
			"desired_revision":       revision.SourceRevision,
			"desired_app_version":    revision.AppVersion,
			"desired_config_version": revision.ConfigVersion,
			"delivery_status":       domain.DeliveryUnitStatusReady,
			"last_synced_at":        time.Now(),
			"last_release_id":       revision.ID,
		},
		OperationPatch: payload,
	}, nil
}

func (e *Executor) executePause(in ExecuteInput, payload map[string]any) ExecuteResult {
	return ExecuteResult{
		Status:        domain.DeliveryOperationStatusSucceeded,
		Progress:      100,
		ResultSummary: "pause completed",
		UnitUpdates: map[string]any{
			"paused":          true,
			"delivery_status": domain.DeliveryUnitStatusPaused,
		},
		OperationPatch: payload,
	}
}

func (e *Executor) executeResume(in ExecuteInput, payload map[string]any) ExecuteResult {
	return ExecuteResult{
		Status:        domain.DeliveryOperationStatusSucceeded,
		Progress:      100,
		ResultSummary: "resume completed",
		UnitUpdates: map[string]any{
			"paused":          false,
			"delivery_status": domain.DeliveryUnitStatusReady,
		},
		OperationPatch: payload,
	}
}

func (e *Executor) executeUninstall(in ExecuteInput, payload map[string]any) ExecuteResult {
	return ExecuteResult{
		Status:        domain.DeliveryOperationStatusSucceeded,
		Progress:      100,
		ResultSummary: "uninstall completed",
		UnitUpdates: map[string]any{
			"paused":                false,
			"delivery_status":       domain.DeliveryUnitStatusUnknown,
			"last_release_id":       nil,
			"desired_app_version":    "",
			"desired_config_version": "",
		},
		OperationPatch: payload,
	}
}

func (e *Executor) executePromote(ctx context.Context, in ExecuteInput, payload map[string]any) (ExecuteResult, error) {
	if e == nil || e.promotions == nil {
		return ExecuteResult{}, ErrGitOpsNotConfigured
	}
	stages, err := e.promotions.Promote(ctx, in.Unit.ID, readString(payload, "environment"), payload)
	if err != nil {
		return ExecuteResult{}, err
	}
	status, summary, failure := normalizeByStageResults(stages)
	unitStatus := domain.DeliveryUnitStatusReady
	if status == domain.DeliveryOperationStatusPartiallySucceeded {
		unitStatus = domain.DeliveryUnitStatusProgressing
	}
	if status == domain.DeliveryOperationStatusFailed {
		unitStatus = domain.DeliveryUnitStatusDegraded
	}
	payload["stages"] = stages
	return ExecuteResult{
		Status:        status,
		Progress:      100,
		ResultSummary: firstNonEmpty(summary, "promote completed"),
		FailureReason: failure,
		StageResults:  stages,
		UnitUpdates: map[string]any{
			"delivery_status": unitStatus,
			"last_synced_at":  time.Now(),
		},
		OperationPatch: payload,
	}, nil
}

func (e *Executor) executeRollback(ctx context.Context, in ExecuteInput, payload map[string]any) (ExecuteResult, error) {
	if e == nil || e.revisions == nil || e.promotions == nil {
		return ExecuteResult{}, ErrGitOpsNotConfigured
	}
	revision, err := e.revisions.RollbackToRevision(ctx, in.Unit.ID, in.Operation.TargetReleaseID, in.Operation.OperatorID)
	if err != nil {
		return ExecuteResult{}, err
	}
	stages, err := e.promotions.RollbackStages(ctx, in.Unit.ID, readString(payload, "environment"), payload)
	if err != nil {
		return ExecuteResult{}, err
	}
	status, summary, failure := normalizeByStageResults(stages)
	if len(stages) == 0 {
		status = domain.DeliveryOperationStatusSucceeded
		summary = "rollback completed"
	}
	payload["revisionId"] = revision.ID
	payload["stages"] = stages
	return ExecuteResult{
		Status:        status,
		Progress:      100,
		ResultSummary: firstNonEmpty(summary, "rollback completed"),
		FailureReason: failure,
		StageResults:  stages,
		UnitUpdates: map[string]any{
			"delivery_status":       mapRollbackUnitStatus(status),
			"last_synced_at":        time.Now(),
			"last_release_id":       revision.ID,
			"desired_revision":       revision.SourceRevision,
			"desired_app_version":    revision.AppVersion,
			"desired_config_version": revision.ConfigVersion,
		},
		OperationPatch: payload,
	}, nil
}

func mapRollbackUnitStatus(status domain.DeliveryOperationStatus) domain.DeliveryUnitStatus {
	switch status {
	case domain.DeliveryOperationStatusPartiallySucceeded:
		return domain.DeliveryUnitStatusProgressing
	case domain.DeliveryOperationStatusFailed:
		return domain.DeliveryUnitStatusDegraded
	default:
		return domain.DeliveryUnitStatusReady
	}
}

func normalizeByStageResults(stages []StageExecutionResult) (domain.DeliveryOperationStatus, string, string) {
	if len(stages) == 0 {
		return domain.DeliveryOperationStatusSucceeded, "", ""
	}
	total := 0
	succeeded := 0
	failed := 0
	failureReasons := make([]string, 0)
	for i := range stages {
		total += stages[i].TargetCount
		succeeded += stages[i].SucceededCount
		failed += stages[i].FailedCount
		if strings.TrimSpace(stages[i].FailureReason) != "" {
			failureReasons = append(failureReasons, stages[i].FailureReason)
		}
	}
	summary := fmt.Sprintf("targets=%d, succeeded=%d, failed=%d", total, succeeded, failed)
	if failed == 0 {
		return domain.DeliveryOperationStatusSucceeded, summary, ""
	}
	if succeeded == 0 {
		return domain.DeliveryOperationStatusFailed, summary, strings.Join(failureReasons, "; ")
	}
	return domain.DeliveryOperationStatusPartiallySucceeded, summary, strings.Join(failureReasons, "; ")
}

func normalizeDeliveryActionType(action domain.DeliveryActionType) domain.DeliveryActionType {
	trimmed := strings.ToLower(strings.TrimSpace(string(action)))
	switch domain.DeliveryActionType(trimmed) {
	case domain.DeliveryActionTypeInstall,
		domain.DeliveryActionTypeSync,
		domain.DeliveryActionTypeResync,
		domain.DeliveryActionTypeUpgrade,
		domain.DeliveryActionTypePromote,
		domain.DeliveryActionTypeRollback,
		domain.DeliveryActionTypePause,
		domain.DeliveryActionTypeResume,
		domain.DeliveryActionTypeUninstall:
		return domain.DeliveryActionType(trimmed)
	default:
		return ""
	}
}

func readString(payload map[string]any, key string) string {
	if len(payload) == 0 {
		return ""
	}
	return strings.TrimSpace(fmt.Sprint(payload[key]))
}

func clonePayload(payload map[string]any) map[string]any {
	if len(payload) == 0 {
		return map[string]any{}
	}
	out := make(map[string]any, len(payload))
	for k, v := range payload {
		out[k] = v
	}
	return out
}

func firstNonEmpty(values ...string) string {
	for _, value := range values {
		if strings.TrimSpace(value) != "" {
			return strings.TrimSpace(value)
		}
	}
	return ""
}
