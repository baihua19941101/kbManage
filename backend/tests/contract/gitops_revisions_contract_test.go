package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestGitOpsRevisionsContract_ListReleases(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-revisions-contract",
		Password: "GitOpsRevisions@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-revisions-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	unitID := createGitOpsUS2ContractDeliveryUnit(t, app.Router, token, access.WorkspaceID, access.ProjectID, access.ClusterID, "revisions")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10)+"/releases", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected list releases status=200 got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}

	payload := mustDecodeGitOpsUS2ContractObject(t, resp.Body.Bytes())
	items, ok := payload["items"].([]any)
	if !ok {
		t.Fatalf("expected releases payload contains items array, payload=%v", payload)
	}
	if len(items) != 0 {
		_, _ = json.Marshal(items)
	}
}
