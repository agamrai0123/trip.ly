-- 000002_add_search_vector.up.sql
-- Adds full-text search vector to trips and a GIN index on places_cache results.
-- search_vector is a stored generated column so it stays in sync with title, destination, and description.

BEGIN;

-- Add generated tsvector column for full-text search on trips.
ALTER TABLE wanderplan.trips
    ADD COLUMN IF NOT EXISTS search_vector tsvector
        GENERATED ALWAYS AS (
            to_tsvector(
                'english',
                coalesce(title, '') || ' ' ||
                coalesce(destination, '') || ' ' ||
                coalesce(description, '')
            )
        ) STORED;

-- GIN index for fast plainto_tsquery / ts_rank lookups on trips.
CREATE INDEX IF NOT EXISTS idx_trips_search_vector
    ON wanderplan.trips USING GIN (search_vector);

-- GIN index for fast JSONB querying on places_cache results.
CREATE INDEX IF NOT EXISTS idx_places_cache_results
    ON wanderplan.places_cache USING GIN (results jsonb_path_ops);

COMMIT;
