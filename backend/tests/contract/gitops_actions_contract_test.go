package contract_test

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

func TestGitOpsActionsContract_SubmitAndGetOperation(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-actions-contract",
		Password: "GitOpsActions@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-actions-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	unitID := createGitOpsUS2ContractDeliveryUnit(t, app.Router, token, access.WorkspaceID, access.ProjectID, access.ClusterID, "actions")

	submitReq := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10)+"/actions", strings.NewReader(`{
		"actionType":"sync"
	}`))
	submitReq.Header.Set("Authorization", "Bearer "+token)
	submitReq.Header.Set("Content-Type", "application/json")
	submitResp := httptest.NewRecorder()
	app.Router.ServeHTTP(submitResp, submitReq)
	if submitResp.Code != http.StatusAccepted {
		t.Fatalf("expected submit action status=202 got=%d body=%s", submitResp.Code, strings.TrimSpace(submitResp.Body.String()))
	}
	created := mustDecodeGitOpsUS2ContractObject(t, submitResp.Body.Bytes())
	operationID := mustReadGitOpsUS2ContractID(t, created, "id")
	if operationID == 0 {
		t.Fatalf("expected operation id > 0 payload=%v", created)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/operations/"+strconv.FormatUint(operationID, 10), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResp := httptest.NewRecorder()
	app.Router.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get operation status=200 got=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
}

func createGitOpsUS2ContractDeliveryUnit(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
	suffix string,
) uint64 {
	t.Helper()

	sourceID := createGitOpsUS2ContractSource(t, r, token, workspaceID, projectID, suffix)
	targetGroupID := createGitOpsUS2ContractTargetGroup(t, r, token, workspaceID, projectID, clusterID, suffix)

	body := fmt.Sprintf(`{
		"name":"orders-unit-%s",
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
	payload := mustDecodeGitOpsUS2ContractObject(t, resp.Body.Bytes())
	return mustReadGitOpsUS2ContractID(t, payload, "id")
}

func createGitOpsUS2ContractSource(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
	suffix string,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"name":"orders-source-%s",
		"sourceType":"git",
		"endpoint":"https://git.example.test/demo/orders-%s.git",
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
	payload := mustDecodeGitOpsUS2ContractObject(t, resp.Body.Bytes())
	return mustReadGitOpsUS2ContractID(t, payload, "id")
}

func createGitOpsUS2ContractTargetGroup(
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
		"name":"orders-target-%s",
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
	payload := mustDecodeGitOpsUS2ContractObject(t, resp.Body.Bytes())
	return mustReadGitOpsUS2ContractID(t, payload, "id")
}

func mustDecodeGitOpsUS2ContractObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadGitOpsUS2ContractID(t *testing.T, payload map[string]any, key string) uint64 {
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
