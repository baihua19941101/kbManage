package integration_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kbmanage/backend/internal/api/router"
	"kbmanage/backend/internal/repository"
)

func TestClusterOverview_MultiClusterFlowSkeleton(t *testing.T) {
	t.Parallel()

	r := newTestRouter()

	listResp := httptest.NewRecorder()
	listReq := httptest.NewRequest(http.MethodGet, "/api/v1/clusters", nil)
	r.ServeHTTP(listResp, listReq)

	if listResp.Code == http.StatusNotFound {
		t.Skip("GET /api/v1/clusters not implemented yet")
	}

	assertOverviewAllowedStatus(t, listResp.Code)

	clusterIDs := extractClusterIDs(listResp.Body.Bytes())
	if len(clusterIDs) == 0 {
		clusterIDs = []string{"cluster-a", "cluster-b"}
	}
	if len(clusterIDs) > 2 {
		clusterIDs = clusterIDs[:2]
	}

	for _, clusterID := range clusterIDs {
		clusterID := clusterID
		t.Run(fmt.Sprintf("resources-%s", clusterID), func(t *testing.T) {
			t.Parallel()

			resourceResp := httptest.NewRecorder()
			resourceReq := httptest.NewRequest(http.MethodGet, "/api/v1/clusters/"+clusterID+"/resources", nil)
			r.ServeHTTP(resourceResp, resourceReq)

			if resourceResp.Code == http.StatusNotFound {
				t.Skipf("GET /api/v1/clusters/%s/resources not implemented yet", clusterID)
			}

			assertOverviewAllowedStatus(t, resourceResp.Code)
			if resourceResp.Code == http.StatusOK {
				assertHasItemsField(t, resourceResp.Body.Bytes())
			}
		})
	}
}

func newTestRouter() http.Handler {
	cfg := repository.Config{
		JWTSecret:       "test-secret",
		AccessTokenTTL:  15 * time.Minute,
		RefreshTokenTTL: 24 * time.Hour,
	}
	return router.NewRouter(nil, nil, cfg)
}

func assertOverviewAllowedStatus(t *testing.T, code int) {
	t.Helper()

	switch code {
	case http.StatusOK,
		http.StatusBadRequest,
		http.StatusUnauthorized,
		http.StatusForbidden,
		http.StatusMethodNotAllowed,
		http.StatusNotImplemented:
		return
	default:
		t.Fatalf("unexpected status code: %d", code)
	}
}

func extractClusterIDs(body []byte) []string {
	var payload struct {
		Items []struct {
			ID string `json:"id"`
		} `json:"items"`
	}
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil
	}

	ids := make([]string, 0, len(payload.Items))
	for _, item := range payload.Items {
		if item.ID != "" {
			ids = append(ids, item.ID)
		}
	}
	return ids
}

func assertHasItemsField(t *testing.T, body []byte) {
	t.Helper()

	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("response is not valid JSON: %v", err)
	}
	if _, ok := payload["items"]; !ok {
		t.Fatalf("expected response to contain items field, got: %v", payload)
	}
}
