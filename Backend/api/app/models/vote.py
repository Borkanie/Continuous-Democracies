import uuid
from sqlalchemy import Column, String, Integer, ForeignKey
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship

from app.database import Base


class Vote(Base):
    __tablename__ = "Votes"

    Id = Column("Id", UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    Name = Column("Name", String, nullable=False, default="")
    PoliticianId = Column("PoliticianId", UUID(as_uuid=True), ForeignKey("Politicians.Id"), nullable=False)
    Position = Column("Position", Integer, nullable=False, default=3)
    RoundId = Column("RoundId", UUID(as_uuid=True), ForeignKey("VotingRounds.Id"), nullable=False)

    politician_rel = relationship("Politician", back_populates="votes", lazy="joined")
    round_rel = relationship("Round", back_populates="votes", lazy="joined")
