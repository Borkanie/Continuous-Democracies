import logging
from uuid import UUID
from sqlalchemy.orm import Session, joinedload
from app.models import Politician, Party, Gender, WorkLocation

logger = logging.getLogger(__name__)


class PoliticianService:
    def __init__(self, db: Session):
        self.db = db

    def get_by_id(self, politician_id: UUID) -> Politician | None:
        return (
            self.db.query(Politician)
            .options(joinedload(Politician.party_rel))
            .filter(Politician.Id == politician_id)
            .first()
        )

    def get_by_name(self, name: str) -> Politician | None:
        return (
            self.db.query(Politician)
            .options(joinedload(Politician.party_rel))
            .filter(Politician.Name == name)
            .first()
        )

    def get_all_politicians(
        self,
        party_acronym: str | None = None,
        party_name: str | None = None,
        is_active: bool | None = None,
        location: WorkLocation | None = None,
        gender: Gender | None = None,
        number: int = 100,
    ) -> list[Politician]:
        query = self.db.query(Politician).options(joinedload(Politician.party_rel))

        if party_acronym or party_name:
            query = query.join(Party)
            if party_acronym:
                query = query.filter(Party.Acronym == party_acronym)
            if party_name:
                query = query.filter(Party.Name == party_name)

        if is_active is not None:
            query = query.filter(Politician.Active == is_active)
        if location is not None:
            query = query.filter(Politician.WorkLocation == location.value)
        if gender is not None:
            query = query.filter(Politician.Gender == gender.value)

        return query.limit(number).all()
