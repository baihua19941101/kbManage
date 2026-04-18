package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestClusterLifecycleIntegration_RegisterFlow(t *testing.T) {
	ctx := newClusterLifecycleIntegrationCtx(t, "workspace-owner")
	resp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/register", `{
		"name":"register-int",
		"displayName":"Register Int",
		"workspaceId":`+strconv.FormatUint(ctx.access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(ctx.access.ProjectID, 10)+`,
		"infrastructureType":"registered",
		"driverKey":"register-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("register failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	detailResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/clusters/1", "")
	if detailResp.Code != http.StatusOK || !strings.Contains(detailResp.Body.String(), `"registrationStatus":"issued"`) {
		t.Fatalf("detail after register failed status=%d body=%s", detailResp.Code, strings.TrimSpace(detailResp.Body.String()))
	}
}
