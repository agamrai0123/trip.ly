package internal

import (
	"time"
)

// User represents a WanderPlan user account.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// RefreshToken is an opaque server-side token used to issue new access tokens.
type RefreshToken struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

// LoginResponse is returned to the frontend after a successful OAuth callback.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	User        *User  `json:"user"`
}

// RefreshResponse is returned after a successful token refresh.
type RefreshResponse struct {
	AccessToken string `json:"access_token"`
}
