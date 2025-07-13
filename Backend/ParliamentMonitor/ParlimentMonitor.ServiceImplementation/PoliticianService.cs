using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using StackExchange.Redis;

namespace ParliamentMonitor.ServiceImplementation
{
    public class PoliticianService : IPoliticianService<Politician>
    {
        internal readonly IConnectionMultiplexer _redis;
        internal IDatabase _cache;
        private AppDBContext _dbContext;

        public PoliticianService(AppDBContext context, IConnectionMultiplexer redis)
        {
            _dbContext = context;
            _redis = redis;
            _cache = redis.GetDatabase();
        }

        /*<inheritdoc/>*/
        public async Task<Politician?> CreatePolitican(string name, Party party, WorkLocation location, Gender gender, bool isCurrentlyActive = true, string? imageUrl = null)
        {
            try
            {
                Politician politician = new Politician();
                politician.Name = name;
                politician.Party = party;
                politician.WorkLocation = location;
                politician.Gender = gender;
                politician.Active = isCurrentlyActive;
                if (imageUrl != null)
                {
                    politician.ImageUrl = imageUrl;
                }
                if (_dbContext.Politicians.Any(x => x == politician))
                    throw new Exception("Already existing politician");
                _dbContext.Politicians.Add(politician);
                _dbContext.SaveChanges();
                return politician;
            }
            catch(Exception ex)
            {
                Console.WriteLine($"Error when creating Politician:{ex.Message}");
                return null; 
            }
            
        }

        public async Task<IList<Politician>> GetAllPoliticians(Party? party = null, bool? isActive = null, WorkLocation? location = null, Gender? gender = null, int number = 100)
        {
            var query = _dbContext.Politicians.AsQueryable();

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

            return query.Take(number).ToList(); // single SQL query executed here

        }

        public async Task<Politician?> GetPolitician(Guid id)
        {
            return _dbContext.Politicians.Find(id);
        }

        public async Task<bool> Update(Politician entity)
        {
            var old = GetPolitician(entity.Id).Result;
            if (old!=null)
            {
                
                _dbContext.Update(old);
                old.Party = entity.Party;
                old.Gender = entity.Gender;
                old.WorkLocation = entity.WorkLocation;
                old.Name = entity.Name;
                old.Active = entity.Active;
                _dbContext.SaveChanges();
                return true;
            }
            return false;
        }

        /*<inheritdoc/>*/
        public async Task<Politician?> UpdatePoliticianActivity(Guid id, bool isCurrentlyActive)
        {
            var item = GetPolitician(id).Result;
            if(item != null)
            {
                _dbContext.Update(item);
                item.Active = isCurrentlyActive;
                _dbContext.SaveChanges();
            }
            return item;
        }

        public async Task<Politician?> UpdatePolitician(Guid id, string? name = null, Party? party = null, WorkLocation? location = null, Gender? gender = null, bool? isCurrentlyActive = null, string? imageUrl = null)
        {
            var item = GetPolitician(id).Result;
            if (item != null)
            {
                _dbContext.Update(item);
                item.Name = name ?? item.Name;
                item.Party = party ?? item.Party;
                item.WorkLocation = location ?? item.WorkLocation;
                item.Gender = gender ?? item.Gender;
                item.Active = isCurrentlyActive ?? item.Active;
                item.ImageUrl = imageUrl ?? item.ImageUrl;
                _dbContext.SaveChanges();
            }
            return item;
        }

        public async Task<bool> Delete(Politician entity)
        {
            if (_dbContext.Politicians.Find(entity) != null)
            {
                _dbContext.Politicians.Remove(entity);
                _dbContext.SaveChanges();
                return true;
            }
            return false;
        }

        public async Task<Politician?> GetPolitician(string name)
        {
            return _dbContext.Politicians.FirstOrDefault(x => String.Equals(x.Name.ToLower(), name.ToLower()));
        }
    }
}
