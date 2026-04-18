package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleIntegration_DriverManagement(t *testing.T) {
	ctx := newClusterLifecycleIntegrationCtx(t, "workspace-owner")
	createDriverForIntegration(t, ctx)
	driverResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/drivers", "")
	if driverResp.Code != http.StatusOK || !strings.Contains(driverResp.Body.String(), "generic-driver") {
		t.Fatalf("driver list failed status=%d body=%s", driverResp.Code, strings.TrimSpace(driverResp.Body.String()))
	}
	capResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/drivers/1/capabilities", "")
	if capResp.Code != http.StatusOK || !strings.Contains(capResp.Body.String(), "network") {
		t.Fatalf("capability list failed status=%d body=%s", capResp.Code, strings.TrimSpace(capResp.Body.String()))
	}
}
