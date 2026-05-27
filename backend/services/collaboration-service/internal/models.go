package internal

import "time"

// Collaborator represents a user who has been granted access to a trip.
type Collaborator struct {
	ID         string     `json:"id"`
	TripID     string     `json:"trip_id"`
	UserID     string     `json:"user_id"`
	Role       string     `json:"role"` // viewer|editor|admin
	InvitedAt  time.Time  `json:"invited_at"`
	AcceptedAt *time.Time `json:"accepted_at"`
}

// InviteCollaboratorRequest is the payload for POST /trips/:id/collaborators.
type InviteCollaboratorRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Role   string `json:"role"    binding:"required"`
}

// UpdateCollaboratorRequest is the payload for PATCH /trips/:id/collaborators/:userId.
type UpdateCollaboratorRequest struct {
	Role *string `json:"role"`
}
