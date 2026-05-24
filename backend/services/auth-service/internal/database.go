package internal

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// ErrNotFound is returned when a row is not found.
var ErrNotFound = errors.New("not found")

// UserRepo handles user persistence.
type UserRepo struct{ pool *pgxpool.Pool }

func NewUserRepo(pool *pgxpool.Pool) *UserRepo { return &UserRepo{pool: pool} }

// Upsert inserts or updates a user row and returns the current record.
func (r *UserRepo) Upsert(ctx context.Context, u *User) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		INSERT INTO wanderplan.users (id, email, name, avatar_url, provider, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
		ON CONFLICT (email) DO UPDATE
			SET name       = EXCLUDED.name,
			    avatar_url = EXCLUDED.avatar_url,
			    provider   = EXCLUDED.provider,
			    updated_at = NOW()
		RETURNING id, email, name, avatar_url, provider, created_at, updated_at`,
		u.ID, u.Email, u.Name, u.AvatarURL, u.Provider,
	)
	out := &User{}
	err := row.Scan(&out.ID, &out.Email, &out.Name, &out.AvatarURL, &out.Provider, &out.CreatedAt, &out.UpdatedAt)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// GetByID fetches a user by their UUID.
func (r *UserRepo) GetByID(ctx context.Context, id string) (*User, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT id, email, name, avatar_url, provider, created_at, updated_at
		FROM wanderplan.users WHERE id = $1`, id)
	u := &User{}
	err := row.Scan(&u.ID, &u.Email, &u.Name, &u.AvatarURL, &u.Provider, &u.CreatedAt, &u.UpdatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return nil, ErrNotFound
	}
	if err != nil {
		return nil, err
	}
	return u, nil
}

// RefreshTokenRepo handles refresh token persistence.
type RefreshTokenRepo struct{ pool *pgxpool.Pool }

func NewRefreshTokenRepo(pool *pgxpool.Pool) *RefreshTokenRepo {
	return &RefreshTokenRepo{pool: pool}
}

func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

// Create stores a new refresh token and returns the raw opaque token string.
func (r *RefreshTokenRepo) Create(ctx context.Context, userID string, ttl time.Duration) (string, error) {
	raw := uuid.New().String() + uuid.New().String()
	tokenHash := hashToken(raw)
	_, err := r.pool.Exec(ctx, `
		INSERT INTO wanderplan.refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		VALUES ($1, $2, $3, $4, NOW())`,
		uuid.New().String(), userID, tokenHash, time.Now().Add(ttl),
	)
	if err != nil {
		return "", err
	}
	return raw, nil
}

// Rotate validates a token, deletes it (rotation), and creates a new one.
func (r *RefreshTokenRepo) Rotate(ctx context.Context, rawToken string, ttl time.Duration) (string, string, error) {
	tokenHash := hashToken(rawToken)
	var userID string
	var expiresAt time.Time
	err := r.pool.QueryRow(ctx, `
		SELECT user_id, expires_at FROM wanderplan.refresh_tokens
		WHERE token_hash = $1`, tokenHash,
	).Scan(&userID, &expiresAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return "", "", ErrNotFound
	}
	if err != nil {
		return "", "", err
	}
	if time.Now().After(expiresAt) {
		_, _ = r.pool.Exec(ctx, `DELETE FROM wanderplan.refresh_tokens WHERE token_hash = $1`, tokenHash)
		return "", "", ErrNotFound
	}
	if _, err = r.pool.Exec(ctx, `DELETE FROM wanderplan.refresh_tokens WHERE token_hash = $1`, tokenHash); err != nil {
		return "", "", err
	}
	newRaw, err := r.Create(ctx, userID, ttl)
	if err != nil {
		return "", "", err
	}
	return userID, newRaw, nil
}

// DeleteByUserID revokes all refresh tokens for a user (logout).
func (r *RefreshTokenRepo) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.pool.Exec(ctx, `DELETE FROM wanderplan.refresh_tokens WHERE user_id = $1`, userID)
	return err
}
