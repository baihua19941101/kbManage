package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreContract_RestoreAndMigration(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")

	performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies", `{
		"name":"namespace-orders",
		"scopeType":"namespace",
		"scopeRef":"orders-prod",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"executionMode":"manual",
		"retentionRule":"14d",
		"consistencyLevel":"application-consistent",
		"status":"active"
	}`)
	performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies/1/run", "")

	restoreResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs", `{
		"restorePointId":1,
		"jobType":"cross-cluster-restore",
		"sourceEnvironment":"prod",
		"targetEnvironment":"dr-site",
		"scopeSelection":{"namespaces":["orders-prod"]}
	}`)
	if restoreResp.Code != http.StatusAccepted {
		t.Fatalf("create restore job failed status=%d body=%s", restoreResp.Code, strings.TrimSpace(restoreResp.Body.String()))
	}

	validateResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs/1/validate", "")
	if validateResp.Code != http.StatusOK || !strings.Contains(validateResp.Body.String(), "consistencyNotice") {
		t.Fatalf("validate restore job failed status=%d body=%s", validateResp.Code, strings.TrimSpace(validateResp.Body.String()))
	}

	migrationResp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/migrations", `{
		"name":"orders-dr-cutover",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"sourceClusterId":11,
		"targetClusterId":22,
		"scopeSelection":{"namespaces":["orders-prod"]},
		"mappingRules":{"targetNamespace":"orders-dr"},
		"cutoverSteps":["freeze writes","switch traffic"]
	}`)
	if migrationResp.Code != http.StatusCreated {
		t.Fatalf("create migration plan failed status=%d body=%s", migrationResp.Code, strings.TrimSpace(migrationResp.Body.String()))
	}
}
