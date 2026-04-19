package contract_test

import (
	"net/http"
	"strings"
	"testing"
)

func TestPlatformMarketplaceContract_TemplateCatalogDetail(t *testing.T) {
	ctx := newMarketplaceContractCtx(t, "workspace-owner")
	sourceID := createMarketplaceCatalogSourceContract(t, ctx, "template-catalog")
	syncMarketplaceCatalogSourceContract(t, ctx, sourceID)

	resp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodGet, "/api/v1/marketplace/templates/1", "")
	if resp.Code != http.StatusOK || !strings.Contains(resp.Body.String(), `"parameterSchemaSummary":"replicas,image"`) {
		t.Fatalf("template catalog detail failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
