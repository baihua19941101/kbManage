package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleIntegration_CreateFlow(t *testing.T) {
	ctx := newClusterLifecycleIntegrationCtx(t, "workspace-owner")
	createClusterForIntegration(t, ctx)
	detailResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/clusters/1", "")
	if detailResp.Code != http.StatusOK {
		t.Fatalf("detail failed status=%d body=%s", detailResp.Code, strings.TrimSpace(detailResp.Body.String()))
	}
	if !strings.Contains(detailResp.Body.String(), `"nodePools"`) || !strings.Contains(detailResp.Body.String(), `"status":"active"`) {
		t.Fatalf("unexpected detail body=%s", strings.TrimSpace(detailResp.Body.String()))
	}
}
