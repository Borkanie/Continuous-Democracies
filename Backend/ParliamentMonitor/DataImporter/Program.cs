// See https://aka.ms/new-console-template for more information
using DataImporter;
using Microsoft.EntityFrameworkCore;
using Microsoft.Extensions.Logging;
using ParliamentMonitor.DataBaseConnector;
using ParliamentMonitor.ServiceImplementation;
using System.Text;
using Microsoft.Extensions.Logging.Console;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;

Encoding.RegisterProvider(CodePagesEncodingProvider.Instance);

Console.WriteLine("Choose intended behaviour: for image mapper type I for anything else press any other key");
var key = Console.ReadKey();

if (key.Key == ConsoleKey.I)
{
    UpdateImagesBasedOnImagefolder();
}
else
{
    await ImportDataFromVoteFolder();
}

static void UpdateImagesBasedOnImagefolder()
{

    var imagefolderPath = "D:\\Continuous-Democracies\\Resources\\DeputyPortraits\\";
    Console.WriteLine($"Currently importing: {imagefolderPath}");
    using (ILoggerFactory factory = LoggerFactory.Create(builder => builder.AddConsole()))
    {
        Console.WriteLine("Starting services");
        var dbContext = new AppDBContext();
        Console.WriteLine("Loaded DB Context");

        var polLogger = factory.CreateLogger<IPoliticianService<Politician>>();
        var politicianService = new PoliticianService(dbContext, polLogger);
        Console.WriteLine("Started Politican service");
        var directory = Directory.GetFiles(imagefolderPath);
        foreach(var politician in politicianService.GetAllPoliticiansAsync(number: 1000).Result)
        {
            var fileUrl = imagefolderPath + UrlBuilder.CleanName(politician.Name);
            if (File.Exists(fileUrl))
            {
                Console.WriteLine($"{fileUrl} exsits and will update: {politician.Name}");
                if(politician.ImageUrl != null && String.Equals(politician.ImageUrl, fileUrl))
                {
                    Console.WriteLine("Image url already up to date skipping!");
                    continue;
                }
                politicianService.UpdatePoliticianAsync(politician.Id, imageUrl: fileUrl);
            }
        }
    }
        

}

static void ImportFile(string pathToXmlFile, VotingDataimporter dataImporter)
{
    if (File.Exists(pathToXmlFile))
    {
        try
        {
            dataImporter.ImportVotingRoundFromXml(pathToXmlFile);
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Exception when writing into db {ex.Message}");
        }

    }
    else
    {
        Console.WriteLine("Fisierul nu exista");
    }
}

static async Task ImportDataFromVoteFolder()
{
    var pathToXmlDir = "D:\\Voturi\\Votes";
    Console.WriteLine($"Currently importing: {pathToXmlDir}");

    if (Directory.Exists(pathToXmlDir))
    {
        try
        {
            using (ILoggerFactory factory = LoggerFactory.Create(builder => builder.AddConsole()))
            {
                Console.WriteLine("Starting services");
                var dbContext = new AppDBContext();
                Console.WriteLine("Loaded DB Context");

                var polLogger = factory.CreateLogger<IPoliticianService<Politician>>();
                var politicianService = new PoliticianService(dbContext, polLogger);
                Console.WriteLine("Started Politican service");

                var votingService = new VotingService(dbContext, politicianService, polLogger);
                Console.WriteLine("Started Voting Service");

                var partyService = new PartyService(dbContext, factory.CreateLogger<IPartyService<Party>>());
                Console.WriteLine("Started Party Service");

                var roundService = new VotingRoundService(dbContext, votingService, (IPoliticianService<Politician>)partyService, factory.CreateLogger<IVotingRoundService<Round>>());
                Console.WriteLine("Started roundService");

                var dataImport = new VotingDataimporter(votingService, politicianService, partyService, roundService);
                /*
                 * foreach (var file in Directory.GetFiles(pathToXmlDir))
                {
                    ImportFile(file,dataImport);
                }
                 * */

                //var imageImporter = new ImageImporter(politicianService);
                //imageImporter.SearchForNewUrls(true);

                var scrapper = new ParliamentScraper(politicianService, partyService);
                await scrapper.ScrapeDeputiesAsync();
            }
        }
        catch (Exception ex)
        {
            Console.WriteLine($"Failed to load service because:{ex.Message}");
        }

    }
}