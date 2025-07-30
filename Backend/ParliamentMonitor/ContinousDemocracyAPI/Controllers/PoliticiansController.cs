using Microsoft.AspNetCore.Mvc;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;

namespace ContinousDemocracyAPI.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class PoliticiansController : ControllerBase
    {
        private readonly IPoliticianService<Politician> politicianService;
        private readonly IPartyService<Party> partyService;

        public PoliticiansController(IPoliticianService<Politician> politicianService, IPartyService<Party> partyService)
        {
            this.politicianService = politicianService;
            this.partyService = partyService;
        }

        [HttpGet("GetById/")]
        public ActionResult<string> GetPoliticianById([FromQuery] Guid id)
        {
            var politician = politicianService.GetAsync(id).Result;
            if(politician == null)
            {
                return NotFound("Politician not found");
            }
            return Ok(politician!);
        }

        [HttpGet("GetByName/")]
        public ActionResult<string> GetPoliticianByName([FromQuery] string name)
        {
            var politician = politicianService.GetPoliticianAsync(name);
            if (politician == null)
            {
                return NotFound("No Politician with this name found");
            }
            return Ok(politician!);
        }

        // GET api/politicians/query
        // Example: api/politicians/query?id1=1&id2=2
        [HttpGet("getAllPoliticians")]
        public ActionResult<string> GetAllPoliticians(
            [FromQuery] string? partyAcronym = null,
            [FromQuery] string? partyName = null,
            [FromQuery] bool? isActive = null,
            [FromQuery] WorkLocation? location = null, 
            [FromQuery] Gender? gender = null,
            [FromQuery] int number = 100)
        {
            var party = partyService.GetPartyAsync(partyName, partyAcronym).Result;

            var politicians = politicianService.GetAllPoliticiansAsync(party, isActive, location, gender, number).Result;
            if(politicians.Count == 0)
            {
                return NotFound("No politicians found with the specified criteria.");
            }
            return Ok(politicians);
        }

    }
}
