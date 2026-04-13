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

func TestWorkloadOpsTerminalContract_CreateGetClose(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "workloadops-terminal-contract",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-terminal-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/terminal/sessions", strings.NewReader(`{
		"clusterId": `+strconv.FormatUint(access.ClusterID, 10)+`,
		"namespace": "default",
		"podName": "demo-api-pod-0",
		"containerName": "app",
		"workloadKind": "Deployment",
		"workloadName": "demo-api"
	}`))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	var created map[string]any
	if err := json.Unmarshal(createResp.Body.Bytes(), &created); err != nil {
		t.Fatalf("invalid create response: %v", err)
	}
	sessionID := intValue(created["id"])
	if sessionID <= 0 {
		t.Fatalf("invalid session id: %v", created)
	}
	if strings.TrimSpace(stringValue(created["status"])) != "active" {
		t.Fatalf("expected active status: %v", created)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/terminal/sessions/"+strconv.Itoa(sessionID), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResp := httptest.NewRecorder()
	app.Router.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get status=200, got=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}

	deleteReq := httptest.NewRequest(http.MethodDelete, "/api/v1/workload-ops/terminal/sessions/"+strconv.Itoa(sessionID), nil)
	deleteReq.Header.Set("Authorization", "Bearer "+token)
	deleteResp := httptest.NewRecorder()
	app.Router.ServeHTTP(deleteResp, deleteReq)
	if deleteResp.Code != http.StatusNoContent {
		t.Fatalf("expected delete status=204, got=%d body=%s", deleteResp.Code, strings.TrimSpace(deleteResp.Body.String()))
	}
}

func intValue(v any) int {
	switch x := v.(type) {
	case float64:
		return int(x)
	case int:
		return x
	default:
		return 0
	}
}
