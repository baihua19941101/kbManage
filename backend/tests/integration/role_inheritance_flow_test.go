package integration_test

import (
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestIdentityTenancyIntegration_RoleInheritanceFlow(t *testing.T) {
	ctx := newIdentityTenancyIntegrationCtx(t, "workspace-owner")
	roleResp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/roles", `{"name":"inherit-role","roleLevel":"project","permissionSummary":"read","inheritancePolicy":"downward-allowed","delegable":true}`)
	if roleResp.Code != http.StatusCreated {
		t.Fatalf("create role failed status=%d body=%s", roleResp.Code, strings.TrimSpace(roleResp.Body.String()))
	}
	assignResp := performIdentityTenancyIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/identity/assignments", `{
		"subjectType":"user",
		"subjectRef":"1",
		"roleDefinitionId":1,
		"scopeType":"project",
		"scopeRef":"1",
		"sourceType":"temporary",
		"validUntil":"`+time.Now().UTC().Add(time.Hour).Format(time.RFC3339)+`"
	}`)
	if assignResp.Code != http.StatusCreated {
		t.Fatalf("create assignment failed status=%d body=%s", assignResp.Code, strings.TrimSpace(assignResp.Body.String()))
	}
}
