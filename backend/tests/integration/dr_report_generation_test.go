package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_DRReportGeneration(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	planID := createBackupRestoreIntegrationDrillPlan(t, ctx, "drill-report-int")
	runResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/plans/"+strconv.FormatUint(planID, 10)+"/run", "")
	recordID := mustReadBackupRestoreIntegrationID(t, runResp.Body.Bytes(), "id")

	reportResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/records/"+strconv.FormatUint(recordID, 10)+"/report", "")
	if reportResp.Code != http.StatusCreated || !strings.Contains(reportResp.Body.String(), "\"goalAssessment\":") {
		t.Fatalf("create drill report failed status=%d body=%s", reportResp.Code, strings.TrimSpace(reportResp.Body.String()))
	}
}
