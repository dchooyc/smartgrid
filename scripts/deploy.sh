#!/usr/bin/env bash
# Updates the hosted app to the latest commit. Run it on the server, from
# anywhere inside the repository:
#
#   ./scripts/deploy.sh           # deploy if there is a new commit
#   ./scripts/deploy.sh --force   # deploy again even if nothing is new
#
# It pulls, runs the tests, builds the new version and only then swaps it in.
# If the tests or the build fail, the running site is left as it was.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
cd "$root"

compose=(docker compose -f compose.yaml -f compose.prod.yaml)

fail() {
  echo "deploy: $1" >&2
  exit 1
}

# --- 1. Safety checks ---
[ -f .env ] || fail "there is no .env file. Copy .env.example to .env and fill it in."
[ -z "$(git status --porcelain --untracked-files=no)" ] ||
  fail "files were changed on the server by hand. Undo them (git checkout .) and run this again."

# --- 2. Look for a new commit ---
branch="$(git rev-parse --abbrev-ref HEAD)"
git fetch --quiet origin "$branch"
old="$(git rev-parse HEAD)"
new="$(git rev-parse "origin/$branch")"

if [ "$old" = "$new" ] && [ "${1:-}" != "--force" ]; then
  echo "Already up to date at $(git log -1 --format='%h %s')."
  exit 0
fi

# --- 3. Pull it ---
git merge --quiet --ff-only "origin/$branch"
echo "Deploying $(git log -1 --format='%h %s')"

# Puts the code back as it was, so that it matches what is still running.
go_back() {
  git reset --quiet --hard "$old"
  fail "$1 The site was not changed and is still at $(git log -1 --format='%h %s')."
}

# --- 4. Run the tests ---
echo
echo "Running the tests..."
"${compose[@]}" run --rm --build test || go_back "The tests failed."

# --- 5. Build the new version, while the old one keeps running ---
echo
echo "Building..."
"${compose[@]}" build || go_back "The build failed."

# --- 6. Swap it in: applies new migrations and waits until it is healthy ---
echo
echo "Starting the new version..."
if ! "${compose[@]}" up --detach --wait; then
  fail "the new version did not start. See why with: ${compose[*]} logs
To go back to the previous version:
  git reset --hard $old && ${compose[*]} up --detach --build"
fi

# The web server only reads its configuration when it starts.
if ! git diff --quiet "$old" "$new" -- Caddyfile; then
  "${compose[@]}" restart caddy
fi

# --- 7. Confirm the site answers ---
# Without a password the site should ask for one, which is a 401.
site="$(grep '^SITE_ADDRESS=' .env | cut -d= -f2-)"
status=000
for _ in 1 2 3 4 5 6; do
  status="$(curl --silent --output /dev/null --max-time 10 --write-out '%{http_code}' "https://$site/" || true)"
  [ "$status" = "401" ] && break
  sleep 5
done
[ "$status" = "401" ] ||
  fail "https://$site/ answered $status, not the password prompt (401). See: ${compose[*]} logs caddy"

# --- 8. Tidy up ---
docker image prune --force > /dev/null

echo
echo "Deployed $(git log -1 --format='%h %s')"
echo "Live at https://$site/"
