// See https://aka.ms/new-console-template for more information
using DataImporter;
using Microsoft.EntityFrameworkCore;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using ParliamentMonitor.ServiceImplementation;
using System.Text;

Encoding.RegisterProvider(CodePagesEncodingProvider.Instance);
var pathToXmlDir = "D:\\Voturi\\Votes";
Console.WriteLine($"Currently importing: {pathToXmlDir}");

if (Directory.Exists(pathToXmlDir))
{
    try
    {
        Console.WriteLine("Starting services");
        var dbContext = new AppDBContext();
        Console.WriteLine("Loaded DB Context");

        var votingService = new VotingService(dbContext);
        Console.WriteLine("Started Voting Service");

        var partyService = new PartyService(dbContext);
        Console.WriteLine("Started Party Service");

        var politicianService = new PoliticianService(dbContext);
        Console.WriteLine("Started Politican service");


        var dataImport = new VotingDataimporter(votingService, politicianService, partyService);
        foreach (var file in Directory.GetFiles(pathToXmlDir))
        {
            ImportFile(file,dataImport);
        }
    }
    catch(Exception ex)
    {
        Console.WriteLine($"Failed to load service because:{ex.Message}");
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