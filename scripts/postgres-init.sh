#!/bin/bash
# Runs only *.up.sql migration files in order.
# Mounted at /docker-entrypoint-initdb.d/init.sh in the postgres container.
set -e

echo "==> Running WanderPlan migrations..."
for f in $(ls /migrations/*up.sql | sort); do
    echo "  --> $f"
    psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -f "$f"
done
echo "==> Migrations complete."
