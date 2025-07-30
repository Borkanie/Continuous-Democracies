using Microsoft.AspNetCore.Mvc;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.ServiceImplementation;

namespace ContinousDemocracyAPI.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class PartyController : Controller
    {
        private IPartyService<Party> partyService;

        public PartyController(IPartyService<Party> partyService)
        {
            this.partyService = partyService;
        }

        [HttpGet("all")]
        public ActionResult<string> GetAllPartys(
            [FromQuery] bool active,
            [FromQuery] int number = 100)
        {
            var parties = partyService.GetAllPartiesAsync(active, number).Result;
            if(parties.Count == 0)
            {
                return NotFound("No parties found.");
            }
            return Ok(parties);
        }

        [HttpGet("GetById/")]
        public ActionResult<string> GetPartyById(Guid id)
        {
            var party = partyService.GetAsync(id).Result;
            if (party == null)
            {
                return NotFound("Party not found.");
            }
            return Ok(party);
        }


        // GET api/party/query
        // Example: api/party/query?id1=1&id2=2
        [HttpGet("query")]
        public ActionResult<string> GetParty(
            [FromQuery] string? name = null,
            [FromQuery] string? acronym = null)
        {
            var party = partyService.GetPartyAsync(name, acronym).Result;
            if (party == null)
            {
                return NotFound("No party found with the specified criteria.");
            }
            return Ok(party);
        }
    }
}
