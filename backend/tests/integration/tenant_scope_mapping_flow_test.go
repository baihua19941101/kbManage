package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestIdentityTenancyIntegration_TenantScopeMappingFlow(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	createOrg := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/organizations", `{"unitType":"organization","name":"mapping-org"}`)
	if createOrg.Code != http.StatusCreated {
		t.Fatalf("create org failed status=%d body=%s", createOrg.Code, strings.TrimSpace(createOrg.Body.String()))
	}
	resp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/organizations/1/mappings", `{
		"scopeType":"workspace",
		"scopeRef":"1",
		"inheritanceMode":"direct"
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create mapping failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
