class VotingRounds:
    def __init__(self, Id: str, VoteId: int, VoteDate: str):
        self.Id = Id
        self.VoteId = VoteId
        self.Title = f"Vot electronic:{VoteId}"
        self.voteDate = VoteDate    