package contract_test

import (
	"net/http"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestObservabilityContract_NotificationTargetsRoute(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "obs-targets-contract-user",
		Password:    "ObsTargets@123",
		DisplayName: "Obs Targets Contract User",
		Email:       "obs-targets-contract-user@example.test",
	})
	_ = testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "obs-targets-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)

	createResp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodPost,
		"/api/v1/observability/notification-targets",
		`{"name":"oncall-webhook","targetType":"webhook","configRef":"secret://ops/oncall"}`,
	)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	listResp := performObservabilityAuthedRequest(
		t,
		app.Router,
		token,
		http.MethodGet,
		"/api/v1/observability/notification-targets",
		"",
	)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	payload := mustDecodeObservabilityObject(t, listResp.Body.Bytes())
	assertObservabilityArrayField(t, payload, "items")
}
