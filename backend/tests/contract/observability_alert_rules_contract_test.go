package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_AlertRulesRoute(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-rules-contract-user",
		Password:    "ObsRules@123",
		DisplayName: "Obs Rules Contract User",
		Email:       "obs-rules-contract-user@example.test",
	})
	seed := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-rules-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	createResp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodPost,
		"/api/v1/observability/alert-rules",
		`{"name":"cpu high","severity":"critical","conditionExpression":"cpu_usage > 80","scopeSnapshot":"{\"workspaceIds\":[`+strconv.FormatUint(seed.WorkspaceID, 10)+`]}"} `,
	)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	listResp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodGet,
		"/api/v1/observability/alert-rules",
		"",
	)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	payload := mustDecodeObservabilityObject(t, listResp.Body.Bytes())
	assertObservabilityArrayField(t, payload, "items")
}
