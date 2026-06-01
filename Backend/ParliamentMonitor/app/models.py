import uuid
import enum
from sqlalchemy import Column, String, Boolean, DateTime, Integer, ForeignKey
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship, DeclarativeBase


class Base(DeclarativeBase):
    pass


class Gender(int, enum.Enum):
    Male = 0
    Female = 1
    Other = 2


class WorkLocation(int, enum.Enum):
    Parliament = 0
    Senate = 1
    Other = 2


class VotePosition(int, enum.Enum):
    Yes = 0
    No = 1
    Abstain = 2
    Absent = 3


class Party(Base):
    __tablename__ = "Parties"

    Id = Column("Id", UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    Name = Column("Name", String, nullable=False, default="")
    Acronym = Column("Acronym", String, nullable=True)
    LogoUrl = Column("LogoUrl", String, nullable=True)
    Color = Column("Color", String, nullable=False, default="#D3D3D3")
    Active = Column("Active", Boolean, nullable=False, default=False)

    politicians = relationship("Politician", back_populates="party_rel")


class Politician(Base):
    __tablename__ = "Politicians"

    Id = Column("Id", UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    Name = Column("Name", String, nullable=False, default="")
    Gender = Column("Gender", Integer, nullable=False, default=Gender.Male.value)
    ImageUrl = Column("ImageUrl", String, nullable=True)
    PartyId = Column("PartyId", UUID(as_uuid=True), ForeignKey("Parties.Id"), nullable=False)
    Active = Column("Active", Boolean, nullable=False, default=False)
    WorkLocation = Column("WorkLocation", Integer, nullable=False, default=WorkLocation.Parliament.value)

    party_rel = relationship("Party", back_populates="politicians", lazy="joined")
    votes = relationship("Vote", back_populates="politician_rel")


class Round(Base):
    __tablename__ = "VotingRounds"

    Id = Column("Id", UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    Name = Column("Name", String, nullable=False, default="")
    Title = Column("Title", String, nullable=False, default="")
    Description = Column("Description", String, nullable=False, default="")
    VoteDate = Column("VoteDate", DateTime(timezone=True), nullable=False)
    VoteId = Column("VoteId", Integer, nullable=False, default=0)

    votes = relationship("Vote", back_populates="round_rel")


class Vote(Base):
    __tablename__ = "Votes"

    Id = Column("Id", UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    Name = Column("Name", String, nullable=False, default="")
    PoliticianId = Column("PoliticianId", UUID(as_uuid=True), ForeignKey("Politicians.Id"), nullable=False)
    Position = Column("Position", Integer, nullable=False, default=VotePosition.Absent.value)
    RoundId = Column("RoundId", UUID(as_uuid=True), ForeignKey("VotingRounds.Id"), nullable=False)

    politician_rel = relationship("Politician", back_populates="votes", lazy="joined")
    round_rel = relationship("Round", back_populates="votes", lazy="joined")
