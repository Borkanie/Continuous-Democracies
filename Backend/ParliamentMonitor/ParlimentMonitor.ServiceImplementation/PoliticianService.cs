using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using System.Drawing;
using System.Reflection;

namespace ParliamentMonitor.ServiceImplementation
{
    public class PoliticianService : IPoliticianService<Politician>
    {
        private AppDBContext dBContext;

        public PoliticianService(AppDBContext context)
        {
            dBContext = context;
            // Cache entries form DB.
            GetAllPoliticians();
        }

        /*<inheritdoc/>*/
        public Politician CreatePolitican(string name, Party party, WorkLocation location, Gender gender, bool isCurrentlyActive = true, string? imageUrl = null)
        {
            Politician politician = new Politician();
            politician.Name = name;
            politician.Party = party;
            politician.WorkLocation = location;
            politician.Gender = gender;
            politician.Active = isCurrentlyActive;
            if(imageUrl != null)
            {
                politician.ImageUrl = imageUrl;
            }
            if (dBContext.Politicians.Any(x => x == politician))
                throw new Exception("Already existing politician");
            dBContext.Politicians.Add(politician);
            dBContext.SaveChanges();
            return politician;
        }

        public IList<Politician> GetAllPoliticians(Party? party = null, bool? isActive = null, WorkLocation? location = null, Gender? gender = null)
        {
            var query = dBContext.Politicians.AsQueryable();

            if (party != null)
            {
                query = query.Where(x => x.Party == party);
            }

            if (isActive != null)
            {
                query = query.Where(x => x.Active == isActive);
            }

            if (location != null)
            {
                query = query.Where(x => x.WorkLocation == location);
            }

            if (gender != null)
            {
                query = query.Where(x => x.Gender == gender);
            }

            return query.ToList(); // single SQL query executed here

        }

        public Politician? GetPolitician(Guid id)
        {
            return dBContext.Politicians.First(x => x.Id == id);
        }

        public void Update(Politician entity)
        {
            if (dBContext.Politicians.Contains(entity))
            {
                dBContext.Update(entity);
                dBContext.SaveChanges();
            }
        }

        /*<inheritdoc/>*/
        public Politician? UpdatePoliticianActivity(Guid id, bool isCurrentlyActive)
        {
            var item = GetPolitician(id);
            if(item != null)
            {
                dBContext.Update(item);
                item.Active = isCurrentlyActive;
                dBContext.SaveChanges();
            }
            return item;
        }

        public Politician? UpdatePolitician(Guid id, string? name = null, Party? party = null, WorkLocation? location = null, Gender? gender = null, bool? isCurrentlyActive = null, string? imageUrl = null)
        {
            var item = GetPolitician(id);
            if (item != null)
            {
                dBContext.Update(item);
                item.Name = name ?? item.Name;
                item.Party = party ?? item.Party;
                item.WorkLocation = location ?? item.WorkLocation;
                item.Gender = gender ?? item.Gender;
                item.Active = isCurrentlyActive ?? item.Active;
                item.ImageUrl = imageUrl ?? item.ImageUrl;
                dBContext.SaveChanges();
            }
            return item;
        }

        public void Delete(Politician entity)
        {
            if (dBContext.Politicians.Contains(entity))
            {
                dBContext.Politicians.Remove(entity);
                dBContext.SaveChanges();
            }
        }

        public Politician? GetPolitician(string name)
        {
            return dBContext.Politicians.FirstOrDefault(x => x.Name == name);
        }
    }
}
