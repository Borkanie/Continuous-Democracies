using Microsoft.AspNetCore.Mvc;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;

namespace ContinousDemocracyAPI.Controllers
{
    [Route("api/[controller]")]
    [ApiController]
    public class PoliticiansController(
        IPoliticianService<Politician> politicianService,
        IPartyService<Party> partyService,
        ILogger<PoliticiansController> logger) : ControllerBase
    {
        private readonly IPoliticianService<Politician> politicianService = politicianService;
        private readonly IPartyService<Party> partyService = partyService;
        private readonly ILogger<PoliticiansController> logger = logger;

        /// <summary>
        /// Retrieves a politician's details by their unique identifier.
        /// </summary>
        /// <remarks>This method uses the <see cref="HttpGetAttribute"/> to handle HTTP GET requests. 
        /// Ensure the provided <paramref name="id"/> corresponds to an existing politician in the system.</remarks>
        /// <param name="id">The unique identifier of the politician to retrieve.</param>
        /// <returns>An HTTP 200 OK response containing the politician's details as a string if found;  otherwise, an HTTP 404
        /// Not Found response with an error message.</returns>
        [HttpGet("GetById/")]
        public ActionResult<string> GetPoliticianById([FromQuery] Guid id)
        {
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var timestamp = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/politicians/GetById?id={Id}",
                timestamp, ip, id);

            var politician = politicianService.GetAsync(id).Result;
            if (politician == null)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (Politician {Id} not found)",
                    timestamp, ip, id);
                return Ok("Politician not found");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK (Politician {Id})",
                timestamp, ip, id);
            return Ok(politician!);
        }

        /// <summary>
        /// Retrieves a politician's details by their name.
        /// </summary>
        /// <remarks>This method queries the underlying data source for a politician matching the
        /// specified name. If no match is found, a 404 Not Found response is returned.</remarks>
        /// <param name="name">The name of the politician to search for. This parameter is case-insensitive.</param>
        /// <returns>An HTTP 200 OK response containing the politician's details as a string if found;  otherwise, an HTTP 404
        /// Not Found response with an error message.</returns>
        [HttpGet("GetByName/")]
        public ActionResult<string> GetPoliticianByName([FromQuery] string name)
        {
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var timestamp = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/politicians/GetByName?name={Name}",
                timestamp, ip, name);

            var politician = politicianService.GetPoliticianAsync(name);
            if (politician == null)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (No Politician with name={Name})",
                    timestamp, ip, name);
                return Ok("No Politician with this name found");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK (Politician {Name})",
                timestamp, ip, name);
            return Ok(politician!);
        }

        /// <summary>
        /// Retrieves a list of politicians based on the specified filter criteria.
        /// </summary>
        /// <remarks>This method allows filtering politicians by various criteria, such as party
        /// affiliation, activity status,  work location, and gender. If no filter criteria are provided, the method
        /// returns up to the specified  number of politicians (default is 100).</remarks>
        /// <param name="partyAcronym">The acronym of the political party to filter by. Optional.</param>
        /// <param name="partyName">The name of the political party to filter by. Optional.</param>
        /// <param name="isActive">Specifies whether to filter by active politicians. If <see langword="true"/>, only active politicians are
        /// returned; if <see langword="false"/>, only inactive politicians are returned. Optional.</param>
        /// <param name="location">The work location to filter politicians by. Optional.</param>
        /// <param name="gender">The gender to filter politicians by. Optional.</param>
        /// <param name="number">The maximum number of politicians to retrieve. Must be a positive integer. Defaults to 100.</param>
        /// <returns>An HTTP 200 OK response containing a JSON-encoded list of politicians matching the specified criteria,  or
        /// an HTTP 404 Not Found response if no politicians match the criteria.</returns>
        [HttpGet("getAllPoliticians")]
        public ActionResult<string> GetAllPoliticians(
            [FromQuery] string? partyAcronym = null,
            [FromQuery] string? partyName = null,
            [FromQuery] bool? isActive = null,
            [FromQuery] WorkLocation? location = null, 
            [FromQuery] Gender? gender = null,
            [FromQuery] int number = 100)
        {
            var ip = HttpContext.Connection.RemoteIpAddress?.ToString();
            var timestamp = DateTime.UtcNow;

            logger.LogInformation("Request at {Timestamp} from {IP} -> GET /api/politicians/getAllPoliticians " +
                                  "(partyAcronym={PartyAcronym}, partyName={PartyName}, isActive={IsActive}, location={Location}, gender={Gender}, number={Number})",
                timestamp, ip, partyAcronym, partyName, isActive, location, gender, number);

            var party = partyService.GetPartyAsync(partyName, partyAcronym).Result;
            var politicians = politicianService.GetAllPoliticiansAsync(party, isActive, location, gender, number).Result;

            if (politicians.Count == 0)
            {
                logger.LogWarning("Response at {Timestamp} to {IP} -> 404 Not Found (No politicians matching criteria)",
                    timestamp, ip);
                return Ok("No politicians found with the specified criteria.");
            }

            logger.LogInformation("Response at {Timestamp} to {IP} -> 200 OK ({Count} politicians)",
                timestamp, ip, politicians.Count);
            return Ok(politicians);
        }

    }
}
