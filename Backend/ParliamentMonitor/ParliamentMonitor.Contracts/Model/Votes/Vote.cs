using System;
using System.Collections.Generic;
using System.ComponentModel.DataAnnotations;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentMonitor.Contracts.Model.Votes
{
    public class Vote : Entity
    {

        public string Title { get => Name; set => Name = value; }
        public string Description { get; set; }
        public DateTime VoteDate { get; set; }

        public List<VoteResult> VoteResults { get; set; }
    }
}
