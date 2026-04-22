package contract_test

import (
	"net/http"
	"strconv"
	"testing"
)

func TestEnterprisePolishContract_DeliveryBundles(t *testing.T) {
	ctx := newEnterpriseContractCtx(t, "workspace-owner")
	bundleID := seedEnterpriseAuditData(t, ctx)
	resp := performEnterpriseContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/delivery/bundles", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("delivery bundles status=%d body=%s", resp.Code, resp.Body.String())
	}
	resp = performEnterpriseContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/delivery/bundles/"+strconv.FormatUint(bundleID, 10)+"/checklists", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("delivery checklist status=%d body=%s", resp.Code, resp.Body.String())
	}
	if len(mustDecodeEnterpriseItems(t, resp.Body.Bytes())) == 0 {
		t.Fatal("expected checklist items")
	}
}
