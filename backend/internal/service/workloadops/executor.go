package workloadops

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"kbmanage/backend/internal/domain"
)

type ActionExecuteResult struct {
	ProgressMessage string
	ResultMessage   string
	FailureReason   string
}

type ActionExecutor interface {
	Execute(ctx context.Context, item *domain.WorkloadActionRequest) (ActionExecuteResult, error)
}

type simulatedActionExecutor struct{}

func NewActionExecutor() ActionExecutor {
	return simulatedActionExecutor{}
}

func (simulatedActionExecutor) Execute(_ context.Context, item *domain.WorkloadActionRequest) (ActionExecuteResult, error) {
	if item == nil {
		return ActionExecuteResult{}, fmt.Errorf("action request is required")
	}
	target := fmt.Sprintf("%s/%s", strings.TrimSpace(item.Namespace), strings.TrimSpace(item.ResourceName))
	switch item.ActionType {
	case domain.WorkloadActionTypeScale:
		return ActionExecuteResult{
			ProgressMessage: "scale action executed",
			ResultMessage:   fmt.Sprintf("scaled workload %s", target),
		}, nil
	case domain.WorkloadActionTypeRestart:
		return ActionExecuteResult{
			ProgressMessage: "restart action executed",
			ResultMessage:   fmt.Sprintf("restarted workload %s", target),
		}, nil
	case domain.WorkloadActionTypeRedeploy:
		return ActionExecuteResult{
			ProgressMessage: "redeploy action executed",
			ResultMessage:   fmt.Sprintf("redeployed workload %s", target),
		}, nil
	case domain.WorkloadActionTypeReplaceInstance:
		return ActionExecuteResult{
			ProgressMessage: "instance replacement executed",
			ResultMessage:   fmt.Sprintf("replaced instance for workload %s", target),
		}, nil
	case domain.WorkloadActionTypeRollback:
		revision := parseRollbackRevision(item.PayloadJSON)
		return ActionExecuteResult{
			ProgressMessage: "rollback action executed",
			ResultMessage:   fmt.Sprintf("rolled back workload %s to revision %s", target, revision),
		}, nil
	default:
		return ActionExecuteResult{}, fmt.Errorf("unsupported action type: %s", item.ActionType)
	}
}

func parseRollbackRevision(payloadJSON string) string {
	if strings.TrimSpace(payloadJSON) == "" {
		return "unknown"
	}
	var payload map[string]any
	if err := json.Unmarshal([]byte(payloadJSON), &payload); err != nil {
		return "unknown"
	}
	if revision, ok := payload["revision"]; ok {
		return strings.TrimSpace(fmt.Sprint(revision))
	}
	return "unknown"
}
