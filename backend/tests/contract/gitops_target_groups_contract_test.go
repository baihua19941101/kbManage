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

func TestGitOpsTargetGroupsContract_CreateAndList(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "gitops-target-groups-contract",
		Password: "GitOpsTargets@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "gitops-target-groups-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	body := fmt.Sprintf(`{
		"name":"orders-targets",
		"workspaceId":%d,
		"projectId":%d,
		"clusterRefs":[%d],
		"selectorSummary":"env in (test,prod)",
		"description":"orders delivery target group"
	}`, access.WorkspaceID, access.ProjectID, access.ClusterID)

	createReq := httptest.NewRequest(http.MethodPost, "/api/v1/gitops/target-groups", strings.NewReader(body))
	createReq.Header.Set("Authorization", "Bearer "+token)
	createReq.Header.Set("Content-Type", "application/json")
	createResp := httptest.NewRecorder()
	app.Router.ServeHTTP(createResp, createReq)

	if createResp.Code != http.StatusCreated {
		t.Fatalf("expected create target-group status=201, got=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	created := mustDecodeGitOpsTargetGroupsContractObject(t, createResp.Body.Bytes())
	targetGroupID := mustReadGitOpsTargetGroupsContractID(t, created, "id")
	if strings.TrimSpace(mustReadGitOpsTargetGroupsContractString(t, created, "name")) != "orders-targets" {
		t.Fatalf("expected target-group name persisted, payload=%v", created)
	}
	if mustReadGitOpsTargetGroupsContractID(t, created, "workspaceId") != access.WorkspaceID {
		t.Fatalf("expected workspaceId=%d payload=%v", access.WorkspaceID, created)
	}
	if strings.TrimSpace(mustReadGitOpsTargetGroupsContractString(t, created, "status")) != "active" {
		t.Fatalf("expected initial target-group status=active payload=%v", created)
	}
	clusterRefs := mustReadGitOpsTargetGroupsContractArray(t, created, "clusterRefs")
	if len(clusterRefs) != 1 || uint64(clusterRefs[0].(float64)) != access.ClusterID {
		t.Fatalf("expected clusterRefs contains cluster=%d payload=%v", access.ClusterID, created)
	}

	listReq := httptest.NewRequest(
		http.MethodGet,
		"/api/v1/gitops/target-groups?workspaceId="+strconv.FormatUint(access.WorkspaceID, 10)+
			"&projectId="+strconv.FormatUint(access.ProjectID, 10),
		nil,
	)
	listReq.Header.Set("Authorization", "Bearer "+token)
	listResp := httptest.NewRecorder()
	app.Router.ServeHTTP(listResp, listReq)

	if listResp.Code != http.StatusOK {
		t.Fatalf("expected list target-groups status=200, got=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}
	listPayload := mustDecodeGitOpsTargetGroupsContractObject(t, listResp.Body.Bytes())
	items := mustReadGitOpsTargetGroupsContractArray(t, listPayload, "items")
	if len(items) == 0 {
		t.Fatalf("expected non-empty target-group items payload=%v", listPayload)
	}

	found := false
	for _, raw := range items {
		item, _ := raw.(map[string]any)
		if item == nil {
			continue
		}
		if mustReadGitOpsTargetGroupsContractID(t, item, "id") == targetGroupID {
			found = true
			if strings.TrimSpace(mustReadGitOpsTargetGroupsContractString(t, item, "name")) != "orders-targets" {
				t.Fatalf("expected listed target-group name=orders-targets item=%v", item)
			}
		}
	}
	if !found {
		t.Fatalf("expected list includes created target-group id=%d payload=%v", targetGroupID, listPayload)
	}
}

func mustDecodeGitOpsTargetGroupsContractObject(t *testing.T, body []byte) map[string]any {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON object: %v body=%s", err, strings.TrimSpace(string(body)))
	}
	return payload
}

func mustReadGitOpsTargetGroupsContractString(t *testing.T, payload map[string]any, key string) string {
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

func mustReadGitOpsTargetGroupsContractID(t *testing.T, payload map[string]any, key string) uint64 {
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

func mustReadGitOpsTargetGroupsContractArray(t *testing.T, payload map[string]any, key string) []any {
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
