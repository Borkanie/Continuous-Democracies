using Microsoft.Extensions.Logging;
using ParliamentMonitor.Contracts.Model;

namespace ParliamentMonitor.Contracts.Services
{
    public interface IDBMService<T> where T : Entity 
    {
        /// <summary>
        /// Get's logger isntance related to the service.
        /// </summary>
        public ILogger Logger { get; }

        /// <summary>
        /// Up[dates an entity in the db based on Key.
        /// </summary>
        public Task<bool> UpdateAsync(T entity);

        /// <summary>
        /// Removes an entity form the DB based on Key.
        /// </summary>
        public Task<bool> DeleteAsync(T entity);

        /// <summary>
        /// Returns an entity if it finds any in cache/db.
        /// </summary>
        public Task<T?> GetAsync(Guid id);

    }
}
