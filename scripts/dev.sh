#!/usr/bin/env bash
# Runs the Go backend (port 8080) and the Next.js frontend (port 3000) together.
# Press Ctrl+C to stop both.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"

# The backend needs the database. This starts the Docker one if it is not
# already running, e.g. after the machine was restarted. With DATABASE_URL
# set, the database is another one and is left alone.
if [ -z "${DATABASE_URL:-}" ]; then
  (cd "$root" && docker compose up --detach --wait db)
fi

# Stop everything this script started when it exits.
trap 'kill 0' EXIT

(cd "$root/backend" && go run .) &
(cd "$root/frontend" && npm run dev) &

# Exit as soon as either one stops, which also stops the other.
wait -n
