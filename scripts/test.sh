#!/usr/bin/env bash
# Runs the tests for the frontend and the backend and prints a summary.
set -uo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"

echo "Frontend: no tests yet."
echo

echo "Running backend tests..."
output="$(cd "$root/backend" && go test -v ./... 2>&1)"
status=$?

# Count the top-level Test functions; each one may contain several subtests.
passed="$(grep -c '^--- PASS' <<<"$output" || true)"
failed="$(grep -c '^--- FAIL' <<<"$output" || true)"

if [ "$status" -eq 0 ]; then
  grep '^--- ' <<<"$output"
else
  echo "$output" # show everything so the failure can be diagnosed
fi

echo
echo "Summary"
echo "  Frontend: 0/0 tests passed"
echo "  Backend:  $passed/$((passed + failed)) tests passed"

exit "$status"
