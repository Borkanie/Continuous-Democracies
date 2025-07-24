using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using StackExchange.Redis;
using System.Drawing;
using static System.Runtime.InteropServices.JavaScript.JSType;

namespace ParliamentMonitor.ServiceImplementation
{
    public class PartyService : RedisService, IPartyService<Party>
    {
        private readonly AppDBContext _dbContext;

        public PartyService(AppDBContext context, IConnectionMultiplexer redis) : base(redis,"party")
        {
            _dbContext = context;
        }

        /*<inheritdoc/>*/
        public async Task<Party?> CreatePartyAsync(string name, string? acronym = null, string? logoUrl = null, Color? color = null)
        {
            try
            {
                var newParty = new Party();
                newParty.Id = Guid.NewGuid();
                newParty.Name = name;
                if (acronym != null)
                    newParty.Acronym = acronym;
                if (logoUrl != null)
                    newParty.LogoUrl = logoUrl;
                if (color != null)
                    newParty.Color = (Color)color;
                if (newParty == null || _dbContext.Parties.Any(x => x == newParty))
                {
                    throw new Exception("Already existing party");
                }
                _dbContext.Parties.Add(newParty);
                _dbContext.SaveChanges();
                SetAsync(MakeKey(newParty.Id.ToString()), newParty);
                return newParty;
            }catch(Exception ex)
            {
                Console.WriteLine($"Error when creating new party:{ex.Message}");
                return null; 
            }
            
        }

        /*<inheritdoc/>*/
        public async Task<Party?> GetAsync(Guid id)
        {
            var cached = await GetAsync<Party>(MakeKey(id.ToString()));
            if (cached != null)
                return cached;
            return _dbContext.Parties.Find(id);
        }

        /*<inheritdoc/>*/
        public async Task<bool> UpdateAsync(Party item)
        {
            var oldItem = GetAsync(item.Id).Result;
            if (oldItem != null)
            {
                _dbContext.Update(item);
                oldItem.Active = item.Active;
                oldItem.Acronym = item.Acronym;
                oldItem.LogoUrl = item.LogoUrl;
                oldItem.Politicians = item.Politicians;
                oldItem.Color = item.Color;
                _dbContext.SaveChanges();
                SetAsync(MakeKey(item.Id.ToString()), item);
                return true;
            }
            return false;
        }

        /*<inheritdoc/>*/
        public async Task<IList<Party>> GetAllPartiesAsync(bool isActive = true, int number = 100)
        {
            var parties = getPartiesAsync(isActive, number).Result;
            if (parties.Count == number)
                return parties;
            return _dbContext.Parties.Where(x => x.Active == isActive).Take(number).ToList();
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
                SetAsync(MakeKey(item.Id.ToString()), item);
            }
            return Task.FromResult(item);
        }

        /*<inheritdoc/>*/
        public async Task<bool> DeleteAsync(Party entity)
        {
            if (_dbContext.Parties.Find(entity.Id) != null)
            {
                _dbContext.Parties.Remove(entity);
                _dbContext.SaveChanges();
                RemoveAsync(MakeKey(entity.Id.ToString()));
                return true;
            }
            return false;
        }

        /*<inheritdoc/>*/
        public async Task<Party?> GetPartyAsync(string? name = null, string? acronym = null)
        {
            var party = getPartiesAsync(true).Result.FirstOrDefault(x => (name != null && string.Equals(x.Name, name))
                                                    || ( acronym != null && string.Equals(x.Acronym, acronym)));
            if (party != null)
                return party;
            if (name != null)
            {
                var value = _dbContext.Parties.FirstOrDefault(x => string.Equals(x.Name, name));
                SetAsync(MakeKey(value.Id.ToString()), value);
                return value;
            }
            if (acronym != null)
            {
                var value = _dbContext.Parties.FirstOrDefault(x => string.Equals(x.Acronym, acronym));
                SetAsync(MakeKey(value.Id.ToString()), value);
                return value;

            }
            return null;
        }

        private async Task<IList<Party>> getPartiesAsync(bool isActive = true, int number = 100)
        {
            var parties = new List<Party>();
            var keys = await _cache.SetMembersAsync(serviceKey);
            foreach (var key in keys)
            {
                var party = await GetAsync<Party>(key);
                if (party != null)
                {
                    if (party.Active != isActive)
                    {
                        continue;
                    }
                    parties.Add(party);
                    if (parties.Count == number)
                        return parties;
                }
            }
            return new List<Party>();

        }
    }
}
