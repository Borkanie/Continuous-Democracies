namespace ParliamentMonitor.Contracts.Model.Votes
{
    public class Round : Entity
    {
        public string Title { get => Name; set => Name = value; }
        public string Description { get; set; } = string.Empty;
        public DateTime VoteDate { get; set; } = DateTime.Now;
        public int VoteId { get; set; }
    }
}
