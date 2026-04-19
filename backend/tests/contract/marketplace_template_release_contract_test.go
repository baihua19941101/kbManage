package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformMarketplaceContract_TemplateReleaseAndInstallation(t *testing.T) {
	ctx := newMarketplaceContractCtx(t, "workspace-owner")
	sourceID := createMarketplaceCatalogSourceContract(t, ctx, "release-contract")
	syncMarketplaceCatalogSourceContract(t, ctx, sourceID)

	releaseResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/templates/1/releases", `{
		"versionId":2,
		"scopeType":"workspace",
		"scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"visibilityMode":"scope"
	}`)
	if releaseResp.Code != http.StatusCreated {
		t.Fatalf("create template release failed status=%d body=%s", releaseResp.Code, strings.TrimSpace(releaseResp.Body.String()))
	}

	listReleaseResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/templates/1/releases", "")
	if listReleaseResp.Code != http.StatusOK || !strings.Contains(listReleaseResp.Body.String(), `"scopeType":"workspace"`) {
		t.Fatalf("list template releases failed status=%d body=%s", listReleaseResp.Code, strings.TrimSpace(listReleaseResp.Body.String()))
	}

	installationResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/installations?scopeType=workspace&scopeId="+strconv.FormatUint(ctx.Access.WorkspaceID, 10), "")
	if installationResp.Code != http.StatusOK || !strings.Contains(installationResp.Body.String(), `"currentInstalledVersion":"1.1.0"`) {
		t.Fatalf("list installations failed status=%d body=%s", installationResp.Code, strings.TrimSpace(installationResp.Body.String()))
	}
}
