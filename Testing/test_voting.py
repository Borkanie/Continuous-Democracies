"""
Tests for VotingController
  GET /api/voting/getAllRounds?startDate=&endDate=&keywords=&maxNumberOfEntries=
  GET /api/voting/getRoundById/?voteNumber={int}
  GET /api/voting/GetResultForVote/?number={int}&partyId={guid}&partyAcronim={str}
  GET /api/voting/GetAllVotesForARoundById/?roundId={guid}

Round shape:  id, name/title, description, voteDate, voteId
Vote shape:   id, name, position (Yes|No|Abstain|Absent), politician, round
"""

import uuid
import requests
import pytest


# ── /api/voting/getAllRounds ──────────────────────────────────────────────────

class TestGetAllRounds:

    def test_no_params_returns_200(self, api_url):
        r = requests.get(f"{api_url}/api/voting/getAllRounds")
        assert r.status_code == 200

    def test_response_is_list_or_not_found_string(self, api_url):
        r = requests.get(f"{api_url}/api/voting/getAllRounds")
        assert isinstance(r.json(), (list, str))

    def test_max_entries_limits_results(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        r = requests.get(f"{api_url}/api/voting/getAllRounds", params={"maxNumberOfEntries": 2})
        body = r.json()
        assert isinstance(body, list)
        assert len(body) <= 2

    def test_round_shape(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        r = requests.get(f"{api_url}/api/voting/getAllRounds", params={"maxNumberOfEntries": 1})
        rounds = r.json()
        assert isinstance(rounds, list)
        rnd = rounds[0]
        assert "id" in rnd
        assert "name" in rnd
        assert "description" in rnd
        assert "voteDate" in rnd
        assert "voteId" in rnd

    def test_date_range_filter(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        # Use a wide window that should include the existing round
        r = requests.get(
            f"{api_url}/api/voting/getAllRounds",
            params={
                "startDate": "2000-01-01",
                "endDate": "2100-12-31",
                "maxNumberOfEntries": 10,
            },
        )
        assert r.status_code == 200
        body = r.json()
        assert isinstance(body, (list, str))

    def test_future_date_range_returns_200(self, api_url):
        r = requests.get(
            f"{api_url}/api/voting/getAllRounds",
            params={"startDate": "2099-01-01", "endDate": "2099-12-31"},
        )
        assert r.status_code == 200

    def test_keyword_filter_returns_200(self, api_url):
        r = requests.get(
            f"{api_url}/api/voting/getAllRounds",
            params={"keywords": "lege", "maxNumberOfEntries": 5},
        )
        assert r.status_code == 200

    def test_multiple_keywords_returns_200(self, api_url):
        # ASP.NET binds repeated params to string?[] keywords
        r = requests.get(
            f"{api_url}/api/voting/getAllRounds",
            params=[("keywords", "lege"), ("keywords", "buget"), ("maxNumberOfEntries", "5")],
        )
        assert r.status_code == 200


# ── /api/voting/getRoundById ─────────────────────────────────────────────────

class TestGetRoundById:

    def test_valid_vote_number_returns_round(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        vote_id = existing_round["voteId"]
        r = requests.get(f"{api_url}/api/voting/getRoundById/", params={"voteNumber": vote_id})
        assert r.status_code == 200
        data = r.json()
        assert isinstance(data, dict)
        assert data["voteId"] == vote_id

    def test_nonexistent_vote_number_returns_200_with_string(self, api_url):
        r = requests.get(f"{api_url}/api/voting/getRoundById/", params={"voteNumber": 999999999})
        assert r.status_code == 200
        assert isinstance(r.json(), str)


# ── /api/voting/GetResultForVote ─────────────────────────────────────────────

class TestGetResultForVote:

    def test_valid_number_returns_200(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        r = requests.get(f"{api_url}/api/voting/GetResultForVote/", params={"number": existing_round["voteId"]})
        assert r.status_code == 200

    def test_response_is_list_or_not_found_string(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        r = requests.get(f"{api_url}/api/voting/GetResultForVote/", params={"number": existing_round["voteId"]})
        assert isinstance(r.json(), (list, str))

    def test_vote_shape(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        r = requests.get(f"{api_url}/api/voting/GetResultForVote/", params={"number": existing_round["voteId"]})
        body = r.json()
        if not isinstance(body, list) or not body:
            pytest.skip("No votes for this round")
        vote = body[0]
        assert "id" in vote
        assert "position" in vote
        assert "politician" in vote

    def test_filter_by_party_id(self, api_url, existing_round, existing_politician):
        if existing_round is None or existing_politician is None:
            pytest.skip("No data in DB")
        party = existing_politician.get("party")
        if not party:
            pytest.skip("Existing politician has no party")
        r = requests.get(
            f"{api_url}/api/voting/GetResultForVote/",
            params={"number": existing_round["voteId"], "partyId": party["id"]},
        )
        assert r.status_code == 200

    def test_filter_by_party_acronym(self, api_url, existing_round, existing_politician):
        if existing_round is None or existing_politician is None:
            pytest.skip("No data in DB")
        acronym = existing_politician.get("party", {}).get("acronym")
        if not acronym:
            pytest.skip("Existing politician has no party acronym")
        r = requests.get(
            f"{api_url}/api/voting/GetResultForVote/",
            params={"number": existing_round["voteId"], "partyAcronim": acronym},
        )
        assert r.status_code == 200

    def test_nonexistent_round_returns_200_with_string(self, api_url):
        r = requests.get(f"{api_url}/api/voting/GetResultForVote/", params={"number": 999999999})
        assert r.status_code == 200
        assert isinstance(r.json(), str)


# ── /api/voting/GetAllVotesForARoundById ─────────────────────────────────────

class TestGetAllVotesForARoundById:

    def test_valid_round_id_returns_200(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        r = requests.get(
            f"{api_url}/api/voting/GetAllVotesForARoundById/",
            params={"roundId": existing_round["id"]},
        )
        assert r.status_code == 200

    def test_response_is_list_or_not_found_string(self, api_url, existing_round):
        if existing_round is None:
            pytest.skip("No voting round data in DB")
        r = requests.get(
            f"{api_url}/api/voting/GetAllVotesForARoundById/",
            params={"roundId": existing_round["id"]},
        )
        assert isinstance(r.json(), (list, str))

    def test_unknown_round_id_returns_200_with_string(self, api_url):
        r = requests.get(
            f"{api_url}/api/voting/GetAllVotesForARoundById/",
            params={"roundId": str(uuid.uuid4())},
        )
        assert r.status_code == 200
        assert isinstance(r.json(), str)
