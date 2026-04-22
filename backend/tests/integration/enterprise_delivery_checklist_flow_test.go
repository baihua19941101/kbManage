package integration_test

import (
	"net/http"
	"strconv"
	"testing"
)

func TestEnterprisePolishIntegration_DeliveryChecklistFlow(t *testing.T) {
	ctx := newEnterpriseIntegrationCtx(t, "workspace-owner")
	bundleID := seedEnterpriseIntegrationData(t, ctx)
	resp := performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/delivery/bundles/"+strconv.FormatUint(bundleID, 10)+"/checklists", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("delivery checklist status=%d body=%s", resp.Code, resp.Body.String())
	}
}
