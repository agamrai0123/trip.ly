-- 000001_init.down.sql
-- Drops all WanderPlan tables in reverse dependency order.

DROP TABLE IF EXISTS wanderplan.audit_log;
DROP TABLE IF EXISTS wanderplan.notifications;
DROP TABLE IF EXISTS wanderplan.places_cache;
DROP TABLE IF EXISTS wanderplan.trip_tags;
DROP TABLE IF EXISTS wanderplan.collaborators;
DROP TABLE IF EXISTS wanderplan.itinerary_items;
DROP TABLE IF EXISTS wanderplan.itinerary_days;
DROP TABLE IF EXISTS wanderplan.trips;
DROP TABLE IF EXISTS wanderplan.refresh_tokens;
DROP TABLE IF EXISTS wanderplan.users;

DROP SCHEMA IF EXISTS wanderplan;
