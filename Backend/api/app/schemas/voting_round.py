from uuid import UUID
from datetime import datetime
from pydantic import BaseModel


class RoundSchema(BaseModel):
    Id: UUID
    Name: str
    Title: str
    Description: str = ""
    VoteDate: datetime
    VoteId: int = 0

    model_config = {"from_attributes": True}
