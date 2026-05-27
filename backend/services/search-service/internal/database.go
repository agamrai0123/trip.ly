package internal

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

// PlaceCacheRepo handles caching of Google Places API results.
type PlaceCacheRepo struct{ pool *pgxpool.Pool }

// NewPlaceCacheRepo creates a repo backed by the given pool.
func NewPlaceCacheRepo(pool *pgxpool.Pool) *PlaceCacheRepo {
	return &PlaceCacheRepo{pool: pool}
}

// Get retrieves a cached place by query key.
func (r *PlaceCacheRepo) Get(ctx context.Context, cacheKey string) ([]*PlaceResult, error) {
	row := r.pool.QueryRow(ctx, `
		SELECT results FROM wanderplan.places_cache WHERE cache_key=$1`, cacheKey)
	var raw []byte
	if err := row.Scan(&raw); err != nil {
		return nil, err
	}
	// Deserialise from JSONB.
	return unmarshalPlaces(raw)
}

// Set upserts a cached result set for the given key.
func (r *PlaceCacheRepo) Set(ctx context.Context, cacheKey string, places []*PlaceResult) error {
	data, err := marshalPlaces(places)
	if err != nil {
		return err
	}
	_, err = r.pool.Exec(ctx, `
		INSERT INTO wanderplan.places_cache (cache_key, results, cached_at)
		VALUES ($1,$2,NOW())
		ON CONFLICT (cache_key) DO UPDATE SET results=EXCLUDED.results, cached_at=NOW()`,
		cacheKey, data)
	return err
}

// TripSearchRepo searches trips by full-text.
type TripSearchRepo struct{ pool *pgxpool.Pool }

// NewTripSearchRepo creates a repo backed by the given pool.
func NewTripSearchRepo(pool *pgxpool.Pool) *TripSearchRepo {
	return &TripSearchRepo{pool: pool}
}

// Search performs a PostgreSQL full-text search on trips.
func (r *TripSearchRepo) Search(ctx context.Context, query string) ([]*TripSearchResult, error) {
	rows, err := r.pool.Query(ctx, `
		SELECT id, title, destination, status, owner_id
		FROM wanderplan.trips
		WHERE to_tsvector('english', coalesce(title,'') || ' ' || coalesce(destination,'')) @@ plainto_tsquery('english', $1)
		LIMIT 50`, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var results []*TripSearchResult
	for rows.Next() {
		t := &TripSearchResult{}
		if err := rows.Scan(&t.ID, &t.Title, &t.Destination, &t.Status, &t.OwnerID); err != nil {
			return nil, err
		}
		results = append(results, t)
	}
	return results, rows.Err()
}
