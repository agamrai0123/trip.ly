-- 000002_add_search_vector.down.sql
-- Reverses 000002_add_search_vector.up.sql exactly.

BEGIN;

DROP INDEX IF EXISTS wanderplan.idx_places_cache_results;
DROP INDEX IF EXISTS wanderplan.idx_trips_search_vector;
ALTER TABLE wanderplan.trips DROP COLUMN IF EXISTS search_vector;

COMMIT;
