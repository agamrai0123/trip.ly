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
// GET /search/places?q=
// ──────────────────────────────────────────────

// TestSearchPlaces_MissingQuery verifies that omitting ?q= returns HTTP 400.
func TestSearchPlaces_MissingQuery(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.GET("/search/places", h.SearchPlaces)

	req := httptest.NewRequest(http.MethodGet, "/search/places", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSearchPlaces_EmptyQuery verifies that ?q= (empty string) returns HTTP 400.
func TestSearchPlaces_EmptyQuery(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.GET("/search/places", h.SearchPlaces)

	req := httptest.NewRequest(http.MethodGet, "/search/places?q=", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ──────────────────────────────────────────────
// GET /search/trips?q=
// ──────────────────────────────────────────────

// TestSearchTrips_MissingQuery verifies that omitting ?q= returns HTTP 400.
func TestSearchTrips_MissingQuery(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.GET("/search/trips", h.SearchTrips)

	req := httptest.NewRequest(http.MethodGet, "/search/trips", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestSearchTrips_EmptyQuery verifies that ?q= (empty string) returns HTTP 400.
func TestSearchTrips_EmptyQuery(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.GET("/search/trips", h.SearchTrips)

	req := httptest.NewRequest(http.MethodGet, "/search/trips?q=", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ──────────────────────────────────────────────
// Health
// ──────────────────────────────────────────────

func TestSearchHealthz(t *testing.T) {
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
