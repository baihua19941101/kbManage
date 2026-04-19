package integration_test

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestIdentityTenancyIntegration_DelegationLifecycle(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	roleResp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/roles", `{"name":"delegable-role","roleLevel":"project","permissionSummary":"read","inheritancePolicy":"bounded","delegable":true}`)
	if roleResp.Code != http.StatusCreated {
		t.Fatalf("create role failed status=%d body=%s", roleResp.Code, strings.TrimSpace(roleResp.Body.String()))
	}
	grantResp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/delegations", `{
		"grantorRef":"1",
		"delegateRef":"2",
		"allowedRoleLevels":["project"],
		"validFrom":"`+time.Now().UTC().Format(time.RFC3339)+`",
		"validUntil":"`+time.Now().UTC().Add(time.Hour).Format(time.RFC3339)+`",
		"reason":"integration delegation"
	}`)
	if grantResp.Code != http.StatusCreated {
		t.Fatalf("create grant failed status=%d body=%s", grantResp.Code, strings.TrimSpace(grantResp.Body.String()))
	}
	assignResp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/assignments", `{
		"subjectType":"user",
		"subjectRef":"2",
		"roleDefinitionId":1,
		"scopeType":"project",
		"scopeRef":"1",
		"sourceType":"delegated",
		"delegationGrantId":1
	}`)
	if assignResp.Code != http.StatusCreated {
		t.Fatalf("delegated assignment failed status=%d body=%s", assignResp.Code, strings.TrimSpace(assignResp.Body.String()))
	}
}
