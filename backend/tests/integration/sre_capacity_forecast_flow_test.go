package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformSREIntegration_CapacityForecastFlow(t *testing.T) {
	ctx := newSREIntegrationCtx(t, "workspace-owner")
	seedSREIntegrationEvidence(t, ctx)
	resp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/capacity/baselines", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("list capacity baselines failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	evidenceResp := performSREIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodGet, "/api/v1/sre/scale-evidence?evidenceType=loadtest", "")
	if evidenceResp.Code != http.StatusOK {
		t.Fatalf("list scale evidence failed status=%d body=%s", evidenceResp.Code, strings.TrimSpace(evidenceResp.Body.String()))
	}
}
