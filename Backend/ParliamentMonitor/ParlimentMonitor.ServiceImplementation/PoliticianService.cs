using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using StackExchange.Redis;
using System.IO;
using System.Text.Json;

namespace ParliamentMonitor.ServiceImplementation
{
    public class PoliticianService : RedisService, IPoliticianService<Politician>
    {
        private AppDBContext _dbContext;       

        public PoliticianService(AppDBContext context, IConnectionMultiplexer redis)  : base(redis,"politician")
        {
            _dbContext = context;
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
                SetAsync(MakeKey(politician.Id.ToString()), politician);
                return politician;
            }
            catch(Exception ex)
            {
                Console.WriteLine($"Error when creating Politician:{ex.Message}");
                return null; 
            }
            
        }

        private async Task<IList<Politician>> getPoliticians(Party? party = null, bool? isActive = null, WorkLocation? location = null, Gender? gender = null, int number = 100)
        {
            var politicians = new List<Politician>();
            var keys = await _cache.SetMembersAsync(serviceKey);
            foreach(var key in keys)
            {
                var politician = await GetAsync<Politician>(key);
                if(politician != null)
                {
                    if(party!= null && politician.Party!=party
                        || isActive != null && politician.Active != isActive
                        || location!= null && politician.WorkLocation != location
                        || gender!= null && politician.Gender != gender
                        )
                    {
                        continue;
                    }
                    politicians.Add(politician);
                    if (politicians.Count == number)
                        return politicians;
                }
            }
           

            return new List<Politician>();

        }

        public async Task<IList<Politician>> GetAllPoliticians(Party? party = null, bool? isActive = null, WorkLocation? location = null, Gender? gender = null, int number = 100)
        {
            var cachedPoliticians = await getPoliticians(party,isActive,location,gender, number);
            if (cachedPoliticians.Count == number)
                return cachedPoliticians;
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

            var politicians = query.Take(number).ToList();
            foreach(var politician in politicians)
            {
                SetAsync(MakeKey(politician.Id.ToString()), politician);
            }
            return politicians; // single SQL query executed here

        }

        public async Task<Politician?> GetPolitician(Guid id)
        {
            var cached = await GetAsync<Politician>(MakeKey(id.ToString()));
            if (cached != null)
                return cached;
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
                SetAsync(MakeKey(entity.Id.ToString()), entity);
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
                SetAsync(MakeKey(item.Id.ToString()), item);
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
                SetAsync(MakeKey(item.Id.ToString()), item);
            }
            return item;
        }

        public async Task<bool> Delete(Politician entity)
        {
            if (_dbContext.Politicians.Find(entity) != null)
            {
                _dbContext.Politicians.Remove(entity);
                _dbContext.SaveChanges();
                RemoveAsync(MakeKey(entity.Id.ToString()));
                return true;
            }
            return false;
        }

        public async Task<Politician?> GetPolitician(string name)
        {
            var politician = getPoliticians().Result.FirstOrDefault(x => string.Equals(x.Name, name));
            if (politician != null)
                return politician;
            var value = _dbContext.Politicians.FirstOrDefault(x => String.Equals(x.Name.ToLower(), name.ToLower()));
            SetAsync(MakeKey(value.Id.ToString()), value);
            return value;
        }

    }
}
