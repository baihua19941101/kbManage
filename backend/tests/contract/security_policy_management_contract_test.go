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

func TestSecurityPolicyContract_CreateListGetUpdate(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "security-policy-contract",
		Password: "SecurityPolicy@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "security-policy-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	createBody := fmt.Sprintf(`{
		"name":"restrict-latest-tag",
		"workspaceId":%d,
		"projectId":%d,
		"scopeLevel":"project",
		"category":"image",
		"ruleTemplate":{"requireDigest":true},
		"defaultEnforcementMode":"warn",
		"riskLevel":"high"
	}`, access.WorkspaceID, access.ProjectID)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/security-policies", strings.NewReader(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	created := mustDecodeSecurityPolicyContractObject(t, createResp.Body.Bytes())
	policyID := mustReadSecurityPolicyContractID(t, created, "id")
	if strings.TrimSpace(mustReadSecurityPolicyContractString(t, created, "name")) != "restrict-latest-tag" {
		t.Fatalf("unexpected created policy payload=%v", created)
	}

	listReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/security-policies?workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+
			"&projectId="+strconv.FormatUint(access.ProjectID, 10)+
			"&scopeLevel=project&category=image",
		nil,
	)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)
	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	listPayload := mustDecodeSecurityPolicyContractObject(t, listResp.Body.Bytes())
	items := mustReadSecurityPolicyContractArray(t, listPayload, "items")
	if len(items) == 0 {
		t.Fatalf("expected non-empty policy list payload=%v", listPayload)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/security-policies/"+strconv.FormatUint(policyID, 10), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResp := httptest.NewRecorder()
	app.Router.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get status=200, got=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}

	updateReq := httptest.NewRequest(http.MethodPatch, "/api/v1/security-policies/"+strconv.FormatUint(policyID, 10), strings.NewReader(`{
		"status":"active",
		"defaultEnforcementMode":"enforce"
	}`))
	updateReq.Header.Set("Authorization", "Bearer "+token)
	updateReq.Header.Set("Content-Type", "application/json")
	updateResp := httptest.NewRecorder()
	app.Router.ServeHTTP(updateResp, updateReq)
	if updateResp.Code != http.StatusOK {
		t.Fatalf("expected update status=200, got=%d body=%s", updateResp.Code, strings.TrimSpace(updateResp.Body.String()))
	}
	updated := mustDecodeSecurityPolicyContractObject(t, updateResp.Body.Bytes())
	if strings.TrimSpace(mustReadSecurityPolicyContractString(t, updated, "status")) != "active" {
		t.Fatalf("expected status=active payload=%v", updated)
	}
}

func mustDecodeSecurityPolicyContractObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadSecurityPolicyContractString(t *testing.T, payload map[string]any, key string) string {
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

func mustReadSecurityPolicyContractID(t *testing.T, payload map[string]any, key string) uint64 {
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

func mustReadSecurityPolicyContractArray(t *testing.T, payload map[string]any, key string) []any {
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
