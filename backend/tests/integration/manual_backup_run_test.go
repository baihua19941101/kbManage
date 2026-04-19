package integration_test

import "testing"

func TestBackupRestoreIntegration_ManualBackupRun(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	policyID := createBackupRestoreIntegrationPolicy(t, ctx, "manual-run-int")
	restorePointID := runBackupRestoreIntegrationPolicy(t, ctx, policyID)
	if restorePointID == 0 {
		t.Fatal("expected restore point id after manual backup run")
	}
}
