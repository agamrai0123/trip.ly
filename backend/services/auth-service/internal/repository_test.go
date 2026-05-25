//go:build integration

package internal

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	testcontainers "github.com/testcontainers/testcontainers-go"
)

// testSchema matches the columns used by auth-service queries.
const testSchema = `
CREATE SCHEMA IF NOT EXISTS wanderplan;

CREATE TABLE IF NOT EXISTS wanderplan.users (
    id          TEXT        PRIMARY KEY,
    email       TEXT        NOT NULL UNIQUE,
    name        TEXT        NOT NULL DEFAULT '',
    avatar_url  TEXT        NOT NULL DEFAULT '',
    provider    TEXT        NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wanderplan.refresh_tokens (
    id           TEXT        PRIMARY KEY,
    user_id      TEXT        NOT NULL REFERENCES wanderplan.users(id) ON DELETE CASCADE,
    token_hash   TEXT        NOT NULL UNIQUE,
    expires_at   TIMESTAMPTZ NOT NULL,
    revoked_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

// setupTestPool starts a throwaway PostgreSQL container and returns a pgxpool.
// It registers container teardown and pool close via t.Cleanup.
func setupTestPool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	ctx := context.Background()

	pgContainer, err := tcpostgres.Run(ctx,
		"postgres:16-alpine",
		tcpostgres.WithDatabase("testdb"),
		tcpostgres.WithUsername("testuser"),
		tcpostgres.WithPassword("testpass"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(60*time.Second),
		),
	)
	require.NoError(t, err)
	t.Cleanup(func() { _ = pgContainer.Terminate(ctx) })

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	require.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	require.NoError(t, err)
	t.Cleanup(pool.Close)

	_, err = pool.Exec(ctx, testSchema)
	require.NoError(t, err, "run test schema")

	return pool
}

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

// ──────────────────────────────────────────────
// UserRepo tests
// ──────────────────────────────────────────────

func TestUserRepo_Upsert_Insert(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	u := &User{
		ID:        "u-001",
		Email:     "alice@example.com",
		Name:      "Alice",
		AvatarURL: "https://example.com/alice.jpg",
		Provider:  "google",
	}

	got, err := repo.Upsert(ctx, u)
	require.NoError(t, err)
	assert.Equal(t, u.ID, got.ID)
	assert.Equal(t, u.Email, got.Email)
	assert.Equal(t, u.Name, got.Name)
}

func TestUserRepo_Upsert_UpdateOnConflict(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	u := &User{ID: "u-002", Email: "bob@example.com", Name: "Bob", Provider: "github"}
	_, err := repo.Upsert(ctx, u)
	require.NoError(t, err)

	// Update name and avatar on second upsert (same email)
	u2 := &User{ID: "u-002", Email: "bob@example.com", Name: "Bobby", AvatarURL: "https://example.com/bob.jpg", Provider: "github"}
	got, err := repo.Upsert(ctx, u2)
	require.NoError(t, err)
	assert.Equal(t, "Bobby", got.Name)
	assert.Equal(t, "https://example.com/bob.jpg", got.AvatarURL)
}

func TestUserRepo_GetByID_Found(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	u := &User{ID: "u-003", Email: "carol@example.com", Name: "Carol", Provider: "google"}
	_, err := repo.Upsert(ctx, u)
	require.NoError(t, err)

	got, err := repo.GetByID(ctx, "u-003")
	require.NoError(t, err)
	assert.Equal(t, "carol@example.com", got.Email)
}

func TestUserRepo_GetByID_NotFound(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewUserRepo(pool)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent-id")
	assert.ErrorIs(t, err, ErrNotFound)
}

// ──────────────────────────────────────────────
// RefreshTokenRepo tests
// ──────────────────────────────────────────────

func TestRefreshTokenRepo_Create(t *testing.T) {
	pool := setupTestPool(t)
	users := NewUserRepo(pool)
	tokens := NewRefreshTokenRepo(pool)
	ctx := context.Background()

	_, err := users.Upsert(ctx, &User{ID: "u-rt-001", Email: "dan@example.com", Name: "Dan", Provider: "google"})
	require.NoError(t, err)

	rawToken, err := tokens.Create(ctx, "u-rt-001", 30*24*time.Hour)
	require.NoError(t, err)
	assert.NotEmpty(t, rawToken)
}

func TestRefreshTokenRepo_Rotate_Success(t *testing.T) {
	pool := setupTestPool(t)
	users := NewUserRepo(pool)
	tokens := NewRefreshTokenRepo(pool)
	ctx := context.Background()

	_, err := users.Upsert(ctx, &User{ID: "u-rt-002", Email: "eve@example.com", Name: "Eve", Provider: "github"})
	require.NoError(t, err)

	rawToken, err := tokens.Create(ctx, "u-rt-002", 30*24*time.Hour)
	require.NoError(t, err)

	userID, newRaw, err := tokens.Rotate(ctx, rawToken, 30*24*time.Hour)
	require.NoError(t, err)
	assert.Equal(t, "u-rt-002", userID)
	assert.NotEmpty(t, newRaw)
	assert.NotEqual(t, rawToken, newRaw, "rotated token must differ from original")
}

func TestRefreshTokenRepo_Rotate_NotFound(t *testing.T) {
	pool := setupTestPool(t)
	tokens := NewRefreshTokenRepo(pool)
	ctx := context.Background()

	_, _, err := tokens.Rotate(ctx, "does-not-exist", 30*24*time.Hour)
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestRefreshTokenRepo_Rotate_Expired(t *testing.T) {
	pool := setupTestPool(t)
	users := NewUserRepo(pool)
	tokens := NewRefreshTokenRepo(pool)
	ctx := context.Background()

	_, err := users.Upsert(ctx, &User{ID: "u-rt-003", Email: "frank@example.com", Name: "Frank", Provider: "google"})
	require.NoError(t, err)

	// Create a token that is already expired
	rawToken, err := tokens.Create(ctx, "u-rt-003", -1*time.Second)
	require.NoError(t, err)

	_, _, err = tokens.Rotate(ctx, rawToken, 30*24*time.Hour)
	assert.Error(t, err, "rotating an expired token should return an error")
}
