"""
Session-scoped fixtures that mirror .github/workflows/backend-container.yml:
  1. Build the Docker image from Backend/ParliamentMonitor/ (same build context as CI)
  2. Start the container with the same env vars the workflow injects from secrets
  3. Poll /health until the API is ready
  4. Yield the base URL to all tests
  5. Remove the container on teardown

Usage:
  pytest                          # build image, run tests, teardown
  pytest --no-docker              # assume API is already on --port (default 18080)
  pytest --no-docker --port 8080  # point at a running instance on a different port
"""

import os
import time
import subprocess
from pathlib import Path

import requests
import pytest
from dotenv import load_dotenv

TESTS_DIR    = Path(__file__).parent
PROJECT_ROOT = TESTS_DIR.parent
BUILD_CONTEXT = PROJECT_ROOT / "Backend" / "ParliamentMonitor"

CONTAINER_NAME = "cd-api-test"
DEFAULT_PORT   = 18080
IMAGE_TAG      = "cd-api-test:local"
HEALTH_TIMEOUT = 120  # seconds


# ── CLI options ──────────────────────────────────────────────────────────────

def pytest_addoption(parser):
    parser.addoption(
        "--no-docker",
        action="store_true",
        default=False,
        help="Skip Docker build/start; assume the API is already running.",
    )
    parser.addoption(
        "--port",
        type=int,
        default=DEFAULT_PORT,
        help=f"Port the API listens on when --no-docker is set (default {DEFAULT_PORT}).",
    )


# ── Helpers ──────────────────────────────────────────────────────────────────

def _load_env():
    env_file = TESTS_DIR / ".env.test"
    if env_file.exists():
        load_dotenv(env_file)
    missing = [k for k in ("DBSERVERADDRESS", "DBNAME", "DBUSER", "DBPASSWORD") if not os.environ.get(k)]
    if missing:
        pytest.exit(
            f"Missing required env vars: {', '.join(missing)}\n"
            f"Copy Testing/.env.test.example to Testing/.env.test and fill in your values."
        )


def _run(*cmd, check=False):
    return subprocess.run(list(cmd), capture_output=True, text=True, check=check)


def _build_image():
    print(f"\n[docker] Building {IMAGE_TAG} from {BUILD_CONTEXT} ...")
    # Mirrors CI step: cd Backend/ParliamentMonitor && docker build -t ... .
    subprocess.run(
        ["docker", "build", "-t", IMAGE_TAG, "."],
        cwd=BUILD_CONTEXT,
        check=True,
    )
    print(f"[docker] Image {IMAGE_TAG} built.")


def _start_container():
    db_host = os.environ["DBSERVERADDRESS"]
    # On macOS/Linux, containers can't reach 127.0.0.1 on the host directly
    if db_host in ("localhost", "127.0.0.1"):
        db_host = "host.docker.internal"

    _run("docker", "rm", "-f", CONTAINER_NAME)

    subprocess.run(
        [
            "docker", "run", "-d",
            "--name", CONTAINER_NAME,
            "-p",  f"{DEFAULT_PORT}:8080",
            "-e",  f"DBSERVERADDRESS={db_host}",
            "-e",  f"DBNAME={os.environ['DBNAME']}",
            "-e",  f"DBUSER={os.environ['DBUSER']}",
            "-e",  f"DBPASSWORD={os.environ['DBPASSWORD']}",
            "-e",  f"LOG_LEVEL={os.environ.get('LOG_LEVEL', 'Warning')}",
            "-e",  "ASPNETCORE_ENVIRONMENT=Development",
            "--add-host", "host.docker.internal:host-gateway",
            IMAGE_TAG,
        ],
        check=True,
    )
    print(f"[docker] Container {CONTAINER_NAME} started on port {DEFAULT_PORT}.")


# ── Primary fixture ───────────────────────────────────────────────────────────

@pytest.fixture(scope="session")
def api_url(request):
    _load_env()

    no_docker = request.config.getoption("--no-docker")
    port      = request.config.getoption("--port") if no_docker else DEFAULT_PORT

    if not no_docker:
        _build_image()
        _start_container()

    base_url = f"http://localhost:{port}"
    print(f"\n[health] Polling {base_url}/health (timeout {HEALTH_TIMEOUT}s) ...")

    deadline = time.time() + HEALTH_TIMEOUT
    while time.time() < deadline:
        try:
            r = requests.get(f"{base_url}/health", timeout=2)
            if r.status_code == 200:
                elapsed = HEALTH_TIMEOUT - (deadline - time.time())
                print(f"[health] API ready after ~{elapsed:.0f}s")
                break
        except Exception:
            pass
        time.sleep(2)
    else:
        if not no_docker:
            logs = _run("docker", "logs", CONTAINER_NAME)
            _run("docker", "rm", "-f", CONTAINER_NAME)
            pytest.fail(
                f"API did not become healthy within {HEALTH_TIMEOUT}s.\n"
                f"--- container stdout ---\n{logs.stdout}\n"
                f"--- container stderr ---\n{logs.stderr}"
            )
        else:
            pytest.fail(f"API not reachable at {base_url}/health after {HEALTH_TIMEOUT}s.")

    yield base_url

    if not no_docker:
        _run("docker", "rm", "-f", CONTAINER_NAME)
        print(f"\n[docker] Container {CONTAINER_NAME} removed.")


# ── Seed fixtures (one real record each, used by ID-based tests) ─────────────

@pytest.fixture(scope="session")
def existing_party(api_url):
    """First active party in the DB, or None if empty."""
    r = requests.get(f"{api_url}/api/party/all", params={"active": "true", "number": 1})
    data = r.json()
    return data[0] if isinstance(data, list) and data else None


@pytest.fixture(scope="session")
def existing_politician(api_url):
    """First politician in the DB, or None if empty."""
    r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"number": 1})
    data = r.json()
    return data[0] if isinstance(data, list) and data else None


@pytest.fixture(scope="session")
def existing_round(api_url):
    """First voting round in the DB, or None if empty."""
    r = requests.get(f"{api_url}/api/voting/getAllRounds", params={"maxNumberOfEntries": 1})
    data = r.json()
    return data[0] if isinstance(data, list) and data else None
