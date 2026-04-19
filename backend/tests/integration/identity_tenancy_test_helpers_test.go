package integration_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"kbmanage/backend/internal/api/middleware"
	"kbmanage/backend/internal/api/router"
	authSvc "kbmanage/backend/internal/service/auth"
	"kbmanage/backend/tests/testutil"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type identityTenancyIntegrationCtx struct {
	Router *gin.Engine
	DB     *gorm.DB
	Token  string
	Access testutil.ObservabilityAccessSeed
	UserID uint64
}

func newIdentityTenancyIntegrationCtx(t *testing.T, roleKey string) *identityTenancyIntegrationCtx {
	t.Helper()
	app := testutil.NewApp(t)
	user := testutil.SeedUser(t, app.DB, testutil.SeedUserInput{
		Username: "identity-integration-" + strings.ReplaceAll(roleKey, "_", "-"),
		Password: "Integration@123",
	})
	access := testutil.SeedObservabilityAccess(t, app.DB, user.User.ID, "identity-integration", roleKey)
	token := testutil.IssueAccessToken(t, app.Config, user.User.ID)
	tokenSvc := authSvc.NewTokenService(app.Config.JWTSecret, app.Config.AccessTokenTTL, app.Config.RefreshTokenTTL)

	gin.SetMode(gin.TestMode)
	engine := gin.New()
	engine.Use(gin.Recovery())
	v1 := engine.Group("/api/v1")
	authed := v1.Group("/")
	authed.Use(middleware.AuthRequired(tokenSvc))
	router.RegisterIdentityTenancyRoutes(authed, app.DB, nil)

	return &identityTenancyIntegrationCtx{
		Router: engine,
		DB:     app.DB,
		Token:  token,
		Access: access,
		UserID: user.User.ID,
	}
}

func performIdentityTenancyIntegrationRequest(t *testing.T, h http.Handler, token, method, path, body string) *httptest.ResponseRecorder {
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
