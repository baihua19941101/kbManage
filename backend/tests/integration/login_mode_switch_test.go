package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyIntegration_LoginModeSwitch(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/sources", `{"name":"oidc-a","sourceType":"oidc","loginMode":"optional","scopeMode":"platform"}`)
	performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/sources", `{"name":"ldap-a","sourceType":"ldap","loginMode":"fallback","scopeMode":"platform"}`)
	resp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/sources", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"loginModes"`) {
		t.Fatalf("login mode query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
