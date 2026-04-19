package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_SelectiveRestoreScope(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	policyID := createBackupRestoreIntegrationPolicy(t, ctx, "selective-policy")
	runBackupRestoreIntegrationPolicy(t, ctx, policyID)

	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs", `{
		"restorePointId":1,
		"jobType":"selective-restore",
		"sourceEnvironment":"prod",
		"targetEnvironment":"staging",
		"scopeSelection":{"configMaps":["orders-config"]}
	}`)
	if resp.Code != http.StatusAccepted || !strings.Contains(resp.Body.String(), "configMaps") {
		t.Fatalf("selective restore failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
