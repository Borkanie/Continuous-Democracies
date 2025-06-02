using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParlimentMonitor.ServiceImplementation
{
    public class VotingService : IVotingService<Vote, VotingRound>
    {
        private readonly AppDBContext dBContext;

        public VotingService(AppDBContext context)
        {
            dBContext = context;
        }

        public Vote CastVote(VotingRound container, Politician politician, VotePosition position)
        {
            throw new NotImplementedException();
        }

        public VotingRound CreateVotingRound(string title, DateTime time, List<Vote>? votes = null, string? description = null)
        {
            var round = new VotingRound();
            round.Title = title;
            round.VoteDate = time;
            if(votes!= null)
            {
                foreach(var vote in votes)
                    round.VoteResults.Add(vote);
            }
            round.Description = description ?? string.Empty;
            dBContext.Votes.Add(round);
            return round;
        }

        public void Delete(VotingRound entity)
        {
            if(dBContext.Votes.Contains(entity))
            {
                dBContext.Votes.Remove(entity);
                dBContext.SaveChanges();
            }
        }

        public void Delete(Vote entity)
        {
            throw new NotImplementedException();
        }

        public void Update(VotingRound voteResult)
        {
            throw new NotImplementedException();
        }

        public void Update(Vote entity)
        {
            throw new NotImplementedException();
        }

        public Vote? UpdateCastedVote(VotingRound container, Politician politician, VotePosition position)
        {
            throw new NotImplementedException();
        }

        public VotingRound? UpdateVoteResult(Guid id, DateTime? time = null, List<Vote>? votes = null, string? Description = null)
        {
            throw new NotImplementedException();
        }
    }
}
