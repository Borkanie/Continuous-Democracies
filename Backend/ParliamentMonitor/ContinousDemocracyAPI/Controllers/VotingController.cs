using Microsoft.AspNetCore.Mvc;
using ParliamentMonitor.Contracts.Model;
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

        /// <summary>
        /// Retrieves a list of voting rounds from the database based on the specified date range and maximum number of
        /// entries.
        /// </summary>
        /// <remarks>This method queries the database for voting rounds that match the specified criteria.
        /// If no criteria are provided, all available rounds up to the specified maximum are returned.</remarks>
        /// <param name="startDate">The optional start date to filter the voting rounds. Only rounds occurring on or after this date will be
        /// included. If null, no lower date limit is applied.</param>
        /// <param name="endDate">The optional end date to filter the voting rounds. Only rounds occurring on or before this date will be
        /// included. If null, no upper date limit is applied.</param>
        /// <param name="maxNumberOfEntries">The maximum number of voting rounds to retrieve. Must be a positive integer. Defaults to 100.</param>
        /// <returns>An HTTP 200 OK response containing the list of voting rounds if any are found;  otherwise, an HTTP 404 Not
        /// Found response with a message indicating no results were available.</returns>
        [HttpGet("getAllRounds")]
        public ActionResult<string> GetAllRounds(
            [FromQuery] DateTime? startDate,
            [FromQuery] DateTime? endDate,
            [FromQuery] int maxNumberOfEntries = 100)
        {
            
            var result = votingRoundService.GetAllRoundsFromDBAsync(startDate, endDate, maxNumberOfEntries).Result;
            if (result.Count == 0)
            {
                return NotFound("No results where available in the db");
            }
            return Ok(result);
        }

        /// <summary>
        /// Retrieves the voting round associated with the specified vote number.
        /// </summary>
        /// <remarks>This method queries the voting round service to retrieve the voting round associated
        /// with the given vote ID. If no voting round is found for the specified ID, a 404 Not Found response is
        /// returned.</remarks>
        /// <param name="voteNumber">The number of the vote for which the voting round is requested.</param>
        /// <returns>An HTTP 200 OK response containing the voting round as a string if found;  otherwise, an HTTP 404 Not Found
        /// response with an error message.</returns>
        [HttpGet("getRoundById/")]
        public ActionResult<string> GetVotingRoundByVoteId([FromQuery] int voteNumber)
        {
            var round = votingRoundService.GetVotingRoundAsync(voteNumber).Result;
            if(round == null)
            {
                return NotFound("Voting votes not found.");
            }   
            return Ok(round);
        }

        /// <summary>
        /// Gets all the votes for a given <see cref="Round"/>.
        /// </summary>
        /// <param name="number">MANDATORY - The integer that identifies the <see cref="Round"/></param>
        /// <param name="partyId">The <see cref="Party"/> id to filter the votes by</param>
        /// <param name="partyAcronim">The acronym of a given <see cref="Party"/> to filter the votes by</param>
        /// <returns>List of <see cref="Vote"/> for a law in a <see cref="Round"/>.</returns>
        [HttpGet("GetResultForVote/")]
        public ActionResult<string> GetVotesForRound(
            [FromQuery] int number,
            [FromQuery] Guid? partyId,
            [FromQuery] string? partyAcronim)
        {
            var votes = votingService.GetAllVotesForRound(number).Result;
            if (votes == null)
            {
                return NotFound("Voting votes not found.");
            }
            if (partyId != null)
            {
                votes = votes.Where(v => v.Politician.Party != null && v.Politician.Party.Id == partyId).ToList();
                return Ok(votes);
            }
            if (partyAcronim != null)
            {
                votes = votes.Where(v => v.Politician.Party != null && string.Equals(v.Politician.Party.Acronym, partyAcronim)).ToList();
                return Ok(votes);
            }
            return Ok(votes);

        }

        /// <summary>
        /// Retrieves all votes associated with a specific round by its unique identifier.
        /// </summary>
        /// <remarks>This method queries the voting service to retrieve all votes for the specified round.
        /// Ensure that the <paramref name="roundId"/> is a valid and existing round identifier.</remarks>
        /// <param name="roundId">The unique identifier of the round for which votes are to be retrieved.</param>
        /// <returns>An HTTP 200 OK response containing the votes as a JSON string if found;  otherwise, an HTTP 404 Not Found
        /// response with an appropriate error message.</returns>
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
