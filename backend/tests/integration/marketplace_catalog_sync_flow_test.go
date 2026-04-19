package integration_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformMarketplaceIntegration_CatalogSyncFlow(t *testing.T) {
	ctx := newMarketplaceIntegrationCtx(t, "workspace-owner")
	createResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources", `{
		"name":"integration-catalog",
		"sourceType":"helm",
		"endpointRef":"https://catalog.example.test/integration",
		"visibilityScope":"platform",
		"templateSeeds":[
			{
				"name":"redis-stack",
				"slug":"redis-stack",
				"category":"middleware",
				"summary":"标准 redis 模板",
				"publishStatus":"active",
				"supportedScopes":["workspace","project"],
				"versions":[
					{"version":"7.0.0","status":"active","dependencies":["storageclass"],"parameterSchemaSummary":"size","deploymentConstraintSummary":"workspace/project","releaseNotes":"基础版本"}
				]
			}
		]
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create catalog source failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	syncResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources/1/sync", "")
	if syncResp.Code != http.StatusAccepted {
		t.Fatalf("sync catalog source failed status=%d body=%s", syncResp.Code, strings.TrimSpace(syncResp.Body.String()))
	}

	templateResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/templates", "")
	if templateResp.Code != http.StatusOK || !strings.Contains(templateResp.Body.String(), `"slug":"redis-stack"`) {
		t.Fatalf("list templates failed status=%d body=%s", templateResp.Code, strings.TrimSpace(templateResp.Body.String()))
	}
}
