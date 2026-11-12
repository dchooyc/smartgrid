#!/usr/bin/env bash
# Applies the SQL migrations in backend/migrations that have not run yet.
# Uses DATABASE_URL, or the local smartgrid database if it is not set.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
url="${DATABASE_URL:-postgres://localhost:5432/smartgrid?sslmode=disable}"

run_sql() {
  psql "$url" --quiet --no-psqlrc --set ON_ERROR_STOP=1 "$@"
}

# Records which migration files have been applied.
run_sql --command "CREATE TABLE IF NOT EXISTS schema_migrations (
  filename   TEXT PRIMARY KEY,
  applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
)"

for file in "$root"/backend/migrations/*.sql; do
  name="$(basename "$file")"

  applied="$(run_sql --tuples-only --no-align \
    --command "SELECT 1 FROM schema_migrations WHERE filename = '$name'")"
  if [ "$applied" = "1" ]; then
    echo "already applied: $name"
    continue
  fi

  # The migration and its record are applied together or not at all.
  run_sql --single-transaction --file "$file" \
    --command "INSERT INTO schema_migrations (filename) VALUES ('$name')"
  echo "applied: $name"
done
