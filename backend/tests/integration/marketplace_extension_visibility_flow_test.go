package integration_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformMarketplaceIntegration_ExtensionVisibilityBlockedFlow(t *testing.T) {
	ctx := newMarketplaceIntegrationCtx(t, "workspace-owner")
	createResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/extensions", `{
		"name":"blocked-ext",
		"extensionType":"plugin",
		"version":"1.0.0",
		"visibilityScope":"workspace:`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`",
		"entrySummary":"受阻扩展",
		"permissionDeclaration":["platform-admin"],
		"compatibility":[{"targetType":"platform-version","targetRef":"current","result":"blocked","summary":"需要更高平台版本"}]
	}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("register blocked extension failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	enableResp := performMarketplaceIntegrationRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/extensions/1/enable", `{
		"scopeType":"workspace",
		"scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`
	}`)
	if enableResp.Code != http.StatusConflict {
		t.Fatalf("enable blocked extension expected conflict status=%d body=%s", enableResp.Code, strings.TrimSpace(enableResp.Body.String()))
	}
}
