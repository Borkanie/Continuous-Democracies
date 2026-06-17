"""Tests for GET /health"""

import requests


def test_health_status_200(api_url):
    r = requests.get(f"{api_url}/health")
    assert r.status_code == 200


def test_health_body_is_healthy(api_url):
    r = requests.get(f"{api_url}/health")
    assert r.text.strip().lower() == "healthy"
