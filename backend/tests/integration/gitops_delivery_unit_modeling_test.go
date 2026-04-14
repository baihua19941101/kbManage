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

func TestGitOpsDeliveryUnitModelingIntegration_CreateUpdateAndStatus(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-delivery-unit-modeling-int",
		Password: "GitOpsUnitModeling@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-delivery-unit-modeling-int", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	sourceID := createGitOpsDeliveryUnitModelingSource(t, app.Router, token, access.WorkspaceID, access.ProjectID)
	targetGroupID := createGitOpsDeliveryUnitModelingTargetGroup(t, app.Router, token, access.WorkspaceID, access.ProjectID, access.ClusterID)

	createBody := fmt.Sprintf(`{
		"name":"orders-modeling-unit",
		"workspaceId":%d,
		"projectId":%d,
		"sourceId":%d,
		"sourcePath":"services/orders",
		"defaultNamespace":"orders",
		"syncMode":"manual",
		"desiredRevision":"main",
		"desiredAppVersion":"1.0.0",
		"desiredConfigVersion":"cfg-v1",
		"environments":[
			{"name":"test","orderIndex":10,"targetGroupId":%d,"promotionMode":"manual"},
			{"name":"prod","orderIndex":20,"targetGroupId":%d,"promotionMode":"manual"}
		],
		"overlays":[
			{"overlayType":"values","overlayRef":"values/common.yaml","precedence":5,"effectiveScope":"global"},
			{"overlayType":"values","overlayRef":"values/test.yaml","precedence":10,"effectiveScope":"env:test"}
		]
	}`, access.WorkspaceID, access.ProjectID, sourceID, targetGroupID, targetGroupID)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/delivery-units", strings.NewReader(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create delivery-unit status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	created := mustDecodeGitOpsDeliveryUnitModelingObject(t, createResp.Body.Bytes())
	unitID := mustReadGitOpsDeliveryUnitModelingID(t, created, "id")

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResp := httptest.NewRecorder()
	app.Router.ServeHTTP(getResp, getReq)

	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get delivery-unit status=200, got=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
	getPayload := mustDecodeGitOpsDeliveryUnitModelingObject(t, getResp.Body.Bytes())
	if len(mustReadGitOpsDeliveryUnitModelingArray(t, getPayload, "environments")) != 2 {
		t.Fatalf("expected two environments in detail payload=%v", getPayload)
	}
	if len(mustReadGitOpsDeliveryUnitModelingArray(t, getPayload, "overlays")) != 2 {
		t.Fatalf("expected two overlays in detail payload=%v", getPayload)
	}

	updateBody := fmt.Sprintf(`{
		"desiredRevision":"release/2026.04",
		"desiredAppVersion":"1.1.0",
		"desiredConfigVersion":"cfg-v2",
		"environments":[
			{"name":"staging","orderIndex":15,"targetGroupId":%d,"promotionMode":"manual"},
			{"name":"prod","orderIndex":20,"targetGroupId":%d,"promotionMode":"manual"}
		],
		"overlays":[
			{"overlayType":"values","overlayRef":"values/common.yaml","precedence":5,"effectiveScope":"global"},
			{"overlayType":"values","overlayRef":"values/staging.yaml","precedence":12,"effectiveScope":"env:staging"}
		]
	}`, targetGroupID, targetGroupID)

	updateReq := httptest.NewRequest(http.MethodPatch, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10), strings.NewReader(updateBody))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	app.Router.ServeHTTP(updateResp, updateReq)

	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected update delivery-unit status=200, got=%d body=%s", updateResp.Code, strings.TrimSpace(updateResp.Body.String()))
	}
	updated := mustDecodeGitOpsDeliveryUnitModelingObject(t, updateResp.Body.Bytes())
	if strings.TrimSpace(mustReadGitOpsDeliveryUnitModelingString(t, updated, "desiredRevision")) != "release/2026.04" {
		t.Fatalf("expected desiredRevision updated payload=%v", updated)
	}

	statusReq := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10)+"/status?environment=staging", nil)
	statusReq.Header.Set("Authorization", "Bearer "+token)
	statusResp := httptest.NewRecorder()
	app.Router.ServeHTTP(statusResp, statusReq)

	if statusResp.Code != http.StatusOK {
		t.Fatalf("expected delivery-unit status endpoint=200, got=%d body=%s", statusResp.Code, strings.TrimSpace(statusResp.Body.String()))
	}
	statusPayload := mustDecodeGitOpsDeliveryUnitModelingObject(t, statusResp.Body.Bytes())
	if strings.TrimSpace(mustReadGitOpsDeliveryUnitModelingString(t, statusPayload, "deliveryStatus")) == "" {
		t.Fatalf("expected deliveryStatus not empty payload=%v", statusPayload)
	}
	if strings.TrimSpace(mustReadGitOpsDeliveryUnitModelingString(t, statusPayload, "driftStatus")) == "" {
		t.Fatalf("expected driftStatus not empty payload=%v", statusPayload)
	}
	statusItems := mustReadGitOpsDeliveryUnitModelingArray(t, statusPayload, "environments")
	if len(statusItems) != 1 {
		t.Fatalf("expected environment-filtered status result len=1 payload=%v", statusPayload)
	}
	envItem, _ := statusItems[0].(map[string]any)
	if strings.TrimSpace(mustReadGitOpsDeliveryUnitModelingString(t, envItem, "environment")) != "staging" {
		t.Fatalf("expected environment=staging payload=%v", envItem)
	}
}

func createGitOpsDeliveryUnitModelingSource(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"name":"orders-source-modeling",
		"sourceType":"git",
		"endpoint":"https://git.example.test/demo/orders-modeling.git",
		"workspaceId":%d,
		"projectId":%d
	}`, workspaceID, projectID)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/sources", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create source fixture failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeGitOpsDeliveryUnitModelingObject(t, resp.Body.Bytes())
	return mustReadGitOpsDeliveryUnitModelingID(t, payload, "id")
}

func createGitOpsDeliveryUnitModelingTargetGroup(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"name":"orders-target-modeling",
		"workspaceId":%d,
		"projectId":%d,
		"clusterRefs":[%d]
	}`, workspaceID, projectID, clusterID)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/target-groups", strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()
	r.ServeHTTP(resp, req)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create target-group fixture failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	payload := mustDecodeGitOpsDeliveryUnitModelingObject(t, resp.Body.Bytes())
	return mustReadGitOpsDeliveryUnitModelingID(t, payload, "id")
}

func mustDecodeGitOpsDeliveryUnitModelingObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadGitOpsDeliveryUnitModelingString(t *testing.T, payload map[string]any, key string) string {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	text, ok := raw.(string)
	if !ok {
		t.Fatalf("field %q must be string, got=%T value=%v", key, raw, raw)
	}
	return text
}

func mustReadGitOpsDeliveryUnitModelingID(t *testing.T, payload map[string]any, key string) uint64 {
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

func mustReadGitOpsDeliveryUnitModelingArray(t *testing.T, payload map[string]any, key string) []any {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("field %q must be array, got=%T value=%v", key, raw, raw)
	}
	return items
}
