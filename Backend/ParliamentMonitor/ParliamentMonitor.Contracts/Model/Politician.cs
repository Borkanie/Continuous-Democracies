
namespace ParliamentMonitor.Contracts.Model
{
    public enum Gender
    {
        Male,
        Female,
        Other
    }

    public enum WorkLocation
    {
        Parliament,
        Senate,
        Other
    }

    /// <summary>
    /// Abstraction of a politician in a political context.
    /// </summary>
    public class Politician
    {
        /// <summary>
        /// Primary key of the politician, will be used to identify the entity.
        /// </summary>
        public Guid Id { get; }

        /// <summary>
        /// The full name of the politician.
        /// </summary>
        public String Name { get; set; } = String.Empty;

        /// <summary>
        /// Gender of the politician, used for demographic purposes.
        /// </summary>
        public Gender Gender { get; set; } = Gender.Male;

        /// <summary>
        /// Image of the politician, used for visual representation.
        /// </summary>
        public String? ImageUrl { get; set; }

        /// <summary>
        /// The party the politician belongs to.
        /// </summary>
        public Party Party { get; set; }

        /// <summary>
        /// True if the politician is currently in active an official capacity, false otherwise.
        /// </summary>
        public bool Active { get; set; } = false;

        /// <summary>
        /// Worklocation of the politician, indicating where they primarily operate.For now, Parliament or Senate.
        /// </summary>
        public WorkLocation WorkLocation { get; set; } = WorkLocation.Parliament;
    }
}
