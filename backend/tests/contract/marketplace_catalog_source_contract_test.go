package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestPlatformMarketplaceContract_CatalogSourceLifecycle(t *testing.T) {
	ctx := newMarketplaceContractCtx(t, "workspace-owner")
	sourceID := createMarketplaceCatalogSourceContract(t, ctx, "catalog-contract")

	listResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/catalog-sources", "")
	if listResp.Code != http.StatusOK || !strings.Contains(listResp.Body.String(), `"sourceType":"helm"`) {
		t.Fatalf("list catalog sources failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}

	syncMarketplaceCatalogSourceContract(t, ctx, sourceID)

	templateResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/templates", "")
	if templateResp.Code != http.StatusOK || !strings.Contains(templateResp.Body.String(), `"slug":"nginx-stack"`) {
		t.Fatalf("list templates failed status=%d body=%s", templateResp.Code, strings.TrimSpace(templateResp.Body.String()))
	}

	detailResp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/templates/1", "")
	if detailResp.Code != http.StatusOK || !strings.Contains(detailResp.Body.String(), `"version":"1.1.0"`) {
		t.Fatalf("get template detail failed status=%d body=%s", detailResp.Code, strings.TrimSpace(detailResp.Body.String()))
	}

	_ = strconv.FormatUint(sourceID, 10)
}
