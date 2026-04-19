package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreIntegration_MigrationPlanFlow(t *testing.T) {
	ctx := newBackupRestoreIntegrationCtx(t, "workspace-owner")
	resp := performBackupRestoreIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/migrations", `{
		"name":"migration-flow-int",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"sourceClusterId":111,
		"targetClusterId":222,
		"scopeSelection":{"namespaces":["orders-prod"]},
		"mappingRules":{"targetNamespace":"orders-dr"},
		"cutoverSteps":["prepare target","switch traffic"]
	}`)
	if resp.Code != http.StatusCreated || !strings.Contains(resp.Body.String(), "\"status\":\"succeeded\"") {
		t.Fatalf("migration plan flow failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
