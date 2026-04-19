package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_BackupLifecycle(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")

	createResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies", `{
		"name":"audit-records",
		"scopeType":"audit-records",
		"scopeRef":"platform",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"executionMode":"manual",
		"retentionRule":"30d",
		"consistencyLevel":"best-effort",
		"status":"active"
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create policy failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	runResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies/1/run", "")
	if runResp.Code != http.StatusAccepted {
		t.Fatalf("run backup failed status=%d body=%s", runResp.Code, strings.TrimSpace(runResp.Body.String()))
	}

	listResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/backup-restore/restore-points?policyId=1", "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), "\"policyId\":1") {
		t.Fatalf("list restore points failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
