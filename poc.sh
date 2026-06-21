#!/usr/bin/env bash
# POC: import a single law via the Go scraper into an ephemeral MongoDB,
# then run the BE integration tests against that data.
#
# Usage: ./poc.sh [LAW_ID]
#   LAW_ID defaults to 40000 if not supplied.
#   WG_CONFIGS_DIR env var sets the WireGuard config directory (default /etc/wireguard/configs).
#
# Example: WG_CONFIGS_DIR=~/wg-configs ./poc.sh 37113

set -euo pipefail

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
LAW_ID=${1:-40000}

echo "==> POC: importing law $LAW_ID via Go scraper"
echo ""

# Bring up ephemeral MongoDB and run the scraper against it.
# --abort-on-container-exit stops everything once the scraper exits.
# --exit-code-from scraper-poc propagates the scraper's exit code.
IMPORT_LAW_ID="$LAW_ID" docker compose \
  -f docker-compose.poc.yml \
  up \
  --build \
  --abort-on-container-exit \
  --exit-code-from scraper-poc

echo ""
echo "==> Scraper finished. Running BE POC integration tests against mongodb-poc (port 27018)..."
echo ""

cd "$ROOT_DIR/new-backend/api"

TEST_MONGO_URI="mongodb://localhost:27018" \
TEST_LAW_ID="$LAW_ID" \
  go test ./tests/poc/... -v -timeout 60s

echo ""
echo "==> Tests complete. Tearing down poc containers..."
docker compose -f "$ROOT_DIR/docker-compose.poc.yml" down -v

echo ""
echo "==> POC done."
