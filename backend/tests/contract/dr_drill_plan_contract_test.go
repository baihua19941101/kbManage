package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreContract_DRDrillPlan(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")
	createBackupRestoreContractDrillPlan(t, ctx, "drill-plan-contract")

	listResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/backup-restore/drills/plans", "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), "drill-plan-contract") {
		t.Fatalf("list drill plans failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
