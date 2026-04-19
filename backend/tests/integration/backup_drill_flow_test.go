package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_DrillAndAuditFlow(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")

	planResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/plans", `{
		"name":"monthly-drill",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"scopeSelection":{"namespaces":["orders"]},
		"rpoTargetMinutes":10,
		"rtoTargetMinutes":20,
		"roleAssignments":["sre","biz-owner"],
		"cutoverProcedure":["capture point","recover point"],
		"validationChecklist":["check deployment","check service"]
	}`)
	if planResp.Code != http.StatusCreated {
		t.Fatalf("create drill plan failed status=%d body=%s", planResp.Code, strings.TrimSpace(planResp.Body.String()))
	}

	runResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/plans/1/run", "")
	if runResp.Code != http.StatusAccepted || !strings.Contains(runResp.Body.String(), "\"status\":\"succeeded\"") {
		t.Fatalf("run drill failed status=%d body=%s", runResp.Code, strings.TrimSpace(runResp.Body.String()))
	}

	reportResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/records/1/report", "")
	if reportResp.Code != http.StatusCreated || !strings.Contains(reportResp.Body.String(), "RPO/RTO") {
		t.Fatalf("generate report failed status=%d body=%s", reportResp.Code, strings.TrimSpace(reportResp.Body.String()))
	}

	auditResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/audit/backup-restore/events?action=backuprestore.drill.report.generate", "")
	if auditResp.Code != http.StatusOK || !strings.Contains(auditResp.Body.String(), "backuprestore.drill.report.generate") {
		t.Fatalf("audit query failed status=%d body=%s", auditResp.Code, strings.TrimSpace(auditResp.Body.String()))
	}
}
