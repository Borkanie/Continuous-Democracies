import uuid
from unittest.mock import patch

from tests.conftest import make_politician, make_party


class TestGetAllPoliticians:
    def test_returns_list(self, client):
        politician = make_politician()
        with patch("app.routers.politicians.PoliticianService") as MockService:
            MockService.return_value.get_all_politicians.return_value = [politician]
            response = client.get("/api/politicians/getAllPoliticians")

        assert response.status_code == 200
        data = response.json()
        assert isinstance(data, list)
        assert len(data) == 1
        assert data[0]["Name"] == politician.Name

    def test_empty_returns_string(self, client):
        with patch("app.routers.politicians.PoliticianService") as MockService:
            MockService.return_value.get_all_politicians.return_value = []
            response = client.get("/api/politicians/getAllPoliticians")

        assert response.status_code == 200
        assert response.json() == "No politicians found"

    def test_filters_forwarded(self, client):
        with patch("app.routers.politicians.PoliticianService") as MockService:
            MockService.return_value.get_all_politicians.return_value = []
            client.get("/api/politicians/getAllPoliticians?partyAcronym=PSD&isActive=true&number=10")
            MockService.return_value.get_all_politicians.assert_called_once()
            call_kwargs = MockService.return_value.get_all_politicians.call_args.kwargs
            assert call_kwargs["party_acronym"] == "PSD"
            assert call_kwargs["is_active"] is True
            assert call_kwargs["number"] == 10

    def test_party_embedded_in_response(self, client):
        party = make_party(Name="PSD", Acronym="PSD")
        politician = make_politician(party=party, Name="Vasile")
        with patch("app.routers.politicians.PoliticianService") as MockService:
            MockService.return_value.get_all_politicians.return_value = [politician]
            response = client.get("/api/politicians/getAllPoliticians")

        data = response.json()
        assert data[0]["Party"]["Name"] == "PSD"


class TestGetPoliticianById:
    def test_found(self, client):
        pid = uuid.uuid4()
        politician = make_politician(Id=pid, Name="Maria Ionescu")
        with patch("app.routers.politicians.PoliticianService") as MockService:
            MockService.return_value.get_by_id.return_value = politician
            response = client.get(f"/api/politicians/GetById/?id={pid}")

        assert response.status_code == 200
        assert response.json()["Name"] == "Maria Ionescu"

    def test_not_found(self, client):
        with patch("app.routers.politicians.PoliticianService") as MockService:
            MockService.return_value.get_by_id.return_value = None
            response = client.get(f"/api/politicians/GetById/?id={uuid.uuid4()}")

        assert response.status_code == 200
        assert response.json() == "Politician not found"


class TestGetPoliticianByName:
    def test_found(self, client):
        politician = make_politician(Name="Gheorghe Mihai")
        with patch("app.routers.politicians.PoliticianService") as MockService:
            MockService.return_value.get_by_name.return_value = politician
            response = client.get("/api/politicians/GetByName/?name=Gheorghe+Mihai")

        assert response.status_code == 200
        assert response.json()["Name"] == "Gheorghe Mihai"

    def test_not_found(self, client):
        with patch("app.routers.politicians.PoliticianService") as MockService:
            MockService.return_value.get_by_name.return_value = None
            response = client.get("/api/politicians/GetByName/?name=Nobody")

        assert response.status_code == 200
        assert response.json() == "Politician not found"
