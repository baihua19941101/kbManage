package integration_test

import (
	"context"
	"testing"

	complianceSvc "kbmanage/backend/internal/service/compliance"
)

func TestComplianceArchiveExportIntegration_CreateAndProcess(t *testing.T) {
	t.Parallel()
	svc := complianceSvc.NewArchiveExportService()
	item, err := svc.CreateExport(context.Background(), 1, complianceSvc.CreateArchiveExportInput{ExportScope: "bundle"})
	if err != nil {
		t.Fatalf("create export failed: %v", err)
	}
	result, err := svc.ProcessExport(context.Background(), item.ID)
	if err != nil || result.Status != "succeeded" {
		t.Fatalf("process export failed status=%v err=%v", result, err)
	}
}
