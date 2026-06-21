import uuid
from datetime import datetime, timezone
from unittest.mock import MagicMock

import pytest
from fastapi.testclient import TestClient

from app.main import app
from app.database import get_db


def make_party(**kwargs):
    p = MagicMock()
    p.Id = kwargs.get("Id", uuid.uuid4())
    p.Name = kwargs.get("Name", "Test Party")
    p.Acronym = kwargs.get("Acronym", "TP")
    p.LogoUrl = kwargs.get("LogoUrl", None)
    p.Color = kwargs.get("Color", "#123456")
    p.Active = kwargs.get("Active", True)
    return p


def make_politician(party=None, **kwargs):
    party = party or make_party()
    p = MagicMock()
    p.Id = kwargs.get("Id", uuid.uuid4())
    p.Name = kwargs.get("Name", "Ion Popescu")
    p.Gender = kwargs.get("Gender", 0)
    p.ImageUrl = kwargs.get("ImageUrl", None)
    p.PartyId = kwargs.get("PartyId", party.Id)
    p.Active = kwargs.get("Active", True)
    p.WorkLocation = kwargs.get("WorkLocation", 0)
    p.party_rel = party
    return p


def make_round(**kwargs):
    r = MagicMock()
    r.Id = kwargs.get("Id", uuid.uuid4())
    r.Name = kwargs.get("Name", "Test Round")
    r.Title = kwargs.get("Title", "Test Round")
    r.Description = kwargs.get("Description", "")
    r.VoteDate = kwargs.get("VoteDate", datetime(2024, 1, 15, 10, 0, 0, tzinfo=timezone.utc))
    r.VoteId = kwargs.get("VoteId", 100)
    return r


def make_vote(politician=None, round_=None, **kwargs):
    politician = politician or make_politician()
    round_ = round_ or make_round()
    v = MagicMock()
    v.Id = kwargs.get("Id", uuid.uuid4())
    v.Name = kwargs.get("Name", "")
    v.PoliticianId = kwargs.get("PoliticianId", politician.Id)
    v.Position = kwargs.get("Position", 0)
    v.RoundId = kwargs.get("RoundId", round_.Id)
    v.politician_rel = politician
    v.round_rel = round_
    return v


@pytest.fixture
def db_mock():
    return MagicMock()


@pytest.fixture
def client(db_mock):
    app.dependency_overrides[get_db] = lambda: db_mock
    with TestClient(app) as c:
        yield c
    app.dependency_overrides.clear()
