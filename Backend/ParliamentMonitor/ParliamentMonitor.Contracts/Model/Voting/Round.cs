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
    public class Round : Entity
    {
        public string Title { get => Name; set => Name = value; }
        public string Description { get; set; } = string.Empty;
        public DateTime VoteDate { get; set; } = DateTime.Now;

        [JsonIgnore]
        public ISet<Vote> VoteResults { get; set; } = new HashSet<Vote>();
        public int VoteId { get; set; }

        [NotMapped]
        [JsonPropertyName("voteResultIds")]
        public HashSet<Guid> VoteResultIds { get;set; } = new HashSet<Guid>();  
    }
}
