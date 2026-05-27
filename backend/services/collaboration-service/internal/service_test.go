package internal

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──────────────────────────────────────────────────────────────────────────────
// publishEvent — nil-producer safety
// ──────────────────────────────────────────────────────────────────────────────

// TestPublishEvent_NilProducer verifies that publishEvent with a nil Kafka
// producer does not panic. This is a safety-net test for the nil-check guard.
func TestPublishEvent_NilProducer(t *testing.T) {
	svc := &CollaborationService{repo: nil, producer: nil}
	assert.NotPanics(t, func() {
		svc.publishEvent(context.Background(), "test-topic", "test.event", map[string]string{"key": "value"})
	})
}

// ──────────────────────────────────────────────────────────────────────────────
// InviteCollaboratorRequest — JSON round-trip
// ──────────────────────────────────────────────────────────────────────────────

// TestInviteCollaboratorRequest_JSONRoundtrip verifies field mapping.
func TestInviteCollaboratorRequest_JSONRoundtrip(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		userID string
		role   string
	}{
		{"editor", `{"user_id":"u-1","role":"editor"}`, "u-1", "editor"},
		{"viewer", `{"user_id":"u-2","role":"viewer"}`, "u-2", "viewer"},
		{"admin", `{"user_id":"u-3","role":"admin"}`, "u-3", "admin"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req InviteCollaboratorRequest
			require.NoError(t, json.Unmarshal([]byte(tt.input), &req))
			assert.Equal(t, tt.userID, req.UserID)
			assert.Equal(t, tt.role, req.Role)
		})
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// Collaborator model — optional AcceptedAt
// ──────────────────────────────────────────────────────────────────────────────

// TestCollaborator_AcceptedAtNil verifies that AcceptedAt is nil for a new invite.
func TestCollaborator_AcceptedAtNil(t *testing.T) {
	c := &Collaborator{
		ID:     "collab-001",
		TripID: "trip-001",
		UserID: "user-001",
		Role:   "viewer",
	}
	assert.Nil(t, c.AcceptedAt, "AcceptedAt should be nil for a pending invite")
}
