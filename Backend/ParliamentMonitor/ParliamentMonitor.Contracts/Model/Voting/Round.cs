using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentMonitor.Contracts.Model.Votes
{
    public class Round : Entity
    {

        public string Title { get => Name; set => Name = value; }
        public string Description { get; set; } = string.Empty;
        public DateTime VoteDate { get; set; } = DateTime.Now;

        public ISet<Vote> VoteResults { get; set; } = new HashSet<Vote>();

        public int VoteId { get; set; }

        public override bool Equals(object? obj)
        {
            if (obj == null) return false;
            if (obj.GetType() != GetType())
            {
                return false;
            }
            var ob = (Round)obj;
            return ob.Id == Id || (Title == ob.Title && ob.VoteDate == VoteDate);
        }
    }
}
