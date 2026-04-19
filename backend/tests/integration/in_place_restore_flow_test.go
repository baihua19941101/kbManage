package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_InPlaceRestoreFlow(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	policyID := createBackupRestoreIntegrationPolicy(t, ctx, "in-place-policy")
	runBackupRestoreIntegrationPolicy(t, ctx, policyID)

	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs", `{
		"restorePointId":1,
		"jobType":"in-place-restore",
		"sourceEnvironment":"prod",
		"targetEnvironment":"prod",
		"scopeSelection":{"namespaces":["orders-prod"]}
	}`)
	if resp.Code != http.StatusAccepted || !strings.Contains(resp.Body.String(), "\"status\":\"succeeded\"") {
		t.Fatalf("in-place restore failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
