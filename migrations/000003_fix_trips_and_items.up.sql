-- Migration 000003: align trips and itinerary_items columns with service code.
-- Safe to re-run: uses IF EXISTS / IF NOT EXISTS guards.
-- This corrects schema applied by 000001 before the rename was done.

-- 1. trips: rename owner_id → user_id (service code uses user_id)
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'wanderplan' AND table_name = 'trips' AND column_name = 'owner_id'
    ) THEN
        ALTER TABLE wanderplan.trips RENAME COLUMN owner_id TO user_id;
        -- Recreate index under new name.
        DROP INDEX IF EXISTS wanderplan.idx_trips_owner_id;
        CREATE INDEX IF NOT EXISTS idx_trips_user_id ON wanderplan.trips(user_id);
    END IF;
END;
$$;

-- 2. trips: rename budget → budget_total
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'wanderplan' AND table_name = 'trips' AND column_name = 'budget'
    ) THEN
        ALTER TABLE wanderplan.trips RENAME COLUMN budget TO budget_total;
    END IF;
END;
$$;

-- 3. trips: add visibility column (private|shared|public)
ALTER TABLE wanderplan.trips
    ADD COLUMN IF NOT EXISTS visibility TEXT NOT NULL DEFAULT 'private';

-- 4. itinerary_items: rename position → order_index
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'wanderplan' AND table_name = 'itinerary_items' AND column_name = 'position'
    ) THEN
        ALTER TABLE wanderplan.itinerary_items RENAME COLUMN position TO order_index;
    END IF;
END;
$$;

-- 5. itinerary_items: add trip_id FK for direct trip-scoped queries and reorder
ALTER TABLE wanderplan.itinerary_items
    ADD COLUMN IF NOT EXISTS trip_id TEXT REFERENCES wanderplan.trips(id) ON DELETE CASCADE;

-- Backfill trip_id from the parent day (only for existing rows where it is NULL)
UPDATE wanderplan.itinerary_items ii
SET    trip_id = d.trip_id
FROM   wanderplan.itinerary_days d
WHERE  ii.day_id = d.id
  AND  ii.trip_id IS NULL;

-- Make trip_id NOT NULL once backfilled
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'wanderplan' AND table_name = 'itinerary_items'
          AND column_name = 'trip_id' AND is_nullable = 'YES'
    ) THEN
        ALTER TABLE wanderplan.itinerary_items ALTER COLUMN trip_id SET NOT NULL;
    END IF;
END;
$$;

CREATE INDEX IF NOT EXISTS idx_itinerary_items_trip_id ON wanderplan.itinerary_items(trip_id);
