-- 000001_init.up.sql
-- WanderPlan initial schema: all tables under the `wanderplan` schema.

CREATE SCHEMA IF NOT EXISTS wanderplan;

-- ── users ────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.users (
    id          TEXT        PRIMARY KEY,
    email       TEXT        NOT NULL UNIQUE,
    name        TEXT        NOT NULL DEFAULT '',
    avatar_url  TEXT        NOT NULL DEFAULT '',
    provider    TEXT        NOT NULL,               -- 'google' | 'github'
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ── refresh_tokens ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.refresh_tokens (
    id           TEXT        PRIMARY KEY,
    user_id      TEXT        NOT NULL REFERENCES wanderplan.users(id) ON DELETE CASCADE,
    token_hash   TEXT        NOT NULL UNIQUE,       -- bcrypt hash of opaque token
    expires_at   TIMESTAMPTZ NOT NULL,
    revoked_at   TIMESTAMPTZ,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_user_id ON wanderplan.refresh_tokens(user_id);

-- ── trips ─────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.trips (
    id              TEXT        PRIMARY KEY,
    user_id         TEXT        NOT NULL REFERENCES wanderplan.users(id) ON DELETE CASCADE,
    title           TEXT        NOT NULL,
    description     TEXT        NOT NULL DEFAULT '',
    destination     TEXT        NOT NULL DEFAULT '',
    cover_image_url TEXT        NOT NULL DEFAULT '',
    start_date      DATE,
    end_date        DATE,
    status          TEXT        NOT NULL DEFAULT 'draft',     -- draft|planned|completed
    visibility      TEXT        NOT NULL DEFAULT 'private',   -- private|shared|public
    budget_total    NUMERIC(12,2) NOT NULL DEFAULT 0,
    currency        TEXT        NOT NULL DEFAULT 'USD',
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_trips_user_id  ON wanderplan.trips(user_id);
CREATE INDEX IF NOT EXISTS idx_trips_status   ON wanderplan.trips(status);

-- ── itinerary_days ────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.itinerary_days (
    id          TEXT        PRIMARY KEY,
    trip_id     TEXT        NOT NULL REFERENCES wanderplan.trips(id) ON DELETE CASCADE,
    day_number  INT         NOT NULL,
    date        DATE,
    title       TEXT        NOT NULL DEFAULT '',
    notes       TEXT        NOT NULL DEFAULT '',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (trip_id, day_number)
);

CREATE INDEX IF NOT EXISTS idx_itinerary_days_trip_id ON wanderplan.itinerary_days(trip_id);

-- ── itinerary_items ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.itinerary_items (
    id            TEXT        PRIMARY KEY,
    day_id        TEXT        NOT NULL REFERENCES wanderplan.itinerary_days(id) ON DELETE CASCADE,
    trip_id       TEXT        NOT NULL REFERENCES wanderplan.trips(id) ON DELETE CASCADE,
    title         TEXT        NOT NULL,
    description   TEXT        NOT NULL DEFAULT '',
    type          TEXT        NOT NULL DEFAULT 'activity',  -- activity|transport|accommodation|meal
    location      TEXT        NOT NULL DEFAULT '',
    start_time    TEXT        NOT NULL DEFAULT '',
    end_time      TEXT        NOT NULL DEFAULT '',
    cost          NUMERIC(10,2) NOT NULL DEFAULT 0,
    currency      TEXT        NOT NULL DEFAULT 'USD',
    place_id      TEXT        NOT NULL DEFAULT '',
    order_index   INT         NOT NULL DEFAULT 0,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_itinerary_items_day_id  ON wanderplan.itinerary_items(day_id);
CREATE INDEX IF NOT EXISTS idx_itinerary_items_trip_id ON wanderplan.itinerary_items(trip_id);

-- ── collaborators ─────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.collaborators (
    id           TEXT        PRIMARY KEY,
    trip_id      TEXT        NOT NULL REFERENCES wanderplan.trips(id) ON DELETE CASCADE,
    user_id      TEXT        NOT NULL REFERENCES wanderplan.users(id) ON DELETE CASCADE,
    role         TEXT        NOT NULL DEFAULT 'viewer',  -- viewer|editor|admin
    invited_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    accepted_at  TIMESTAMPTZ,
    UNIQUE (trip_id, user_id)
);

CREATE INDEX IF NOT EXISTS idx_collaborators_trip_id ON wanderplan.collaborators(trip_id);
CREATE INDEX IF NOT EXISTS idx_collaborators_user_id ON wanderplan.collaborators(user_id);

-- ── trip_tags ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.trip_tags (
    trip_id TEXT NOT NULL REFERENCES wanderplan.trips(id) ON DELETE CASCADE,
    tag     TEXT NOT NULL,
    PRIMARY KEY (trip_id, tag)
);

-- ── places_cache ──────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.places_cache (
    cache_key   TEXT        PRIMARY KEY,
    results     JSONB       NOT NULL DEFAULT '[]',
    cached_at   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- ── notifications ─────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.notifications (
    id          TEXT        PRIMARY KEY,
    user_id     TEXT        NOT NULL REFERENCES wanderplan.users(id) ON DELETE CASCADE,
    type        TEXT        NOT NULL,
    title       TEXT        NOT NULL,
    body        TEXT        NOT NULL DEFAULT '',
    read_at     TIMESTAMPTZ,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id    ON wanderplan.notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON wanderplan.notifications(created_at DESC);

-- ── audit_log ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS wanderplan.audit_log (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     TEXT        REFERENCES wanderplan.users(id) ON DELETE SET NULL,
    action      TEXT        NOT NULL,
    entity_type TEXT        NOT NULL,
    entity_id   TEXT        NOT NULL,
    metadata    JSONB       NOT NULL DEFAULT '{}',
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_user_id    ON wanderplan.audit_log(user_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_entity     ON wanderplan.audit_log(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON wanderplan.audit_log(created_at DESC);
