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
// POST /trips/:id/collaborators
// ──────────────────────────────────────────────────────────────────────────────

// TestInviteCollaborator_MissingUserID checks binding rejects missing user_id.
func TestInviteCollaborator_MissingUserID(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.POST("/trips/:id/collaborators", injectUserID("owner-001"), h.InviteCollaborator)

	body := `{"role":"editor"}`
	req := httptest.NewRequest(http.MethodPost, "/trips/trip-001/collaborators",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	env := parseEnv(t, w)
	require.NotNil(t, env.Error)
}

// TestInviteCollaborator_MissingRole checks binding rejects missing role.
func TestInviteCollaborator_MissingRole(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.POST("/trips/:id/collaborators", injectUserID("owner-001"), h.InviteCollaborator)

	body := `{"user_id":"user-456"}`
	req := httptest.NewRequest(http.MethodPost, "/trips/trip-001/collaborators",
		bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
	env := parseEnv(t, w)
	require.NotNil(t, env.Error)
}

// TestInviteCollaborator_InvalidJSON checks that malformed JSON returns 400.
func TestInviteCollaborator_InvalidJSON(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.POST("/trips/:id/collaborators", injectUserID("owner-001"), h.InviteCollaborator)

	req := httptest.NewRequest(http.MethodPost, "/trips/trip-001/collaborators",
		bytes.NewBufferString("{not valid json"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// TestUpdateCollaborator_InvalidJSON checks that malformed JSON returns 400.
func TestUpdateCollaborator_InvalidJSON(t *testing.T) {
	h := &Handlers{svc: nil}
	r := gin.New()
	r.PATCH("/trips/:id/collaborators/:userId", injectUserID("owner-001"), h.UpdateCollaborator)

	req := httptest.NewRequest(http.MethodPatch, "/trips/trip-001/collaborators/user-456",
		bytes.NewBufferString("{broken"))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}

// ──────────────────────────────────────────────────────────────────────────────
// Table-driven binding validation tests
// ──────────────────────────────────────────────────────────────────────────────

func TestInviteCollaborator_RequiredFields(t *testing.T) {
	tests := []struct {
		name string
		body string
		want int
	}{
		{"missing both", `{}`, http.StatusBadRequest},
		{"missing role", `{"user_id":"u-1"}`, http.StatusBadRequest},
		{"missing user_id", `{"role":"viewer"}`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &Handlers{svc: nil}
			r := gin.New()
			r.POST("/trips/:id/collaborators", injectUserID("owner"), h.InviteCollaborator)

			req := httptest.NewRequest(http.MethodPost, "/trips/t-1/collaborators",
				bytes.NewBufferString(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			r.ServeHTTP(w, req)

			assert.Equal(t, tt.want, w.Code, "test case: %s", tt.name)
		})
	}
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
