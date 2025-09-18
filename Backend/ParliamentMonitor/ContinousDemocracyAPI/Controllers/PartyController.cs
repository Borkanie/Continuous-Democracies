using Microsoft.AspNetCore.Mvc;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.ServiceImplementation;

namespace ContinousDemocracyAPI.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class PartyController(IPartyService<Party> partyService, ILogger<PartyController> logger) : Controller
    {
        private readonly IPartyService<Party> partyService = partyService;
        private readonly ILogger<PartyController> logger = logger;

        /// <summary>
        /// Retrieves a list of parties based on their active status and a specified limit.
        /// </summary>
        /// <remarks>This method queries the party service to retrieve the requested parties. If no
        /// parties match the criteria, a 404 Not Found response is returned with an appropriate message.</remarks>
        /// <param name="active">A value indicating whether to retrieve only active parties. If <see langword="true"/>, only active parties
        /// are returned; otherwise, inactive parties are included.</param>
        /// <param name="number">The maximum number of parties to retrieve. The default value is 100.</param>
        /// <returns>An <see cref="ActionResult{T}"/> containing the list of parties if found; otherwise, a 404 Not Found
        /// response.</returns>
        [HttpGet("all")]
        public ActionResult<string> GetAllPartys(
            [FromQuery] bool active,
            [FromQuery] int number = 100)
        {
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var timestamp = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/party/all?active={Active}&number={Number}",
                timestamp, ip, active, number);

            var parties = partyService.GetAllPartiesAsync(active, number).Result;
            if (parties.Count == 0)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (No parties)", timestamp, ip);
                return Ok("No parties found.");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK ({Count} parties)", timestamp, ip, parties.Count);
            return Ok(parties);
        }

        /// <summary>
        /// Retrieves a party by its unique identifier.
        /// </summary>
        /// <param name="id">The unique identifier of the party to retrieve.</param>
        /// <returns>An HTTP 200 response containing the party details as a string if the party is found;  otherwise, an HTTP 404
        /// response with a "Party not found" message.</returns>
        [HttpGet("GetById/")]
        public ActionResult<string> GetPartyById(Guid id)
        {
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var timestamp = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/party/GetById?id={Id}", timestamp, ip, id);

            var party = partyService.GetAsync(id).Result;
            if (party == null)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (Party {Id} not found)", timestamp, ip, id);
                return Ok("Party not found.");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK (Party {Id})", timestamp, ip, id);
            return Ok(party);
        }


        /// <summary>
        /// Retrieves party information based on the specified query parameters.
        /// </summary>
        /// <remarks>At least one of the query parameters, <paramref name="name"/> or <paramref
        /// name="acronym"/>, should be provided to perform a meaningful search. If both parameters are null or empty,
        /// the method may return a 404 response.</remarks>
        /// <param name="name">The name of the party to search for. This parameter is optional.</param>
        /// <param name="acronym">The acronym of the party to search for. This parameter is optional.</param>
        /// <returns>An HTTP 200 response containing the party information as a string if a matching party is found; otherwise,
        /// an HTTP 404 response with a message indicating that no party matches the specified criteria.</returns>
        [HttpGet("query")]
        public ActionResult<string> GetParty(
            [FromQuery] string? name = null,
            [FromQuery] string? acronym = null)
        {
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var timestamp = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/party/query?name={Name}&acronym={Acronym}",
                timestamp, ip, name, acronym);

            var party = partyService.GetPartyAsync(name, acronym).Result;
            if (party == null)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (name={Name}, acronym={Acronym})",
                    timestamp, ip, name, acronym);
                return Ok("No party found with the specified criteria.");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK (Party {Party})",
                timestamp, ip, party);
            return Ok(party);
        }
    }
}
