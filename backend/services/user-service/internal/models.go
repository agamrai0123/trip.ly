package internal

import "time"

// User holds the public profile of a WanderPlan user.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	AvatarURL string    `json:"avatar_url"`
	Provider  string    `json:"provider"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdateUserRequest is the payload for PATCH /users/me.
type UpdateUserRequest struct {
	Name      *string `json:"name"`
	AvatarURL *string `json:"avatar_url"`
}
