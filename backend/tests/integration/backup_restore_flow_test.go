package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_RestoreAndMigrationFlow(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")

	performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies", `{
		"name":"cluster-config",
		"scopeType":"cluster-config",
		"scopeRef":"cluster-a",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"executionMode":"manual",
		"retentionRule":"14d",
		"consistencyLevel":"platform-consistent",
		"status":"active"
	}`)
	performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/policies/1/run", "")

	restoreResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/restore-jobs", `{
		"restorePointId":1,
		"jobType":"selective-restore",
		"sourceEnvironment":"prod",
		"targetEnvironment":"staging",
		"scopeSelection":{"configMaps":["app-config"]}
	}`)
	if restoreResp.Code != http.StatusAccepted || !strings.Contains(restoreResp.Body.String(), "\"status\":\"succeeded\"") {
		t.Fatalf("restore flow failed status=%d body=%s", restoreResp.Code, strings.TrimSpace(restoreResp.Body.String()))
	}

	migrationResp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/migrations", `{
		"name":"cluster-config-migrate",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"sourceClusterId":101,
		"targetClusterId":202,
		"scopeSelection":{"configMaps":["app-config"]},
		"mappingRules":{"targetNamespace":"app-dr"},
		"cutoverSteps":["prepare target","switch target"]
	}`)
	if migrationResp.Code != http.StatusCreated || !strings.Contains(migrationResp.Body.String(), "\"status\":\"succeeded\"") {
		t.Fatalf("migration flow failed status=%d body=%s", migrationResp.Code, strings.TrimSpace(migrationResp.Body.String()))
	}
}
