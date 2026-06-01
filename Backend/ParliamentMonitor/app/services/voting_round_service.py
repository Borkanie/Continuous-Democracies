import logging
from datetime import datetime
from uuid import UUID
from sqlalchemy.orm import Session
from app.models import Round

logger = logging.getLogger(__name__)


class VotingRoundService:
    def __init__(self, db: Session):
        self.db = db

    def get_by_id(self, round_id: UUID) -> Round | None:
        return self.db.query(Round).filter(Round.Id == round_id).first()

    def get_by_vote_id(self, vote_id: int) -> Round | None:
        return self.db.query(Round).filter(Round.VoteId == vote_id).first()

    def get_all_rounds(
        self,
        start_date: datetime | None = None,
        end_date: datetime | None = None,
        keywords: list[str] | None = None,
        max_entries: int = 100,
    ) -> list[Round]:
        query = self.db.query(Round)

        if start_date:
            query = query.filter(Round.VoteDate >= start_date)
        if end_date:
            query = query.filter(Round.VoteDate <= end_date)

        if keywords:
            for kw in keywords:
                kw_lower = kw.lower()
                query = query.filter(
                    Round.Title.ilike(f"%{kw_lower}%")
                    | Round.Description.ilike(f"%{kw_lower}%")
                )

        return query.limit(max_entries).all()
