package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformMarketplaceIntegration_TemplateReleaseFlow(t *testing.T) {
	ctx := newMarketplaceIntegrationCtx(t, "workspace-owner")
	_ = performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources", `{
		"name":"release-flow",
		"sourceType":"helm",
		"endpointRef":"https://catalog.example.test/release-flow",
		"visibilityScope":"platform",
		"templateSeeds":[
			{"name":"nginx-stack","slug":"nginx-stack","category":"web","publishStatus":"active","supportedScopes":["workspace"],"versions":[{"version":"1.0.0","status":"active","releaseNotes":"基础版"}]}
		]
	}`)
	_ = performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources/1/sync", "")

	releaseResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/templates/1/releases", `{
		"versionId":1,
		"scopeType":"workspace",
		"scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"visibilityMode":"scope"
	}`)
	if releaseResp.Code != http.StatusCreated {
		t.Fatalf("create template release failed status=%d body=%s", releaseResp.Code, strings.TrimSpace(releaseResp.Body.String()))
	}

	installationResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/installations?scopeType=workspace&scopeId="+strconv.FormatUint(ctx.Access.WorkspaceID, 10), "")
	if installationResp.Code != http.StatusOK || !strings.Contains(installationResp.Body.String(), `"lifecycleStatus":"installed"`) {
		t.Fatalf("list installation records failed status=%d body=%s", installationResp.Code, strings.TrimSpace(installationResp.Body.String()))
	}
}
