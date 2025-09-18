using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using System.Drawing;

namespace ParliamentMonitor.ServiceImplementation
{
    public class PartyService(AppDBContext context, ILogger<IPartyService<Party>> logger) : IPartyService<Party>
    {
        private readonly AppDBContext _dbContext = context;
        private ILogger<IPartyService<Party>> _logger = logger;

        public ILogger Logger { get => _logger; }

        /*<inheritdoc/>*/
        public Task<Party?> CreatePartyAsync(string name, string? acronym = null, string? logoUrl = null, Color? color = null)
        {
            var newParty = new Party() { Id = Guid.NewGuid() };
            newParty.Id = Guid.NewGuid();
            newParty.Name = name;
            if (acronym != null)
                newParty.Acronym = acronym;
            if (logoUrl != null)
                newParty.LogoUrl = logoUrl;
            if (color != null)
                newParty.Color = (Color)color;
            if (newParty == null || _dbContext.Parties.Any(x => x.Name == newParty.Name && x.Acronym == newParty.Acronym))
            {
                throw new Exception($"The party:{newParty.Name} already exists.");
            }
            _dbContext.Parties.Add(newParty);
            _dbContext.SaveChanges();
            _logger.LogInformation($"Created new party:{newParty.Name}");
            return Task.FromResult<Party?>(newParty);
        }

        /*<inheritdoc/>*/
        public Task<Party?> GetAsync(Guid id)
        {         
            var party = _dbContext.Parties.Find(id);
            return Task.FromResult(party);
        }

        /*<inheritdoc/>*/
        public Task<bool> UpdateAsync(Party item)
        {
            var oldItem = GetAsync(item.Id).Result;
            if (oldItem != null)
            {
                _dbContext.Update(item);
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
            if (parties == null)
            {
                return Task.FromResult<IList<Party>>(new List<Party>());
            }
            return Task.FromResult<IList<Party>>(parties);
        }

        /*<inheritdoc/>*/
        public Task<Party?> UpdatePartyAsync(Guid id, string? name = null, string? acronym = null, string? logoUrl = null, Color? color = null)
        {
            var item = GetAsync(id).Result;
            if (item != null)
            {
                _dbContext.Update(item);
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
            if (_dbContext.Parties.Find(entity.Id) != null)
            {
                _dbContext.Parties.Remove(entity);
                _dbContext.SaveChanges();
                _logger.LogInformation($"Deleted party:{entity.Id}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /*<inheritdoc/>*/
        public Task<Party?> GetPartyAsync(string? name = null, string? acronym = null)
        {
            if (name != null)
            {
                var value = _dbContext.Parties.FirstOrDefault(x => string.Equals(x.Name.Trim(), name.Trim()));
                if (value == null)
                    return Task.FromResult<Party?>(null);
                return Task.FromResult<Party?>(value);
            }
            if (acronym != null)
            {
                var value = _dbContext.Parties.FirstOrDefault(x => string.Equals(x.Acronym, acronym));
                if(value == null)
                    return Task.FromResult<Party?>(null);
                return Task.FromResult<Party?>(value);

            }
            return Task.FromResult<Party?>(null);
        }
    }
}
