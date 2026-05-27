package internal

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ──────────────────────────────────────────────────────────────────────────────
// UpdateUserRequest — field optionality
// ──────────────────────────────────────────────────────────────────────────────

// TestUpdateUserRequest_PartialFields verifies that UpdateUserRequest can carry
// partial data (all fields are optional pointers).
func TestUpdateUserRequest_PartialFields(t *testing.T) {
	tests := []struct {
		name      string
		json      string
		wantName  *string
		wantAvURL *string
	}{
		{
			name:      "both set",
			json:      `{"name":"Alice","avatar_url":"https://example.com/a.jpg"}`,
			wantName:  strPtr("Alice"),
			wantAvURL: strPtr("https://example.com/a.jpg"),
		},
		{
			name:      "only name",
			json:      `{"name":"Bob"}`,
			wantName:  strPtr("Bob"),
			wantAvURL: nil,
		},
		{
			name:      "only avatar_url",
			json:      `{"avatar_url":"https://example.com/b.jpg"}`,
			wantName:  nil,
			wantAvURL: strPtr("https://example.com/b.jpg"),
		},
		{
			name:      "empty object",
			json:      `{}`,
			wantName:  nil,
			wantAvURL: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var req UpdateUserRequest
			require.NoError(t, json.Unmarshal([]byte(tt.json), &req))
			assert.Equal(t, tt.wantName, req.Name)
			assert.Equal(t, tt.wantAvURL, req.AvatarURL)
		})
	}
}

// ──────────────────────────────────────────────────────────────────────────────
// User model — JSON serialisation
// ──────────────────────────────────────────────────────────────────────────────

// TestUser_JSONRoundtrip verifies that a User can be marshalled and unmarshalled
// without data loss.
func TestUser_JSONRoundtrip(t *testing.T) {
	u := &User{
		ID:        "user-001",
		Email:     "alice@example.com",
		Name:      "Alice",
		AvatarURL: "https://example.com/alice.jpg",
		Provider:  "google",
	}

	b, err := json.Marshal(u)
	require.NoError(t, err)

	var got User
	require.NoError(t, json.Unmarshal(b, &got))
	assert.Equal(t, u.ID, got.ID)
	assert.Equal(t, u.Email, got.Email)
	assert.Equal(t, u.Name, got.Name)
	assert.Equal(t, u.AvatarURL, got.AvatarURL)
	assert.Equal(t, u.Provider, got.Provider)
}

// strPtr is a helper to get a pointer to a string literal.
func strPtr(s string) *string { return &s }
