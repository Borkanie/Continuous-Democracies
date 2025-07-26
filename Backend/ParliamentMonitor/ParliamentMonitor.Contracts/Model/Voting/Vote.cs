using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.ComponentModel.DataAnnotations.Schema;
using System.Linq;
using System.Text;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace ParliamentMonitor.Contracts.Model.Votes
{
    public enum VotePosition
    {
        Yes,
        No,
        Abstain,
        Absent
    }
    public class Vote : Entity
    {
        [JsonIgnore]
        public Politician Politician { get; set; }

        [NotMapped]
        [JsonPropertyName("politicianId")]
        public Guid? PoliticianId { get; set; }

        public VotePosition Position { get; set; } = VotePosition.Absent;

        [JsonIgnore]
        public Round Round {get;set;}


        [NotMapped]
        [JsonPropertyName("roundId")]
        public Guid? RoundId { get; set; }
    }
}
