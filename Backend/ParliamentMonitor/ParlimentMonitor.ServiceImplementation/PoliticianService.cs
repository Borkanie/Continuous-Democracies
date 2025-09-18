using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;

namespace ParliamentMonitor.ServiceImplementation
{
    public class PoliticianService(AppDBContext context, ILogger<IPoliticianService<Politician>> logger) : IPoliticianService<Politician>
    {
        private AppDBContext _dbContext = context;
        private ILogger<IPoliticianService<Politician>> _logger = logger;

        public ILogger Logger { get => _logger; }

        /*<inheritdoc/>*/
        public Task<Politician?> CreatePoliticanAsync(string name, Party party, WorkLocation location, Gender gender, bool isCurrentlyActive = true, string? imageUrl = null)
        {
            Politician politician = new Politician() { Id = Guid.NewGuid(), Party = party };
            politician.Name = name;
            politician.WorkLocation = location;
            politician.Gender = gender;
            politician.Active = isCurrentlyActive;
            if (imageUrl != null)
            {
                politician.ImageUrl = imageUrl;
            }
            if (_dbContext.Politicians.Any(x => x == politician))
                throw new Exception($"Already existing politician:{politician.Name}");
            _dbContext.Politicians.Add(politician);
            _dbContext.SaveChanges();
            _logger.LogInformation($"Created new politician:{politician.Name}");
            return Task.FromResult<Politician?>(politician);

        }


        /*<inheritdoc/>*/
        public async Task<IList<Politician>> GetAllPoliticiansAsync(Party? party = null, bool? isActive = null, WorkLocation? location = null, Gender? gender = null, int number = 100)
        {
            var query = _dbContext.Politicians.Include(x => x.Party).AsQueryable();

            if (party != null)
            {
                _logger.LogInformation($"Fitlering by Part: {party.Name}");
                query = query.Where(x => x.Party == party);
            }

            if (isActive != null)
            {
                _logger.LogInformation($"Fitlering by Activity: {isActive}");
                query = query.Where(x => x.Active == isActive);
            }

            if (location != null)
            {
                _logger.LogInformation($"Fitlering by Location: {location}");
                query = query.Where(x => x.WorkLocation == location);
            }

            if (gender != null)
            {

                _logger.LogInformation($"Fitlering by Gender: {gender}");
                query = query.Where(x => x.Gender == gender);
            }

            var politicians = query.Take(number).ToList();

            return politicians;

        }

        /*<inheritdoc/>*/
        public async Task<Politician?> GetAsync(Guid id)
        {
            var politician = _dbContext.Politicians.Find(id);
            if (politician == null)
            {
                throw new Exception($"Politician with Id: {id} does not exist in DB");
            }
            return politician;
        }

        /*<inheritdoc/>*/
        public Task<bool> UpdateAsync(Politician entity)
        {
            var old = GetAsync(entity.Id).Result;
            if (old!=null)
            {
                _dbContext.Update(old);
                old.Party = entity.Party;
                old.Gender = entity.Gender;
                old.WorkLocation = entity.WorkLocation;
                old.Name = entity.Name;
                old.Active = entity.Active;
                _dbContext.SaveChanges();
                _logger.LogInformation($"Updated politician:{entity.Id}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /*<inheritdoc/>*/
        public Task<Politician?> UpdatePoliticianActivityAsync(Guid id, bool isCurrentlyActive)
        {
            var item = GetAsync(id).Result;
            if(item != null)
            {
                _dbContext.Update(item);
                item.Active = isCurrentlyActive;
                _dbContext.SaveChanges();
                _logger.LogInformation($"Updated politician activity:{item.Id}");
            }
            return Task.FromResult(item);
        }

        /*<inheritdoc/>*/
        public Task<Politician?> UpdatePoliticianAsync(Guid id, string? name = null, Party? party = null, WorkLocation? location = null, Gender? gender = null, bool? isCurrentlyActive = null, string? imageUrl = null)
        {
            var item = GetAsync(id).Result;
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
                _logger.LogInformation($"Updated politician:{item.Id}");
            }
            return Task.FromResult(item);
        }

        /*<inheritdoc/>*/
        public Task<bool> DeleteAsync(Politician entity)
        {
            if (_dbContext.Politicians.Find(entity) != null)
            {
                _dbContext.Politicians.Remove(entity);
                _dbContext.SaveChanges();
                _logger.LogInformation($"Deleted politician:{entity.Id}");
                return Task.FromResult(true);
            }
            return Task.FromResult(false);
        }

        /*<inheritdoc/>*/
        public Task<Politician?> GetPoliticianAsync(string name)
        {
            var politician = _dbContext.Politicians.Include(x => x.Party).FirstOrDefault(x => String.Equals(x.Name.ToLower(), name.ToLower()));
            if(politician == null)
            {
                return Task.FromResult<Politician?>(null);
            }
            return Task.FromResult<Politician?>(politician);
        }

    }
}
