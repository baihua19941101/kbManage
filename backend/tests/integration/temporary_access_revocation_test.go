package integration_test

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestIdentityTenancyIntegration_TemporaryAccessRevocation(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/roles", `{"name":"temp-role","roleLevel":"project","permissionSummary":"read","inheritancePolicy":"bounded","delegable":true}`)
	performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/assignments", `{
		"subjectType":"user",
		"subjectRef":"1",
		"roleDefinitionId":1,
		"scopeType":"project",
		"scopeRef":"1",
		"sourceType":"temporary",
		"validUntil":"`+time.Now().UTC().Add(-time.Minute).Format(time.RFC3339)+`"
	}`)
	resp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/identity/assignments", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"status":"expired"`) {
		t.Fatalf("expected expired assignment status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
