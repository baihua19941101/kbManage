package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_BackupAuditQuery(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	policyID := createBackupRestoreIntegrationPolicy(t, ctx, "audit-query-policy")
	runBackupRestoreIntegrationPolicy(t, ctx, policyID)

	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/audit/backup-restore/events?action=backuprestore.policy.run", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "backuprestore.policy.run") {
		t.Fatalf("query backup audit failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
