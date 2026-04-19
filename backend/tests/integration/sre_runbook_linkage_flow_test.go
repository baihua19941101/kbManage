package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformSREIntegration_RunbookLinkageFlow(t *testing.T) {
	ctx := newSREIntegrationCtx(t, "workspace-owner")
	seedSREIntegrationEvidence(t, ctx)
	resp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/runbooks", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list runbooks failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	auditResp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/audit/sre/events", "")
	if auditResp.Code != http.StatusOK {
		t.Fatalf("list sre audit events failed status=%d body=%s", auditResp.Code, strings.TrimSpace(auditResp.Body.String()))
	}
}
