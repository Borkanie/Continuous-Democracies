import logging
from uuid import UUID
from fastapi import APIRouter, Depends, Query
from fastapi.responses import JSONResponse
from sqlalchemy.orm import Session

from app.database import get_db
from app.models import Gender, WorkLocation
from app.schemas import PoliticianSchema
from app.services.politician_service import PoliticianService

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/politicians", tags=["Politicians"])


@router.get("/GetById/")
def get_politician_by_id(
    id: UUID = Query(...),
    db: Session = Depends(get_db),
):
    service = PoliticianService(db)
    politician = service.get_by_id(id)
    if not politician:
        return JSONResponse(status_code=404, content="Politician not found")
    return PoliticianSchema.from_orm_model(politician)


@router.get("/GetByName/")
def get_politician_by_name(
    name: str = Query(...),
    db: Session = Depends(get_db),
):
    service = PoliticianService(db)
    politician = service.get_by_name(name)
    if not politician:
        return JSONResponse(
            status_code=404,
            content="No Politician with this name found",
        )
    return PoliticianSchema.from_orm_model(politician)


@router.get("/getAllPoliticians")
def get_all_politicians(
    partyAcronym: str | None = Query(default=None),
    partyName: str | None = Query(default=None),
    isActive: bool | None = Query(default=None),
    location: WorkLocation | None = Query(default=None),
    gender: Gender | None = Query(default=None),
    number: int = Query(default=100),
    db: Session = Depends(get_db),
):
    service = PoliticianService(db)
    politicians = service.get_all_politicians(
        party_acronym=partyAcronym,
        party_name=partyName,
        is_active=isActive,
        location=location,
        gender=gender,
        number=number,
    )
    if not politicians:
        return JSONResponse(
            status_code=404,
            content="No politicians found with the specified criteria",
        )
    return [PoliticianSchema.from_orm_model(p) for p in politicians]
