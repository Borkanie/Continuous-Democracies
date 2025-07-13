using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using StackExchange.Redis;
using System.Drawing;

namespace ParliamentMonitor.ServiceImplementation
{
    public class PartyService : IPartyService<Party>
    {
        internal readonly IConnectionMultiplexer _redis;
        internal IDatabase _cache;
        private readonly AppDBContext _dbContext;

        public PartyService(AppDBContext context, IConnectionMultiplexer redis)
        {
            _dbContext = context;
            _redis = redis;
            _cache = redis.GetDatabase();
        }

        /*<inheritdoc/>*/
        public async Task<Party?> CreateParty(string name, string? acronym = null, string? logoUrl = null, Color? color = null)
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
                return newParty;
            }catch(Exception ex)
            {
                Console.WriteLine($"Error when creating new party:{ex.Message}");
                return null; 
            }
            
        }


        /*<inheritdoc/>*/
        public async Task<Party?> GetParty(Guid id)
        {
            return _dbContext.Parties.Find(id);
        }

        public async Task<bool> Update(Party item)
        {
            var oldItem = _dbContext.Parties.Find(item.Id);
            if (oldItem != null)
            {
                _dbContext.Update(item);
                oldItem.Active = item.Active;
                oldItem.Acronym = item.Acronym;
                oldItem.LogoUrl = item.LogoUrl;
                oldItem.Politicians = item.Politicians;
                oldItem.Color = item.Color;
                _dbContext.SaveChanges();
                return true;
            }
            return false;
        }

        public async Task<IList<Party>> GetAllParties(bool isActive = true, int number = 100)
        {
            return _dbContext.Parties.Where(x => x.Active == isActive).Take(number).ToList();
        }

        public Task<Party?> UpdateParty(Guid id, string? name = null, string? acronym = null, string? logoUrl = null, Color? color = null)
        {
            var item = GetParty(id).Result;
            if (item != null)
            {
                _dbContext.Update(item);
                item.Name = name ?? item.Name;
                item.Acronym = acronym ?? item.Acronym;
                item.LogoUrl = logoUrl ?? item.LogoUrl;
                item.Color = color ?? item.Color;
                _dbContext.SaveChanges();
            }
            return Task.FromResult(item);
        }

        public async Task<bool> Delete(Party entity)
        {
            if (_dbContext.Parties.Find(entity.Id) != null)
            {
                _dbContext.Parties.Remove(entity);
                _dbContext.SaveChanges();
                return true;
            }
            return false;
        }

        public async Task<Party?> GetParty(string? name = null, string? acronym = null)
        {
            if(name != null)
            {
                return _dbContext.Parties.FirstOrDefault(x => x.Name == name);
            }
            if (acronym != null)
            {
                return _dbContext.Parties.FirstOrDefault(x => x.Acronym == acronym);
            }
            return null;
        }
    }
}
