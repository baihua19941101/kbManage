package integration_test

import (
	"net/http"
	"testing"
)

func TestEnterprisePolishIntegration_PermissionAuditFlow(t *testing.T) {
	ctx := newEnterpriseIntegrationCtx(t, "workspace-owner")
	seedEnterpriseIntegrationData(t, ctx)
	resp := performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/audit/permission-trails", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("permission trails status=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/audit/key-operations", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("key operations status=%d body=%s", resp.Code, resp.Body.String())
	}
}
