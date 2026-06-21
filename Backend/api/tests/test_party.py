import uuid
from unittest.mock import patch

from tests.conftest import make_party


class TestGetAllParties:
    def test_returns_party_list(self, client):
        party = make_party(Name="PSD", Acronym="PSD", Active=True)
        with patch("app.routers.party.PartyService") as MockService:
            MockService.return_value.get_all_parties.return_value = [party]
            response = client.get("/api/party/all")

        assert response.status_code == 200
        data = response.json()
        assert isinstance(data, list)
        assert len(data) == 1
        assert data[0]["Name"] == "PSD"
        assert data[0]["Acronym"] == "PSD"

    def test_active_false_param_passed(self, client):
        with patch("app.routers.party.PartyService") as MockService:
            MockService.return_value.get_all_parties.return_value = []
            client.get("/api/party/all?active=false")
            MockService.return_value.get_all_parties.assert_called_once_with(is_active=False, number=100)

    def test_empty_result_returns_string(self, client):
        with patch("app.routers.party.PartyService") as MockService:
            MockService.return_value.get_all_parties.return_value = []
            response = client.get("/api/party/all")

        assert response.status_code == 200
        assert response.json() == "No parties found"

    def test_number_param(self, client):
        with patch("app.routers.party.PartyService") as MockService:
            MockService.return_value.get_all_parties.return_value = []
            client.get("/api/party/all?number=5")
            MockService.return_value.get_all_parties.assert_called_once_with(is_active=True, number=5)


class TestGetPartyById:
    def test_found(self, client):
        party_id = uuid.uuid4()
        party = make_party(Id=party_id, Name="PNL")
        with patch("app.routers.party.PartyService") as MockService:
            MockService.return_value.get_by_id.return_value = party
            response = client.get(f"/api/party/GetById/?id={party_id}")

        assert response.status_code == 200
        assert response.json()["Name"] == "PNL"

    def test_not_found_returns_string(self, client):
        with patch("app.routers.party.PartyService") as MockService:
            MockService.return_value.get_by_id.return_value = None
            response = client.get(f"/api/party/GetById/?id={uuid.uuid4()}")

        assert response.status_code == 200
        assert response.json() == "Party not found"


class TestQueryParty:
    def test_by_acronym(self, client):
        party = make_party(Acronym="USR")
        with patch("app.routers.party.PartyService") as MockService:
            MockService.return_value.get_party.return_value = party
            response = client.get("/api/party/query?acronym=USR")

        assert response.status_code == 200
        assert response.json()["Acronym"] == "USR"

    def test_not_found(self, client):
        with patch("app.routers.party.PartyService") as MockService:
            MockService.return_value.get_party.return_value = None
            response = client.get("/api/party/query?name=NOBODY")

        assert response.status_code == 200
        assert response.json() == "Party not found"
