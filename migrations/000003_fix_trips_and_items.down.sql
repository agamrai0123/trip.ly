-- Rollback for 000003_fix_trips_and_items.
-- Reverses renames and drops added columns.

-- 5. Drop trip_id index and column
DROP INDEX IF EXISTS wanderplan.idx_itinerary_items_trip_id;
ALTER TABLE wanderplan.itinerary_items DROP COLUMN IF EXISTS trip_id;

-- 4. Rename order_index → position
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'wanderplan' AND table_name = 'itinerary_items' AND column_name = 'order_index'
    ) THEN
        ALTER TABLE wanderplan.itinerary_items RENAME COLUMN order_index TO position;
    END IF;
END;
$$;

-- 3. Drop visibility column
ALTER TABLE wanderplan.trips DROP COLUMN IF EXISTS visibility;

-- 2. Rename budget_total → budget
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'wanderplan' AND table_name = 'trips' AND column_name = 'budget_total'
    ) THEN
        ALTER TABLE wanderplan.trips RENAME COLUMN budget_total TO budget;
    END IF;
END;
$$;

-- 1. Rename user_id → owner_id and restore old index
DO $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'wanderplan' AND table_name = 'trips' AND column_name = 'user_id'
    ) THEN
        ALTER TABLE wanderplan.trips RENAME COLUMN user_id TO owner_id;
        DROP INDEX IF EXISTS wanderplan.idx_trips_user_id;
        CREATE INDEX IF NOT EXISTS idx_trips_owner_id ON wanderplan.trips(owner_id);
    END IF;
END;
$$;
