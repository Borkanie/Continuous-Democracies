using ParliamentMonitor.Contracts.Model;
using System.Drawing;

namespace ParliamentMonitor.Contracts.Services
{
    /// <summary>
    /// Deals with operations related to political parties.
    /// </summary>
    internal interface IPartyService
    {
        /// <summary>
        /// Returns a party by its unique identifier.
        /// </summary>
        /// <returns>Null if not found.</returns>
        public Party? GetParty(Guid id);

        /// <summary>
        /// Creates a new party with the specified parameters.
        /// </summary>
        /// <returns>The created instance.</returns>
        public Party CreateParty(string name, string? acronym = null, string? logoUrl = null, Color? color = null);

        /// <summary>
        /// Updates the name of a given party identified by its unique identifier.
        /// </summary>
        /// <returns>The updated party.</returns>
        public Party UpdatePartyName(Guid id, string name);

        /// <summary>
        /// Updates the acronym of a given party identified by its unique identifier.
        /// </summary>
        /// <returns>The updated party.</returns>
        public Party UpdatePartyAcronym(Guid id, string acronym);

        /// <summary>
        /// Updates the logo URL of a given party identified by its unique identifier.
        /// </summary>
        public Party UpdatePartyLogo(Guid id, string logoUrl);

        /// <summary>
        /// Updates the color of a given party identified by its unique identifier.
        /// </summary>
        public Party UpdatePartyColor(Guid id, Color color);

        /// <summary>
        /// Returns all parties in the system.
        /// </summary>
        public IList<Party> GetAllParties();

        /// <summary>
        /// Returns all parties that are currently active.
        /// </summary>
        public IList<Party> GetActiveParties();

        /// <summary>
        /// Returns all parties that are currently inactive.
        /// </summary>
        public IList<Party> GetInactiveParties();
    }
}
