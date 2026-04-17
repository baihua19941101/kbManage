package contract_test

import (
	"net/http"
	"strconv"
	"strings"
	"testing"
)

func TestComplianceArchiveExportsContract_CreateListGet(t *testing.T) {
	t.Parallel()
	env := newComplianceContractEnv(t, "compliance-exports-contract")
	createResp := complianceRequest(t, env, http.MethodPost, "/api/v1/compliance/archive-exports", `{"workspaceId":1,"projectId":1,"exportScope":"bundle"}`)
	if createResp.Code != http.StatusCreated {
		t.Fatalf("create export failed status=%d body=%s", createResp.Code, strings.TrimSpace(createResp.Body.String()))
	}
	payload := complianceDecodeObject(t, createResp.Body.Bytes())
	exportID := payload["id"].(string)

	listResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/archive-exports?status=pending", "")
	if listResp.Code != http.StatusOK {
		t.Fatalf("list exports failed status=%d body=%s", listResp.Code, strings.TrimSpace(listResp.Body.String()))
	}

	getResp := complianceRequest(t, env, http.MethodGet, "/api/v1/compliance/archive-exports/"+exportID, "")
	if getResp.Code != http.StatusOK {
		t.Fatalf("get export failed status=%d body=%s", getResp.Code, strings.TrimSpace(getResp.Body.String()))
	}
	_ = strconv.Itoa(1)
}
