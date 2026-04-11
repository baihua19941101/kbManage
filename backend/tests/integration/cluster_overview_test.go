package integration_test

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"

	"kbmanage/backend/internal/repository"
	"kbmanage/backend/tests/testutil"
)

func TestClusterOverview_MultiClusterFlowSkeleton(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username:    "cluster-integration-user",
		Password:    "Cluster@123456",
		DisplayName: "Cluster Integration User",
		Email:       "cluster-integration-user@example.test",
	})
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	clusterIDs := make([]string, 0, 2)
	for _, clusterName := range []string{"cluster-a", "cluster-b"} {
		resp := performClusterAuthedRequest(t, app.Router, token, http.MethodPost, "/api/v1/clusters", `{
			"name":"`+clusterName+`",
			"apiServer":"https://`+clusterName+`.example.test",
			"authType":"kubeconfig",
			"kubeConfig":"apiVersion: v1"
		}`)
		if resp.Code != http.StatusCreated {
			t.Fatalf("expected create cluster status=201, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
		}
		clusterIDs = append(clusterIDs, extractIDField(t, resp.Body.Bytes(), "id"))
	}

	listClustersResp := performClusterAuthedRequest(t, app.Router, token, http.MethodGet, "/api/v1/clusters", "")
	if listClustersResp.Code != http.StatusOK {
		t.Fatalf("expected list clusters status=200, got status=%d body=%s", listClustersResp.Code, strings.TrimSpace(listClustersResp.Body.String()))
	}
	clusterItems := extractJSONArrayField(t, listClustersResp.Body.Bytes(), "items")
	assertClustersPresent(t, clusterItems, clusterIDs)

	for _, clusterID := range clusterIDs {
		clusterID := clusterID
		t.Run("health-summary-"+clusterID, func(t *testing.T) {
			resp := performClusterAuthedRequest(t, app.Router, token, http.MethodGet, "/api/v1/clusters/"+clusterID+"/health-summary", "")
			if resp.Code != http.StatusOK {
				t.Fatalf("expected health summary status=200, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
			}
			assertHasClusterField(t, resp.Body.Bytes(), "clusterId")
			assertHasClusterField(t, resp.Body.Bytes(), "total")
		})

		t.Run("connectivity-"+clusterID, func(t *testing.T) {
			resp := performClusterAuthedRequest(t, app.Router, token, http.MethodPost, "/api/v1/clusters/"+clusterID+"/connectivity", "")
			if resp.Code != http.StatusOK {
				t.Fatalf("expected connectivity status=200, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
			}
			assertConnectivityPayload(t, resp.Body.Bytes())
		})

		t.Run("sync-"+clusterID, func(t *testing.T) {
			resp := performClusterAuthedRequest(t, app.Router, token, http.MethodPost, "/api/v1/clusters/"+clusterID+"/sync", "")
			if resp.Code == http.StatusNotFound || resp.Code == http.StatusNotImplemented {
				t.Fatalf("expected sync route implemented (not 404/501), got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
			}
			if resp.Code != http.StatusAccepted && resp.Code != http.StatusOK {
				t.Fatalf("expected sync status=202(or explicit success status), got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
			}
			assertSyncPayload(t, resp.Body.Bytes(), clusterID)
		})
	}

	targetClusterID := clusterIDs[0]
	beforeItems := listResourcesByCluster(t, app.Router, token, targetClusterID)
	syncedName := persistSyncedResource(t, app, targetClusterID, "integration")
	afterItems := waitForResourceCountAtLeast(t, app.Router, token, targetClusterID, len(beforeItems)+1)
	assertResourceInList(t, afterItems, syncedName)
}

func performClusterAuthedRequest(t *testing.T, h http.Handler, token, method, target, body string) *httptest.ResponseRecorder {
	t.Helper()

	req := httptest.NewRequest(method, target, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	return resp
}

func extractIDField(t *testing.T, body []byte, field string) string {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}
	val, ok := payload[field]
	if !ok && strings.ToLower(field) == "id" {
		val, ok = payload["ID"]
	}
	if !ok {
		t.Fatalf("response missing field %q", field)
	}
	switch v := val.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		if v <= 0 || v != math.Trunc(v) {
			t.Fatalf("field %q has invalid numeric value: %v", field, v)
		}
		return strconv.FormatInt(int64(v), 10)
	default:
		t.Fatalf("field %q has unsupported type %T", field, val)
		return ""
	}
}

func assertHasClusterField(t *testing.T, body []byte, field string) {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}
	if _, ok := payload[field]; !ok {
		t.Fatalf("expected field %q in response, got: %v", field, payload)
	}
}

func extractJSONArrayField(t *testing.T, body []byte, field string) []any {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}

	raw, ok := payload[field]
	if !ok {
		t.Fatalf("response missing field %q: %v", field, payload)
	}
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("field %q should be array, got %T", field, raw)
	}
	return items
}

func assertClustersPresent(t *testing.T, items []any, expectedClusterIDs []string) {
	t.Helper()

	if len(items) == 0 {
		t.Fatal("expected non-empty cluster list")
	}

	seen := make(map[string]struct{}, len(items))
	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		id := jsonScalarToString(obj["id"])
		if id == "" {
			id = jsonScalarToString(obj["ID"])
		}
		if id != "" {
			seen[id] = struct{}{}
		}
	}

	for _, expectedID := range expectedClusterIDs {
		if _, ok := seen[expectedID]; !ok {
			t.Fatalf("expected cluster id=%s in list, seen=%v", expectedID, seen)
		}
	}
}

func assertConnectivityPayload(t *testing.T, body []byte) {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}

	if _, ok := payload["success"].(bool); !ok {
		t.Fatalf("connectivity response field success should be bool, payload=%v", payload)
	}

	message, ok := payload["message"].(string)
	if !ok {
		t.Fatalf("connectivity response field message should be string, payload=%v", payload)
	}
	msg := strings.TrimSpace(message)
	if msg == "" {
		t.Fatalf("connectivity response message should not be empty, payload=%v", payload)
	}

	msgLower := strings.ToLower(msg)
	if strings.Contains(msgLower, "stub") {
		t.Fatalf("connectivity response message should not be stub text, message=%q", msg)
	}
	if !strings.Contains(msgLower, "connectivity") &&
		!strings.Contains(msgLower, "kubeconfig") &&
		!strings.Contains(msgLower, "credential") &&
		!strings.Contains(msgLower, "api server") {
		t.Fatalf("connectivity response message should be diagnostic, message=%q", msg)
	}
}

func assertSyncPayload(t *testing.T, body []byte, expectedClusterID string) {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("sync response is not valid JSON object: %v", err)
	}

	clusterID := jsonScalarToString(payload["clusterId"])
	if clusterID == "" {
		clusterID = jsonScalarToString(payload["clusterID"])
	}
	if clusterID == "" {
		t.Fatalf("sync response field clusterId should be non-empty, payload=%v", payload)
	}
	if clusterID != expectedClusterID {
		t.Fatalf("sync response clusterId mismatch: expected=%s got=%s payload=%v", expectedClusterID, clusterID, payload)
	}

	status := strings.ToLower(strings.TrimSpace(jsonScalarToString(payload["status"])))
	if status == "" {
		t.Fatalf("sync response field status should be string, payload=%v", payload)
	}
	if status != "accepted" && status != "queued" && status != "success" {
		t.Fatalf("sync response status should be explicit success state, got=%q payload=%v", status, payload)
	}

	message, ok := payload["message"].(string)
	if !ok {
		t.Fatalf("sync response field message should be string, payload=%v", payload)
	}
	msg := strings.TrimSpace(message)
	if msg == "" {
		t.Fatalf("sync response message should not be empty, payload=%v", payload)
	}
	msgLower := strings.ToLower(msg)
	if strings.Contains(msgLower, "stub") {
		t.Fatalf("sync response message should not be stub text, message=%q", msg)
	}
	if !strings.Contains(msgLower, "sync") && !strings.Contains(msgLower, "enqueue") {
		t.Fatalf("sync response message should describe sync dispatch, message=%q", msg)
	}
}

func listResourcesByCluster(t *testing.T, h http.Handler, token, clusterID string) []any {
	t.Helper()

	resp := performClusterAuthedRequest(t, h, token, http.MethodGet, "/api/v1/resources?clusterId="+clusterID+"&limit=200", "")
	if resp.Code != http.StatusOK {
		t.Fatalf("expected resource list status=200, got status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return extractJSONArrayField(t, resp.Body.Bytes(), "items")
}

func persistSyncedResource(t *testing.T, app *testutil.App, clusterID, prefix string) string {
	t.Helper()

	id, err := strconv.ParseUint(clusterID, 10, 64)
	if err != nil {
		t.Fatalf("invalid cluster id for resource seed: %v", err)
	}

	name := prefix + "-synced-resource-" + clusterID
	record := repository.ResourceInventory{
		ClusterID: id,
		Namespace: "default",
		Kind:      "Pod",
		Name:      name,
		Health:    "healthy",
	}
	if err := app.DB.Create(&record).Error; err != nil {
		t.Fatalf("seed resource inventory failed: %v", err)
	}
	return name
}

func waitForResourceCountAtLeast(t *testing.T, h http.Handler, token, clusterID string, expected int) []any {
	t.Helper()

	deadline := time.Now().Add(2 * time.Second)
	var last []any
	for time.Now().Before(deadline) {
		last = listResourcesByCluster(t, h, token, clusterID)
		if len(last) >= expected {
			return last
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatalf("resource list did not grow to expected count=%d, got=%d", expected, len(last))
	return nil
}

func assertResourceInList(t *testing.T, items []any, expectedName string) {
	t.Helper()

	for _, item := range items {
		obj, ok := item.(map[string]any)
		if !ok {
			continue
		}
		name := strings.TrimSpace(jsonScalarToString(obj["name"]))
		if name == "" {
			name = strings.TrimSpace(jsonScalarToString(obj["Name"]))
		}
		if name == expectedName {
			return
		}
	}
	t.Fatalf("expected resource %q in list items=%v", expectedName, items)
}

func jsonScalarToString(v any) string {
	switch value := v.(type) {
	case string:
		return strings.TrimSpace(value)
	case float64:
		if value == math.Trunc(value) {
			return strconv.FormatInt(int64(value), 10)
		}
		return strconv.FormatFloat(value, 'f', -1, 64)
	default:
		return ""
	}
}
