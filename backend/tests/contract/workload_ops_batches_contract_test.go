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

func TestWorkloadOpsBatchesContract_SubmitAndQuery(t *testing.T) {
	t.Parallel()

	app := testutil.NewApp(t)
	seeded := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "wops-batches-contract",
		Password: "WorkloadOps@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, seeded.User.ID, "wops-batches-contract", "workspace-owner")
	token := testutil.IssueAccessToken(t, app.Config, seeded.User.ID)

	submitReq := httptest.NewRequest(http.MethodPost, "/api/v1/workload-ops/batches", strings.NewReader(`{
		"actionType": "restart",
		"riskConfirmed": true,
		"targets": [
			{"clusterId":`+strconv.FormatUint(access.ClusterID, 10)+`,"namespace":"default","resourceKind":"Deployment","resourceName":"demo-api-1"},
			{"clusterId":`+strconv.FormatUint(access.ClusterID, 10)+`,"namespace":"default","resourceKind":"Deployment","resourceName":"demo-api-2"}
		]
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
	if strings.TrimSpace(fmt.Sprint(created["status"])) == "" {
		t.Fatalf("expected status in batch payload=%v", created)
	}
	id := int(created["id"].(float64))

	getReq := httptest.NewRequest(http.MethodGet, "/api/v1/workload-ops/batches/"+strconv.Itoa(id), nil)
	getReq.Header.Set("Authorization", "Bearer "+token)
	getResp := httptest.NewRecorder()
	app.Router.ServeHTTP(getResp, getReq)
	if getResp.Code != http.StatusOK {
		t.Fatalf("expected get status=200 got=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
	var fetched map[string]any
	_ = json.Unmarshal(getResp.Body.Bytes(), &fetched)
	items, ok := fetched["items"].([]any)
	if !ok || len(items) != 2 {
		t.Fatalf("expected 2 batch items payload=%v", fetched)
	}
}
