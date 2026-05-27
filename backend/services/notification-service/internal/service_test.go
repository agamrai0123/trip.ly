package internal

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──────────────────────────────────────────────────────────────────────────────
// Hub — pure concurrency-safe bookkeeping
// ──────────────────────────────────────────────────────────────────────────────

// TestHub_NewHub_IsEmpty verifies that a fresh Hub has no clients.
func TestHub_NewHub_IsEmpty(t *testing.T) {
	hub := NewHub()
	hub.mu.RLock()
	defer hub.mu.RUnlock()
	assert.Empty(t, hub.clients)
}

// TestHub_Register_AddsEntry verifies that Register increments the client list.
func TestHub_Register_AddsEntry(t *testing.T) {
	hub := NewHub()
	hub.Register("user-A", nil) // nil conn is fine for bookkeeping
	hub.mu.RLock()
	count := len(hub.clients["user-A"])
	hub.mu.RUnlock()
	assert.Equal(t, 1, count)
}

// TestHub_Register_MultipleConns verifies that the same user can have multiple connections.
func TestHub_Register_MultipleConns(t *testing.T) {
	hub := NewHub()
	hub.Register("user-B", nil)
	hub.Register("user-B", nil)
	hub.mu.RLock()
	count := len(hub.clients["user-B"])
	hub.mu.RUnlock()
	assert.Equal(t, 2, count)
}

// TestHub_Unregister_UnknownUser_NoPanic verifies graceful handling of unknown users.
func TestHub_Unregister_UnknownUser_NoPanic(t *testing.T) {
	hub := NewHub()
	assert.NotPanics(t, func() {
		hub.Unregister("no-such-user", nil)
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// Notification model
// ──────────────────────────────────────────────────────────────────────────────

// TestNotification_ReadAt_NilForUnread verifies that ReadAt is nil by default.
func TestNotification_ReadAt_NilForUnread(t *testing.T) {
	n := Notification{ID: "n-001", UserID: "u-001", Type: "trip.update"}
	assert.Nil(t, n.ReadAt)
}

// TestNotification_JSONSerialisation verifies the JSON field names used by the API.
func TestNotification_JSONSerialisation(t *testing.T) {
	now := time.Date(2026, 1, 15, 12, 0, 0, 0, time.UTC)
	n := Notification{
		ID:        "n-100",
		UserID:    "u-100",
		Type:      "collab.invite",
		Title:     "New invitation",
		Body:      "You were invited to a trip",
		ReadAt:    &now,
		CreatedAt: now,
	}

	b, err := json.Marshal(n)
	require.NoError(t, err)

	var m map[string]interface{}
	require.NoError(t, json.Unmarshal(b, &m))

	// Verify JSON key names match what the API returns
	assert.Contains(t, m, "id")
	assert.Contains(t, m, "user_id")
	assert.Contains(t, m, "type")
	assert.Contains(t, m, "title")
	assert.Contains(t, m, "body")
	assert.Contains(t, m, "read_at")
	assert.Contains(t, m, "created_at")
	assert.Equal(t, "n-100", m["id"])
}
