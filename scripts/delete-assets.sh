#!/usr/bin/env bash
# Deletes every asset from the database, for starting over during development.
# Uses DATABASE_URL, or the Docker Compose database if it is not set.
set -euo pipefail

url="${DATABASE_URL:-postgres://smartgrid:smartgrid@localhost:5432/smartgrid?sslmode=disable}"

# Without psql on this machine, the one inside the database container is used.
run_sql() {
  local options=(--quiet --no-psqlrc --set ON_ERROR_STOP=1 --tuples-only --no-align)
  if command -v psql >/dev/null 2>&1; then
    psql "$url" "${options[@]}" "$@"
  elif [ -z "${DATABASE_URL:-}" ]; then
    (cd "$(dirname "$0")/.." && docker compose exec -T db psql --username smartgrid --dbname smartgrid "${options[@]}" "$@")
  else
    echo "psql is not installed, so DATABASE_URL cannot be used." >&2
    exit 1
  fi
}

red=$'\033[1;31m'
plain=$'\033[0m'

# Shown without any password the URL may contain.
database="$(sed -E 's#//[^@/]*@#//#' <<<"$url")"
count="$(run_sql --command "SELECT count(*) FROM assets")"

echo "${red}WARNING: this script is for development only.${plain}"
echo "${red}It permanently deletes all $count assets from the database below.${plain}"
echo "${red}There is no undo and no backup is taken.${plain}"
echo
echo "Database: $database"
echo

# Anything except Y or y, including just pressing Enter, cancels.
read -r -p "Delete all assets? [y/N] " answer || answer=""
if [[ "$answer" != "Y" && "$answer" != "y" ]]; then
  echo "Cancelled. Nothing was deleted."
  exit 0
fi

run_sql --command "DELETE FROM assets"

echo "Deleted $count assets."
echo "The asset pages may still show them: click \"Refresh data\" there to clear the browser's copy."
