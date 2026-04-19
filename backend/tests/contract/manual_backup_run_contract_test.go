package contract_test

import "testing"

func TestBackupRestoreContract_ManualBackupRun(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")
	policyID := createBackupRestoreContractPolicy(t, ctx, "manual-backup-policy")
	restorePointID := runBackupRestoreContractPolicy(t, ctx, policyID)
	if restorePointID == 0 {
		t.Fatal("expected restore point id after manual backup run")
	}
}
