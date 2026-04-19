package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestBackupRestoreContract_MigrationPlan(t *testing.T) {
	ctx := newBackupRestoreContractCtx(t, "workspace-owner")
	resp := performBackupRestoreContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/backup-restore/migrations", `{
		"name":"migration-contract",
		"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.Access.ProjectID, 10)+`,
		"sourceClusterId":101,
		"targetClusterId":202,
		"scopeSelection":{"namespaces":["orders-prod"]},
		"mappingRules":{"targetNamespace":"orders-dr"},
		"cutoverSteps":["freeze writes","switch traffic"]
	}`)
	if resp.Code != http.StatusCreated || !strings.Contains(resp.Body.String(), "\"name\":\"migration-contract\"") {
		t.Fatalf("create migration plan failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
