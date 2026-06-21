import uuid
from sqlalchemy import Column, String, Boolean, Integer, ForeignKey
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship

from app.database import Base


class Politician(Base):
    __tablename__ = "Politicians"

    Id = Column("Id", UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    Name = Column("Name", String, nullable=False, default="")
    Gender = Column("Gender", Integer, nullable=False, default=0)
    ImageUrl = Column("ImageUrl", String, nullable=True)
    PartyId = Column("PartyId", UUID(as_uuid=True), ForeignKey("Parties.Id"), nullable=False)
    Active = Column("Active", Boolean, nullable=False, default=False)
    WorkLocation = Column("WorkLocation", Integer, nullable=False, default=0)

    party_rel = relationship("Party", back_populates="politicians", lazy="joined")
    votes = relationship("Vote", back_populates="politician_rel")
