using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using StackExchange.Redis;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using static System.Runtime.InteropServices.JavaScript.JSType;

namespace ParliamentMonitor.ServiceImplementation
{
    internal class VotingRoundService : RedisService, IVotingRoundService<Round>
    {
        private readonly AppDBContext _dbContext;
        private readonly IVotingService<Vote> _votingService;

        public VotingRoundService(AppDBContext context, IVotingService<Vote> votingService, IConnectionMultiplexer redis) : base(redis, "round")
        {
            _dbContext = context;
            _votingService = votingService;
        }

        /// <inheritdoc/>
        public async Task<Round?> CreateVotingRoundAsync(string title, DateTime time, int id = 0, List<Vote>? votes = null, string? description = null)
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
                SetAsync(MakeKey(round.VoteId.ToString()), round);
                return round;
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error when creating Voting Round:{ex.Message}");
                return null;
            }
        }

        /// <inheritdoc/>
        public async Task<bool> DeleteAsync(Round entity)
        {
            if (_dbContext.VotingRounds.Contains(entity))
            {
                _dbContext.VotingRounds.Remove(entity);
                _dbContext.SaveChanges();
                RemoveAsync(MakeKey(entity.VoteId.ToString()));
                return true;
            }
            return false;
        }

        public async Task<Round?> GetVotingRoundAsync(int votingRoundId)
        {
            var round = GetAsync<Round>(MakeKey(votingRoundId.ToString())).Result;
            if(round == null)
                 round = _dbContext.VotingRounds.Find(votingRoundId);
            if (round == null)
                Console.WriteLine($"No round with Id{votingRoundId} found in db");
            else
            {
                Console.WriteLine($"Returning round:{votingRoundId}");
                if (round.VoteResults.Count == 0)
                {
                    Console.WriteLine($"Loading vote results for round{round.VoteId}");
                    foreach(var id in round.VoteResults)
                    round.VoteResults = _votingService.GetAsync()
                      //_dbContext.Votes.Include(x => x.Politician).Include(x => x.Politician.Party).ToHashSet();
                }
            }
            return round;
        }

        public async Task<bool> UpdateAsync(Round voteResult)
        {
            var vote = _dbContext.VotingRounds.FirstOrDefault(x => x.Id == voteResult.Id);
            if (vote != null)
            {
                _dbContext.Update(vote);
                vote.Description = voteResult.Description;
                vote.VoteResults = voteResult.VoteResults;
                vote.Title = voteResult.Title;
                _dbContext.SaveChanges();
                SetAsync(MakeKey(vote.Id.ToString()), vote);
                return true;
            }
            return false;
        }

        public async Task<Round?> UpdateVoteResultAsync(Guid id, DateTime? time = null, List<Vote>? votes = null, string? Description = null)
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
                foreach (var vote in round.VoteResults)
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
            SetAsync(MakeKey(round.Id.ToString()), round);
            return round;
        }

        /// <inheritdoc/>
        public async Task<ISet<Round>> GetAllRoundsFromDBAsync(int number = 100)
        {
            return _dbContext.VotingRounds.Take(number).ToHashSet();
        }

        public async Task<Vote?> CastVoteAsync(Round container, Politician politician, VotePosition position)
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
            }
            catch (Exception ex)
            {
                Console.WriteLine($"Error when casting vote for {politician.Name} for law {container.Name} : {ex.Message}");
                return null;
            }

        }

        public async Task<Vote?> UpdateCastedVoteAsync(Round container, Politician politician, VotePosition position)
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
            SetAsync(MakeKey(vote.Id.ToString()), vote);
            return vote;
        }

        /// <inheritdoc/>
        public async Task<Round?> GetAsync(Guid id)
        {
            var cached = GetAsync<Round>(id.ToString()).Result;
            if (cached != null)
                return cached;
            return _dbContext.VotingRounds.FirstOrDefault(x => x.Id == id);
        }
    }
}
