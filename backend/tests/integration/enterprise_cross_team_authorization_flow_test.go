package integration_test

import (
	"net/http"
	"testing"
)

func TestEnterprisePolishIntegration_CoverageAndActionFlow(t *testing.T) {
	ctx := newEnterpriseIntegrationCtx(t, "workspace-owner")
	seedEnterpriseIntegrationData(t, ctx)
	resp := performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/governance/coverage", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("coverage status=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/governance/action-items", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("action items status=%d body=%s", resp.Code, resp.Body.String())
	}
}
