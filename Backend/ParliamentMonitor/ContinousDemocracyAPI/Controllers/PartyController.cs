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

        // GET api/party/all
        [HttpGet("all")]
        public ActionResult<string> GetAllPartys(
            [FromQuery] bool active,
            [FromQuery] int number = 100)
        {

            return Ok(partyService.GetAllParties(active, number));
        }

        // GET api/voting/GetById/{id}
        [HttpGet("GetById/{id}")]
        public ActionResult<string> GetPartyById(Guid id)
        {
            return Ok(partyService.GetParty(id));
        }


        // GET api/party/query
        // Example: api/party/query?id1=1&id2=2
        [HttpGet("query")]
        public ActionResult<string> GetParty(
            [FromQuery] string? name = null,
            [FromQuery] string? acronym = null)
        {

            return Ok(partyService.GetParty(name,acronym));
        }
    }
}
