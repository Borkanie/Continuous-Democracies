"""
Tests for PoliticiansController
  GET /api/politicians/GetById/?id={guid}
  GET /api/politicians/GetByName/?name={str}
  GET /api/politicians/getAllPoliticians?partyAcronym=&partyName=&isActive=&location=&gender=&number=

Enums (passed as string names in query params):
  Gender:       Male | Female | Other
  WorkLocation: Parliament | Senate | Other
"""

import uuid
import requests
import pytest


# ── /api/politicians/getAllPoliticians ───────────────────────────────────────

class TestGetAllPoliticians:

    def test_no_filters_returns_200(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians")
        assert r.status_code == 200

    def test_response_is_list_or_not_found_string(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians")
        assert isinstance(r.json(), (list, str))

    def test_number_limit_applied(self, api_url, existing_politician):
        if existing_politician is None:
            pytest.skip("No politician data in DB")
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"number": 3})
        body = r.json()
        assert isinstance(body, list)
        assert len(body) <= 3

    def test_politician_shape(self, api_url, existing_politician):
        if existing_politician is None:
            pytest.skip("No politician data in DB")
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"number": 1})
        politicians = r.json()
        assert isinstance(politicians, list)
        p = politicians[0]
        assert "id" in p
        assert "name" in p
        assert "active" in p
        assert "gender" in p
        assert "workLocation" in p

    def test_filter_active_true(self, api_url, existing_politician):
        if existing_politician is None:
            pytest.skip("No politician data in DB")
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"isActive": "true", "number": 10})
        body = r.json()
        if isinstance(body, list):
            for p in body:
                assert p["active"] is True

    def test_filter_active_false(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"isActive": "false", "number": 10})
        body = r.json()
        if isinstance(body, list):
            for p in body:
                assert p["active"] is False

    def test_filter_by_gender_male(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"gender": "Male", "number": 5})
        assert r.status_code == 200

    def test_filter_by_gender_female(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"gender": "Female", "number": 5})
        assert r.status_code == 200

    def test_filter_by_location_parliament(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"location": "Parliament", "number": 5})
        assert r.status_code == 200

    def test_filter_by_location_senate(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"location": "Senate", "number": 5})
        assert r.status_code == 200

    def test_filter_by_party_acronym(self, api_url, existing_politician):
        if existing_politician is None:
            pytest.skip("No politician data in DB")
        acronym = existing_politician.get("party", {}).get("acronym")
        if not acronym:
            pytest.skip("Existing politician has no party acronym")
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"partyAcronym": acronym, "number": 10})
        assert r.status_code == 200
        body = r.json()
        if isinstance(body, list):
            for p in body:
                assert p["party"]["acronym"] == acronym

    def test_filter_by_party_name(self, api_url, existing_politician):
        if existing_politician is None:
            pytest.skip("No politician data in DB")
        party_name = existing_politician.get("party", {}).get("name")
        if not party_name:
            pytest.skip("Existing politician has no party name")
        r = requests.get(f"{api_url}/api/politicians/getAllPoliticians", params={"partyName": party_name, "number": 10})
        assert r.status_code == 200

    def test_nonexistent_party_returns_200_with_string(self, api_url):
        r = requests.get(
            f"{api_url}/api/politicians/getAllPoliticians",
            params={"partyAcronym": "ZZZ_NONEXISTENT_999"},
        )
        assert r.status_code == 200
        assert isinstance(r.json(), str)


# ── /api/politicians/GetById ─────────────────────────────────────────────────

class TestGetPoliticianById:

    def test_valid_id_returns_politician(self, api_url, existing_politician):
        if existing_politician is None:
            pytest.skip("No politician data in DB")
        r = requests.get(f"{api_url}/api/politicians/GetById/", params={"id": existing_politician["id"]})
        assert r.status_code == 200
        data = r.json()
        assert isinstance(data, dict)
        assert data["id"] == existing_politician["id"]

    def test_unknown_id_returns_200_with_not_found_string(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/GetById/", params={"id": str(uuid.uuid4())})
        assert r.status_code == 200
        assert isinstance(r.json(), str)


# ── /api/politicians/GetByName ───────────────────────────────────────────────

class TestGetPoliticianByName:

    def test_valid_name_returns_politician(self, api_url, existing_politician):
        if existing_politician is None:
            pytest.skip("No politician data in DB")
        r = requests.get(f"{api_url}/api/politicians/GetByName/", params={"name": existing_politician["name"]})
        assert r.status_code == 200

    def test_unknown_name_returns_200(self, api_url):
        r = requests.get(f"{api_url}/api/politicians/GetByName/", params={"name": "ZZZ_Nonexistent_Person_999"})
        assert r.status_code == 200
