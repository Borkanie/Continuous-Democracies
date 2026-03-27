using Microsoft.EntityFrameworkCore;
using Microsoft.EntityFrameworkCore.Storage.ValueConversion;
using Microsoft.Extensions.Configuration;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using System;
using System.Collections.Generic;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentMonitor.DataBaseConnector 
{
    public class AppDBContext : DbContext
    {
        public AppDBContext(DbContextOptions<AppDBContext> options) : base(options) { }

        public AppDBContext(string connectionString) : base(new DbContextOptionsBuilder<AppDBContext>().UseNpgsql(connectionString).Options) {}

        public static string getConnStringFromEnvVariables()
        {
            var dbServer = Environment.GetEnvironmentVariable("DBSERVERADDRESS");
            var dbName = Environment.GetEnvironmentVariable("DBNAME");
            var dbUser = Environment.GetEnvironmentVariable("DBUSER");
            var dbPassword = Environment.GetEnvironmentVariable("DBPASSWORD");

            string connectionString;
            if (!string.IsNullOrEmpty(dbServer) && !string.IsNullOrEmpty(dbName) && !string.IsNullOrEmpty(dbUser))
            {
                connectionString = $"Server={dbServer};Port=5432;Database={dbName};Username={dbUser};Password={dbPassword};";
                Console.WriteLine("Using connection string from environment variables");
            }
            else
            {
                throw new Exception("Database environment variables are not set. Please set DBSERVERADDRESS, DBNAME, DBUSER, and DBPASSWORD.");
            }
            return connectionString;
        }

        protected override void OnModelCreating(ModelBuilder modelBuilder)
        {
            // Color <-> string converter (ARGB hex string)
            var colorConverter = new ValueConverter<Color, string>(
                v => ColorTranslator.ToHtml(v),
                v => ColorTranslator.FromHtml(v));

            modelBuilder.Entity<Party>()
                .Property(p => p.Color)
                .HasConversion(colorConverter);

            base.OnModelCreating(modelBuilder);
        }

        public DbSet<Politician> Politicians { get; set; }

        public DbSet<Party> Parties { get; set; }

        public DbSet<Round> VotingRounds { get; set; }

        public DbSet<Vote> Votes { get; set; }
    }
}
