package integration_test

import (
	"net/http"
	"testing"
)

func TestEnterprisePolishIntegration_DeliveryBundleFlow(t *testing.T) {
	ctx := newEnterpriseIntegrationCtx(t, "workspace-owner")
	seedEnterpriseIntegrationData(t, ctx)
	resp := performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/delivery/artifacts", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("delivery artifacts status=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/delivery/bundles", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("delivery bundles status=%d body=%s", resp.Code, resp.Body.String())
	}
}
