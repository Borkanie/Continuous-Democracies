import logging
from uuid import UUID
from datetime import datetime
from fastapi import APIRouter, Depends, Query
from sqlalchemy.orm import Session

from app.database import get_db
from app.schemas import RoundSchema, VoteSchema
from app.services.voting import VotingService
from app.services.voting_round import VotingRoundService

logger = logging.getLogger(__name__)
router = APIRouter(prefix="/api/voting", tags=["Voting"])


@router.get("/getAllRounds")
def get_all_rounds(
    startDate: datetime | None = Query(default=None),
    endDate: datetime | None = Query(default=None),
    keywords: list[str] | None = Query(default=None),
    maxNumberOfEntries: int = Query(default=100),
    db: Session = Depends(get_db),
):
    service = VotingRoundService(db)
    rounds = service.get_all_rounds(
        start_date=startDate,
        end_date=endDate,
        keywords=keywords,
        max_entries=maxNumberOfEntries,
    )
    if not rounds:
        return "No results were available in the db"
    return [RoundSchema.model_validate(r) for r in rounds]


@router.get("/getRoundById/")
def get_round_by_id(
    voteNumber: int = Query(...),
    db: Session = Depends(get_db),
):
    service = VotingRoundService(db)
    round_ = service.get_by_vote_id(voteNumber)
    if not round_:
        return "No voting round found with that ID"
    return RoundSchema.model_validate(round_)


@router.get("/GetResultForVote/")
def get_result_for_vote(
    number: int = Query(...),
    partyId: UUID | None = Query(default=None),
    partyAcronim: str | None = Query(default=None),
    db: Session = Depends(get_db),
):
    service = VotingService(db)
    votes = service.get_all_votes_for_round_by_vote_id(number)

    if partyId:
        votes = [v for v in votes if v.politician_rel and v.politician_rel.PartyId == partyId]
    if partyAcronim:
        votes = [
            v for v in votes
            if v.politician_rel
            and v.politician_rel.party_rel
            and v.politician_rel.party_rel.Acronym == partyAcronim
        ]

    return [VoteSchema.from_orm_model(v) for v in votes]


@router.get("/GetAllVotesForARoundById/")
def get_all_votes_for_round_by_id(
    roundId: UUID = Query(...),
    db: Session = Depends(get_db),
):
    service = VotingService(db)
    votes = service.get_all_votes_for_round_by_guid(roundId)
    return [VoteSchema.from_orm_model(v) for v in votes]
