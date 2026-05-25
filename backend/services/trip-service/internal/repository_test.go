//go:build integration

package internal

import (
	"context"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	testcontainers "github.com/testcontainers/testcontainers-go"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// testSchema mirrors the columns expected by trip-service database.go queries.
// NOTE: uses user_id, budget_total, visibility — NOT owner_id/budget from migration file.
const testSchema = `
CREATE SCHEMA IF NOT EXISTS wanderplan;

CREATE TABLE IF NOT EXISTS wanderplan.trips (
    id              TEXT        PRIMARY KEY,
    user_id         TEXT        NOT NULL,
    title           TEXT        NOT NULL,
    destination     TEXT        NOT NULL DEFAULT '',
    cover_image_url TEXT        NOT NULL DEFAULT '',
    start_date      DATE,
    end_date        DATE,
    status          TEXT        NOT NULL DEFAULT 'draft',
    visibility      TEXT        NOT NULL DEFAULT 'private',
    budget_total    NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency        TEXT        NOT NULL DEFAULT 'USD',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wanderplan.itinerary_days (
    id         TEXT        PRIMARY KEY,
    trip_id    TEXT        NOT NULL REFERENCES wanderplan.trips(id) ON DELETE CASCADE,
    day_number INT         NOT NULL,
    date       DATE,
    notes      TEXT        NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS wanderplan.itinerary_items (
    id          TEXT        PRIMARY KEY,
    day_id      TEXT        NOT NULL REFERENCES wanderplan.itinerary_days(id) ON DELETE CASCADE,
    trip_id     TEXT        NOT NULL,
    title       TEXT        NOT NULL,
    description TEXT        NOT NULL DEFAULT '',
    location    TEXT        NOT NULL DEFAULT '',
    place_id    TEXT        NOT NULL DEFAULT '',
    type        TEXT        NOT NULL DEFAULT 'activity',
    start_time  TEXT        NOT NULL DEFAULT '',
    end_time    TEXT        NOT NULL DEFAULT '',
    cost        NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency    TEXT        NOT NULL DEFAULT 'USD',
    order_index INT         NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
`

// setupTestPool starts a throwaway PostgreSQL container and returns a pgxpool.
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

// ──────────────────────────────────────────────
// TripRepo integration tests
// ──────────────────────────────────────────────

func TestTripRepo_Create_And_Get(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTripRepo(pool)
	ctx := context.Background()

	req := &CreateTripRequest{
		Title:       "Japan Adventure",
		Destination: "Tokyo",
		Status:      "draft",
		Visibility:  "private",
		Currency:    "USD",
	}

	trip, err := repo.Create(ctx, "user-001", req)
	require.NoError(t, err)
	require.NotEmpty(t, trip.ID)
	assert.Equal(t, "Japan Adventure", trip.Title)
	assert.Equal(t, "user-001", trip.UserID)
	assert.Equal(t, "draft", trip.Status)
	assert.Equal(t, "private", trip.Visibility)

	got, err := repo.GetByID(ctx, trip.ID)
	require.NoError(t, err)
	assert.Equal(t, trip.ID, got.ID)
	assert.Equal(t, "Tokyo", got.Destination)
}

func TestTripRepo_GetByID_NotFound(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTripRepo(pool)
	ctx := context.Background()

	_, err := repo.GetByID(ctx, "non-existent")
	assert.ErrorIs(t, err, ErrNotFound)
}

func TestTripRepo_ListByUser(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTripRepo(pool)
	ctx := context.Background()

	// Create 2 trips for user-A and 1 for user-B
	for _, title := range []string{"Trip 1", "Trip 2"} {
		_, err := repo.Create(ctx, "user-A", &CreateTripRequest{
			Title: title, Destination: "Paris", Status: "draft",
			Visibility: "private", Currency: "USD",
		})
		require.NoError(t, err)
	}
	_, err := repo.Create(ctx, "user-B", &CreateTripRequest{
		Title: "Other Trip", Destination: "Berlin", Status: "draft",
		Visibility: "private", Currency: "EUR",
	})
	require.NoError(t, err)

	trips, err := repo.ListByUser(ctx, "user-A")
	require.NoError(t, err)
	assert.Len(t, trips, 2)
	for _, tr := range trips {
		assert.Equal(t, "user-A", tr.UserID)
	}
}

func TestTripRepo_Delete(t *testing.T) {
	pool := setupTestPool(t)
	repo := NewTripRepo(pool)
	ctx := context.Background()

	trip, err := repo.Create(ctx, "user-del", &CreateTripRequest{
		Title: "Delete Me", Destination: "Nowhere", Status: "draft",
		Visibility: "private", Currency: "USD",
	})
	require.NoError(t, err)

	err = repo.Delete(ctx, trip.ID)
	require.NoError(t, err)

	_, err = repo.GetByID(ctx, trip.ID)
	assert.ErrorIs(t, err, ErrNotFound)
}

// ──────────────────────────────────────────────
// DayRepo integration tests
// ──────────────────────────────────────────────

func TestDayRepo_CreateAndList(t *testing.T) {
	pool := setupTestPool(t)
	tripRepo := NewTripRepo(pool)
	dayRepo := NewDayRepo(pool)
	ctx := context.Background()

	trip, err := tripRepo.Create(ctx, "user-day", &CreateTripRequest{
		Title: "Day Test Trip", Destination: "Rome", Status: "draft",
		Visibility: "private", Currency: "EUR",
	})
	require.NoError(t, err)

	day, err := dayRepo.Create(ctx, trip.ID, &CreateDayRequest{DayNumber: 1, Notes: "Day 1"})
	require.NoError(t, err)
	require.NotEmpty(t, day.ID)
	assert.Equal(t, 1, day.DayNumber)

	days, err := dayRepo.ListByTrip(ctx, trip.ID)
	require.NoError(t, err)
	assert.Len(t, days, 1)
	assert.Equal(t, "Day 1", days[0].Notes)
}
