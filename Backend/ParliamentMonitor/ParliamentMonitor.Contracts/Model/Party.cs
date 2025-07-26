using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using System.Drawing;
using System.Text.Json.Serialization;

namespace ParliamentMonitor.Contracts.Model
{
    /// <summary>
    /// Abstraction of a party in a political context.
    /// </summary>
    public class Party : Entity
    {
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
        [JsonIgnore]
        public List<Politician> Politicians { get; set; } = new List<Politician>();

        [NotMapped]
        [JsonPropertyName("politicianIds")]
        public HashSet<Guid> PoliticianIds { get; set; } = new HashSet<Guid>();

        public bool Active {  get; set; }

    }
}
