package integration_test

import (
	"context"
	"testing"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

func TestComplianceRecheckLifecycleIntegration_CreateAndComplete(t *testing.T) {
	t.Parallel()
	svc := complianceSvc.NewRecheckService()
	item, err := svc.CreateTask(context.Background(), 1, "finding-1", complianceSvc.CreateRecheckInput{TriggerSource: "manual"})
	if err != nil {
		t.Fatalf("create recheck failed: %v", err)
	}
	completed, err := svc.CompleteTask(context.Background(), 1, item.ID, complianceSvc.CompleteRecheckInput{Passed: true, Summary: "done"})
	if err != nil || completed.Status != "passed" {
		t.Fatalf("complete recheck failed status=%v err=%v", completed, err)
	}
}
