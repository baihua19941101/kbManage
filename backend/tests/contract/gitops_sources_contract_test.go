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

func TestGitOpsSourcesContract_CreateAndList(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-sources-contract",
		Password: "GitOpsSources@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-sources-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/sources", strings.NewReader(`{
		"name":"orders-git",
		"sourceType":"git",
		"endpoint":"https://git.example.test/demo/orders.git",
		"defaultRef":"main",
		"credentialRef":"secret://git/orders",
		"workspaceId":`+strconv.FormatUint(access.WorkspaceID, 10)+`,
		"projectId":`+strconv.FormatUint(access.ProjectID, 10)+`
	}`))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	created := mustDecodeGitOpsSourcesContractObject(t, createResp.Body.Bytes())
	sourceID := mustReadGitOpsSourcesContractID(t, created, "id")
	if strings.TrimSpace(mustReadGitOpsSourcesContractString(t, created, "name")) != "orders-git" {
		t.Fatalf("expected name=orders-git, payload=%v", created)
	}
	if strings.TrimSpace(mustReadGitOpsSourcesContractString(t, created, "sourceType")) != "git" {
		t.Fatalf("expected sourceType=git, payload=%v", created)
	}
	if strings.TrimSpace(mustReadGitOpsSourcesContractString(t, created, "endpoint")) != "https://git.example.test/demo/orders.git" {
		t.Fatalf("expected endpoint persisted, payload=%v", created)
	}
	if strings.TrimSpace(mustReadGitOpsSourcesContractString(t, created, "status")) != "pending" {
		t.Fatalf("expected initial status=pending, payload=%v", created)
	}

	listReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/gitops/sources?workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+
			"&projectId="+strconv.FormatUint(access.ProjectID, 10)+
			"&sourceType=git&status=pending",
		nil,
	)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)

	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	listPayload := mustDecodeGitOpsSourcesContractObject(t, listResp.Body.Bytes())
	items := mustReadGitOpsSourcesContractArray(t, listPayload, "items")
	if len(items) == 0 {
		t.Fatalf("expected non-empty source items, payload=%v", listPayload)
	}

	found := false
	for _, raw := range items {
		item, _ := raw.(map[string]any)
		if item == nil {
			continue
		}
		if mustReadGitOpsSourcesContractID(t, item, "id") == sourceID {
			found = true
			if strings.TrimSpace(mustReadGitOpsSourcesContractString(t, item, "name")) != "orders-git" {
				t.Fatalf("expected listed item name=orders-git, item=%v", item)
			}
			if strings.TrimSpace(mustReadGitOpsSourcesContractString(t, item, "status")) != "pending" {
				t.Fatalf("expected listed item status=pending, item=%v", item)
			}
		}
	}
	if !found {
		t.Fatalf("expected list contains created source id=%d payload=%v", sourceID, listPayload)
	}
}

func mustDecodeGitOpsSourcesContractObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadGitOpsSourcesContractString(t *testing.T, payload map[string]any, key string) string {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q in payload=%v", key, payload)
	}
	val, ok := raw.(string)
	if !ok {
		t.Fatalf("field %q must be string, got=%T value=%v", key, raw, raw)
	}
	return val
}

func mustReadGitOpsSourcesContractID(t *testing.T, payload map[string]any, key string) uint64 {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q in payload=%v", key, payload)
	}
	number, ok := raw.(float64)
	if !ok || number <= 0 {
		t.Fatalf("field %q must be positive number, got=%T value=%v", key, raw, raw)
	}
	return uint64(number)
}

func mustReadGitOpsSourcesContractArray(t *testing.T, payload map[string]any, key string) []any {
	t.Helper()
	raw, ok := payload[key]
	if !ok {
		t.Fatalf("missing field %q in payload=%v", key, payload)
	}
	items, ok := raw.([]any)
	if !ok {
		t.Fatalf("field %q must be array, got=%T value=%v", key, raw, raw)
	}
	return items
}
