package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Create(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	createClusterLifecycleDriverForContract(t, ctx, "")
	resp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters", `{
		"name":"create-contract",
		"displayName":"Create Contract",
		"workspaceId":1,
		"projectId":1,
		"infrastructureType":"generic",
		"driverKey":"generic-driver",
		"driverVersion":"v1",
		"kubernetesVersion":"v1.30.1"
	}`)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("create failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
