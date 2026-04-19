package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreContract_RestorePoints(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")
	policyID := createBackupRestoreContractPolicy(t, ctx, "restore-point-policy")
	restorePointID := runBackupRestoreContractPolicy(t, ctx, policyID)

	listResp := performBackupRestoreContractRequest(
		t,
		ctx.Router,
		ctx.Token,
		http.MethodGet,
		"/api/v1/backup-restore/restore-points?policyId="+strconv.FormatUint(policyID, 10),
		"",
	)
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), "restore-point-policy") && !strings.Contains(listResp.Body.String(), "\"items\"") {
		t.Fatalf("list restore points failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}

	detailResp := performBackupRestoreContractRequest(
		t,
		ctx.Router,
		ctx.Token,
		http.MethodGet,
		"/api/v1/backup-restore/restore-points/"+strconv.FormatUint(restorePointID, 10),
		"",
	)
	if detailResp.Code != http.StatusOK || !strings.Contains(detailResp.Body.String(), "\"id\":") {
		t.Fatalf("get restore point failed status=%d body=%s", detailResp.Code, strings.TrimSpace(detailResp.Body.String()))
	}
}
