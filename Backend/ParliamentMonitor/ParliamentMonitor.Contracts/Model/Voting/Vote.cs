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
        public Politician Politician { get; set; }

        public VotePosition Position { get; set; } = VotePosition.Absent;

        public Round Round {get;set;}
    }
}
