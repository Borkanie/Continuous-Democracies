#!/usr/bin/env bash
# ============================================================================
# Local CI equivalent of .github/workflows/backend-container.yml
#
# Steps mirrored from the workflow:
#   1. Build Docker image  (workflow: "Build and push Docker image")
#   2. Run the container   (workflow: "Deploy to Docker")
#   3. Run tests           (new — the workflow does not run tests yet)
#   4. Teardown            (handled automatically by conftest.py)
#
# Usage:
#   ./run_tests.sh                     # full run: build + test
#   ./run_tests.sh --no-docker         # skip build/start (API must already be up)
#   ./run_tests.sh -k test_health      # pass any pytest args through
# ============================================================================
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR"

echo "================================================================"
echo "  Continuous Democracies — API Integration Tests"
echo "  Mirrors: .github/workflows/backend-container.yml"
echo "================================================================"

# ── Prerequisite: .env.test ──────────────────────────────────────────────────
if [ ! -f ".env.test" ]; then
  echo ""
  echo "ERROR: Testing/.env.test not found."
  echo "       Copy env.test.example to .env.test and fill in your DB credentials."
  exit 1
fi

# ── Step 1: Install Python test dependencies ─────────────────────────────────
echo ""
echo "--- Installing dependencies ---"
pip3 install -q -r requirements.txt

# ── Step 2-4: Build image, start container, run tests, teardown ───────────────
# conftest.py handles Docker lifecycle automatically unless --no-docker is passed.
echo ""
echo "--- Running tests ---"
python3 -m pytest "$@"

echo ""
echo "Done."
