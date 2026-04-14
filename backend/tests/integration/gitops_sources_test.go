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

func TestGitOpsSourcesIntegration_CreateUpdateVerifyAndFilter(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-sources-integration",
		Password: "GitOpsSourcesIntegration@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-sources-integration", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	createBody := fmt.Sprintf(`{
		"name":"orders-source-int",
		"sourceType":"git",
		"endpoint":"https://git.example.test/demo/orders-int.git",
		"defaultRef":"main",
		"workspaceId":%d,
		"projectId":%d
	}`, access.WorkspaceID, access.ProjectID)
	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/sources", strings.NewReader(createBody))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create source status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	created := mustDecodeGitOpsSourcesIntegrationObject(t, createResp.Body.Bytes())
	sourceID := mustReadGitOpsSourcesIntegrationID(t, created, "id")

	patchDisableReq := httptest.NewRequest(http.MethodPatch, "/api/v1/gitops/sources/"+strconv.FormatUint(sourceID, 10), strings.NewReader(`{"disabled":true}`))
	patchDisableReq.Header.Set("Authorization", "Bearer "+token)
	patchDisableReq.Header.Set("Content-Type", "application/json")
	patchDisableResp := httptest.NewRecorder()
	app.Router.ServeHTTP(patchDisableResp, patchDisableReq)

	if patchDisableResp.Code != http.StatusOK {
		t.Fatalf("expected patch disabled status=200, got=%d body=%s", patchDisableResp.Code, strings.TrimSpace(patchDisableResp.Body.String()))
	}
	disabledSource := mustDecodeGitOpsSourcesIntegrationObject(t, patchDisableResp.Body.Bytes())
	if strings.TrimSpace(mustReadGitOpsSourcesIntegrationString(t, disabledSource, "status")) != "disabled" {
		t.Fatalf("expected source status=disabled payload=%v", disabledSource)
	}

	patchEnableReq := httptest.NewRequest(http.MethodPatch, "/api/v1/gitops/sources/"+strconv.FormatUint(sourceID, 10), strings.NewReader(`{"defaultRef":"release/v1","disabled":false}`))
	patchEnableReq.Header.Set("Authorization", "Bearer "+token)
	patchEnableReq.Header.Set("Content-Type", "application/json")
	patchEnableResp := httptest.NewRecorder()
	app.Router.ServeHTTP(patchEnableResp, patchEnableReq)

	if patchEnableResp.Code != http.StatusOK {
		t.Fatalf("expected patch enable status=200, got=%d body=%s", patchEnableResp.Code, strings.TrimSpace(patchEnableResp.Body.String()))
	}
	enabledSource := mustDecodeGitOpsSourcesIntegrationObject(t, patchEnableResp.Body.Bytes())
	if strings.TrimSpace(mustReadGitOpsSourcesIntegrationString(t, enabledSource, "defaultRef")) != "release/v1" {
		t.Fatalf("expected defaultRef updated payload=%v", enabledSource)
	}

	verifyReq := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/sources/"+strconv.FormatUint(sourceID, 10)+"/verify", nil)
	verifyReq.Header.Set("Authorization", "Bearer "+token)
	verifyResp := httptest.NewRecorder()
	app.Router.ServeHTTP(verifyResp, verifyReq)

	if verifyResp.Code != http.StatusAccepted {
		t.Fatalf("expected verify status=202, got=%d body=%s", verifyResp.Code, strings.TrimSpace(verifyResp.Body.String()))
	}
	verifyPayload := mustDecodeGitOpsSourcesIntegrationObject(t, verifyResp.Body.Bytes())
	if strings.TrimSpace(mustReadGitOpsSourcesIntegrationString(t, verifyPayload, "actionType")) != "verify-source" {
		t.Fatalf("expected verify actionType=verify-source payload=%v", verifyPayload)
	}
	if strings.TrimSpace(mustReadGitOpsSourcesIntegrationString(t, verifyPayload, "status")) != "pending" {
		t.Fatalf("expected verify operation status=pending payload=%v", verifyPayload)
	}

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/gitops/sources/"+strconv.FormatUint(sourceID, 10), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResp := httptest.NewRecorder()
	app.Router.ServeHTTP(getResp, getReq)

	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get source status=200, got=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
	getPayload := mustDecodeGitOpsSourcesIntegrationObject(t, getResp.Body.Bytes())
	if strings.TrimSpace(mustReadGitOpsSourcesIntegrationString(t, getPayload, "status")) != "ready" {
		t.Fatalf("expected verified source status=ready payload=%v", getPayload)
	}
	if strings.TrimSpace(mustReadGitOpsSourcesIntegrationString(t, getPayload, "lastVerifiedAt")) == "" {
		t.Fatalf("expected lastVerifiedAt not empty payload=%v", getPayload)
	}

	listReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/gitops/sources?workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+
			"&projectId="+strconv.FormatUint(access.ProjectID, 10)+
			"&sourceType=git&status=ready",
		nil,
	)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)

	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list source status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	listPayload := mustDecodeGitOpsSourcesIntegrationObject(t, listResp.Body.Bytes())
	items := mustReadGitOpsSourcesIntegrationArray(t, listPayload, "items")
	if len(items) == 0 {
		t.Fatalf("expected ready source list not empty payload=%v", listPayload)
	}
}

func mustDecodeGitOpsSourcesIntegrationObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadGitOpsSourcesIntegrationString(t *testing.T, payload map[string]any, key string) string {
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

func mustReadGitOpsSourcesIntegrationID(t *testing.T, payload map[string]any, key string) uint64 {
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

func mustReadGitOpsSourcesIntegrationArray(t *testing.T, payload map[string]any, key string) []any {
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
