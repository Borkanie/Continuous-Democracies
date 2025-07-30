using ParliamentMonitor.Contracts.Model;

namespace ParliamentMonitor.Contracts.Services
{
    /// <summary>
    /// Deals with operations related to politicians.
    /// </summary>
    public interface IPoliticianService<T> : IDBMService<T> where T : Politician
    {
        /// <summary>
        /// Returns a politician by its unique identifier.
        /// </summary>
        public Task<T?> GetPoliticianAsync(String name);

        /// <summary>
        /// Creates a new politician with the specified parameters.
        /// </summary>
        /// <returns>The new instance of a politician.</returns>
        public Task<T?> CreatePoliticanAsync(string name, Party party, WorkLocation location, Gender gender, bool isCurrentlyActive = true, string? imageUrl = null);

        /// <summary>
        /// Updates a given polititacns information based on they're Id. All null informations will be ignored.
        /// </summary>
        /// <param name="id"></param>
        /// <param name="name"></param>
        /// <param name="party"></param>
        /// <param name="location"></param>
        /// <param name="gender"></param>
        /// <param name="isCurrentlyActive"></param>
        /// <param name="imageUrl"></param>
        /// <returns></returns>
        public Task<T?> UpdatePoliticianAsync(Guid id, string? name = null, Party? party = null, WorkLocation? location = null, Gender? gender = null, bool? isCurrentlyActive = null, string? imageUrl = null);

        /// <summary>
        /// Returns all the avialable politicians in the system.
        /// </summary>
        public Task<IList<T>> GetAllPoliticiansAsync(Party? party = null,bool? isActive = null, WorkLocation? location = null,Gender? gender = null, int number = 100);

    }
}
