﻿using Microsoft.EntityFrameworkCore;
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

        public AppDBContext() : base(new DbContextOptionsBuilder<AppDBContext>()
                .UseNpgsql("Server=localhost;Port=5432;Database=parliment_8;Username=ps_user;Password=password123;Client Encoding=UTF8;")
                .Options)
        {
         

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
