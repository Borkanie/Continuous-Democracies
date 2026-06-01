from pydantic import BaseModel
from uuid import UUID
from datetime import datetime
from app.models import Gender, WorkLocation, VotePosition


class PartySchema(BaseModel):
    Id: UUID
    Name: str
    Acronym: str | None = None
    LogoUrl: str | None = None
    Color: str = "#D3D3D3"
    Active: bool = False

    model_config = {"from_attributes": True}


class PoliticianSchema(BaseModel):
    Id: UUID
    Name: str
    Gender: int
    ImageUrl: str | None = None
    PartyId: UUID
    Active: bool = False
    WorkLocation: int
    Party: PartySchema | None = None

    model_config = {"from_attributes": True}

    @classmethod
    def from_orm_model(cls, obj):
        party = None
        if obj.party_rel:
            party = PartySchema.model_validate(obj.party_rel)
        return cls(
            Id=obj.Id,
            Name=obj.Name,
            Gender=obj.Gender,
            ImageUrl=obj.ImageUrl,
            PartyId=obj.PartyId,
            Active=obj.Active,
            WorkLocation=obj.WorkLocation,
            Party=party,
        )


class RoundSchema(BaseModel):
    Id: UUID
    Name: str
    Title: str
    Description: str = ""
    VoteDate: datetime
    VoteId: int = 0

    model_config = {"from_attributes": True}


class VoteSchema(BaseModel):
    Id: UUID
    Name: str
    PoliticianId: UUID
    Position: int
    RoundId: UUID
    Politician: PoliticianSchema | None = None
    Round: RoundSchema | None = None

    model_config = {"from_attributes": True}

    @classmethod
    def from_orm_model(cls, obj):
        politician = None
        if obj.politician_rel:
            politician = PoliticianSchema.from_orm_model(obj.politician_rel)
        round_ = None
        if obj.round_rel:
            round_ = RoundSchema.model_validate(obj.round_rel)
        return cls(
            Id=obj.Id,
            Name=obj.Name,
            PoliticianId=obj.PoliticianId,
            Position=obj.Position,
            RoundId=obj.RoundId,
            Politician=politician,
            Round=round_,
        )
