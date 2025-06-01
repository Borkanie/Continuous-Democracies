using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
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
    public class VoteResult : Entity
    {
        public Vote Vote { get; set; }

        public Politician Politician { get; set; }

        public VotePosition Position { get; set; } = VotePosition.Absent;
    }
}
