from uuid import UUID
from pydantic import BaseModel


class PartySchema(BaseModel):
    Id: UUID
    Name: str
    Acronym: str | None = None
    LogoUrl: str | None = None
    Color: str = "#D3D3D3"
    Active: bool = False

    model_config = {"from_attributes": True}
