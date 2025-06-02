using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using System;
using System.Collections.Generic;
using System.Drawing;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Xml.Linq;

namespace ParlimentMonitor.ServiceImplementation
{
    internal class PartyService : IPartyService<Party>
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
                    return dbContext.Parties.First(x => x == newParty);
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
                item.Name = name ?? item.Name;
                item.Acronym = acronym ?? item.Acronym;
                item.LogoUrl = logoUrl ?? item.LogoUrl;
                item.Color = color ?? item.Color;
                Update(item);
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
    }
}
