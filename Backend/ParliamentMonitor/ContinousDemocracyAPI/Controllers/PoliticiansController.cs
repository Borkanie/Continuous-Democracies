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

        public PoliticiansController(IPoliticianService<Politician> politicianService)
        {
            this.politicianService = politicianService;
        }

        // GET api/politicians/GetById/{id}
        [HttpGet("GetById/{id}")]
        public ActionResult<string> GetPoliticianById(Guid id)
        {
            var politician = politicianService.GetPolitician(id);
            if(politician == null)
            {
                return Ok("");
            }
            return Ok(politician!);
        }


        // GET api/politicians/GetByName/{name}
        [HttpGet("GetByName/{name}")]
        public ActionResult<string> GetPoliticianByName(string name)
        {
            var politician = politicianService.GetPolitician(name);
            if (politician == null)
            {
                return Ok("");
            }
            return Ok(politician!);
        }

        // GET api/politicians/query
        // Example: api/politicians/query?id1=1&id2=2
        [HttpGet("query")]
        public ActionResult<string> GetAllPoliticians(
            [FromQuery] Party? party = null, 
            [FromQuery] bool? isActive = null,
            [FromQuery] WorkLocation? location = null, 
            [FromQuery] Gender? gender = null,
            [FromQuery] int number = 100)
        {

            return Ok(politicianService.GetAllPoliticians(party,isActive,location,gender,number));
        }

    }
}
