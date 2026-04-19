package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreContract_RestoreJobs(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")
	policyID := createBackupRestoreContractPolicy(t, ctx, "restore-job-policy")
	runBackupRestoreContractPolicy(t, ctx, policyID)

	createResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs", `{
		"restorePointId":1,
		"jobType":"in-place-restore",
		"sourceEnvironment":"prod",
		"targetEnvironment":"prod",
		"scopeSelection":{"namespaces":["orders-prod"]}
	}`)
	if createResp.Code != http.StatusAccepted {
		t.Fatalf("create restore job failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	listResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/backup-restore/restore-jobs", "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), "\"items\"") {
		t.Fatalf("list restore jobs failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
