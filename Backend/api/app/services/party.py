import logging
from uuid import UUID
from sqlalchemy.orm import Session

from app.models import Party

logger = logging.getLogger(__name__)


class PartyService:
    def __init__(self, db: Session):
        self.db = db

    def get_all_parties(self, is_active: bool = True, number: int = 100) -> list[Party]:
        return (
            self.db.query(Party)
            .filter(Party.Active == is_active)
            .limit(number)
            .all()
        )

    def get_by_id(self, party_id: UUID) -> Party | None:
        return self.db.query(Party).filter(Party.Id == party_id).first()

    def get_party(self, name: str | None = None, acronym: str | None = None) -> Party | None:
        query = self.db.query(Party)
        if name:
            query = query.filter(Party.Name == name)
        if acronym:
            query = query.filter(Party.Acronym == acronym)
        return query.first()
