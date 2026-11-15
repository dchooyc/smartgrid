#!/usr/bin/env bash
# Prepares this machine for local development: checks that the tools are
# installed, starts the Docker database, applies the migrations and installs
# the dependencies. Safe to run again. Afterwards, start the app with
# scripts/dev.sh.
#
# It installs nothing on the system itself: if a tool is missing it says how
# to install it and stops.
set -euo pipefail

root="$(cd "$(dirname "$0")/.." && pwd)"
min_node=20 # the oldest Node.js major version that Next.js supports

# install_hint TOOL: how to install it on this kind of machine.
install_hint() {
  local tool="$1" family=""
  if [ "$(uname)" = "Darwin" ]; then
    family="mac"
  elif [ -r /etc/os-release ]; then
    family="$(. /etc/os-release && echo "${ID:-} ${ID_LIKE:-}")"
  fi

  case "$tool:$family" in
    docker:mac) echo "Install Docker Desktop: https://docs.docker.com/desktop/setup/install/mac-install/" ;;
    go:mac) echo "brew install go" ;;
    node:mac) echo "brew install node" ;;
    docker:*) echo "Install Docker Engine with the Compose plugin: https://docs.docker.com/engine/install/" ;;
    go:*arch*) echo "sudo pacman -S go" ;;
    node:*arch*) echo "sudo pacman -S nodejs npm" ;;
    go:*debian* | go:*ubuntu* | go:*fedora* | go:*rhel*) echo "Install Go from https://go.dev/doc/install (the system package is often too old)" ;;
    node:*debian* | node:*ubuntu*) echo "Install Node.js $min_node or newer: https://nodejs.org/en/download" ;;
    node:*fedora* | node:*rhel*) echo "sudo dnf install nodejs" ;;
    go:*) echo "Install Go: https://go.dev/doc/install" ;;
    node:*) echo "Install Node.js $min_node or newer: https://nodejs.org/en/download" ;;
  esac
}

missing=0
need() {
  local tool="$1" command="$2"
  if command -v "$command" >/dev/null 2>&1; then
    return 0
  fi
  echo "Missing: $tool. $(install_hint "$tool")"
  missing=1
  return 1
}

echo "Checking the tools..."
need docker docker || true
need go go || true
if need node node; then
  node_major="$(node --version | sed -E 's/^v([0-9]+).*/\1/')"
  if [ "$node_major" -lt "$min_node" ]; then
    echo "Node.js $(node --version) is too old. $(install_hint node)"
    missing=1
  fi
  need node npm || true
fi
if [ "$missing" -ne 0 ]; then
  echo
  echo "Install what is missing, then run this script again."
  exit 1
fi
if ! docker compose version >/dev/null 2>&1; then
  echo "Docker is installed but its Compose plugin is not. $(install_hint docker)"
  exit 1
fi
if ! docker info >/dev/null 2>&1; then
  echo "Docker is installed but cannot be reached. Start Docker, and on Linux make sure"
  echo "your user may use it (https://docs.docker.com/engine/install/linux-postinstall/)."
  exit 1
fi
echo "  docker $(docker --version | sed -E 's/^Docker version ([^,]+).*/\1/'), $(go version | cut -d' ' -f3), node $(node --version)"

cd "$root"

echo
echo "Starting the database..."
if ! docker compose up --detach --wait db; then
  echo
  echo "The database did not start. If port 5432 is already in use, e.g. by a local"
  echo "PostgreSQL, stop that, or pick another port for this one:"
  echo "  echo DB_PORT=5433 > .env    (Docker Compose reads this file)"
  echo "  export DATABASE_URL='postgres://smartgrid:smartgrid@localhost:5433/smartgrid?sslmode=disable'"
  echo "then run this script again, and keep DATABASE_URL set when running the others."
  exit 1
fi

echo
echo "Applying the migrations..."
docker compose run --rm migrate

echo
echo "Installing the dependencies..."
(cd backend && go mod download)
(cd frontend && npm ci --no-audit --no-fund)

echo
echo "Ready. Start the app with: scripts/dev.sh"
echo "Then open http://localhost:3000"
