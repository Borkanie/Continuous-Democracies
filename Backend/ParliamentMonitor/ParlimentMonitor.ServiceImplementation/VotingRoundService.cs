using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingRoundService(
        AppDBContext context,
        IVotingService<Vote> votingService,
        IPoliticianService<Politician> politicianService,
        ILogger<IVotingRoundService<Round>> logger) : IVotingRoundService<Round>
    {
        private readonly AppDBContext _dbContext = context;
        private readonly IVotingService<Vote> _votingService = votingService;
        private readonly IPoliticianService<Politician> _politicianService = politicianService;
        private readonly ILogger<IVotingRoundService<Round>> _logger = logger;

        public ILogger Logger { get => _logger; }

        /// <inheritdoc/>
        public Task<Round?> CreateVotingRoundAsync(string title, DateTime time, int id = 0, string? description = null)
        {
            var round = new Round() { Id = Guid.NewGuid() };
            round.Title = title;
            round.VoteDate = time;
            round.VoteId = id;
            round.Description = description ?? string.Empty;
            if (_dbContext.VotingRounds.Any(x => x == round))
                throw new Exception($"Already existing voting round:{round.VoteId}");
            _dbContext.VotingRounds.Add(round);
            _dbContext.SaveChanges();
            _logger.LogInformation($"Created new voting round:{round.VoteId}");
            return Task.FromResult<Round?>(round);
        }

        /// <inheritdoc/>
        public Task<bool> DeleteAsync(Round entity)
        {
            if (_dbContext.VotingRounds.Contains(entity))
            {
                _dbContext.VotingRounds.Remove(entity);
                _dbContext.SaveChanges();
                _logger.LogInformation($"Deleted voting round:{entity.Id}");
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
                Console.WriteLine($"No round with Id {votingRoundId} found in db");
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
                _logger.LogInformation($"Updated voting round:{vote.Id}");
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
            _logger.LogInformation($"Updated voting round:{round.Id}");
            return Task.FromResult<Round?>(round);
        }

        /// <inheritdoc/>
        public Task<ISet<Round>> GetAllRoundsFromDBAsync(DateTime? startDate, DateTime? endDate, string?[] keywords, int number = 100)
        {
            var query = _dbContext.VotingRounds.AsQueryable();
            if (keywords != null && keywords.Length > 0)
            {
                // TODO: increase efficiency here because with a big DB this will be slow
                foreach (var keyword in keywords)
                {
                    if(keyword == null) 
                        continue;
                    var tempKeyword = keyword.ToLower();
                    query = query.Where(x => x.Title.ToLower().Contains(tempKeyword) || (x.Description != null && x.Description.ToLower().Contains(tempKeyword)));
                }
            }
            if(startDate != null && endDate == null)
            {
                if (DateTime.Now < startDate)
                {
                    throw new Exception("End date cannot be in the future");
                }
                endDate = DateTime.Now;
            }
            if (startDate != null && endDate != null)
            {
                return Task.FromResult<ISet<Round>>(query.Where(x => x.VoteDate >= startDate && x.VoteDate <= endDate).Take(number).ToHashSet());
            }
            return Task.FromResult<ISet<Round>>(query.Take(number).ToHashSet());
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
                    _logger.LogInformation($"Registered vote:{vote.Position} for politician:{politician.Name} in voting round:{container.VoteId}");
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
                throw new Exception($"Voting round {container.VoteId} not registered yet");
        }

        /// <inheritdoc/>
        public Task<Round?> GetAsync(Guid id)
        {
            return Task.FromResult(_dbContext.VotingRounds.FirstOrDefault(x => x.Id == id));
        }
    }
}
