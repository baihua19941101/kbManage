package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_CrossClusterRestoreFlow(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	policyID := createBackupRestoreIntegrationPolicy(t, ctx, "cross-cluster-policy")
	runBackupRestoreIntegrationPolicy(t, ctx, policyID)

	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs", `{
		"restorePointId":1,
		"jobType":"cross-cluster-restore",
		"sourceEnvironment":"prod",
		"targetEnvironment":"dr-site",
		"scopeSelection":{"namespaces":["orders-prod"]}
	}`)
	if resp.Code != http.StatusAccepted || !strings.Contains(resp.Body.String(), "\"targetEnvironment\":\"dr-site\"") {
		t.Fatalf("cross-cluster restore failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
