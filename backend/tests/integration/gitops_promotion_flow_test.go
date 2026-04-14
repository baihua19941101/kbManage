package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestGitOpsPromotionFlowIntegration_SubmitAndQueryOperation(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-promotion-integration",
		Password: "GitOpsPromotion@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-promotion-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	unitID := createGitOpsUS2IntegrationDeliveryUnit(t, app.Router, token, access.WorkspaceID, access.ProjectID, access.ClusterID, "promote")
	operationID := submitGitOpsUS2IntegrationAction(t, app.Router, token, unitID, "promote")
	getGitOpsUS2IntegrationOperation(t, app.Router, token, operationID)

	releasesReq := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10)+"/releases", nil)
	releasesReq.Header.Set("Authorization", "Bearer "+token)
	releasesResp := httptest.NewRecorder()
	app.Router.ServeHTTP(releasesResp, releasesReq)
	if releasesResp.Code != http.StatusOK {
		t.Fatalf("expected releases status=200 got=%d body=%s", releasesResp.Code, strings.TrimSpace(releasesResp.Body.String()))
	}
}
