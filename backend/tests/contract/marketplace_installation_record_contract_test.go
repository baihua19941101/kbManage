package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformMarketplaceContract_InstallationRecordQuery(t *testing.T) {
	ctx := newMarketplaceContractCtx(t, "workspace-owner")
	sourceID := createMarketplaceCatalogSourceContract(t, ctx, "installation-contract")
	syncMarketplaceCatalogSourceContract(t, ctx, sourceID)
	_ = performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/templates/1/releases", `{
		"versionId":2,
		"scopeType":"workspace",
		"scopeId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,
		"visibilityMode":"scope"
	}`)

	resp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/installations?scopeType=workspace&scopeId="+strconv.FormatUint(ctx.Access.WorkspaceID, 10), "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"scopeType":"workspace"`) {
		t.Fatalf("installation record query failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
