package integration_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestWorkloadTerminalSessionIntegration_CreateCloseLifecycle(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "workloadops-terminal-integration",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-terminal-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/terminal/sessions", strings.NewReader(`{
		"clusterId": `+strconv.FormatUint(access.ClusterID, 10)+`,
		"namespace": "default",
		"podName": "api-server-pod-0",
		"containerName": "app",
		"workloadKind": "Deployment",
		"workloadName": "api-server"
	}`))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}

	var created map[string]any
	_ = json.Unmarshal(createResp.Body.Bytes(), &created)
	sessionID := int(created["id"].(float64))
	if sessionID <= 0 {
		t.Fatalf("invalid session id: %v", created)
	}

	closeReq := httptest.NewRequest(http.MethodDelete, "/api/v1/workload-ops/terminal/sessions/"+strconv.Itoa(sessionID), nil)
	closeReq.Header.Set("Authorization", "Bearer "+token)
	closeResp := httptest.NewRecorder()
	app.Router.ServeHTTP(closeResp, closeReq)
	if closeResp.Code != http.StatusNoContent {
		t.Fatalf("expected close status=204, got=%d body=%s", closeResp.Code, strings.TrimSpace(closeResp.Body.String()))
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/terminal/sessions/"+strconv.Itoa(sessionID), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResp := httptest.NewRecorder()
	app.Router.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get status=200, got=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
	var payload map[string]any
	_ = json.Unmarshal(getResp.Body.Bytes(), &payload)
	if strings.TrimSpace(payload["status"].(string)) != "closed" {
		t.Fatalf("expected session status closed, payload=%v", payload)
	}
}
