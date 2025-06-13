using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentMonitor.Contracts.Services
{
    public interface IVotingService<T,Y> : IDBMService<T> where T : Vote where Y : Round
    {
        /// <summary>
        /// Udpates a VotingRound result. 
        /// It will replace the previous vote round result with the udpated item based on Guid.
        /// </summary>
        /// <param name="voteResult">Updated vote result instance.</param>
        public void Update(Y voteResult);

        public void Delete(Y entity);

        /// <summary>
        /// Creates a new voting round in the DB.
        /// </summary>
        /// <param name="title">Title of the vote that happened.</param>
        /// <param name="time">When it took ended.</param>
        /// <param name="votes">The list of votes for all present individuals.</param>
        /// <param name="Description">Short descriptiona bout the vote.</param>
        public Y CreateVotingRound(string title, DateTime time, List<T>? votes = null, string? Description = null);

        /// <summary>
        /// Update a vote result in the db.
        /// Null values will be ignored.
        /// </summary>
        /// <param name="title">Title of the vote that happened.</param>
        /// <param name="time">When it took ended.</param>
        /// <param name="votes">The list of votes for all present individuals.</param>
        /// <param name="Description">Short descriptiona bout the vote.</param>
        /// <returns></returns>
        public Y? UpdateVoteResult(Guid id, DateTime? time = null, List<T>? votes = null, string? Description = null);


        /// <summary>
        /// Casts a vote for a given politican in a given voting round.
        /// A politician can only give a single vote per round.
        /// </summary>
        public T CastVote(Round container, Politician politician, VotePosition position);

        /// <summary>
        /// Changes the position of apolitican in a given round.
        /// </summary>
        public T? UpdateCastedVote(Round container, Politician politician, VotePosition position);

        /// <summary>
        /// Return voting round based off voting number.
        /// </summary>
        /// <param name="votingRoundId"></param>
        /// <returns></returns>
        public Y? GetVotingRound(int votingRoundId);
    }
}
