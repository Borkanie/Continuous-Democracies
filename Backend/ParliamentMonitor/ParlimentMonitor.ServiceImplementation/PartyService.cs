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
        private readonly AppDBContext dbContext;

        public PartyService(AppDBContext context)
        {
            dbContext = context;
            // Cache parties on startup
            GetAllParties();
        }

        /*<inheritdoc/>*/
        public Party? CreateParty(string name, string? acronym = null, string? logoUrl = null, Color? color = null)
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
                if (newParty == null || dbContext.Parties.Any(x => x == newParty))
                {
                    throw new Exception("Already existing party");
                }
                dbContext.Parties.Add(newParty);
                dbContext.SaveChanges();
                return newParty;
            }catch(Exception ex)
            {
                Console.WriteLine($"Error when creating new party:{ex.Message}");
                return null; 
            }
            
        }


        /*<inheritdoc/>*/
        public Party? GetParty(Guid id)
        {
            return dbContext.Parties.Find(id);
        }

        public void Update(Party item)
        {
            var oldItem = dbContext.Parties.Find(item.Id);
            if (oldItem != null)
            {
                dbContext.Update(item);
                oldItem.Active = item.Active;
                oldItem.Acronym = item.Acronym;
                oldItem.LogoUrl = item.LogoUrl;
                oldItem.Politicians = item.Politicians;
                oldItem.Color = item.Color;
                dbContext.SaveChanges();
            }  
        }

        public IList<Party> GetAllParties(bool isActive = true, int number = 100)
        {
            return dbContext.Parties.Where(x => x.Active == isActive).Take(number).ToList();
        }

        public Party? UpdateParty(Guid id, string? name = null, string? acronym = null, string? logoUrl = null, Color? color = null)
        {
            var item = GetParty(id);
            if (item != null)
            {
                dbContext.Update(item);
                item.Name = name ?? item.Name;
                item.Acronym = acronym ?? item.Acronym;
                item.LogoUrl = logoUrl ?? item.LogoUrl;
                item.Color = color ?? item.Color;
                dbContext.SaveChanges();
            }
            return item;
        }

        public void Delete(Party entity)
        {
            if (dbContext.Parties.Find(entity.Id) != null)
            {
                dbContext.Parties.Remove(entity);
                dbContext.SaveChanges();
            }
        }

        public Party? GetParty(string? name = null, string? acronym = null)
        {
            if(name != null)
            {
                return dbContext.Parties.FirstOrDefault(x => x.Name == name);
            }
            if (acronym != null)
            {
                return dbContext.Parties.FirstOrDefault(x => x.Acronym == acronym);
            }
            return null;
        }
    }
}
