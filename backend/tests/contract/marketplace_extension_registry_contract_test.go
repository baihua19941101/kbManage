package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformMarketplaceContract_ExtensionLifecycle(t *testing.T) {
	ctx := newMarketplaceContractCtx(t, "workspace-owner")
	createResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/extensions", `{
		"name":"cost-insight",
		"extensionType":"plugin",
		"version":"1.0.0",
		"visibilityScope":"workspace:`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`",
		"entrySummary":"成本洞察扩展",
		"permissionDeclaration":["marketplace:read","marketplace:manage-extension"],
		"compatibility":[
			{"targetType":"platform-version","targetRef":"current","result":"compatible","summary":"可兼容当前平台"}
		]
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("register extension failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	compatResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/extensions/1/compatibility", "")
	if compatResp.Code != http.StatusOK || !strings.Contains(compatResp.Body.String(), `"result":"compatible"`) {
		t.Fatalf("get extension compatibility failed status=%d body=%s", compatResp.Code, strings.TrimSpace(compatResp.Body.String()))
	}

	enableResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/extensions/1/enable", `{
		"scopeType":"workspace",
		"scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`
	}`)
	if enableResp.Code != http.StatusAccepted || !strings.Contains(enableResp.Body.String(), `"status":"enabled"`) {
		t.Fatalf("enable extension failed status=%d body=%s", enableResp.Code, strings.TrimSpace(enableResp.Body.String()))
	}

	disableResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/extensions/1/disable", `{
		"scopeType":"workspace",
		"scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`
	}`)
	if disableResp.Code != http.StatusAccepted || !strings.Contains(disableResp.Body.String(), `"status":"disabled"`) {
		t.Fatalf("disable extension failed status=%d body=%s", disableResp.Code, strings.TrimSpace(disableResp.Body.String()))
	}
}
