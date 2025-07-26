using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using ParliamentMonitor.ServiceImplementation.Utils;
using StackExchange.Redis;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingService : RedisService, IVotingService<Vote>
    {
        private readonly AppDBContext _dbContext;
        private readonly IPoliticianService<Politician> _politicianService;
        private Lazy<IVotingRoundService<Round>> _votingRoundService;

        public VotingService(AppDBContext context, IConnectionMultiplexer redis, IPoliticianService<Politician> politicianService, Lazy<IVotingRoundService<Round>> votingRoundService) : base(redis, "vote", new VotingJsonSerializer())
        {
            _dbContext = context;
            _politicianService = politicianService;
            _votingRoundService = votingRoundService;
        }

        public string MakeKey(Vote vote)
        {
            return $"{serviceKey}:{vote.Round.VoteId}:{vote.Id}";
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
            _ = SetAsync<Vote>(MakeKey(vote), vote);
            return Task.FromResult(vote);

        }

        private void AddVoteToRound(Round round, Vote vote)
        {
            round.VoteResults.Add(vote);
            if (_votingRoundService.IsValueCreated)
            {
                _ = _votingRoundService.Value.UpdateVoteResultAsync(round.Id, votes: round.VoteResults);
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
            if (_votingRoundService.IsValueCreated)
            {
                _ = _votingRoundService.Value.UpdateVoteResultAsync(entity.Round.Id, votes: entity.Round.VoteResults);
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
                _ = RemoveAsync(MakeKey(entity));
                return Task.FromResult(true);
            }
            return Task.FromResult(true);
        }

        

        /// <inheritdoc/>
        public Task<Vote?> GetAsync(Guid id)
        {
            var cached = GetAsync<Vote>(MakeKey(id.ToString())).Result;
            if(cached != null) 
            {
                var politician = _politicianService.GetAsync(cached.Politician.Id).Result;
                if(politician == null)
                {
                    throw new Exception($"Politician with ID {cached.Politician.Id} not found when  retrieving vote: {cached.Id}");
                }
                if(_votingRoundService.IsValueCreated )
                {
                    if(cached.RoundId == null)
                    {
                        throw new Exception($"Round ID is null for vote: {cached.Id}");
                    }
                    var round = _votingRoundService.Value.GetAsync(cached.RoundId!.Value).Result;
                    if(round == null)
                    {
                        throw new Exception($"Voting round with ID {cached.Round.Id} not found when retrieving vote: {cached.Id}");
                    }
                    cached.Round = round;
                }
                cached.Politician = politician;
                return Task.FromResult<Vote?>(cached);
            }
            return Task.FromResult(_dbContext.Votes.FirstOrDefault(x => x.Id == id));
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
                _ = SetAsync(MakeKey(vote.Id.ToString()), vote);
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }
    }
}
