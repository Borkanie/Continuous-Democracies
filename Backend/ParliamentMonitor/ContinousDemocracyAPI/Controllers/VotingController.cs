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
        private IVotingService<Vote> votingService;

        private IVotingRoundService<Round> votingRoundService;

        public VotingController(IVotingService<Vote> votingService, IVotingRoundService<Round> votingRoundService)
        {
            this.votingService = votingService;
            this.votingRoundService = votingRoundService;
        }

        // GET api/voting/all
        [HttpGet("getAllRound")]
        public ActionResult<string> GetAllRounds([FromQuery] int number = 100)
        {

            return Ok(votingRoundService.GetAllRoundsFromDBAsync());
        }

        // GET api/voting/GetById/{id}
        [HttpGet("getRoundById/{id}")]
        public ActionResult<string> GetVotingRoundByVoteId(int id)
        {
            var round = votingRoundService.GetVotingRoundAsync(id).Result;
            if(round == null)
            {
                return NotFound("Voting round not found.");
            }   
            return Ok(round);
        }

        // GET api/voting/GetById/{id}
        [HttpGet("GetAllVotesForARoundById/{id}")]
        public ActionResult<string> GetVotesForRound(Guid roundId)
        {
            var round = votingRoundService.GetAsync(roundId).Result;
            if(round == null)
            {
                return NotFound("Voting round not found.");
            }
            return Ok(round.VoteResults);
        }
    }
}
