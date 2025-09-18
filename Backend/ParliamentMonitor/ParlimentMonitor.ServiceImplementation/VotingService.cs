using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingService(AppDBContext context, IPoliticianService<Politician> politicianService, ILogger<IPoliticianService<Politician>> logger) : IVotingService<Vote>
    {
        private readonly AppDBContext _dbContext = context;
        private readonly ILogger<IPoliticianService<Politician>> _logger = logger;

        public ILogger Logger { get => _logger; }

        /// <inheritdoc/>
        public Task<Vote> CreateNewVote(Round round, Politician politician, VotePosition votePosition)
        {
            var oldvote = _dbContext.Votes.FirstOrDefault(x => x.Round.Id == round.Id && x.Politician.Id == politician.Id);
            if (oldvote != null)
            {
                throw new InvalidOperationException($"Politician:{politician} has already voted in this round:{round} and voted:{oldvote}");
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
            _logger.LogInformation($"Created new vote:{vote.Position} for politician:{politician.Name} in round:{round.VoteId}");
            return Task.FromResult(vote);

        }

        /// <inheritdoc/>
        public Task<List<Vote>> GetAllVotesForRound(Guid roundId)
        {
            return Task.FromResult(_dbContext.Votes.Include(x => x.Politician).Include(x => x.Politician.Party).Include(x => x.Round).Where(x => x.Round.Id == roundId).ToList());
        }

        /// <inheritdoc/>
        public Task<bool> DeleteAsync(Vote entity)
        {
            if (_dbContext.Votes.Contains(entity))
            {
                _dbContext.Votes.Remove(entity);
                _dbContext.SaveChanges();
                _logger.LogInformation($"Deleted vote:{entity.Id}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /// <inheritdoc/>
        public Task<Vote?> GetAsync(Guid id)
        {
            return Task.FromResult(_dbContext.Votes.Include(x => x.Politician).Include(x => x.Politician.Party).Include(x => x.Round).FirstOrDefault(x => x.Id == id));
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
                _logger.LogInformation($"Updated vote:{vote.Id}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        public Task<List<Vote>> GetAllVotesForRound(int roundVoteId)
        {
            return Task.FromResult(_dbContext.Votes.Include(x => x.Politician).Include(x => x.Politician.Party).Include(x => x.Round).Where(x => x.Round.VoteId == roundVoteId).ToList());
        }
    }
}
