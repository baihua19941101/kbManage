package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformSREContract_ListScaleEvidence(t *testing.T) {
	ctx := newSREContractCtx(t, "workspace-owner")
	seedSRECapacityEvidence(t, ctx)
	resp := performSREContractRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/scale-evidence?evidenceType=loadtest", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list scale evidence failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	if len(mustDecodeSREItems(t, resp.Body.Bytes())) == 0 {
		t.Fatal("expected non-empty scale evidence")
	}
}
