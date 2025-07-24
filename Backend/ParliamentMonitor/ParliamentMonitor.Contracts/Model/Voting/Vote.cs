using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
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
        public required Politician Politician { get; set; }

        public VotePosition Position { get; set; } = VotePosition.Absent;

        [JsonIgnore]
        public required Round Round {get;set;}

    }
}
