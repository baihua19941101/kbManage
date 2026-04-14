package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestGitOpsSyncExecutionIntegration_SubmitAndQueryOperation(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-sync-exec-integration",
		Password: "GitOpsSync@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-sync-exec-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	unitID := createGitOpsUS2IntegrationDeliveryUnit(t, app.Router, token, access.WorkspaceID, access.ProjectID, access.ClusterID, "sync")
	operationID := submitGitOpsUS2IntegrationAction(t, app.Router, token, unitID, "sync")
	getGitOpsUS2IntegrationOperation(t, app.Router, token, operationID)
}

func createGitOpsUS2IntegrationDeliveryUnit(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
	suffix string,
) uint64 {
	t.Helper()

	sourceID := createGitOpsUS2IntegrationSource(t, r, token, workspaceID, projectID, suffix)
	targetGroupID := createGitOpsUS2IntegrationTargetGroup(t, r, token, workspaceID, projectID, clusterID, suffix)

	body := fmt.Sprintf(`{
		"name":"orders-int-unit-%s",
		"workspaceId":%d,
		"projectId":%d,
		"sourceId":%d,
		"sourcePath":"services/orders",
		"defaultNamespace":"orders",
		"syncMode":"manual",
		"desiredRevision":"main",
		"desiredAppVersion":"1.0.0",
		"desiredConfigVersion":"cfg-v1",
		"environments":[{"name":"test","orderIndex":10,"targetGroupId":%d,"promotionMode":"manual"}]
	}`, suffix, workspaceID, projectID, sourceID, targetGroupID)

	req := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/delivery-units", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create delivery-unit fixture failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeGitOpsUS2IntegrationObject(t, resp.Body.Bytes())
	return mustReadGitOpsUS2IntegrationID(t, payload, "id")
}

func createGitOpsUS2IntegrationSource(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
	suffix string,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"name":"orders-int-source-%s",
		"sourceType":"git",
		"endpoint":"https://git.example.test/demo/orders-int-%s.git",
		"workspaceId":%d,
		"projectId":%d
	}`, suffix, suffix, workspaceID, projectID)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/sources", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create source fixture failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeGitOpsUS2IntegrationObject(t, resp.Body.Bytes())
	return mustReadGitOpsUS2IntegrationID(t, payload, "id")
}

func createGitOpsUS2IntegrationTargetGroup(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
	suffix string,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"name":"orders-int-target-%s",
		"workspaceId":%d,
		"projectId":%d,
		"clusterRefs":[%d]
	}`, suffix, workspaceID, projectID, clusterID)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/target-groups", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create target-group fixture failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeGitOpsUS2IntegrationObject(t, resp.Body.Bytes())
	return mustReadGitOpsUS2IntegrationID(t, payload, "id")
}

func submitGitOpsUS2IntegrationAction(
	t *testing.T,
	r http.Handler,
	token string,
	unitID uint64,
	actionType string,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{"actionType":"%s"}`, actionType)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10)+"/actions", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusAccepted {
		t.Fatalf("submit action=%s expected status=202 got=%d body=%s", actionType, resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeGitOpsUS2IntegrationObject(t, resp.Body.Bytes())
	return mustReadGitOpsUS2IntegrationID(t, payload, "id")
}

func getGitOpsUS2IntegrationOperation(
	t *testing.T,
	r http.Handler,
	token string,
	operationID uint64,
) {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/operations/"+strconv.FormatUint(operationID, 10), nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusOK {
		t.Fatalf("get operation expected status=200 got=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}

func mustDecodeGitOpsUS2IntegrationObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadGitOpsUS2IntegrationID(t *testing.T, payload map[string]any, key string) uint64 {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	number, ok := raw.(float64)
	if !ok || number <= 0 {
		t.Fatalf("field %q must be positive number, got=%T value=%v", key, raw, raw)
	}
	return uint64(number)
}
