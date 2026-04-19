package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyIntegration_IdentitySourceLifecycle(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	resp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/sources", `{
		"name":"corp-oidc",
		"sourceType":"oidc",
		"loginMode":"optional",
		"scopeMode":"platform"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create source failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	listResp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/sources", "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), `"sourceType":"oidc"`) {
		t.Fatalf("list source failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
}
