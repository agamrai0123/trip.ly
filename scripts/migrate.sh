#!/usr/bin/env bash
# scripts/migrate.sh — Run golang-migrate against DATABASE_URL.
#
# Usage:
#   ./scripts/migrate.sh up          # apply all pending migrations
#   ./scripts/migrate.sh down        # roll back one step
#   ./scripts/migrate.sh down N      # roll back N steps
#   ./scripts/migrate.sh version     # print current schema version
#   ./scripts/migrate.sh force N     # force schema version (use after manual fixes)
#
# Requires DATABASE_URL or individual DB_* env vars.
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
MIGRATIONS_DIR="$ROOT_DIR/migrations"

# Build DATABASE_URL from parts if not set explicitly.
if [[ -z "${DATABASE_URL:-}" ]]; then
  : "${DB_HOST:=localhost}"
  : "${DB_PORT:=5432}"
  : "${DB_NAME:=wanderplan}"
  : "${DB_USER:=postgres}"
  : "${DB_PASSWORD:=postgres}"
  DATABASE_URL="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
fi

COMMAND="${1:-up}"
EXTRA="${2:-}"

echo "▶  migrate $COMMAND $EXTRA  →  $MIGRATIONS_DIR"
echo "   target: ${DATABASE_URL//:*@/:***@}"   # mask password in output

migrate \
  -path "$MIGRATIONS_DIR" \
  -database "$DATABASE_URL" \
  "$COMMAND" $EXTRA

echo "✅  Migration $COMMAND complete."
