from uuid import UUID
from pydantic import BaseModel

from app.schemas.party import PartySchema


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
        party = PartySchema.model_validate(obj.party_rel) if obj.party_rel else None
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
