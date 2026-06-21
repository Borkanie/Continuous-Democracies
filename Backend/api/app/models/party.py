import uuid
from sqlalchemy import Column, String, Boolean
from sqlalchemy.dialects.postgresql import UUID
from sqlalchemy.orm import relationship

from app.database import Base


class Party(Base):
    __tablename__ = "Parties"

    Id = Column("Id", UUID(as_uuid=True), primary_key=True, default=uuid.uuid4)
    Name = Column("Name", String, nullable=False, default="")
    Acronym = Column("Acronym", String, nullable=True)
    LogoUrl = Column("LogoUrl", String, nullable=True)
    Color = Column("Color", String, nullable=False, default="#D3D3D3")
    Active = Column("Active", Boolean, nullable=False, default=False)

    politicians = relationship("Politician", back_populates="party_rel")
