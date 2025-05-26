using System;
using System.Collections.Generic;
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
    public class VoteResult
    {
        public Guid Id { get; set; }

        public Vote Vote { get; set; }

        public Politician Politician { get; set; }

        public VotePosition Position { get; set; } = VotePosition.Absent;
    }
}
