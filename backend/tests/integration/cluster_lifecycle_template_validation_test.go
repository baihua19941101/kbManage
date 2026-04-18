package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestClusterLifecycleIntegration_TemplateValidation(t *testing.T) {
	ctx := newClusterLifecycleIntegrationCtx(t, "workspace-owner")
	createDriverForIntegration(t, ctx)
	templateResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/templates", `{
		"name":"template-int",
		"description":"integration template",
		"infrastructureType":"generic",
		"driverKey":"generic-driver",
		"driverVersionRange":"v1",
		"requiredCapabilities":["network"],
		"workspaceId":1,
		"projectId":1
	}`)
	if templateResp.Code != http.StatusCreated {
		t.Fatalf("template create failed status=%d body=%s", templateResp.Code, strings.TrimSpace(templateResp.Body.String()))
	}
	validateResp := performClusterLifecycleIntegrationRequest(t, ctx.app.Router, ctx.token, http.MethodPost, "/api/v1/cluster-lifecycle/templates/1/validate", `{
		"infrastructureType":"generic",
		"driverVersion":"v1"
	}`)
	if validateResp.Code != http.StatusOK || !strings.Contains(validateResp.Body.String(), `"overallStatus":"passed"`) {
		t.Fatalf("template validate failed status=%d body=%s", validateResp.Code, strings.TrimSpace(validateResp.Body.String()))
	}
}
