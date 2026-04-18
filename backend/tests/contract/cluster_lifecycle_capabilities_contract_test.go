package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleContract_Capabilities(t *testing.T) {
	ctx := newClusterLifecycleContractCtx(t, "workspace-owner")
	createClusterLifecycleDriverForContract(t, ctx, `,
		"capabilities":[{"capabilityDomain":"network","supportLevel":"native","compatibilityStatus":"compatible"}]`)
	resp := performClusterLifecycleContractRequest(t, ctx.app.Router, ctx.token, http.MethodGet, "/api/v1/cluster-lifecycle/drivers/1/capabilities", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list capabilities failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if !strings.Contains(resp.Body.String(), "network") {
		t.Fatalf("expected capability body=%s", strings.TrimSpace(resp.Body.String()))
	}
}
