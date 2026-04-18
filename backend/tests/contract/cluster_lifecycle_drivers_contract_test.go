package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Drivers(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	createClusterLifecycleDriverForContract(t, ctx, "")
	listResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/drivers", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list drivers failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	if !strings.Contains(listResp.Body.String(), "generic-driver") {
		t.Fatalf("expected driver in response body=%s", strings.TrimSpace(listResp.Body.String()))
	}
}
