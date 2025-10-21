using System;
using System.Collections.Generic;
using System.Linq;
using System.Threading.Tasks;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using Moq;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;// adjust namespace if needed
using ParliamentMonitor.ServiceImplementation;
using Xunit;

namespace ParliamentMonitor.Tests.Unit
{
    public class PartyServiceUnitTests
    {
        private Mock<IAppDbContext> _mockContext = null!;
        private Mock<DbSet<Party>> _mockSet = null!;
        private ILogger<IPartyService<Party>> _logger = null!;

        private List<Party> _data = null!;

        public PartyServiceUnitTests()
        {
            _data = new List<Party>
            {
                new Party { Id = Guid.NewGuid(), Name = "Green Party", Acronym = "GP", Active = true },
                new Party { Id = Guid.NewGuid(), Name = "Democratic Party", Acronym = "DP", Active = true },
                new Party { Id = Guid.NewGuid(), Name = "Liberal Party", Acronym = "LP", Active = false }
            };

            _mockSet = CreateMockDbSet(_data.AsQueryable());

            _mockContext = new Mock<IAppDbContext>();
            _mockContext.Setup(c => c.Parties).Returns(_mockSet.Object);
            _mockContext.Setup(c => c.SaveChanges()).Returns(1);

            _logger = new LoggerFactory().CreateLogger<IPartyService<Party>>();
        }

        private static Mock<DbSet<T>> CreateMockDbSet<T>(IQueryable<T> data) where T : class
        {
            var mockSet = new Mock<DbSet<T>>();
            mockSet.As<IQueryable<T>>().Setup(m => m.Provider).Returns(data.Provider);
            mockSet.As<IQueryable<T>>().Setup(m => m.Expression).Returns(data.Expression);
            mockSet.As<IQueryable<T>>().Setup(m => m.ElementType).Returns(data.ElementType);
            mockSet.As<IQueryable<T>>().Setup(m => m.GetEnumerator()).Returns(data.GetEnumerator());

            mockSet.Setup(d => d.Add(It.IsAny<T>())).Callback<T>((s) => {
                // cannot modify underlying IQueryable; tests manipulate backing list directly
            });

            mockSet.Setup(d => d.Remove(It.IsAny<T>())).Callback<T>((s) => {
                // placeholder
            });

            return mockSet;
        }

        [Fact]
        public async Task CreatePartyAsync_NullOrWhitespaceName_ReturnsNull()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var res1 = await service.CreatePartyAsync("   ");
            var res2 = await service.CreatePartyAsync(null!);
            Assert.Null(res1);
            Assert.Null(res2);
        }

        [Fact]
        public async Task CreatePartyAsync_Duplicate_ReturnsNull()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var res = await service.CreatePartyAsync("Green Party", "GP");
            Assert.Null(res);
        }

        [Fact]
        public async Task CreatePartyAsync_Valid_AddsAndReturns()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var newParty = new Party {Id = Guid.NewGuid(), Name = "New Party", Acronym = "NP", Active = true };
            // Ensure backing list is mutable and the mock DbSet points to it
            var backing = _data;
            var mockSet = CreateMockDbSet(backing.AsQueryable());
            _mockContext.Setup(c => c.Parties).Returns(mockSet.Object);
            mockSet.Setup(d => d.Add(It.IsAny<Party>())).Callback<Party>(p => backing.Add(p));

            var res = await service.CreatePartyAsync(newParty.Name, newParty.Acronym);
            Assert.NotNull(res);
            Assert.Contains(backing, p => p.Name == "New Party" && p.Acronym == "NP");
        }

        [Fact]
        public async Task GetAsync_InvalidId_ReturnsNull()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var res = await service.GetAsync(default);
            Assert.Null(res);
        }

        [Fact]
        public async Task GetAsync_ValidId_ReturnsParty()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var id = _data[0].Id;
            var res = await service.GetAsync(id);
            Assert.NotNull(res);
            Assert.Equal("Green Party", res!.Name);
        }

        [Fact]
        public async Task GetAllPartiesAsync_ReturnsFilteredAndLimited()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var res = await service.GetAllPartiesAsync(isActive: true, number: 10);
            Assert.All(res, p => Assert.True(p.Active));
            Assert.Contains(res, p => p.Acronym == "GP");
        }

        [Fact]
        public async Task UpdateAsync_NullInput_ReturnsFalse()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var res = await service.UpdateAsync(null!);
            Assert.False(res);
        }

        [Fact]
        public async Task UpdateAsync_Nonexistent_ReturnsFalse()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var fake = new Party { Id = Guid.NewGuid(), Name = "Nope" };
            var res = await service.UpdateAsync(fake);
            Assert.False(res);
        }

        [Fact]
        public async Task UpdateAsync_Valid_Updates()
        {
            var backing = _data;
            var mockSet = CreateMockDbSet(backing.AsQueryable());
            _mockContext.Setup(c => c.Parties).Returns(mockSet.Object);
            mockSet.Setup(d => d.Remove(It.IsAny<Party>())).Callback<Party>(p => backing.Remove(p));

            var service = new PartyService(_mockContext.Object, _logger);
            var toUpdate = backing[0];
            toUpdate.Name = "Updated Name";

            var res = await service.UpdateAsync(toUpdate);
            Assert.True(res);
            Assert.Equal("Updated Name", backing[0].Name);
        }

        [Fact]
        public async Task DeleteAsync_NullOrMissing_ReturnsFalse()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var res1 = await service.DeleteAsync(null!);
            var res2 = await service.DeleteAsync(new Party { Id = Guid.NewGuid() });
            Assert.False(res1);
            Assert.False(res2);
        }

        [Fact]
        public async Task DeleteAsync_Existing_RemovesAndReturnsTrue()
        {
            var backing = _data;
            var mockSet = CreateMockDbSet(backing.AsQueryable());
            _mockContext.Setup(c => c.Parties).Returns(mockSet.Object);
            mockSet.Setup(d => d.Remove(It.IsAny<Party>())).Callback<Party>(p => backing.RemoveAll(x => x.Id == p.Id));

            var service = new PartyService(_mockContext.Object, _logger);
            var toDelete = backing[1];
            var res = await service.DeleteAsync(toDelete);
            Assert.True(res);
            Assert.DoesNotContain(backing, p => p.Id == toDelete.Id);
        }

        [Fact]
        public async Task GetPartyAsync_ByNameOrAcronym_Works()
        {
            var service = new PartyService(_mockContext.Object, _logger);
            var byName = await service.GetPartyAsync("Green Party", null);
            Assert.NotNull(byName);
            Assert.Equal("GP", byName!.Acronym);

            var byAcr = await service.GetPartyAsync(null, "DP");
            Assert.NotNull(byAcr);
            Assert.Equal("Democratic Party", byAcr!.Name);
        }
    }
}
