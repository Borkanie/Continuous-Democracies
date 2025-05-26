using ParliamentMonitor.Contracts.Model;

namespace ParliamentMonitor.Contracts.Services
{
    /// <summary>
    /// Deals with operations related to politicians.
    /// </summary>
    public interface IPoliticianService
    {
        /// <summary>
        /// Returns a politician by its unique identifier.
        /// </summary>
        Politician? GetPolitician(Guid id);

        /// <summary>
        /// Creates a new politician with the specified parameters.
        /// </summary>
        /// <returns>The new instance of a politician.</returns>
        Politician CreatePolitican(string name, Party party, WorkLocation location, Gender gender, bool isCurrentlyActive = true);

        /// <summary>
        /// Updates the name of a given politician identified by its unique identifier.
        /// </summary>
        /// <returns>Returns the new instance.</returns>
        Politician UpdatePoliticianName(Guid id, string name);

        /// <summary>
        /// Updates the party of a given politician identified by its unique identifier.
        /// </summary>
        /// <returns>Returns the new instance.</returns>
        Politician UpdatePoliticianParty(Guid id, Party party);

        /// <summary>
        /// Updates the activity status of a given politician identified by its unique identifier.
        /// </summary>
        /// <returns>Returns the new instance.</returns>
        Politician UpdatePoliticianActivity(Guid id, bool isCurrentlyActive);

        /// <summary>
        /// Updates the gender of a given politician identified by its unique identifier.
        /// </summary>
        /// <returns>Returns the new instance.</returns>
        Politician UpdatePoliticianGender(Guid id, Gender gender);

        /// <summary>
        /// updates the working location of a given politician identified by its unique identifier.
        /// </summary>
        /// <returns>Returns the new instance.</returns>
        Politician UpdatePoliticianLocation(Guid id, WorkLocation location);

        /// <summary>
        /// Returns all the avialable politicians in the system.
        /// </summary>
        IList<Politician> GetAllPoliticians();

        /// <summary>
        /// Returns all the politicans from a given party.
        /// </summary>
        IList<Politician> GetPoliticiansByParty(Party party);

        /// <summary>
        /// Returns all the politicians that are currently active.
        /// </summary>
        IList<Politician> GetActivePoliticians();

        /// <summary>
        /// Returns all the politicians that are currently inactive.
        /// </summary>
        IList<Politician> GetInactivePoliticians();

        /// <summary>
        /// Returns all the politicians that are located in a specific work location.
        /// </summary>
        IList<Politician> GetPoliticiansByLocation(WorkLocation location);

        /// <summary>
        /// Returns all the politicians of a specific gender.
        /// </summary>
        IList<Politician> GetAllPoliticians(Gender gender);

        /// <summary>
        /// Returns all the politicians from a given gender affiliated with a specific party.
        /// </summary>
        IList<Politician> GetPoliticiansByParty(Party party, Gender gender);

        /// <summary>
        /// Returns all the politicians that are currently active of a given gender.
        /// </summary>
        IList<Politician> GetActivePoliticians(Gender gender);

        /// <summary>
        /// Returns all the politicians that are currently inactive of a given gender.
        /// </summary>
        IList<Politician> GetInactivePoliticians(Gender gender);

        /// <summary>
        /// Return all the politicians that are located in a specific work location and are of a given gender.
        /// </summary>
        IList<Politician> GetPoliticiansByLocation(WorkLocation location, Gender gender);
    }
}
