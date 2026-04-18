package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_NodePool(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	createClusterLifecycleClusterForContract(t, ctx)
	listResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/clusters/1/node-pools", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list node pools failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	scaleResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/node-pools/1/scale", `{"desiredCount":4}`)
	if scaleResp.Code != http.StatusAccepted {
		t.Fatalf("scale node pool failed status=%d body=%s", scaleResp.Code, strings.TrimSpace(scaleResp.Body.String()))
	}
}
