package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformMarketplaceIntegration_TemplateVersionFlow(t *testing.T) {
	ctx := newMarketplaceIntegrationCtx(t, "workspace-owner")
	_ = performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources", `{
		"name":"version-flow",
		"sourceType":"helm",
		"endpointRef":"https://catalog.example.test/version-flow",
		"visibilityScope":"platform",
		"templateSeeds":[
			{
				"name":"app-stack",
				"slug":"app-stack",
				"category":"app",
				"publishStatus":"active",
				"supportedScopes":["workspace"],
				"versions":[
					{"version":"1.0.0","status":"active","dependencies":["config"],"parameterSchemaSummary":"replicas","deploymentConstraintSummary":"workspace","releaseNotes":"基础版"},
					{"version":"1.1.0","status":"active","dependencies":["config","secret"],"parameterSchemaSummary":"replicas,image","deploymentConstraintSummary":"workspace","releaseNotes":"增强版","supersedesVersion":"1.0.0"}
				]
			}
		]
	}`)
	_ = performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources/1/sync", "")

	resp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/templates/1", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"supersedesVersionId":1`) {
		t.Fatalf("template version flow failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
