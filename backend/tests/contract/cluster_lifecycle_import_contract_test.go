package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Import(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	resp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/import", `{
		"name":"imported-contract",
		"displayName":"Imported Contract",
		"workspaceId":`+strconv.FormatUint(ctx.access.WorkspaceID, 10)+`,
		"projectId":0,
		"infrastructureType":"existing",
		"driverKey":"import-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"apiServer":"https://imported.example.test"
	}`)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("import failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if !strings.Contains(resp.Body.String(), `"operation"`) || !strings.Contains(resp.Body.String(), `"cluster"`) {
		t.Fatalf("unexpected import body=%s", strings.TrimSpace(resp.Body.String()))
	}
}
