using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using System.Drawing;

namespace ParlimentMonitor.ServiceImplementation
{
    public class PartyService : IPartyService<Party>
    {
        private readonly AppDBContext dbContext;

        public PartyService(AppDBContext context)
        {
            dbContext = context;
            // Cache parties on startup
            GetAllParties();
        }

        /*<inheritdoc/>*/
        public Party CreateParty(string name, string? acronym = null, string? logoUrl = null, Color? color = null)
        {

                var newParty = new Party();
                newParty.Id = Guid.NewGuid();
                newParty.Name = name;
                if(acronym != null) 
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
            
        }


        /*<inheritdoc/>*/
        public Party? GetParty(Guid id)
        {
            return dbContext.Parties.First(x => x.Id == id);
        }

        public void Update(Party item)
        {
            if(dbContext.Parties.Contains(item))
            {
                dbContext.Update(item);
                dbContext.SaveChanges();
            }  
        }

        public IList<Party> GetAllParties(bool? isActive = null)
        {
            return isActive==null ? dbContext.Parties.Where(x => true).ToList() : 
                dbContext.Parties.Where(x => x.Active == isActive).ToList();
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
            if (dbContext.Parties.Contains(entity))
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
