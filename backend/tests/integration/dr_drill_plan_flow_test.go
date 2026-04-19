package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_DRDrillPlanFlow(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	createBackupRestoreIntegrationDrillPlan(t, ctx, "drill-plan-int")

	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/backup-restore/drills/plans", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), "drill-plan-int") {
		t.Fatalf("drill plan flow failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
