package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestClusterLifecycleIntegration_ImportFlow(t *testing.T) {
	ctx := newClusterLifecycleIntegrationCtx(t, "workspace-owner")
	resp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/import", `{
		"name":"import-int",
		"displayName":"Import Int",
		"workspaceId":`+strconv.FormatUint(ctx.access.WorkspaceID, 10)+`,
		"projectId":0,
		"infrastructureType":"existing",
		"driverKey":"import-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1",
		"apiServer":"https://import.example.test"
	}`)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("import failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	detailResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/clusters/1", "")
	if detailResp.Code != http.StatusOK || !strings.Contains(detailResp.Body.String(), `"status":"active"`) {
		t.Fatalf("detail after import failed status=%d body=%s", detailResp.Code, strings.TrimSpace(detailResp.Body.String()))
	}
}
