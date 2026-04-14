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

func TestGitOpsDeliveryUnitsContract_CreateListAndStatus(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-delivery-units-contract",
		Password: "GitOpsUnits@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-delivery-units-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	sourceID := createGitOpsDeliveryUnitsContractSource(t, app.Router, token, access.WorkspaceID, access.ProjectID)
	targetGroupID := createGitOpsDeliveryUnitsContractTargetGroup(t, app.Router, token, access.WorkspaceID, access.ProjectID, access.ClusterID)

	createBody := fmt.Sprintf(`{
		"name":"orders-unit",
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
	created := mustDecodeGitOpsDeliveryUnitsContractObject(t, createResp.Body.Bytes())
	unitID := mustReadGitOpsDeliveryUnitsContractID(t, created, "id")
	if strings.TrimSpace(mustReadGitOpsDeliveryUnitsContractString(t, created, "name")) != "orders-unit" {
		t.Fatalf("expected name=orders-unit payload=%v", created)
	}
	if mustReadGitOpsDeliveryUnitsContractID(t, created, "sourceId") != sourceID {
		t.Fatalf("expected sourceId=%d payload=%v", sourceID, created)
	}
	if strings.TrimSpace(mustReadGitOpsDeliveryUnitsContractString(t, created, "deliveryStatus")) != "unknown" {
		t.Fatalf("expected initial deliveryStatus=unknown payload=%v", created)
	}
	envs := mustReadGitOpsDeliveryUnitsContractArray(t, created, "environments")
	if len(envs) != 2 {
		t.Fatalf("expected two environments payload=%v", created)
	}
	overlays := mustReadGitOpsDeliveryUnitsContractArray(t, created, "overlays")
	if len(overlays) != 1 {
		t.Fatalf("expected one overlay payload=%v", created)
	}

	listReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/gitops/delivery-units?workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+"&projectId="+strconv.FormatUint(access.ProjectID, 10),
		nil,
	)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)

	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list delivery-units status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	listPayload := mustDecodeGitOpsDeliveryUnitsContractObject(t, listResp.Body.Bytes())
	items := mustReadGitOpsDeliveryUnitsContractArray(t, listPayload, "items")
	if len(items) == 0 {
		t.Fatalf("expected non-empty delivery-unit items payload=%v", listPayload)
	}

	found := false
	for _, raw := range items {
		item, _ := raw.(map[string]any)
		if item == nil {
			continue
		}
		if mustReadGitOpsDeliveryUnitsContractID(t, item, "id") == unitID {
			found = true
			if strings.TrimSpace(mustReadGitOpsDeliveryUnitsContractString(t, item, "name")) != "orders-unit" {
				t.Fatalf("expected listed delivery-unit name=orders-unit item=%v", item)
			}
		}
	}
	if !found {
		t.Fatalf("expected list includes created delivery-unit id=%d payload=%v", unitID, listPayload)
	}

	statusReq := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/delivery-units/"+strconv.FormatUint(unitID, 10)+"/status", nil)
	statusReq.Header.Set("Authorization", "Bearer "+token)
	statusResp := httptest.NewRecorder()
	app.Router.ServeHTTP(statusResp, statusReq)

	if statusResp.Code != http.StatusOK {
		t.Fatalf("expected status endpoint status=200, got=%d body=%s", statusResp.Code, strings.TrimSpace(statusResp.Body.String()))
	}
	statusPayload := mustDecodeGitOpsDeliveryUnitsContractObject(t, statusResp.Body.Bytes())
	if strings.TrimSpace(mustReadGitOpsDeliveryUnitsContractString(t, statusPayload, "deliveryStatus")) == "" {
		t.Fatalf("expected deliveryStatus not empty payload=%v", statusPayload)
	}
	if strings.TrimSpace(mustReadGitOpsDeliveryUnitsContractString(t, statusPayload, "driftStatus")) == "" {
		t.Fatalf("expected driftStatus not empty payload=%v", statusPayload)
	}
}

func createGitOpsDeliveryUnitsContractSource(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"name":"orders-source-contract",
		"sourceType":"git",
		"endpoint":"https://git.example.test/demo/orders-contract.git",
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
	payload := mustDecodeGitOpsDeliveryUnitsContractObject(t, resp.Body.Bytes())
	return mustReadGitOpsDeliveryUnitsContractID(t, payload, "id")
}

func createGitOpsDeliveryUnitsContractTargetGroup(
	t *testing.T,
	r http.Handler,
	token string,
	workspaceID uint64,
	projectID uint64,
	clusterID uint64,
) uint64 {
	t.Helper()
	body := fmt.Sprintf(`{
		"name":"orders-target-contract",
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
	payload := mustDecodeGitOpsDeliveryUnitsContractObject(t, resp.Body.Bytes())
	return mustReadGitOpsDeliveryUnitsContractID(t, payload, "id")
}

func mustDecodeGitOpsDeliveryUnitsContractObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadGitOpsDeliveryUnitsContractString(t *testing.T, payload map[string]any, key string) string {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q payload=%v", key, payload)
	}
	val, ok := raw.(string)
	if !ok {
		t.Fatalf("field %q must be string, got=%T value=%v", key, raw, raw)
	}
	return val
}

func mustReadGitOpsDeliveryUnitsContractID(t *testing.T, payload map[string]any, key string) uint64 {
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

func mustReadGitOpsDeliveryUnitsContractArray(t *testing.T, payload map[string]any, key string) []any {
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
