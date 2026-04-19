package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformSREContract_ListCapacityBaselines(t *testing.T) {
	ctx := newSREContractCtx(t, "workspace-owner")
	seedSRECapacityEvidence(t, ctx)
	resp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/capacity/baselines", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list capacity baselines failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if len(mustDecodeSREItems(t, resp.Body.Bytes())) == 0 {
		t.Fatal("expected non-empty capacity baselines")
	}
}
