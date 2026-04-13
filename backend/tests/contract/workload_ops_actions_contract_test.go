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

func TestWorkloadOpsActionsContract_SubmitAndQuery(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-actions-contract",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-actions-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	submitReq := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/actions", strings.NewReader(`{
		"clusterId": `+strconv.FormatUint(access.ClusterID, 10)+`,
		"namespace": "default",
		"resourceKind": "Deployment",
		"resourceName": "demo-api",
		"actionType": "restart",
		"riskConfirmed": true
	}`))
	submitReq.Header.Set("Authorization", "Bearer "+token)
	submitReq.Header.Set("Content-Type", "application/json")
	submitResp := httptest.NewRecorder()
	app.Router.ServeHTTP(submitResp, submitReq)
	if submitResp.Code != http.StatusAccepted {
		t.Fatalf("expected submit status=202 got=%d body=%s", submitResp.Code, strings.TrimSpace(submitResp.Body.String()))
	}
	var created map[string]any
	_ = json.Unmarshal(submitResp.Body.Bytes(), &created)
	if strings.TrimSpace(strValue(created["status"])) != "succeeded" {
		t.Fatalf("expected succeeded status, payload=%v", created)
	}
	id := int(created["id"].(float64))

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/actions/"+strconv.Itoa(id), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResp := httptest.NewRecorder()
	app.Router.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get status=200 got=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
	var fetched map[string]any
	_ = json.Unmarshal(getResp.Body.Bytes(), &fetched)
	if strings.TrimSpace(strValue(fetched["actionType"])) != "restart" {
		t.Fatalf("expected restart actionType payload=%v", fetched)
	}
}

func strValue(v any) string {
	s, _ := v.(string)
	return s
}
