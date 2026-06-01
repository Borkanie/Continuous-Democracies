import logging
from uuid import UUID
from sqlalchemy.orm import Session, joinedload
from app.models import Vote, Round

logger = logging.getLogger(__name__)


class VotingService:
    def __init__(self, db: Session):
        self.db = db

    def get_all_votes_for_round_by_guid(self, round_id: UUID) -> list[Vote]:
        return (
            self.db.query(Vote)
            .options(
                joinedload(Vote.politician_rel).joinedload("party_rel"),
                joinedload(Vote.round_rel),
            )
            .filter(Vote.RoundId == round_id)
            .all()
        )

    def get_all_votes_for_round_by_vote_id(self, vote_id: int) -> list[Vote]:
        round_ = self.db.query(Round).filter(Round.VoteId == vote_id).first()
        if not round_:
            return []
        return self.get_all_votes_for_round_by_guid(round_.Id)
