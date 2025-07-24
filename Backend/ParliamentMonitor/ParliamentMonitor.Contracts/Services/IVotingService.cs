using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentMonitor.Contracts.Services
{
    public interface IVotingService<T> : IDBMService<T> where T : Vote 
    {
        /// <summary>
        /// Creates a new vote for a given politician with a given position.
        /// </summary>
        public Task<T> CreateNewVote(Round round, Politician politician, VotePosition votePosition);

    }
}
