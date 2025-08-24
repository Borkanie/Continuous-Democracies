using Microsoft.Extensions.Caching.Distributed;
using Microsoft.Extensions.Logging;
using ParliamentMonitor.Contracts.Model;
using System.Drawing;

namespace ParliamentMonitor.Contracts.Services
{
    /// <summary>
    /// Deals with operations related to political parties.
    /// </summary>
    public interface IPartyService<T> : IDBMService<T> where T: Party
    {
        /// <summary>
        /// Returns a party by its unique identifier.
        /// </summary>
        /// <returns>Null if not found.</returns>
        public Task<T?> GetPartyAsync(string? name = null,string? acronym = null);

        /// <summary>
        /// Creates a new party with the specified parameters.
        /// </summary>
        /// <returns>The created instance.</returns>
        public Task<T?> CreatePartyAsync(string name, string? acronym = null, string? logoUrl = null, Color? color = null);
        
        /// <summary>
        /// Updates a party identified by its unique identifier.
        /// All null parameters will be ignored.
        /// </summary>
        /// <returns>The updated party.</returns>
        public Task<T?> UpdatePartyAsync(Guid id, string? name = null, string? acronym = null, string? logoUrl = null, Color? color = null);

        /// <summary>
        /// Returns all parties in the system.
        /// </summary>
        public Task<IList<T>> GetAllPartiesAsync(bool isActive = true, int number = 100);
    }
}
