package contract_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/api/router"
	authSvc "kbmanage/backend/internal/service/auth"
	"kbmanage/backend/tests/testutil"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type marketplaceContractCtx struct {
	Router *gin.Engine
	DB     *gorm.DB
	Token  string
	Access testutil.ObservabilityAccessSeed
	UserID uint64
}

func newMarketplaceContractCtx(t *testing.T, roleKey string) *marketplaceContractCtx {
	t.Helper()
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "marketplace-contract-" + strings.ReplaceAll(roleKey, "_", "-"),
		Password: "Contract@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "marketplace-contract", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	tokenSvc := authSvc.NewTokenService(app.Config.JWTSecret, app.Config.AccessTokenTTL, app.Config.RefreshTokenTTL)

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	v1 := engine.Group("/api/v1")
	authed := v1.Group("/")
	authed.Use(middleware.AuthRequired(tokenSvc))
	router.RegisterMarketplaceRoutes(authed, app.DB, nil)

	return &marketplaceContractCtx{
		Router: engine,
		DB:     app.DB,
		Token:  token,
		Access: access,
		UserID: user.User.ID,
	}
}

func performMarketplaceContractRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(method, path, strings.NewReader(body))
	req.Header.Set("Authorization", "Bearer "+token)
	if body != "" {
		req.Header.Set("Content-Type", "application/json")
	}
	resp := httptest.NewRecorder()
	h.ServeHTTP(resp, req)
	return resp
}

func mustReadMarketplaceContractID(t *testing.T, body []byte, key string) uint64 {
	t.Helper()
	var payload map[string]any
	if err := json.Unmarshal(body, &payload); err != nil {
		t.Fatalf("invalid JSON response: %v body=%s", err, strings.TrimSpace(string(body)))
	}
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

func createMarketplaceCatalogSourceContract(t *testing.T, ctx *marketplaceContractCtx, name string) uint64 {
	t.Helper()
	resp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources", `{
		"name":"`+name+`",
		"sourceType":"helm",
		"endpointRef":"https://catalog.example.test/`+name+`",
		"visibilityScope":"platform",
		"templateSeeds":[
			{
				"name":"nginx-stack",
				"slug":"nginx-stack",
				"category":"web",
				"summary":"标准 nginx 模板",
				"publishStatus":"active",
				"supportedScopes":["workspace","project"],
				"releaseNotesSummary":"首发",
				"versions":[
					{
						"version":"1.0.0",
						"status":"active",
						"dependencies":["ingress"],
						"parameterSchemaSummary":"replicas,image",
						"deploymentConstraintSummary":"workspace/project",
						"releaseNotes":"首个稳定版"
					},
					{
						"version":"1.1.0",
						"status":"active",
						"dependencies":["ingress","metrics"],
						"parameterSchemaSummary":"replicas,image,resources",
						"deploymentConstraintSummary":"workspace/project",
						"releaseNotes":"增加 metrics",
						"supersedesVersion":"1.0.0"
					}
				]
			}
		]
	}`)
	if resp.Code != http.StatusCreated {
		t.Fatalf("create catalog source failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
	return mustReadMarketplaceContractID(t, resp.Body.Bytes(), "id")
}

func syncMarketplaceCatalogSourceContract(t *testing.T, ctx *marketplaceContractCtx, sourceID uint64) {
	t.Helper()
	resp := performMarketplaceContractRequest(t, ctx.Router, ctx.Token, http.MethodPost, "/api/v1/marketplace/catalog-sources/"+strconv.FormatUint(sourceID, 10)+"/sync", "")
	if resp.Code != http.StatusAccepted {
		t.Fatalf("sync catalog source failed status=%d body=%s", resp.Code, strings.TrimSpace(resp.Body.String()))
	}
}
