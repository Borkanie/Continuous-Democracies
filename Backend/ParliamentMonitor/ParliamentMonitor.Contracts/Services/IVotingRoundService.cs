using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentMonitor.Contracts.Services
{
    public interface IVotingRoundService<T> : IDBMService<T> where T : Round
    {
        /// <summary>
        /// Casts a vote for a given politican in a given voting round.
        /// A politician can only give a single vote per round.
        /// </summary>
        public Task<Vote?> RegisterVoteAsync(T container, Politician politician, VotePosition position);

        /// <summary>
        /// Creates a new voting round in the DB.
        /// </summary>
        /// <param name="title">Title of the vote that happened.</param>
        /// <param name="time">When it took ended.</param>
        /// <param name="votes">The list of votes for all present individuals.</param>
        /// <param name="Description">Short descriptiona bout the vote.</param>
        public Task<T?> CreateVotingRoundAsync(string title, DateTime time, int id = 0, string? Description = null);

        /// <summary>
        /// Return voting round based off voting number.
        /// </summary>
        /// <param name="votingRoundId"></param>
        /// <returns></returns>
        public Task<T?> GetVotingRoundAsync(int votingRoundId);

        /// <summary>
        /// Returns all the rounds from DB.
        /// Unpopulated with votes they need to be fetched separately for speed.
        /// This will return all the rounds in a given time interval.
        /// </summary>
        /// <returns></returns>
        public Task<ISet<T>> GetAllRoundsFromDBAsync(DateTime? startDate, DateTime? endDate, string?[] keywords, int number = 100);
    }
}
