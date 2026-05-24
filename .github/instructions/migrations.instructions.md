---
applyTo: "migrations/**/*.sql"
---

# Database Migration Rules

## Naming
- Files use sequential four-digit prefixes: `000001_init.up.sql` / `000001_init.down.sql`.
- The description part is snake_case and describes the change, not the table (e.g. `add_places_cache`, not `places_cache_table`).
- Every `.up.sql` must have a matching `.down.sql` that fully reverses the change.

## Schema
- All tables live in the `wanderplan` schema. Prefix every statement: `CREATE TABLE wanderplan.table_name`.
- Primary keys are `UUID` generated with `gen_random_uuid()`. Never use `SERIAL` or `BIGSERIAL`.
- Timestamps: `created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`, `updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()`.
- Add a trigger to auto-update `updated_at` on every table that has it.
- JSONB columns (`payload`, `photos`, `metadata`) must have a `jsonb_path_ops` GIN index if they will be queried.

## Safety
- Never `DROP TABLE` or `DROP COLUMN` in an up migration. Mark columns as deprecated with a comment and remove in a later dedicated migration after confirming no code reads them.
- Always wrap destructive changes (`ALTER TABLE`, `DROP INDEX`) in a transaction.
- Add indexes for all foreign key columns and all columns used in `WHERE` clauses.

## Down migrations
- The `.down.sql` must be the exact inverse of `.up.sql` — drop what was created, restore what was altered.
- Test every down migration locally with `make migrate-down` before committing.
