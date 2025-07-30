using Microsoft.AspNetCore.Mvc;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;

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

        [HttpGet("getAllRounds")]
        public ActionResult<string> GetAllRounds()
        {
            
            var result = votingRoundService.GetAllRoundsFromDBAsync().Result;
            if (result.Count == 0)
            {
                return NotFound("No results where available in the db");
            }
            return Ok(result);
        }

        [HttpGet("getTheLastRounds/")]
        public ActionResult<string> GetAllRounds([FromQuery] int number = 10)
        {
            var result = votingRoundService.GetAllRoundsFromDBAsync(number).Result;
            if(result.Count == 0)
            {
                return NotFound("No results where available in the db");
            }
            return Ok(result);
        }

        [HttpGet("getRoundById/")]
        public ActionResult<string> GetVotingRoundByVoteId([FromQuery] int voteId)
        {
            var round = votingRoundService.GetVotingRoundAsync(voteId).Result;
            if(round == null)
            {
                return NotFound("Voting votes not found.");
            }   
            return Ok(round);
        }

        [HttpGet("GetResultForVote/")]
        public ActionResult<string> GetVotesForRound([FromQuery] int number)
        {
            var votes = votingService.GetAllVotesForRound(number).Result;
            if (votes == null)
            {
                return NotFound("Voting votes not found.");
            }
            return Ok(votes);
        }

        [HttpGet("GetAllVotesForARoundById/")]
        public ActionResult<string> GetAllVotesForARoundById([FromQuery] Guid roundId)
        {
            var votes = votingService.GetAllVotesForRound(roundId).Result;
            if(votes == null)
            {
                return NotFound("Voting votes not found.");
            }
            return Ok(votes);
        }
    }
}
