package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Templates(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	resp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/templates", `{
		"name":"template-contract",
		"description":"contract template",
		"infrastructureType":"generic",
		"driverKey":"generic-driver",
		"driverVersionRange":"v1",
		"requiredCapabilities":["network"],
		"workspaceId":1,
		"projectId":1
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create template failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	listResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/templates", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list templates failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
