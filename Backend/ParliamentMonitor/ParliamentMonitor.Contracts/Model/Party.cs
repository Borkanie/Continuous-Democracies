using System.Drawing;

namespace ParliamentMonitor.Contracts.Model
{
    /// <summary>
    /// Abstraction of a party in a political context.
    /// </summary>
    public class Party
    {
        /// <summary>
        /// Unique ID of the party, will be used to identyfy the entity.
        /// </summary>
        public Guid Id { get; }

        /// <summary>
        /// The Full Name of the Party.
        /// </summary>
        public String Name { get; set; } = String.Empty;

        /// <summary>
        /// The Acronym of the Party, used for quick identification.
        /// </summary>
        public String? Acronym { get; set; }

        /// <summary>
        /// path to the Party Logo, used for visual representation.
        /// </summary>
        public String? LogoUrl { get; set; }

        /// <summary>
        /// Defaul party color for UI representation.
        /// </summary>
        public Color Color { get; set; } = Color.LightGray;

        /// <summary>
        /// List of politicans affiliated with this party.
        /// </summary>
        public List<Politician> Politicians { get; set; }
    }
}
