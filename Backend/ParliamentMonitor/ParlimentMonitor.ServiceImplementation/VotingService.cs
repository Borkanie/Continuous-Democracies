using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using StackExchange.Redis;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingService : IVotingService<Vote, Round>
    {
        internal readonly IConnectionMultiplexer _redis;
        internal IDatabase _cache;
        private readonly AppDBContext _dbContext;

        public VotingService(AppDBContext context, IConnectionMultiplexer redis)
        {
            _dbContext = context;
            _redis = redis;
            _cache = redis.GetDatabase();

        }

        public async Task<Round?> CreateVotingRound(string title, DateTime time, int id = 0, List<Vote>? votes = null, string? description = null)
        {
            try
            {
                var round = new Round();
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
                return round;
            }catch(Exception ex)
            {
                Console.WriteLine($"Error when creating Voting Round:{ex.Message}");
                return null;
            }
        }

        public async Task<bool> Delete(Round entity)
        {
            if(_dbContext.VotingRounds.Contains(entity))
            {
                _dbContext.VotingRounds.Remove(entity);
                _dbContext.SaveChanges();
                return true;
            }
            return false;
        }

        public async Task<bool> Delete(Vote entity)
        {
            if (_dbContext.Votes.Contains(entity))
            {
                _dbContext.Votes.Remove(entity);
                _dbContext.SaveChanges();
                return true;
            }
            return false;
        }

        public async Task<Round?> GetVotingRound(int votingRoundId)
        {
           var round = _dbContext.VotingRounds.Find(votingRoundId);
            if (round == null)
                Console.WriteLine($"No round with Id{votingRoundId} found in db");
            else
            {
                Console.WriteLine($"Returning round:{votingRoundId}");
                if(round.VoteResults.Count == 0)
                {
                    Console.WriteLine($"Loading vote results for round{round.VoteId}");
                    round.VoteResults = _dbContext.Votes.Include(x => x.Politician).Include(x => x.Politician.Party).Where(x => x.Round == round).ToHashSet();
                }
            }
           return round;
        }

        public async Task<bool> Update(Round voteResult)
        {
            var vote = _dbContext.VotingRounds.FirstOrDefault(x => x.Id == voteResult.Id); 
            if (vote != null)
            {
                _dbContext.Update(vote);
                vote.Description = voteResult.Description;
                vote.VoteResults = voteResult.VoteResults;
                vote.Title = voteResult.Title;
                _dbContext.SaveChanges();
                return true;
            }
            return false;
        }

        public async Task<bool> Update(Vote entity)
        {
            var vote = _dbContext.Votes.FirstOrDefault(x => x.Id == entity.Id);
            if (vote != null)
            {
                _dbContext.Update(vote);
                vote.Position = entity.Position;
                vote.Politician = entity.Politician;
                vote.Name = entity.Name;
                _dbContext.SaveChanges();
                return true;
            }
            return false;
        }
        public async Task<Vote?> CastVote(Round container, Politician politician, VotePosition position)
        {
            try
            {
                if (_dbContext.VotingRounds.Find(container) == null)
                    throw new Exception($"Voting round {container} not registered yet");
                _dbContext.Update(container);
                var vote = container.VoteResults.FirstOrDefault(x => x.Politician == politician);
                if (vote == null)
                {
                    vote = new Vote();
                    vote.Position = position;
                    vote.Politician = politician;
                    vote.Name = $"Vote from {politician.Name} during {container.Name}";
                    container.VoteResults.Add(vote);
                }
                else
                {
                    vote.Position = position;
                }
                _dbContext.Votes.Add(vote);
                _dbContext.SaveChanges();
                return vote;
            }catch(Exception ex)
            {
                Console.WriteLine($"Error when casting vote for {politician.Name} for law {container.Name} : {ex.Message}");
                return null;
            }
            
        }

        public async Task<Vote?> UpdateCastedVote(Round container, Politician politician, VotePosition position)
        {
            if (_dbContext.VotingRounds.Find(container) == null)
                throw new Exception($"Voting round {container} not registered yet");
            var vote = container.VoteResults.FirstOrDefault(x => x.Politician == politician);
            if (vote == null)
            {
                return null;
            }
            _dbContext.Update(vote);
            vote.Position = position;
            _dbContext.SaveChanges();
            return vote;
        }

        public async Task<Round?> UpdateVoteResult(Guid id, DateTime? time = null, List<Vote>? votes = null, string? Description = null)
        {
            var round = _dbContext.VotingRounds.Find(id);
            if (round == null)
            {
                return null;
            }
            _dbContext.Update(round);
            if (time != null)
            {
                round.VoteDate = time.Value;
            }
            if (votes != null)
            {
                foreach( var vote in round.VoteResults)
                    _dbContext.Remove(vote);
                round.VoteResults.Clear();
                foreach (var vote in votes)
                    round.VoteResults.Add(vote);
            }
            if (Description != null)
            {
                round.Description = Description;
            }
            _dbContext.SaveChanges();
            return round;
        }

        /// <inheritdoc/>
        public async Task<ISet<Round>> GetAllRoundsFromDB(int number = 100)
        {
            return _dbContext.VotingRounds.Take(number).ToHashSet();
        }
    }
}
