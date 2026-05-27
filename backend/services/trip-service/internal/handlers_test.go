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

// parseEnv unmarshals the standard envelope.
func parseEnv(t *testing.T, w *httptest.ResponseRecorder) response.Envelope {
	t.Helper()
	var env response.Envelope
	require.NoError(t, json.NewDecoder(w.Body).Decode(&env))
	return env
}

// setUserID is a middleware that injects a user_id into the Gin context.
func setUserID(id string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", id)
		c.Next()
	}
}

// ──────────────────────────────────────────────
// CreateTrip handler
// ──────────────────────────────────────────────

// TestCreateTrip_MissingTitle checks that binding validation rejects
// a request without the required "title" field (returns 400).
func TestCreateTrip_MissingTitle(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.POST("/trips", setUserID("user-123"), h.CreateTrip)

	body := `{"destination":"Paris"}`
	req := httptest.NewRequest(http.MethodPost, "/trips", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	env := parseEnv(t, w)
	require.NotNil(t, env.Error)
}

// TestCreateTrip_MissingDestination checks that binding validation rejects
// a request without the required "destination" field.
func TestCreateTrip_MissingDestination(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.POST("/trips", setUserID("user-123"), h.CreateTrip)

	body := `{"title":"My Trip"}`
	req := httptest.NewRequest(http.MethodPost, "/trips", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	env := parseEnv(t, w)
	require.NotNil(t, env.Error)
}

// TestCreateTrip_InvalidJSON checks that malformed JSON returns 400.
func TestCreateTrip_InvalidJSON(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.POST("/trips", setUserID("user-123"), h.CreateTrip)

	req := httptest.NewRequest(http.MethodPost, "/trips", bytes.NewBufferString("{invalid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ──────────────────────────────────────────────
// CreateDay handler
// ──────────────────────────────────────────────

// TestCreateDay_MissingDayNumber verifies that binding rejects missing day_number.
func TestCreateDay_MissingDayNumber(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.POST("/trips/:id/days", setUserID("user-123"), h.CreateDay)

	body := `{"notes":"Day notes"}`
	req := httptest.NewRequest(http.MethodPost, "/trips/trip-001/days", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ──────────────────────────────────────────────
// CreateItem handler
// ──────────────────────────────────────────────

// TestCreateItem_MissingTitle verifies that binding rejects missing title.
func TestCreateItem_MissingTitle(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.POST("/trips/:id/days/:dayId/items", setUserID("user-123"), h.CreateItem)

	body := `{"type":"activity"}`
	req := httptest.NewRequest(http.MethodPost, "/trips/trip-001/days/day-001/items",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ──────────────────────────────────────────────
// ReorderItems handler
// ──────────────────────────────────────────────

// TestReorderItems_InvalidJSON checks that malformed reorder payloads return 400.
func TestReorderItems_InvalidJSON(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.PATCH("/trips/:id/items/reorder", setUserID("user-123"), h.ReorderItems)

	req := httptest.NewRequest(http.MethodPatch, "/trips/trip-001/items/reorder",
		bytes.NewBufferString("not-json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ──────────────────────────────────────────────
// Health endpoint
// ──────────────────────────────────────────────

func TestHealthEndpoint(t *testing.T) {
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
