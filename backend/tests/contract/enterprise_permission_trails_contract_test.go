package contract_test

import (
	"net/http"
	"testing"
)

func TestEnterprisePolishContract_PermissionTrailsAndKeyOperations(t *testing.T) {
	ctx := newEnterpriseContractCtx(t, "workspace-owner")
	seedEnterpriseAuditData(t, ctx)
	resp := performEnterpriseContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/audit/permission-trails", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("permission trails status=%d body=%s", resp.Code, resp.Body.String())
	}
	if len(mustDecodeEnterpriseItems(t, resp.Body.Bytes())) == 0 {
		t.Fatal("expected permission trails")
	}
	resp = performEnterpriseContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/enterprise/audit/key-operations", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("key operations status=%d body=%s", resp.Code, resp.Body.String())
	}
	if len(mustDecodeEnterpriseItems(t, resp.Body.Bytes())) == 0 {
		t.Fatal("expected key operations")
	}
}
