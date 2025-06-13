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
    public class Vote : Entity
    {
        public Round Round { get; set; }

        public Politician Politician { get; set; }

        public VotePosition Position { get; set; } = VotePosition.Absent;

        public override bool Equals(object? obj)
        {
            if (obj == null) return false;
            if (obj.GetType() != GetType())
            {
                return false;
            }
            var ob = (Vote)obj;
            return ob.Id == Id || (ob.Round == Round && ob.Politician == Politician);
        }
    }
}
