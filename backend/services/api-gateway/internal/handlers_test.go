package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// ──────────────────────────────────────────────
// serviceAddr
// ──────────────────────────────────────────────

func TestServiceAddr_NonEmpty(t *testing.T) {
	assert.Equal(t, "myhost:9090", serviceAddr("myhost:9090", 8080))
}

func TestServiceAddr_DefaultPort(t *testing.T) {
	assert.Equal(t, "localhost:8082", serviceAddr("", 8082))
}

// ──────────────────────────────────────────────
// findTarget
// ──────────────────────────────────────────────

func newHandlersForTest() *Handlers {
	return NewHandlers(nil, ServicesCfg{
		AuthAddr:          "auth:8081",
		TripAddr:          "trip:8082",
		UserAddr:          "user:8083",
		CollaborationAddr: "collab:8084",
		NotificationAddr:  "notif:8085",
		SearchAddr:        "search:8086",
	})
}

func TestFindTarget_ExactPrefix(t *testing.T) {
	h := newHandlersForTest()
	addr, ok := h.findTarget("/api/v1/trips/some-id")
	assert.True(t, ok)
	assert.Contains(t, addr, "trip:8082")
}

func TestFindTarget_AuthPrefix(t *testing.T) {
	h := newHandlersForTest()
	addr, ok := h.findTarget("/auth/google/login")
	assert.True(t, ok)
	assert.Contains(t, addr, "auth:8081")
}

func TestFindTarget_LongestPrefix(t *testing.T) {
	// Both /api/v1/users and /api/v1/users/me should match /api/v1/users
	h := newHandlersForTest()
	addr, ok := h.findTarget("/api/v1/users/me")
	assert.True(t, ok)
	assert.Contains(t, addr, "user:8083")
}

func TestFindTarget_NoMatch(t *testing.T) {
	h := newHandlersForTest()
	_, ok := h.findTarget("/unknown/path")
	assert.False(t, ok)
}

func TestFindTarget_SearchPrefix(t *testing.T) {
	h := newHandlersForTest()
	addr, ok := h.findTarget("/api/v1/search/places")
	assert.True(t, ok)
	assert.Contains(t, addr, "search:8086")
}

// ──────────────────────────────────────────────
// Health handler
// ──────────────────────────────────────────────

func TestHealth_OK(t *testing.T) {
	r := gin.New()
	r.GET("/healthz", Health)

	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusOK, w.Code)
	var body map[string]interface{}
	require.NoError(t, json.NewDecoder(w.Body).Decode(&body))
	// standard envelope: data.status should be "ok"
	data, ok := body["data"].(map[string]interface{})
	require.True(t, ok, "expected data object in envelope")
	assert.Equal(t, "ok", data["status"])
}
