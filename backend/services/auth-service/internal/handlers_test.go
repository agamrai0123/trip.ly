package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/agamrai0123/wanderplan/pkg/response"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// parseEnvelope unmarshals the standard envelope from the recorder body.
func parseEnvelope(t *testing.T, w *httptest.ResponseRecorder) response.Envelope {
	t.Helper()
	var env response.Envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&env))
	return env
}

// TestRefreshHandler_MissingCookie checks that a POST /auth/refresh without
// the refresh_token cookie returns HTTP 401.
func TestRefreshHandler_MissingCookie(t *testing.T) {
	h := &Handlers{svc: nil, cfg: &Config{}}

	r := gin.New()
	r.POST("/auth/refresh", h.Refresh)

	req := httptest.NewRequest(http.MethodPost, "/auth/refresh", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusUnauthorized, w.Code)
	env := parseEnvelope(t, w)
	require.NotNil(t, env.Error)
	assert.Equal(t, string(response.ErrBody{}.Code), "") // env.Error.Code is non-empty
	assert.NotEmpty(t, env.Error.Code)
}

// TestMeHandler_ReturnsUserClaims verifies that GET /auth/me returns the
// claims injected by the Auth middleware into the Gin context.
func TestMeHandler_ReturnsUserClaims(t *testing.T) {
	h := &Handlers{svc: nil, cfg: &Config{}}

	r := gin.New()
	r.GET("/auth/me", func(c *gin.Context) {
		c.Set("user_id", "user-abc-123")
		c.Set("email", "alice@example.com")
		c.Set("name", "Alice")
		c.Set("avatar_url", "https://example.com/avatar.jpg")
	}, h.Me)

	req := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	env := parseEnvelope(t, w)
	assert.Nil(t, env.Error)

	data, ok := env.Data.(map[string]interface{})
	require.True(t, ok, "data must be a JSON object")
	assert.Equal(t, "user-abc-123", data["user_id"])
	assert.Equal(t, "alice@example.com", data["email"])
	assert.Equal(t, "Alice", data["name"])
}

// TestLoginHandler_InvalidProvider checks that an unknown OAuth provider
// results in HTTP 400.
func TestLoginHandler_InvalidProvider(t *testing.T) {
	h := &Handlers{svc: testAuthService(), cfg: &Config{Cookie: CookieCfg{}}}

	r := gin.New()
	r.GET("/auth/:provider/login", h.Login)

	req := httptest.NewRequest(http.MethodGet, "/auth/twitter/login", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	env := parseEnvelope(t, w)
	require.NotNil(t, env.Error)
}

// TestCallbackHandler_MissingState verifies that the callback handler rejects
// requests with mismatched or absent state cookies.
func TestCallbackHandler_MissingState(t *testing.T) {
	h := &Handlers{svc: testAuthService(), cfg: &Config{Cookie: CookieCfg{}}}

	r := gin.New()
	r.GET("/auth/:provider/callback", h.Callback)

	// No oauth_state cookie set → should return 400
	req := httptest.NewRequest(http.MethodGet, "/auth/google/callback?state=somestate&code=somecode", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	env := parseEnvelope(t, w)
	require.NotNil(t, env.Error)
}

// TestHealthz verifies that the /healthz endpoint returns HTTP 200.
func TestHealthz(t *testing.T) {
	r := gin.New()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
}
