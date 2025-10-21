using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using System.Drawing;
using System.Threading.Tasks;
using System.Linq;

namespace ParliamentMonitor.ServiceImplementation
{
    public class PartyService : IPartyService<Party>
    {
        private readonly IAppDbContext _dbContext;
        private readonly ILogger<IPartyService<Party>> _logger;

        public PartyService(IAppDbContext context, ILogger<IPartyService<Party>> logger)
        {
            _dbContext = context;
            _logger = logger;
        }

        public ILogger Logger { get => _logger; }

        /*<inheritdoc/>*/
        public Task<Party?> CreatePartyAsync(string name, string? acronym = null, string? logoUrl = null, Color? color = null)
        {
            if (string.IsNullOrWhiteSpace(name))
                return Task.FromResult<Party?>(null);

            var newParty = new Party() { Id = Guid.NewGuid() };
            newParty.Name = name;
            if (acronym != null)
                newParty.Acronym = acronym;
            if (logoUrl != null)
                newParty.LogoUrl = logoUrl;
            if (color != null)
                newParty.Color = (Color)color;
            if (_dbContext.Parties.Any(x => x.Name == newParty.Name && x.Acronym == newParty.Acronym))
            {
                return Task.FromResult<Party?>(null);
            }

            _dbContext.Parties.Add(newParty);
            _dbContext.SaveChanges();
            _logger.LogInformation($"Created new party:{newParty.Name}");
            return Task.FromResult<Party?>(newParty);
        }

        /*<inheritdoc/>*/
        public Task<Party?> GetAsync(Guid id)
        {
            if (id == default) return Task.FromResult<Party?>(null);
            var party = _dbContext.Parties.FirstOrDefault(p => p.Id == id);
            return Task.FromResult(party);
        }

        /*<inheritdoc/>*/
        public Task<bool> UpdateAsync(Party item)
        {
            if (item == null) return Task.FromResult(false);
            var oldItem = _dbContext.Parties.FirstOrDefault(p => p.Id == item.Id);
            if (oldItem != null)
            {
                oldItem.Active = item.Active;
                oldItem.Acronym = item.Acronym;
                oldItem.LogoUrl = item.LogoUrl;
                oldItem.Color = item.Color;
                _dbContext.SaveChanges();
                _logger.LogInformation($"Updated party:{item.Id}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /*<inheritdoc/>*/
        public Task<IList<Party>> GetAllPartiesAsync(bool isActive = true, int number = 100)
        {
            var parties = _dbContext.Parties.Where(x => x.Active == isActive).Take(number).ToList();
            return Task.FromResult<IList<Party>>(parties ?? new List<Party>());
        }

        /*<inheritdoc/>*/
        public Task<Party?> UpdatePartyAsync(Guid id, string? name = null, string? acronym = null, string? logoUrl = null, Color? color = null)
        {
            var item = _dbContext.Parties.FirstOrDefault(p => p.Id == id);
            if (item != null)
            {
                item.Name = name ?? item.Name;
                item.Acronym = acronym ?? item.Acronym;
                item.LogoUrl = logoUrl ?? item.LogoUrl;
                item.Color = color ?? item.Color;
                _dbContext.SaveChanges();
                _logger.LogInformation($"Updated party:{item.Id}");
            }
            return Task.FromResult(item);
        }

        /*<inheritdoc/>*/
        public Task<bool> DeleteAsync(Party entity)
        {
            if (entity == null) return Task.FromResult(false);
            var existing = _dbContext.Parties.FirstOrDefault(p => p.Id == entity.Id);
            if (existing != null)
            {
                _dbContext.Parties.Remove(existing);
                _dbContext.SaveChanges();
                _logger.LogInformation($"Deleted party:{entity.Id}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /*<inheritdoc/>*/
        public Task<Party?> GetPartyAsync(string? name = null, string? acronym = null)
        {
            if (!string.IsNullOrWhiteSpace(name))
            {
                var value = _dbContext.Parties.FirstOrDefault(x => string.Equals(x.Name.Trim(), name.Trim()));
                return Task.FromResult(value);
            }
            if (!string.IsNullOrWhiteSpace(acronym))
            {
                var value = _dbContext.Parties.FirstOrDefault(x => string.Equals(x.Acronym, acronym));
                return Task.FromResult(value);

            }
            return Task.FromResult<Party?>(null);
        }
    }
}
