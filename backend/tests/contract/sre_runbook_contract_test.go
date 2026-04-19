package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformSREContract_ListRunbooks(t *testing.T) {
	ctx := newSREContractCtx(t, "workspace-owner")
	seedSRECapacityEvidence(t, ctx)
	resp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/runbooks", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list runbooks failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if len(mustDecodeSREItems(t, resp.Body.Bytes())) == 0 {
		t.Fatal("expected non-empty runbooks")
	}
}
