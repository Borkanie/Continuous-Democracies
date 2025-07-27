using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingService : IVotingService<Vote>
    {
        private readonly AppDBContext _dbContext;
        private IVotingRoundService<Round>? _votingRoundService;

        public VotingService(AppDBContext context,IPoliticianService<Politician> politicianService) 
        {
            _dbContext = context;
        }

        public void SetRoundService(IVotingRoundService<Round> votingRoundService)
        {
            _votingRoundService = votingRoundService;
        }

        /// <inheritdoc/>
        public Task<Vote> CreateNewVote(Round round, Politician politician, VotePosition votePosition)
        {
            if (round.VoteResults.Any(x => x.Politician.Id == politician.Id))
            {
                throw new InvalidOperationException("Politician has already voted in this round.");
            }
            var vote = new Vote()
            {
                Round = round,
                Politician = politician,
                Id = Guid.NewGuid(),
                Position = votePosition
            };

            _dbContext.Votes.Add(vote);
            _dbContext.SaveChanges();
            AddVoteToRound(round, vote);
            return Task.FromResult(vote);

        }

        private void AddVoteToRound(Round round, Vote vote)
        {
            round.VoteResults.Add(vote);
            if (_votingRoundService!=null)
            {
                _ = _votingRoundService.UpdateVoteResultAsync(round.Id, votes: round.VoteResults);
            } 
        }

        /// <inheritdoc/>
        public Task<bool> DeleteAsync(Vote entity)
        {
            if (_dbContext.Votes.Contains(entity))
            {
                RemoveEntityFromRoundDbAndCache(entity);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        private void RemoveEntityFromRoundDbAndCache(Vote entity)
        {
            entity.Round.VoteResults.Remove(entity);
            if (_votingRoundService != null)
            {
                _ = _votingRoundService.UpdateVoteResultAsync(entity.Round.Id, votes: entity.Round.VoteResults);
            }
        }

        /// <inheritdoc/>
        public Task<bool> DeleteWithouthRemovingFromRound(Guid id)
        {
            var entity = GetAsync(id).Result;
            if (entity != null)
            {
                _dbContext.Votes.Remove(entity);
                _dbContext.VotingRounds.Update(entity.Round);
                return Task.FromResult(true);
            }
            return Task.FromResult(true);
        }

        

        /// <inheritdoc/>
        public Task<Vote?> GetAsync(Guid id)
        {
            return Task.FromResult(_dbContext.Votes.Include(x => x.Politician).Include(x => x.Round).FirstOrDefault(x => x.Id == id));
        }

        /// <inheritdoc/>
        public Task<bool> UpdateAsync(Vote entity)
        {
            var vote = GetAsync(entity.Id).Result;
            if (vote != null)
            {
                _dbContext.Update(vote);
                vote.Position = entity.Position;
                vote.Politician = entity.Politician;
                vote.Name = entity.Name;
                if(entity.Round != vote.Round)
                {
                    RemoveEntityFromRoundDbAndCache(vote);
                    vote.Round = entity.Round;
                    AddVoteToRound(vote.Round, vote);
                }
                _dbContext.SaveChanges();
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }
    }
}
