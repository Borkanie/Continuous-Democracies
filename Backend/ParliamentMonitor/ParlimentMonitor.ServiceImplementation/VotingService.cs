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

        public VotingService(AppDBContext context,IPoliticianService<Politician> politicianService) 
        {
            _dbContext = context;
        }

        /// <inheritdoc/>
        public Task<Vote> CreateNewVote(Round round, Politician politician, VotePosition votePosition)
        {
            if (_dbContext.Votes.Any(x => x.Round.Id == round.Id && x.Politician.Id == politician.Id))
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
            return Task.FromResult(vote);

        }

        /// <inheritdoc/>
        public Task<List<Vote>> GetAllVotesForRound(Guid roundId)
        {
            return Task.FromResult(_dbContext.Votes.Include(x => x.Politician).Include(x => x.Round).Where(x => x.Round.Id == roundId).ToList());
        }

        /// <inheritdoc/>
        public Task<bool> DeleteAsync(Vote entity)
        {
            if (_dbContext.Votes.Contains(entity))
            {
                _dbContext.Votes.Remove(entity);
                _dbContext.SaveChanges();
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
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
                vote.Round = entity.Round;
                _dbContext.SaveChanges();
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        public Task<List<Vote>> GetAllVotesForRound(int roundVoteId)
        {
            return Task.FromResult(_dbContext.Votes.Include(x => x.Politician).Include(x => x.Round).Where(x => x.Round.VoteId == roundVoteId).ToList());
        }
    }
}
