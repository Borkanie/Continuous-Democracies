using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;

namespace ParliamentMonitor.Contracts.Services
{
    public interface IVotingService<T> : IDBMService<T> where T : Vote 
    {
        /// <summary>
        /// Creates a new vote for a given politician with a given position.
        /// </summary>
        public Task<T> CreateNewVote(Round round, Politician politician, VotePosition votePosition);

        /// <summary>
        /// Removes reference to vote from cache and db withouth touching round.
        /// </summary>
        /// <param name="id"></param>
        /// <returns></returns>
        public Task<bool> DeleteWithouthRemovingFromRound(Guid id);

        public void SetRoundService(IVotingRoundService<Round> votingRoundService);
    }
}
