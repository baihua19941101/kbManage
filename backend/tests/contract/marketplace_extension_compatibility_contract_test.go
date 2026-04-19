package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformMarketplaceContract_ExtensionCompatibilityQuery(t *testing.T) {
	ctx := newMarketplaceContractCtx(t, "workspace-owner")
	_ = performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/extensions", `{
		"name":"compat-query",
		"extensionType":"plugin",
		"version":"1.0.0",
		"visibilityScope":"platform",
		"entrySummary":"兼容性查询扩展",
		"permissionDeclaration":["marketplace:read"],
		"compatibility":[{"targetType":"platform-version","targetRef":"current","result":"compatible","summary":"可兼容"}]
	}`)

	resp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/extensions/1/compatibility", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"summary":"可兼容"`) {
		t.Fatalf("extension compatibility query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
