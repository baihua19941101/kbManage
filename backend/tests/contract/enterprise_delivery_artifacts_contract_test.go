package contract_test

import (
	"net/http"
	"testing"
)

func TestEnterprisePolishContract_DeliveryArtifacts(t *testing.T) {
	ctx := newEnterpriseContractCtx(t, "workspace-owner")
	seedEnterpriseAuditData(t, ctx)
	resp := performEnterpriseContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/delivery/artifacts", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("delivery artifacts status=%d body=%s", resp.Code, resp.Body.String())
	}
	if len(mustDecodeEnterpriseItems(t, resp.Body.Bytes())) == 0 {
		t.Fatal("expected delivery artifacts")
	}
}
