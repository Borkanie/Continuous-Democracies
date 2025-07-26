using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using ParliamentMonitor.ServiceImplementation.Utils;
using StackExchange.Redis;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingRoundService : RedisService, IVotingRoundService<Round>
    {
        private readonly AppDBContext _dbContext;
        private readonly IVotingService<Vote> _votingService;
        private readonly IPoliticianService<Politician> _politicianService;
        public VotingRoundService(
            AppDBContext context,
            IVotingService<Vote> votingService,
            IPoliticianService<Politician> politicianService,
            IConnectionMultiplexer redis
        )
            : base(redis, "round", new VotingRoundJsonSerializer()) // Pass your custom converter here
        {
            _dbContext = context;
            _votingService = votingService;
            _politicianService = politicianService;
        }

        /// <inheritdoc/>
        public Task<Round?> CreateVotingRoundAsync(string title, DateTime time, int id = 0, List<Vote>? votes = null, string? description = null)
        {
            try
            {
                var round = new Round() { Id = Guid.NewGuid() };
                round.Title = title;
                round.VoteDate = time;
                round.VoteId = id;
                if (votes != null)
                {
                    foreach (var vote in votes)
                        round.VoteResults.Add(vote);
                }
                round.Description = description ?? string.Empty;
                _dbContext.VotingRounds.Add(round);
                _dbContext.SaveChanges();
                _ = SetAsync(MakeKey(round.VoteId.ToString()), round);
                return Task.FromResult<Round?>(round);
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error when creating Voting Round:{ex.Message}");
                return Task.FromResult<Round?>(null);
            }
        }

        /// <inheritdoc/>
        public Task<bool> DeleteAsync(Round entity)
        {
            if (_dbContext.VotingRounds.Contains(entity))
            {
                foreach (var vote in entity.VoteResults)
                {
                    _votingService.DeleteWithouthRemovingFromRound(vote.Id).Wait();
                }
                _dbContext.VotingRounds.Remove(entity);
                _dbContext.SaveChanges();
                _ = RemoveAsync(MakeKey(entity.VoteId.ToString()));
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
            Round? round = GetRoundWithVotesFromCache(votingRoundId);
            if (round != null)
            {
                return Task.FromResult<Round?>(round);
            }
            round = GetRoundWithVotesFromDataBase(votingRoundId);
            return Task.FromResult(round);
        }

        private Round? GetRoundWithVotesFromDataBase(int votingRoundId)
        {
            Round? round = _dbContext.VotingRounds.Find(votingRoundId);
            if (round == null)
                Console.WriteLine($"No round with Id{votingRoundId} found in db");
            else
            {
                Console.WriteLine($"Returning round:{votingRoundId}");
                if (round.VoteResults.Count == 0)
                {
                    Console.WriteLine($"Loading vote results for round{round.VoteId}");
                    foreach (var id in round.VoteResults)
                        round.VoteResults = _dbContext.Votes.Include(x => x.Politician).Include(x => x.Politician.Party).ToHashSet();
                }
            }

            return round;
        }

        private Round? GetRoundWithVotesFromCache(int votingRoundId)
        {
            var round = GetAsync<Round>(MakeKey(votingRoundId.ToString())).Result;
            if (round != null)
            {
                Console.WriteLine($"Returning round:{votingRoundId}");
                if (round.VoteResults.Count < round.VoteResultIds.Count)
                {
                    foreach (var id in round.VoteResultIds)
                    {
                        if (round.VoteResults.Any(x => x.Id == id))
                            continue;
                        var vote = _votingService.GetAsync(id).Result;
                        if(vote == null)
                        {
                            Console.WriteLine($"Vote with Id {id} not found in cache");
                            continue;
                        }
                        round.VoteResults.Add(vote);
                    }
                }

            }

            return round;
        }

        /// <inheritdoc/>
        public Task<bool> UpdateAsync(Round voteResult)
        {
            var vote = _dbContext.VotingRounds.FirstOrDefault(x => x.Id == voteResult.Id);
            if (vote != null)
            {
                _dbContext.Update(vote);
                vote.Description = voteResult.Description;
                vote.VoteResults = voteResult.VoteResults;
                vote.Title = voteResult.Title;
                _dbContext.SaveChanges();
                _ = SetAsync(MakeKey(vote.Id.ToString()), vote);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /// <inheritdoc/>
        public Task<Round?> UpdateVoteResultAsync(Guid id, DateTime? time = null, ISet<Vote>? votes = null, string? Description = null)
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
            if (votes != null)
            {
                UpdateExistingVoteOrRemoveIfNoLongerPresentAndDeleteFromCacheAndDB(votes, round);
                round.VoteResults.Clear();
                foreach (var vote in votes)
                {
                    if (!round.VoteResults.Contains(vote))
                        round.VoteResults.Add(vote);
                }
            }
            if (Description != null)
            {
                round.Description = Description;
            }
            _dbContext.SaveChanges();
            _ = SetAsync(MakeKey(round.Id.ToString()), round);
            return Task.FromResult<Round?>(round);
        }

        private void UpdateExistingVoteOrRemoveIfNoLongerPresentAndDeleteFromCacheAndDB(ISet<Vote> votes, Round round)
        {
            foreach (var vote in round.VoteResults)
            {
                if (votes.Contains(vote))
                {
                    // If the vote already exists, update it
                    var existingVote = round.VoteResults.FirstOrDefault(x => x.Id == vote.Id);
                    if (existingVote != null)
                    {
                        existingVote.Position = vote.Position;
                        existingVote.Politician = vote.Politician; // Update politician if needed
                    }
                }
                else
                {
                    // If the vote does not exist, add it
                    _votingService.DeleteWithouthRemovingFromRound(vote.Id).Wait();
                    round.VoteResults.Remove(vote);
                }
            }
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
                    container.VoteResults.Add(vote);
                    _dbContext.SaveChanges();
                    _ = SetAsync(MakeKey(container.Id.ToString()), container);
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

        public Task<Vote?> UpdateCastedVoteAsync(Round container, Politician politician, VotePosition position)
        {
            if (_dbContext.VotingRounds.Find(container) == null)
                throw new Exception($"Voting round {container} not registered yet");
            var vote = container.VoteResults.FirstOrDefault(x => x.Politician == politician);
            if (vote == null)
            {
                return Task.FromResult<Vote?>(null);
            }
            _dbContext.Update(vote);
            vote.Position = position;
            _dbContext.SaveChanges();
            _ = SetAsync(MakeKey(vote.Id.ToString()), vote);
            return Task.FromResult<Vote?>(vote);
        }

        /// <inheritdoc/>
        public Task<Round?> GetAsync(Guid id)
        {
            var cached = GetAsync<Round>(id.ToString()).Result;
            if (cached != null)
                return Task.FromResult<Round?>(cached);
            return Task.FromResult(_dbContext.VotingRounds.FirstOrDefault(x => x.Id == id));
        }
    }
}
