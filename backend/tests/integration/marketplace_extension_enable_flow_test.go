package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformMarketplaceIntegration_ExtensionEnableFlow(t *testing.T) {
	ctx := newMarketplaceIntegrationCtx(t, "workspace-owner")
	createResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/extensions", `{
		"name":"ops-console",
		"extensionType":"plugin",
		"version":"1.0.0",
		"visibilityScope":"workspace:`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`",
		"entrySummary":"运维控制台",
		"permissionDeclaration":["marketplace:read"],
		"compatibility":[{"targetType":"platform-version","targetRef":"current","result":"compatible","summary":"通过兼容性校验"}]
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("register extension failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	enableResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/extensions/1/enable", `{
		"scopeType":"workspace",
		"scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`
	}`)
	if enableResp.Code != http.StatusAccepted || !strings.Contains(enableResp.Body.String(), `"status":"enabled"`) {
		t.Fatalf("enable extension failed status=%d body=%s", enableResp.Code, strings.TrimSpace(enableResp.Body.String()))
	}

	compatResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/extensions/1/compatibility", "")
	if compatResp.Code != http.StatusOK || !strings.Contains(compatResp.Body.String(), `"blockedReasons":[]`) {
		t.Fatalf("get extension compatibility failed status=%d body=%s", compatResp.Code, strings.TrimSpace(compatResp.Body.String()))
	}
}
