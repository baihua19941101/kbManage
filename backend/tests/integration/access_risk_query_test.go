package integration_test

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestIdentityTenancyIntegration_AccessRiskQuery(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/roles", `{"name":"risk-role","roleLevel":"project","permissionSummary":"read","inheritancePolicy":"downward-allowed","delegable":true}`)
	performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/assignments", `{
		"subjectType":"user",
		"subjectRef":"1",
		"roleDefinitionId":1,
		"scopeType":"project",
		"scopeRef":"1",
		"sourceType":"temporary",
		"validUntil":"`+time.Now().UTC().Add(time.Hour).Format(time.RFC3339)+`"
	}`)
	resp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/access-risks", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"severity"`) {
		t.Fatalf("risk query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
