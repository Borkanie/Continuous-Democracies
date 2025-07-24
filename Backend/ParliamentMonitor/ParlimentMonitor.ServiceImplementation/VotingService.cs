using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using StackExchange.Redis;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingService : RedisService, IVotingService<Vote>
    {
        private readonly AppDBContext _dbContext;
       

        public VotingService(AppDBContext context, IConnectionMultiplexer redis) : base(redis, "vote")
        {
            _dbContext = context;
            var endpoint = _redis.GetEndPoints().First();
            var server = _redis.GetServer(endpoint)
        }

        public string MakeKey(Vote vote)
        {
            return $"{serviceKey}:{vote.Round.VoteId}:{vote.Id}";
        }

        /*<inheritdoc/>*/
        public Task<Vote> CreateNewVote(Round round, Politician politician, VotePosition votePosition)
        {
            var vote = new Vote()
            {
                Round = round,
                Politician = politician,
                Id = Guid.NewGuid(),
                Position = votePosition
            };
            
            SetAsync<Vote>(MakeKey(vote), vote);
            return Task.FromResult(vote);

        }

        /*<inheritdoc/>*/
        public async Task<bool> DeleteAsync(Vote entity)
        {
            if (_dbContext.Votes.Contains(entity))
            {
                _dbContext.Votes.Remove(entity);
                _dbContext.SaveChanges();
                RemoveAsync(MakeKey(entity));
                return true;
            }
            return false;
        }

        /*<inheritdoc/>*/
        public async Task<Vote?> GetAsync(Guid id)
        {
            var keys = 
            var cached = GetAsync<Vote>(MakeKey(id.ToString())).Result;
            if(cached != null) 
            {
                return cached;
            }
            return _dbContext.Votes.FirstOrDefault(x => x.Id == id);
        }

        /*<inheritdoc/>*/
        public async Task<bool> UpdateAsync(Vote entity)
        {
            var vote = GetAsync(entity.Id).Result;
            if (vote != null)
            {
                _dbContext.Update(vote);
                vote.Position = entity.Position;
                vote.Politician = entity.Politician;
                vote.Name = entity.Name;
                _dbContext.SaveChanges();
                SetAsync(MakeKey(vote.Id.ToString()), vote);
                return true;
            }
            return false;
        }
    }
}
