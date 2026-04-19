package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreContract_BackupPolicyAndRestorePoint(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")

	createResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies", `{
		"name":"platform-meta",
		"description":"platform metadata backup",
		"scopeType":"platform-metadata",
		"scopeRef":"global",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"executionMode":"manual",
		"retentionRule":"7d",
		"consistencyLevel":"platform-consistent",
		"status":"active"
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create policy failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	runResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies/1/run", "")
	if runResp.Code != http.StatusAccepted {
		t.Fatalf("run policy failed status=%d body=%s", runResp.Code, strings.TrimSpace(runResp.Body.String()))
	}

	listResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/backup-restore/restore-points", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list restore points failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if !strings.Contains(listResp.Body.String(), "platform-metadata") {
		t.Fatalf("expected restore points to include consistency summary body=%s", strings.TrimSpace(listResp.Body.String()))
	}

	detailResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/backup-restore/restore-points/1", "")
	if detailResp.Code != http.StatusOK {
		t.Fatalf("get restore point failed status=%d body=%s", detailResp.Code, strings.TrimSpace(detailResp.Body.String()))
	}
}
