package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Upgrade(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	createClusterLifecycleClusterForContract(t, ctx)
	planResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/upgrade-plans", `{
		"targetVersion":"v1.31.0",
		"impactSummary":"contract upgrade"
	}`)
	if planResp.Code != http.StatusCreated {
		t.Fatalf("create plan failed status=%d body=%s", planResp.Code, strings.TrimSpace(planResp.Body.String()))
	}
	executeResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/upgrade-plans/1/execute", "")
	if executeResp.Code != http.StatusAccepted {
		t.Fatalf("execute plan failed status=%d body=%s", executeResp.Code, strings.TrimSpace(executeResp.Body.String()))
	}
}
