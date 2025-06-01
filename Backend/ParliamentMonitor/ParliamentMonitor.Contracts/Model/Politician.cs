
using System.ComponentModel.DataAnnotations;

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
    public class Politician: Entity
    {


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

        public override bool Equals(object? obj)
        {
            if (obj == null) return false;
            if(obj.GetType()!= typeof(Politician)) return false;
            return ((Politician)obj).Id == Id;
        }
    }
}
