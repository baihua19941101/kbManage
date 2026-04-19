package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyIntegration_OrganizationTreeFlow(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	rootResp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/organizations", `{"unitType":"organization","name":"root-org"}`)
	if rootResp.Code != http.StatusCreated {
		t.Fatalf("create root failed status=%d body=%s", rootResp.Code, strings.TrimSpace(rootResp.Body.String()))
	}
	resp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/organizations", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"name":"root-org"`) {
		t.Fatalf("organization tree failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
