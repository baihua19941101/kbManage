package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreContract_RestorePrecheck(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")
	policyID := createBackupRestoreContractPolicy(t, ctx, "precheck-policy")
	runBackupRestoreContractPolicy(t, ctx, policyID)
	performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs", `{
		"restorePointId":1,
		"jobType":"cross-cluster-restore",
		"sourceEnvironment":"prod",
		"targetEnvironment":"dr-site",
		"scopeSelection":{"namespaces":["orders-prod"]}
	}`)

	validateResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs/1/validate", "")
	if validateResp.Code != http.StatusOK || !strings.Contains(validateResp.Body.String(), "blockers") {
		t.Fatalf("restore precheck failed status=%d body=%s", validateResp.Code, strings.TrimSpace(validateResp.Body.String()))
	}
}
