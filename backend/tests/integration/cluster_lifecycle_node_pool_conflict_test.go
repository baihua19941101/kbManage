package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleIntegration_NodePoolConflict(t *testing.T) {
	ctx := newClusterLifecycleIntegrationCtx(t, "workspace-owner")
	createClusterForIntegration(t, ctx)
	seedRunningLifecycleOperation(t, ctx, 1)
	resp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/clusters/1/node-pools/1/scale", `{"desiredCount":4}`)
	if resp.Code != http.StatusConflict {
		t.Fatalf("expected conflict status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
