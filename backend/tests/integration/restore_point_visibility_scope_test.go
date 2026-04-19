package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_RestorePointVisibilityScope(t *testing.T) {
	ownerCtx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	policyID := createBackupRestoreIntegrationPolicy(t, ownerCtx, "scope-policy")
	runBackupRestoreIntegrationPolicy(t, ownerCtx, policyID)

	viewerCtx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	resp := performBackupRestoreIntegrationRequest(t, viewerCtx.Router, viewerCtx.Token, http.MethodGet, "/api/v1/backup-restore/restore-points", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list restore points failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if strings.Contains(resp.Body.String(), "scope-policy") {
		t.Fatalf("expected scoped viewer not to see restore point body=%s", strings.TrimSpace(resp.Body.String()))
	}
}
