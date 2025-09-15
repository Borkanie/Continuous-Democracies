using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingRoundService : IVotingRoundService<Round>
    {
        private readonly AppDBContext _dbContext;
        private readonly IVotingService<Vote> _votingService;
        private readonly IPoliticianService<Politician> _politicianService;
        private readonly ILogger<IVotingRoundService<Round>> _logger;

        public ILogger Logger { get => _logger; }

        public VotingRoundService(
            AppDBContext context,
            IVotingService<Vote> votingService,
            IPoliticianService<Politician> politicianService,
            ILogger<IVotingRoundService<Round>> logger)
        {
            _dbContext = context;
            _votingService = votingService;
            _politicianService = politicianService;
            _logger = logger;
        }

        /// <inheritdoc/>
        public Task<Round?> CreateVotingRoundAsync(string title, DateTime time, int id = 0, string? description = null)
        {
            var round = new Round() { Id = Guid.NewGuid() };
            round.Title = title;
            round.VoteDate = time;
            round.VoteId = id;
            round.Description = description ?? string.Empty;
            if (_dbContext.VotingRounds.Any(x => x == round))
                throw new Exception($"Already existing voting round:{round}");
            _dbContext.VotingRounds.Add(round);
            _dbContext.SaveChanges();
            _logger.LogInformation($"Created new voting round:{round}");
            return Task.FromResult<Round?>(round);
        }

        /// <inheritdoc/>
        public Task<bool> DeleteAsync(Round entity)
        {
            if (_dbContext.VotingRounds.Contains(entity))
            {
                _dbContext.VotingRounds.Remove(entity);
                _dbContext.SaveChanges();
                _logger.LogInformation($"Deleted voting round:{entity}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /// <summary>
        /// Works directly with DB.
        /// Gets the voting round based off the voting number.
        /// </summary>
        /// <param name="votingRoundId"></param>
        /// <returns></returns>
        public Task<Round?> GetVotingRoundAsync(int votingRoundId)
        {
            // Use SingleOrDefault to fetch the round with included VoteResults by VoteId
            Round? round = _dbContext.VotingRounds.FirstOrDefault(x => x.VoteId == votingRoundId);

            if (round == null)
                Console.WriteLine($"No round with Id{votingRoundId} found in db");
            return Task.FromResult(round);
        }

        /// <inheritdoc/>
        public Task<bool> UpdateAsync(Round voteResult)
        {
            var vote = _dbContext.VotingRounds.FirstOrDefault(x => x.Id == voteResult.Id);
            if (vote != null)
            {
                _dbContext.Update(vote);
                vote.Description = voteResult.Description;
                vote.Title = voteResult.Title;
                _dbContext.SaveChanges();
                _logger.LogInformation($"Updated voting round:{vote}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /// <inheritdoc/>
        public Task<Round?> UpdateVoteResultAsync(Guid id, DateTime? time = null, string? Description = null)
        {
            var round = _dbContext.VotingRounds.Find(id);
            if (round == null)
            {
                return Task.FromResult<Round?>(null);
            }
            _dbContext.Update(round);
            if (time != null)
            {
                round.VoteDate = time.Value;
            }
            if (Description != null)
            {
                round.Description = Description;
            }
            _dbContext.SaveChanges();
            _logger.LogInformation($"Updated voting round:{round}");
            return Task.FromResult<Round?>(round);
        }

        /// <inheritdoc/>
        public Task<ISet<Round>> GetAllRoundsFromDBAsync(int number = 100)
        {
            return Task.FromResult<ISet<Round>>(_dbContext.VotingRounds.Take(number).ToHashSet());
        }

        /// <inheritdoc/>
        public Task<Vote?> RegisterVoteAsync(Round container, Politician politician, VotePosition position)
        {
            try
            {
                checkVotingRoundAndPolitciansExist(container, politician);

                var vote = _votingService.CreateNewVote(container, politician, position).Result;
                if (vote != null)
                {
                    _dbContext.Update(container);
                    _dbContext.SaveChanges();
                    _logger.LogInformation($"Registered vote:{vote} for politician:{politician} in voting round:{container}");
                }
                return Task.FromResult(vote);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error when casting vote for {politician.Name} for law {container.Name} : {ex.Message}");
                return Task.FromResult<Vote?>(null);
            }

        }

        private void checkVotingRoundAndPolitciansExist(Round container, Politician politician)
        {
            if (_politicianService.GetAsync(politician.Id).Result == null)
                throw new Exception($"Politician {politician.Name} not registered yet");

            if (GetAsync(container.Id).Result == null)
                throw new Exception($"Voting round {container} not registered yet");
        }

        /// <inheritdoc/>
        public Task<Round?> GetAsync(Guid id)
        {
            return Task.FromResult(_dbContext.VotingRounds.FirstOrDefault(x => x.Id == id));
        }

        /// <inheritdoc/>
        public Task<List<Round>> GetAllVotesInInterval(DateTime startDate, DateTime endDate)
        {
            if(startDate > endDate)
            {
                var dummy = endDate;
                endDate = startDate;
                startDate = dummy;
            }
            return Task.FromResult(_dbContext.VotingRounds.Where(x => x.VoteDate >= startDate && x.VoteDate <= endDate).ToList());
        }
    }
}
