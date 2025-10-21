using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using System.Threading;
using System.Threading.Tasks;

namespace ParliamentMonitor.DataBaseConnector
{
    public interface IAppDbContext
    {
        DbSet<Politician> Politicians { get; }
        DbSet<Party> Parties { get; }
        DbSet<Round> VotingRounds { get; }
        DbSet<Vote> Votes { get; }

        int SaveChanges();
        Task<int> SaveChangesAsync(CancellationToken cancellationToken = default);
    }
}
