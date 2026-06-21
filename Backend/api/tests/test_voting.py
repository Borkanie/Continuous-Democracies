import uuid
from unittest.mock import patch

from tests.conftest import make_round, make_vote, make_politician, make_party


class TestGetAllRounds:
    def test_returns_list(self, client):
        round_ = make_round(Name="Vot lege", VoteId=200)
        with patch("app.routers.voting.VotingRoundService") as MockService:
            MockService.return_value.get_all_rounds.return_value = [round_]
            response = client.get("/api/voting/getAllRounds")

        assert response.status_code == 200
        data = response.json()
        assert isinstance(data, list)
        assert data[0]["VoteId"] == 200

    def test_empty_returns_string(self, client):
        with patch("app.routers.voting.VotingRoundService") as MockService:
            MockService.return_value.get_all_rounds.return_value = []
            response = client.get("/api/voting/getAllRounds")

        assert response.status_code == 200
        assert "No results" in response.json()

    def test_keyword_param_forwarded(self, client):
        with patch("app.routers.voting.VotingRoundService") as MockService:
            MockService.return_value.get_all_rounds.return_value = []
            client.get("/api/voting/getAllRounds?keywords=buget&keywords=pensii")
            call_kwargs = MockService.return_value.get_all_rounds.call_args.kwargs
            assert call_kwargs["keywords"] == ["buget", "pensii"]


class TestGetRoundById:
    def test_found(self, client):
        round_ = make_round(VoteId=123)
        with patch("app.routers.voting.VotingRoundService") as MockService:
            MockService.return_value.get_by_vote_id.return_value = round_
            response = client.get("/api/voting/getRoundById/?voteNumber=123")

        assert response.status_code == 200
        assert response.json()["VoteId"] == 123

    def test_not_found(self, client):
        with patch("app.routers.voting.VotingRoundService") as MockService:
            MockService.return_value.get_by_vote_id.return_value = None
            response = client.get("/api/voting/getRoundById/?voteNumber=9999")

        assert response.status_code == 200
        assert "not found" in response.json().lower()


class TestGetResultForVote:
    def test_returns_votes(self, client):
        vote = make_vote()
        with patch("app.routers.voting.VotingService") as MockService:
            MockService.return_value.get_all_votes_for_round_by_vote_id.return_value = [vote]
            response = client.get("/api/voting/GetResultForVote/?number=100")

        assert response.status_code == 200
        assert isinstance(response.json(), list)
        assert len(response.json()) == 1

    def test_filter_by_party_id(self, client):
        party_id = uuid.uuid4()
        matching_party = make_party(Id=party_id)
        matching_pol = make_politician(party=matching_party)
        matching_pol.PartyId = party_id
        other_party = make_party(Id=uuid.uuid4())
        other_pol = make_politician(party=other_party)
        other_pol.PartyId = other_party.Id

        vote_match = make_vote(politician=matching_pol)
        vote_other = make_vote(politician=other_pol)

        with patch("app.routers.voting.VotingService") as MockService:
            MockService.return_value.get_all_votes_for_round_by_vote_id.return_value = [
                vote_match, vote_other
            ]
            response = client.get(f"/api/voting/GetResultForVote/?number=100&partyId={party_id}")

        assert response.status_code == 200
        assert len(response.json()) == 1

    def test_filter_by_party_acronym(self, client):
        psd_party = make_party(Acronym="PSD")
        psd_pol = make_politician(party=psd_party)
        psd_pol.party_rel = psd_party

        other_party = make_party(Acronym="PNL")
        other_pol = make_politician(party=other_party)
        other_pol.party_rel = other_party

        vote_psd = make_vote(politician=psd_pol)
        vote_other = make_vote(politician=other_pol)

        with patch("app.routers.voting.VotingService") as MockService:
            MockService.return_value.get_all_votes_for_round_by_vote_id.return_value = [
                vote_psd, vote_other
            ]
            response = client.get("/api/voting/GetResultForVote/?number=100&partyAcronim=PSD")

        assert response.status_code == 200
        assert len(response.json()) == 1

    def test_empty_returns_empty_list(self, client):
        with patch("app.routers.voting.VotingService") as MockService:
            MockService.return_value.get_all_votes_for_round_by_vote_id.return_value = []
            response = client.get("/api/voting/GetResultForVote/?number=9999")

        assert response.status_code == 200
        assert response.json() == []


class TestGetAllVotesForARoundById:
    def test_returns_votes(self, client):
        round_id = uuid.uuid4()
        vote = make_vote()
        with patch("app.routers.voting.VotingService") as MockService:
            MockService.return_value.get_all_votes_for_round_by_guid.return_value = [vote]
            response = client.get(f"/api/voting/GetAllVotesForARoundById/?roundId={round_id}")

        assert response.status_code == 200
        assert isinstance(response.json(), list)

    def test_empty_returns_empty_list(self, client):
        round_id = uuid.uuid4()
        with patch("app.routers.voting.VotingService") as MockService:
            MockService.return_value.get_all_votes_for_round_by_guid.return_value = []
            response = client.get(f"/api/voting/GetAllVotesForARoundById/?roundId={round_id}")

        assert response.status_code == 200
        assert response.json() == []
