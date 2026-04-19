package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyIntegration_LocalFallbackAccess(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/sources", `{"name":"sso-a","sourceType":"sso","loginMode":"optional","scopeMode":"platform"}`)
	resp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/sources", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"sourceType":"local"`) {
		t.Fatalf("expected local fallback status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
