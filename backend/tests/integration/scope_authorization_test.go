package integration_test

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/tests/testutil"
)

func TestScopeAuthorization_WorkspaceIsolationFlow(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	admin := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "scope-admin",
		Password: "Scope@123456",
		Email:    "scope-admin@example.test",
	})
	member := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "scope-member",
		Password: "Scope@123456",
		Email:    "scope-member@example.test",
	})
	adminToken := testutil.IssueAccessToken(t, app.Config, admin.User.ID)
	memberToken := testutil.IssueAccessToken(t, app.Config, member.User.ID)

	createWorkspaceResp := performScopeAuthedRequest(t, app.Router, adminToken, http.MethodPost, "/api/v1/workspaces", `{
		"name":"isolation-workspace",
		"description":"workspace isolation integration test"
	}`)
	if createWorkspaceResp.Code != http.StatusCreated {
		t.Fatalf("expected create workspace status=201, got status=%d body=%s", createWorkspaceResp.Code, strings.TrimSpace(createWorkspaceResp.Body.String()))
	}
	workspaceID := extractScopeID(createWorkspaceResp.Body.Bytes())
	if workspaceID == "" {
		t.Fatalf("workspace creation response missing id, body=%s", strings.TrimSpace(createWorkspaceResp.Body.String()))
	}

	createProjectResp := performScopeAuthedRequest(t, app.Router, adminToken, http.MethodPost, "/api/v1/workspaces/"+workspaceID+"/projects", `{
		"name":"billing-api",
		"description":"project for scope integration"
	}`)
	if createProjectResp.Code != http.StatusCreated {
		t.Fatalf("expected create project status=201, got status=%d body=%s", createProjectResp.Code, strings.TrimSpace(createProjectResp.Body.String()))
	}

	memberListResp := performScopeAuthedRequest(t, app.Router, memberToken, http.MethodGet, "/api/v1/workspaces/"+workspaceID+"/projects", "")
	if memberListResp.Code != http.StatusForbidden {
		t.Fatalf("expected member project list status=403 without binding, got status=%d body=%s", memberListResp.Code, strings.TrimSpace(memberListResp.Body.String()))
	}

	createBindingResp := performScopeAuthedRequest(t, app.Router, adminToken, http.MethodPost, "/api/v1/role-bindings", `{
		"subjectType":"user",
		"subjectId":`+strconv.FormatUint(member.User.ID, 10)+`,
		"scopeType":"workspace",
		"scopeId":`+workspaceID+`,
		"roleKey":"workspace-viewer"
	}`)
	if createBindingResp.Code != http.StatusCreated {
		t.Fatalf("expected role binding creation status=201, got status=%d body=%s", createBindingResp.Code, strings.TrimSpace(createBindingResp.Body.String()))
	}

	memberListResp = performScopeAuthedRequest(t, app.Router, memberToken, http.MethodGet, "/api/v1/workspaces/"+workspaceID+"/projects", "")
	if memberListResp.Code != http.StatusOK {
		t.Fatalf("expected member project list status=200 after binding, got status=%d body=%s", memberListResp.Code, strings.TrimSpace(memberListResp.Body.String()))
	}
	assertScopeHasItems(t, memberListResp.Body.Bytes())
}

func performScopeAuthedRequest(t *testing.T, h http.Handler, token, method, target, body string) *httptest.ResponseRecorder {
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

func extractScopeID(body []byte) string {
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		return ""
	}
	val, ok := payload["id"]
	if !ok {
		val, ok = payload["ID"]
	}
	if !ok {
		return ""
	}
	switch v := val.(type) {
	case string:
		return strings.TrimSpace(v)
	case float64:
		if v <= 0 || v != math.Trunc(v) {
			return ""
		}
		return strconv.FormatInt(int64(v), 10)
	default:
		return ""
	}
}

func assertScopeHasItems(t *testing.T, body []byte) {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v", err)
	}
	if _, ok := payload["items"]; !ok {
		t.Fatalf("expected response to contain items field, got: %v", payload)
	}
}
