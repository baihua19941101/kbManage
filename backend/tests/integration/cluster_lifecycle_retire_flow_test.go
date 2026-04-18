package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleIntegration_RetireFlow(t *testing.T) {
	ctx := newClusterLifecycleIntegrationCtx(t, "workspace-owner")
	createClusterForIntegration(t, ctx)
	disableResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/disable", `{"reason":"retire flow"}`)
	if disableResp.Code != http.StatusAccepted {
		t.Fatalf("disable failed status=%d body=%s", disableResp.Code, strings.TrimSpace(disableResp.Body.String()))
	}
	retireResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/retire", `{
		"reason":"retire flow",
		"confirmationScope":"full",
		"conclusion":"complete"
	}`)
	if retireResp.Code != http.StatusAccepted {
		t.Fatalf("retire failed status=%d body=%s", retireResp.Code, strings.TrimSpace(retireResp.Body.String()))
	}
}
