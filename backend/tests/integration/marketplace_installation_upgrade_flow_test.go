package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformMarketplaceIntegration_InstallationUpgradeFlow(t *testing.T) {
	ctx := newMarketplaceIntegrationCtx(t, "workspace-owner")
	_ = performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources", `{
		"name":"install-upgrade",
		"sourceType":"helm",
		"endpointRef":"https://catalog.example.test/install-upgrade",
		"visibilityScope":"platform",
		"templateSeeds":[
			{
				"name":"upgrade-stack",
				"slug":"upgrade-stack",
				"category":"app",
				"publishStatus":"active",
				"supportedScopes":["workspace"],
				"versions":[
					{"version":"1.0.0","status":"active","releaseNotes":"基础版"},
					{"version":"1.1.0","status":"active","releaseNotes":"升级版","supersedesVersion":"1.0.0"}
				]
			}
		]
	}`)
	_ = performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources/1/sync", "")
	_ = performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/templates/1/releases", `{"versionId":1,"scopeType":"workspace","scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,"visibilityMode":"scope"}`)
	upgradeResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/templates/1/releases", `{"versionId":2,"scopeType":"workspace","scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,"visibilityMode":"scope"}`)
	if upgradeResp.Code != http.StatusCreated {
		t.Fatalf("upgrade release failed status=%d body=%s", upgradeResp.Code, strings.TrimSpace(upgradeResp.Body.String()))
	}
	resp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/installations?scopeType=workspace&scopeId="+strconv.FormatUint(ctx.Access.WorkspaceID, 10), "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"lifecycleStatus":"upgraded"`) {
		t.Fatalf("installation upgrade flow failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
