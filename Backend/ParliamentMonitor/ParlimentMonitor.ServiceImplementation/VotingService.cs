using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;

namespace ParliamentMonitor.ServiceImplementation
{
    public class VotingService : IVotingService<Vote, Round>
    {
        private readonly AppDBContext dBContext;

        public VotingService(AppDBContext context)
        {
            dBContext = context;
        }

        public Round CreateVotingRound(string title, DateTime time, int id = 0, List<Vote>? votes = null, string? description = null)
        {
            var round = new Round();
            round.Title = title;
            round.VoteDate = time;
            round.VoteId = id;
            if(votes!= null)
            {
                foreach(var vote in votes)
                    round.VoteResults.Add(vote);
            }
            round.Description = description ?? string.Empty;
            dBContext.VotingRounds.Add(round);
            dBContext.SaveChanges();
            return round;
        }

        public void Delete(Round entity)
        {
            if(dBContext.VotingRounds.Contains(entity))
            {
                dBContext.VotingRounds.Remove(entity);
                dBContext.SaveChanges();
            }
        }

        public void Delete(Vote entity)
        {
            if (dBContext.Votes.Contains(entity))
            {
                dBContext.Votes.Remove(entity);
                dBContext.SaveChanges();
            }
        }

        public Round? GetVotingRound(int votingRoundId)
        {
           var round = dBContext.VotingRounds.FirstOrDefault(x => x.VoteId == votingRoundId);
            if (round == null)
                Console.WriteLine($"No round with Id{votingRoundId} found in db");
            else
            {
                Console.WriteLine($"Returning round:{votingRoundId}");
                if(round.VoteResults.Count == 0)
                {
                    Console.WriteLine($"Loading vote results for round{round.VoteId}");
                    round.VoteResults = dBContext.Votes.Include(x => x.Politician).Include(x => x.Politician.Party).Where(x => x.Round == round).ToHashSet();
                }
            }
           return round;
        }

        public void Update(Round voteResult)
        {
            var vote = dBContext.VotingRounds.FirstOrDefault(x => x.Id == voteResult.Id); 
            if (vote != null)
            {
                dBContext.Update(vote);
                vote.Description = voteResult.Description;
                vote.VoteResults = voteResult.VoteResults;
                vote.Title = voteResult.Title;
                dBContext.SaveChanges();
            }
        }

        public void Update(Vote entity)
        {
            var vote = dBContext.Votes.FirstOrDefault(x => x.Id == entity.Id);
            if (vote != null)
            {
                dBContext.Update(vote);
                vote.Position = entity.Position;
                vote.Politician = entity.Politician;
                vote.Name = entity.Name;
                dBContext.SaveChanges();
            }
        }
        public Vote CastVote(Round container, Politician politician, VotePosition position)
        {
            if (!dBContext.VotingRounds.Contains(container))
                throw new Exception($"Voting round {container} not registered yet");
            dBContext.Update(container);
            var vote = container.VoteResults.FirstOrDefault(x => x.Politician == politician);
            if (vote == null)
            {
                vote = new Vote();
                vote.Position = position;
                vote.Politician = politician;
                dBContext.Votes.Add(vote);
                container.VoteResults.Add(vote);
            }
            else
            {
                vote.Position = position;    
            }
            dBContext.Votes.Add(vote);
            dBContext.SaveChanges();
            return vote;
        }

        public Vote? UpdateCastedVote(Round container, Politician politician, VotePosition position)
        {
            if (!dBContext.VotingRounds.Contains(container))
                throw new Exception($"Voting round {container} not registered yet");
            var vote = container.VoteResults.FirstOrDefault(x => x.Politician == politician);
            if (vote == null)
            {
                return null;
            }
            dBContext.Update(vote);
            vote.Position = position;
            dBContext.SaveChanges();
            return vote;
        }

        public Round? UpdateVoteResult(Guid id, DateTime? time = null, List<Vote>? votes = null, string? Description = null)
        {
            var round = dBContext.VotingRounds.FirstOrDefault(x => x.Id == id);
            if (round == null)
            {
                return null;
            }
            dBContext.Update(round);
            if (time != null)
            {
                round.VoteDate = time.Value;
            }
            if (votes != null)
            {
                foreach( var vote in round.VoteResults)
                    dBContext.Remove(vote);
                round.VoteResults.Clear();
                foreach (var vote in votes)
                    round.VoteResults.Add(vote);
            }
            if (Description != null)
            {
                round.Description = Description;
            }
            dBContext.SaveChanges();
            return round;
        }

        /// <inheritdoc/>
        public ISet<Round> GetAllRoundsFromDB()
        {
            return dBContext.VotingRounds.Where( x => x.VoteId != 0).ToHashSet();
        }
    }
}
