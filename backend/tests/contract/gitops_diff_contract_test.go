package contract_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestGitOpsDiffContract_GetDiff(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-diff-contract",
		Password: "GitOpsDiff@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-diff-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	unitID := createGitOpsUS2ContractDeliveryUnit(t, app.Router, token, access.WorkspaceID, access.ProjectID, access.ClusterID, "diff")

	req := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10)+"/diff", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("expected get diff status=200 got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
