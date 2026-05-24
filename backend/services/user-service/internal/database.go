package internal

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a queried row does not exist.
var ErrNotFound = errors.New("not found")

// UserRepo handles user persistence.
type UserRepo struct{ pool *pgxpool.Pool }

// NewUserRepo creates a UserRepo backed by the supplied pool.
func NewUserRepo(pool *pgxpool.Pool) *UserRepo { return &UserRepo{pool: pool} }

// GetByID fetches a user by primary key.
func (r *UserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, email, name, avatar_url, provider, created_at, updated_at
		FROM wanderplan.users WHERE id=$1`, id)
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.Provider, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	return u, err
}

// Update applies partial changes to a user profile.
func (r *UserRepo) Update(ctx context.Context, id string, req *UpdateUserRequest) (*User, error) {
	existing, err := r.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if req.Name != nil {
		existing.Name = *req.Name
	}
	if req.AvatarURL != nil {
		existing.AvatarURL = *req.AvatarURL
	}
	existing.UpdatedAt = time.Now().UTC()

	row := r.pool.QueryRow(ctx, `
		UPDATE wanderplan.users SET name=$2, avatar_url=$3, updated_at=$4
		WHERE id=$1
		RETURNING id, email, name, avatar_url, provider, created_at, updated_at`,
		existing.ID, existing.Name, existing.AvatarURL, existing.UpdatedAt,
	)
	out := &User{}
	err = row.Scan(&out.ID, &out.Email, &out.Name, &out.AvatarURL, &out.Provider, &out.CreatedAt, &out.UpdatedAt)
	return out, err
}
