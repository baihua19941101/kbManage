package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreContract_DrillAndAuditQuery(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")

	planResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/plans", `{
		"name":"quarterly-drill",
		"description":"quarterly validation",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"scopeSelection":{"namespaces":["orders-prod"]},
		"rpoTargetMinutes":15,
		"rtoTargetMinutes":30,
		"roleAssignments":["sre:oncall","owner:orders"],
		"cutoverProcedure":["enable backup point","restore traffic"],
		"validationChecklist":["check pods","check traffic"]
	}`)
	if planResp.Code != http.StatusCreated {
		t.Fatalf("create drill plan failed status=%d body=%s", planResp.Code, strings.TrimSpace(planResp.Body.String()))
	}

	runResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/plans/1/run", "")
	if runResp.Code != http.StatusAccepted {
		t.Fatalf("run drill failed status=%d body=%s", runResp.Code, strings.TrimSpace(runResp.Body.String()))
	}

	reportResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/drills/records/1/report", "")
	if reportResp.Code != http.StatusCreated {
		t.Fatalf("generate drill report failed status=%d body=%s", reportResp.Code, strings.TrimSpace(reportResp.Body.String()))
	}

	auditResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/audit/backup-restore/events?action=backuprestore.drill.run", "")
	if auditResp.Code != http.StatusOK || !strings.Contains(auditResp.Body.String(), "backuprestore.drill.run") {
		t.Fatalf("query audit failed status=%d body=%s", auditResp.Code, strings.TrimSpace(auditResp.Body.String()))
	}
}
