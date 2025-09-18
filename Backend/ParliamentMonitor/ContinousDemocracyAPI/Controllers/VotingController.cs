using Microsoft.AspNetCore.Mvc;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;

namespace ContinousDemocracyAPI.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class VotingController(
        IVotingService<Vote> votingService,
        IVotingRoundService<Round> votingRoundService,
        ILogger<VotingController> logger) : Controller
    {
        private readonly IVotingService<Vote> votingService = votingService;
        private readonly IVotingRoundService<Round> votingRoundService = votingRoundService;
        private readonly ILogger<VotingController> logger = logger;


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

            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var ts = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/voting/getAllRounds (startDate={StartDate}, endDate={EndDate}, max={Max})",
                ts, ip, startDate, endDate, maxNumberOfEntries);

            var result = votingRoundService.GetAllRoundsFromDBAsync(startDate, endDate, maxNumberOfEntries).Result;

            if (result.Count == 0)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (No rounds)", ts, ip);
                return Ok("No results where available in the db");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK ({Count} rounds)", ts, ip, result.Count);
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
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var ts = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/voting/getRoundById?voteNumber={VoteNumber}", ts, ip, voteNumber);

            var round = votingRoundService.GetVotingRoundAsync(voteNumber).Result;
            if (round == null)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (VoteNumber={VoteNumber})", ts, ip, voteNumber);
                return Ok("Voting votes not found.");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK (Round found for VoteNumber={VoteNumber})", ts, ip, voteNumber);
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
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var ts = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/voting/GetResultForVote (number={Number}, partyId={PartyId}, partyAcronim={PartyAcronim})",
                ts, ip, number, partyId, partyAcronim);

            var votes = votingService.GetAllVotesForRound(number).Result;
            if (votes == null)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (Round {Number})", ts, ip, number);
                return Ok("Voting votes not found.");
            }

            if (partyId != null)
            {
                votes = votes.Where(v => v.Politician.Party != null && v.Politician.Party.Id == partyId).ToList();
                logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK ({Count} votes filtered by PartyId={PartyId})", ts, ip, votes.Count, partyId);
                return Ok(votes);
            }

            if (partyAcronim != null)
            {
                votes = votes.Where(v => v.Politician.Party != null &&
                                         string.Equals(v.Politician.Party.Acronym, partyAcronim, StringComparison.OrdinalIgnoreCase))
                             .ToList();
                logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK ({Count} votes filtered by Acronym={PartyAcronim})", ts, ip, votes.Count, partyAcronim);
                return Ok(votes);
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK ({Count} total votes)", ts, ip, votes.Count);
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
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var ts = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/voting/GetAllVotesForARoundById?roundId={RoundId}", ts, ip, roundId);

            var votes = votingService.GetAllVotesForRound(roundId).Result;
            if (votes == null)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (RoundId={RoundId})", ts, ip, roundId);
                return Ok("Voting votes not found.");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK ({Count} votes for RoundId={RoundId})", ts, ip, votes.Count, roundId);
            return Ok(votes);
        }
    }
}
