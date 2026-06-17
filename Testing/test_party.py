"""
Tests for PartyController
  GET /api/party/all?active={bool}&number={int}
  GET /api/party/GetById/?id={guid}
  GET /api/party/query?name={str}&acronym={str}

Note: all endpoints always return HTTP 200.
"Not found" cases return a plain JSON string, e.g. "No parties found."
Real data cases return a JSON object or array.
"""

import uuid
import requests
import pytest


# ── /api/party/all ───────────────────────────────────────────────────────────

class TestGetAllParties:

    def test_active_true_returns_200(self, api_url):
        r = requests.get(f"{api_url}/api/party/all", params={"active": "true"})
        assert r.status_code == 200

    def test_active_false_returns_200(self, api_url):
        r = requests.get(f"{api_url}/api/party/all", params={"active": "false"})
        assert r.status_code == 200

    def test_response_is_list_or_not_found_string(self, api_url):
        r = requests.get(f"{api_url}/api/party/all", params={"active": "true"})
        body = r.json()
        assert isinstance(body, (list, str))

    def test_number_param_limits_results(self, api_url, existing_party):
        if existing_party is None:
            pytest.skip("No active party data in DB")
        r = requests.get(f"{api_url}/api/party/all", params={"active": "true", "number": 2})
        data = r.json()
        assert isinstance(data, list)
        assert len(data) <= 2

    def test_party_shape(self, api_url, existing_party):
        if existing_party is None:
            pytest.skip("No active party data in DB")
        r = requests.get(f"{api_url}/api/party/all", params={"active": "true", "number": 5})
        parties = r.json()
        assert isinstance(parties, list)
        p = parties[0]
        assert "id" in p
        assert "name" in p
        assert "active" in p
        assert p["active"] is True

    def test_inactive_parties_have_active_false(self, api_url):
        r = requests.get(f"{api_url}/api/party/all", params={"active": "false"})
        body = r.json()
        if isinstance(body, str):
            return  # no inactive parties — acceptable
        for party in body:
            assert party["active"] is False

    def test_default_number_returns_at_most_100(self, api_url, existing_party):
        if existing_party is None:
            pytest.skip("No active party data in DB")
        r = requests.get(f"{api_url}/api/party/all", params={"active": "true"})
        body = r.json()
        if isinstance(body, list):
            assert len(body) <= 100


# ── /api/party/GetById ───────────────────────────────────────────────────────

class TestGetPartyById:

    def test_valid_id_returns_correct_party(self, api_url, existing_party):
        if existing_party is None:
            pytest.skip("No party data in DB")
        r = requests.get(f"{api_url}/api/party/GetById/", params={"id": existing_party["id"]})
        assert r.status_code == 200
        data = r.json()
        assert isinstance(data, dict)
        assert data["id"] == existing_party["id"]
        assert data["name"] == existing_party["name"]

    def test_unknown_id_returns_200_with_not_found_string(self, api_url):
        r = requests.get(f"{api_url}/api/party/GetById/", params={"id": str(uuid.uuid4())})
        assert r.status_code == 200
        assert isinstance(r.json(), str)
        assert "not found" in r.json().lower()


# ── /api/party/query ─────────────────────────────────────────────────────────

class TestQueryParty:

    def test_query_by_exact_name(self, api_url, existing_party):
        if existing_party is None:
            pytest.skip("No party data in DB")
        r = requests.get(f"{api_url}/api/party/query", params={"name": existing_party["name"]})
        assert r.status_code == 200
        data = r.json()
        assert isinstance(data, dict)
        assert data["name"] == existing_party["name"]

    def test_query_by_acronym(self, api_url, existing_party):
        if existing_party is None or not existing_party.get("acronym"):
            pytest.skip("No party with acronym in DB")
        r = requests.get(f"{api_url}/api/party/query", params={"acronym": existing_party["acronym"]})
        assert r.status_code == 200
        data = r.json()
        assert isinstance(data, dict)
        assert data["acronym"] == existing_party["acronym"]

    def test_query_nonexistent_name_returns_200_with_string(self, api_url):
        r = requests.get(f"{api_url}/api/party/query", params={"name": "ZZZ_NONEXISTENT_PARTY_999"})
        assert r.status_code == 200
        assert isinstance(r.json(), str)

    def test_query_both_name_and_acronym(self, api_url, existing_party):
        if existing_party is None or not existing_party.get("acronym"):
            pytest.skip("No party with acronym in DB")
        r = requests.get(
            f"{api_url}/api/party/query",
            params={"name": existing_party["name"], "acronym": existing_party["acronym"]},
        )
        assert r.status_code == 200
        assert isinstance(r.json(), dict)
