#!/usr/bin/env bash
# Runs the Go backend (port 8080) and the Next.js frontend (port 3000) together.
# Press Ctrl+C to stop both.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"

# Stop everything this script started when it exits.
trap 'kill 0' EXIT

(cd "$root/backend" && go run .) &
(cd "$root/frontend" && npm run dev) &

# Exit as soon as either one stops, which also stops the other.
wait -n
