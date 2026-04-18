package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Retire(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	createClusterLifecycleClusterForContract(t, ctx)
	disableResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/disable", `{"reason":"contract disable"}`)
	if disableResp.Code != http.StatusAccepted {
		t.Fatalf("disable failed status=%d body=%s", disableResp.Code, strings.TrimSpace(disableResp.Body.String()))
	}
	retireResp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/retire", `{
		"reason":"contract retire",
		"confirmationScope":"full",
		"conclusion":"complete"
	}`)
	if retireResp.Code != http.StatusAccepted {
		t.Fatalf("retire failed status=%d body=%s", retireResp.Code, strings.TrimSpace(retireResp.Body.String()))
	}
}
