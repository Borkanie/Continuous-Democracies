import uuid
from sqlalchemy import Column, String, Integer, DateTime
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship

from app.database import Base


class Round(Base):
    __tablename__ = "VotingRounds"

    Id = Column("Id", UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    Name = Column("Name", String, nullable=False, default="")
    Title = Column("Title", String, nullable=False, default="")
    Description = Column("Description", String, nullable=False, default="")
    VoteDate = Column("VoteDate", DateTime(timezone=True), nullable=False)
    VoteId = Column("VoteId", Integer, nullable=False, default=0)

    votes = relationship("Vote", back_populates="round_rel")
