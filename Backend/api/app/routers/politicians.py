import logging
from uuid import UUID
from fastapi import APIRouter, Depends, Query
from sqlalchemy.orm import Session

from app.database import get_db
from app.models.enums import Gender, WorkLocation
from app.schemas import PoliticianSchema
from app.services.politician import PoliticianService

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/politicians", tags=["Politicians"])


@router.get("/getAllPoliticians")
def get_all_politicians(
    partyAcronym: str | None = Query(default=None),
    partyName: str | None = Query(default=None),
    isActive: bool | None = Query(default=None),
    location: int | None = Query(default=None),
    gender: int | None = Query(default=None),
    number: int = Query(default=100),
    db: Session = Depends(get_db),
):
    loc = WorkLocation(location) if location is not None else None
    gen = Gender(gender) if gender is not None else None

    service = PoliticianService(db)
    politicians = service.get_all_politicians(
        party_acronym=partyAcronym,
        party_name=partyName,
        is_active=isActive,
        location=loc,
        gender=gen,
        number=number,
    )
    if not politicians:
        return "No politicians found"
    return [PoliticianSchema.from_orm_model(p) for p in politicians]


@router.get("/GetById/")
def get_politician_by_id(
    id: UUID = Query(...),
    db: Session = Depends(get_db),
):
    service = PoliticianService(db)
    politician = service.get_by_id(id)
    if not politician:
        return "Politician not found"
    return PoliticianSchema.from_orm_model(politician)


@router.get("/GetByName/")
def get_politician_by_name(
    name: str = Query(...),
    db: Session = Depends(get_db),
):
    service = PoliticianService(db)
    politician = service.get_by_name(name)
    if not politician:
        return "Politician not found"
    return PoliticianSchema.from_orm_model(politician)
