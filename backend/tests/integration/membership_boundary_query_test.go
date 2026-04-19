package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyIntegration_MembershipBoundaryQuery(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/organizations", `{"unitType":"organization","name":"membership-org"}`)
	resp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/organizations/1/memberships", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"memberType":"user"`) {
		t.Fatalf("membership query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
