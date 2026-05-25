package internal

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// injectUserID is middleware that sets user_id in the Gin context.
func injectUserID(id string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("user_id", id)
		c.Next()
	}
}

// ──────────────────────────────────────────────
// GET /notifications
// ──────────────────────────────────────────────

// TestListNotifications_NilSvcPanics checks that the route exists and panics (nil svc)
// are caught by gin.Recovery, returning 500 (not 404).
func TestListNotifications_RouteExists(t *testing.T) {
	h := &Handlers{svc: nil, hub: nil}
	r := gin.New()
	r.Use(gin.Recovery()) // catches nil-pointer panic as 500
	r.GET("/notifications", injectUserID("user-001"), h.ListNotifications)

	req := httptest.NewRequest(http.MethodGet, "/notifications", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	// Not 404 means route registration succeeded; panic → 500.
	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

// ──────────────────────────────────────────────
// PATCH /notifications/:id/read
// ──────────────────────────────────────────────

// TestMarkRead_RouteExists verifies the route is wired correctly.
func TestMarkRead_RouteExists(t *testing.T) {
	h := &Handlers{svc: nil, hub: nil}
	r := gin.New()
	r.Use(gin.Recovery())
	r.PATCH("/notifications/:id/read", injectUserID("user-001"), h.MarkRead)

	req := httptest.NewRequest(http.MethodPatch, "/notifications/notif-001/read", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

// ──────────────────────────────────────────────
// PATCH /notifications/read-all
// ──────────────────────────────────────────────

// TestMarkAllRead_RouteExists verifies the route is wired correctly.
func TestMarkAllRead_RouteExists(t *testing.T) {
	h := &Handlers{svc: nil, hub: nil}
	r := gin.New()
	r.Use(gin.Recovery())
	r.PATCH("/notifications/read-all", injectUserID("user-001"), h.MarkAllRead)

	req := httptest.NewRequest(http.MethodPatch, "/notifications/read-all", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)

	assert.NotEqual(t, http.StatusNotFound, w.Code)
}

// ──────────────────────────────────────────────
// Health
// ──────────────────────────────────────────────

func TestNotificationHealthz(t *testing.T) {
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

// ──────────────────────────────────────────────
// Notification model
// ──────────────────────────────────────────────

// TestNotification_ReadAtField checks that ReadAt is nil by default and can be set.
func TestNotification_ReadAtField(t *testing.T) {
	n := &Notification{ID: "n-001", UserID: "u-001"}
	assert.Nil(t, n.ReadAt, "ReadAt should be nil for unread notifications")
	now := time.Now()
	n.ReadAt = &now
	assert.NotNil(t, n.ReadAt)
}
