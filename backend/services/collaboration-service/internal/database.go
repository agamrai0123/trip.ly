package internal

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a queried row does not exist.
var ErrNotFound = errors.New("not found")

// CollaboratorRepo handles collaborator persistence.
type CollaboratorRepo struct{ pool *pgxpool.Pool }

// NewCollaboratorRepo creates a repo backed by the supplied pool.
func NewCollaboratorRepo(pool *pgxpool.Pool) *CollaboratorRepo {
	return &CollaboratorRepo{pool: pool}
}

// List returns all collaborators for a trip.
func (r *CollaboratorRepo) List(ctx context.Context, tripID string) ([]*Collaborator, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, trip_id, user_id, role, invited_at, accepted_at
		FROM wanderplan.collaborators WHERE trip_id=$1 ORDER BY invited_at`, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var collabs []*Collaborator
	for rows.Next() {
		c := &Collaborator{}
		if err := rows.Scan(&c.ID, &c.TripID, &c.UserID, &c.Role, &c.InvitedAt, &c.AcceptedAt); err != nil {
			return nil, err
		}
		collabs = append(collabs, c)
	}
	return collabs, rows.Err()
}

// GetByTripAndUser fetches a single collaborator record.
func (r *CollaboratorRepo) GetByTripAndUser(ctx context.Context, tripID, userID string) (*Collaborator, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, trip_id, user_id, role, invited_at, accepted_at
		FROM wanderplan.collaborators WHERE trip_id=$1 AND user_id=$2`, tripID, userID)
	c := &Collaborator{}
	err := row.Scan(&c.ID, &c.TripID, &c.UserID, &c.Role, &c.InvitedAt, &c.AcceptedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return c, err
}

// Invite inserts a new collaborator record.
func (r *CollaboratorRepo) Invite(ctx context.Context, tripID, userID, role string) (*Collaborator, error) {
	id := uuid.New().String()
	now := time.Now().UTC()
	row := r.pool.QueryRow(ctx, `
		INSERT INTO wanderplan.collaborators (id, trip_id, user_id, role, invited_at)
		VALUES ($1,$2,$3,$4,$5)
		RETURNING id, trip_id, user_id, role, invited_at, accepted_at`,
		id, tripID, userID, role, now)
	c := &Collaborator{}
	err := row.Scan(&c.ID, &c.TripID, &c.UserID, &c.Role, &c.InvitedAt, &c.AcceptedAt)
	return c, err
}

// Update changes the role of an existing collaborator.
func (r *CollaboratorRepo) Update(ctx context.Context, tripID, userID string, req *UpdateCollaboratorRequest) (*Collaborator, error) {
	existing, err := r.GetByTripAndUser(ctx, tripID, userID)
	if err != nil {
		return nil, err
	}
	if req.Role != nil {
		existing.Role = *req.Role
	}
	row := r.pool.QueryRow(ctx, `
		UPDATE wanderplan.collaborators SET role=$3
		WHERE trip_id=$1 AND user_id=$2
		RETURNING id, trip_id, user_id, role, invited_at, accepted_at`,
		tripID, userID, existing.Role)
	c := &Collaborator{}
	err = row.Scan(&c.ID, &c.TripID, &c.UserID, &c.Role, &c.InvitedAt, &c.AcceptedAt)
	return c, err
}

// Remove deletes a collaborator record.
func (r *CollaboratorRepo) Remove(ctx context.Context, tripID, userID string) error {
	tag, err := r.pool.Exec(ctx, `
		DELETE FROM wanderplan.collaborators WHERE trip_id=$1 AND user_id=$2`, tripID, userID)
	if err != nil {
		return err
	}
	if tag.RowsAffected() == 0 {
		return ErrNotFound
	}
	return nil
}
