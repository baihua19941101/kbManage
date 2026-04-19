package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreContract_DRDrillReport(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")
	planID := createBackupRestoreContractDrillPlan(t, ctx, "drill-report-contract")
	runResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/plans/"+strconv.FormatUint(planID, 10)+"/run", "")
	recordID := mustReadBackupRestoreContractID(t, runResp.Body.Bytes(), "id")

	reportResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/records/"+strconv.FormatUint(recordID, 10)+"/report", "")
	if reportResp.Code != http.StatusCreated || !strings.Contains(reportResp.Body.String(), "\"drillRecordId\":") {
		t.Fatalf("create drill report failed status=%d body=%s", reportResp.Code, strings.TrimSpace(reportResp.Body.String()))
	}
}
