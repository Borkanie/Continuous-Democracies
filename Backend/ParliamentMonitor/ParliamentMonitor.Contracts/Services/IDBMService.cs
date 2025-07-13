using ParliamentMonitor.Contracts.Model;

namespace ParliamentMonitor.Contracts.Services
{
    public interface IDBMService<T> where T : Entity
    {
        /// <summary>
        /// Up[dates an entity in the db based on Key.
        /// </summary>
        /// <param name="entity"></param>
        public Task<bool> Update(T entity);

        /// <summary>
        /// Removes an entity form the DB based on Key.
        /// </summary>
        /// <param name="entity"></param>
        public Task<bool> Delete(T entity);
    }
}
