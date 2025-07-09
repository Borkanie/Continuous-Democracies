using Microsoft.AspNetCore.Mvc;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.ServiceImplementation;

namespace ContinousDemocracyAPI.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class VotingController : Controller
    {
        private IVotingService<Vote, Round> votingService;

        public VotingController(IVotingService<Vote, Round> votingService)
        {
            this.votingService = votingService;
        }

        // GET api/voting/all
        [HttpGet("all")]
        public ActionResult<string> GetAllPoliticians([FromQuery] int number = 100)
        {

            return Ok(votingService.GetAllRoundsFromDB());
        }

        // GET api/voting/GetById/{id}
        [HttpGet("GetById/{id}")]
        public ActionResult<string> GetPoliticianById(int id)
        {
            return Ok(votingService.GetVotingRound(id));
        }
    }
}
