package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestGitOpsRollbackIntegration_SubmitAndQueryOperation(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-rollback-integration",
		Password: "GitOpsRollback@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-rollback-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	unitID := createGitOpsUS2IntegrationDeliveryUnit(t, app.Router, token, access.WorkspaceID, access.ProjectID, access.ClusterID, "rollback")
	operationID := submitGitOpsUS2IntegrationAction(t, app.Router, token, unitID, "rollback")
	getGitOpsUS2IntegrationOperation(t, app.Router, token, operationID)

	diffReq := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10)+"/diff", nil)
	diffReq.Header.Set("Authorization", "Bearer "+token)
	diffResp := httptest.NewRecorder()
	app.Router.ServeHTTP(diffResp, diffReq)
	if diffResp.Code != http.StatusOK {
		t.Fatalf("expected diff status=200 got=%d body=%s", diffResp.Code, strings.TrimSpace(diffResp.Body.String()))
	}
}
