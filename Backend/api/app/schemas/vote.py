from uuid import UUID
from pydantic import BaseModel

from app.schemas.politician import PoliticianSchema
from app.schemas.voting_round import RoundSchema


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
        politician = PoliticianSchema.from_orm_model(obj.politician_rel) if obj.politician_rel else None
        round_ = RoundSchema.model_validate(obj.round_rel) if obj.round_rel else None
        return cls(
            Id=obj.Id,
            Name=obj.Name,
            PoliticianId=obj.PoliticianId,
            Position=obj.Position,
            RoundId=obj.RoundId,
            Politician=politician,
            Round=round_,
        )
