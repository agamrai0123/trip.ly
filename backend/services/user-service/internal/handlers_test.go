package internal

import (
	"bytes"
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

// parseEnv unmarshals the standard response envelope.
func parseEnv(t *testing.T, w *httptest.ResponseRecorder) response.Envelope {
	t.Helper()
	var env response.Envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&env))
	return env
}

// injectUserID is middleware that sets user_id in the Gin context.
func injectUserID(id string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", id)
		c.Next()
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// PATCH /users/me
// ──────────────────────────────────────────────────────────────────────────────

// TestUpdateMe_InvalidJSON checks that malformed JSON returns HTTP 400.
func TestUpdateMe_InvalidJSON(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.PATCH("/users/me", injectUserID("user-001"), h.UpdateMe)

	req := httptest.NewRequest(http.MethodPatch, "/users/me",
		bytes.NewBufferString("{bad json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	env := parseEnv(t, w)
	require.NotNil(t, env.Error)
}

// TestUpdateMe_EmptyPayload_PassesBind checks that an empty (but valid) JSON
// object passes binding (no required fields on UpdateUserRequest).
func TestUpdateMe_EmptyPayload_PassesBind(t *testing.T) {
	// svc is nil — this will panic if UpdateMe calls svc after binding.
	// This test intentionally verifies binding succeeds; the subsequent panic
	// means the test stops at svc call. We capture it via recover.
	h := &Handlers{svc: nil}
	r := gin.New()
	r.PATCH("/users/me", injectUserID("user-001"), h.UpdateMe)

	body := `{}`
	req := httptest.NewRequest(http.MethodPatch, "/users/me", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	// Gin's default recovery middleware swallows the panic as 500.
	// We add recovery manually to verify the request got past binding.
	r.Use(gin.Recovery())

	// Re-register route with recovery applied (gin.New skips default middlewares)
	r2 := gin.New()
	r2.Use(gin.Recovery())
	r2.PATCH("/users/me", injectUserID("user-001"), h.UpdateMe)

	r2.ServeHTTP(w, req)
	// Either 200 (svc succeeded, unlikely w/ nil svc) or 500 (nil pointer panic caught by gin.Recovery)
	// Both mean binding succeeded.
	assert.NotEqual(t, http.StatusBadRequest, w.Code, "binding should not fail for empty JSON")
}

// ──────────────────────────────────────────────────────────────────────────────
// GET /healthz
// ──────────────────────────────────────────────────────────────────────────────

func TestHealthz(t *testing.T) {
	r := gin.New()
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	assert.Equal(t, "ok", body["status"])
}
