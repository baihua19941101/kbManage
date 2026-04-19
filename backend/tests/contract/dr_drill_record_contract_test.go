package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreContract_DRDrillRecord(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")
	planID := createBackupRestoreContractDrillPlan(t, ctx, "drill-record-contract")

	runResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/plans/"+strconv.FormatUint(planID, 10)+"/run", "")
	if runResp.Code != http.StatusAccepted {
		t.Fatalf("run drill plan failed status=%d body=%s", runResp.Code, strings.TrimSpace(runResp.Body.String()))
	}
	recordID := mustReadBackupRestoreContractID(t, runResp.Body.Bytes(), "id")

	recordResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/backup-restore/drills/records/"+strconv.FormatUint(recordID, 10), "")
	if recordResp.Code != http.StatusOK || !strings.Contains(recordResp.Body.String(), "\"status\":") {
		t.Fatalf("get drill record failed status=%d body=%s", recordResp.Code, strings.TrimSpace(recordResp.Body.String()))
	}
}
