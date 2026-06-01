import logging
from uuid import UUID
from fastapi import APIRouter, Depends, Query
from fastapi.responses import JSONResponse
from sqlalchemy.orm import Session

from app.database import get_db
from app.schemas import PartySchema
from app.services.party_service import PartyService

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/party", tags=["Party"])


@router.get("/all")
def get_all_parties(
    active: bool = True,
    number: int = Query(default=100),
    db: Session = Depends(get_db),
):
    service = PartyService(db)
    parties = service.get_all_parties(is_active=active, number=number)
    if not parties:
        return JSONResponse(status_code=404, content="No parties found")
    return [PartySchema.model_validate(p) for p in parties]


@router.get("/GetById/")
def get_party_by_id(
    id: UUID = Query(...),
    db: Session = Depends(get_db),
):
    service = PartyService(db)
    party = service.get_by_id(id)
    if not party:
        return JSONResponse(status_code=404, content="Party not found")
    return PartySchema.model_validate(party)


@router.get("/query")
def query_party(
    name: str | None = Query(default=None),
    acronym: str | None = Query(default=None),
    db: Session = Depends(get_db),
):
    service = PartyService(db)
    party = service.get_party(name=name, acronym=acronym)
    if not party:
        return JSONResponse(
            status_code=404,
            content="No party found with the specified criteria",
        )
    return PartySchema.model_validate(party)
