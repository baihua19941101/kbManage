package integration_test

import (
	"net/http"
	"strconv"
	"testing"
)

func TestEnterprisePolishIntegration_ExportAuditFlow(t *testing.T) {
	ctx := newEnterpriseIntegrationCtx(t, "workspace-owner")
	seedEnterpriseIntegrationData(t, ctx)
	reportResp := performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/enterprise/reports", `{"workspaceId":`+strconv.FormatUint(ctx.Access.WorkspaceID, 10)+`,"reportType":"audit","title":"审计报告","audienceType":"auditor","timeRange":"30d"}`)
	if reportResp.Code != http.StatusCreated {
		t.Fatalf("create report status=%d body=%s", reportResp.Code, reportResp.Body.String())
	}
	record := mustDecodeEnterpriseIntegrationObject(t, reportResp.Body.Bytes())
	reportID := strconv.FormatUint(uint64(record["id"].(float64)), 10)
	resp := performEnterpriseIntegrationRequest(t, ctx.App.Router, ctx.Token, http.MethodPost, "/api/v1/enterprise/reports/"+reportID+"/exports", `{"audienceScope":"customer","contentLevel":"summary","exportType":"report"}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("export status=%d body=%s", resp.Code, resp.Body.String())
	}
}
