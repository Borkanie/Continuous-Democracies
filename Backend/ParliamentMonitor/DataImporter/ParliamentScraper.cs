using HtmlAgilityPack;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace DataImporter
{
    internal class ParliamentScraper
    {
        IPoliticianService<Politician> politicianService;

        IPartyService<Party> partyService;
        public ParliamentScraper(IPoliticianService<Politician> politicianService, IPartyService<Party> partyService)
        {
            this.politicianService = politicianService;
            this.partyService = partyService;
        }

        private readonly HttpClient client = new HttpClient();

        public static string MoveCapsNameToEnd(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return input;

            var words = input.Split(' ', StringSplitOptions.RemoveEmptyEntries).ToList();

            // Find the index of the fully UPPERCASE word
            var capsWord = words.FirstOrDefault(w => w.All(char.IsUpper));

            if (!string.IsNullOrEmpty(capsWord))
            {
                words.Remove(capsWord);
                words.Add(capsWord);
            }

            return string.Join(" ", words);
        }

        private Party? GetPartyById(int id)
        {
            lock (partyService)
            {
                switch (id)
                {
                    case 1:
                        return partyService.GetPartyAsync(acronym: "PSD").Result;
                    case 2:
                        return partyService.GetPartyAsync(acronym: "AUR").Result;
                    case 3:
                        return partyService.GetPartyAsync(acronym: "PNL").Result;
                    case 4:
                        return partyService.GetPartyAsync(acronym: "USR").Result;
                    case 5:
                        return partyService.GetPartyAsync(acronym: "SOS").Result;
                    case 6:
                        return partyService.GetPartyAsync(acronym: "POT").Result;
                    case 7:
                        return partyService.GetPartyAsync(acronym: "UDMR").Result;
                    default:
                        return partyService.GetPartyAsync(acronym: "neafiliat").Result;
                }
            }
        }

        public async Task ScrapeDeputiesAsync()
        {
            for (int cam = 1; cam <= 2; cam++)
            {
                int maxIdm = (cam == 1) ? 180 : 330;

                for (int idm = 2; idm <= maxIdm; idm++)
                {
                    string url = $"https://www.cdep.ro/pls/parlam/structura2015.mp?idm={idm}&cam={cam}&leg=2024";
                    Console.WriteLine($"navigating to {url}");
                    try
                    {
                        var html = await client.GetStringAsync(url);
                        var doc = new HtmlDocument();
                        doc.LoadHtml(html);

                        // Get name
                        var nameNode = doc.DocumentNode.SelectSingleNode("//div[@class='boxTitle']/h1");
                        string name = nameNode?.InnerText.Trim() ?? String.Empty;


                        // Get image URL
                        var imgNode = doc.DocumentNode.SelectSingleNode("//div[@class='profile-pic-dep']/img");
                        string imageUrl = imgNode?.GetAttributeValue("src", "") ?? "";

                        // Get party ID (from hyperlink's idg parameter)
                        var partyLinkNode = doc.DocumentNode.SelectSingleNode("//div[@class='boxDep clearfix']//a[contains(@href, 'idg=')]\r\n");
                        string partyHref = partyLinkNode?.GetAttributeValue("href", "") ?? "";
                        int partyId = int.Parse(GetQueryParameterValue(partyHref, "idg"));

                        if (!String.IsNullOrEmpty(name))
                        {
                            Console.WriteLine($"Cuurently downloading for: {name}");
                            name = MoveCapsNameToEnd(name);
                            lock (politicianService)
                            {
                                var politician = politicianService.GetPoliticianAsync(name: name).Result;
                                if(politician == null)
                                {
                                    politician = politicianService.CreatePoliticanAsync(name, GetPartyById(partyId), cam == 1 ? WorkLocation.Parliament : WorkLocation.Senate, Gender.Male).Result;
                                }
                                if (!String.IsNullOrEmpty(imageUrl) && String.IsNullOrEmpty(politician.ImageUrl))
                                {
                                    getImageForPolitician(politician, "https://www.cdep.ro" + imageUrl);
                                }
                            }
                        }

                        // Output for now
                        Console.WriteLine($"[{idm}/{cam}] Name: {name}, PartyID: {partyId}, Image: {imageUrl}");
                    }
                    catch (Exception ex)
                    {
                        Console.WriteLine($"Error fetching idm {idm}, cam {cam}: {ex.Message}");
                    }
                }
            }
        }

        public void getImageForPolitician(Politician target, string imageUrl)
        {
            string targetLocationRelative = ImageImporter.CleanName(target.Name.Trim());
            var file = "D:/Continuous-Democracies/Resources/DeputyPortraits/" + targetLocationRelative + ".jpg";
            Console.WriteLine($"Image Path: {file}");
            if (!File.Exists(file))
            {
                try
                {
                    ImageImporter.DownloadFileAsync(imageUrl, file).Wait();
                    politicianService.UpdatePoliticianAsync(target.Id, imageUrl: file);
                }
                catch (Exception ex)
                {
                    Console.WriteLine($" Error when lading image for: {target.Name}. \n Message:{ex.Message}");
                }
            }
        }

        private string GetQueryParameterValue(string url, string parameterName)
        {
            if (string.IsNullOrEmpty(url) || !url.Contains("?"))
                return "";

            var query = url.Split('?')[1];
            var parameters = query.Split('&');

            foreach (var param in parameters)
            {
                var keyValue = param.Split('=');
                if (keyValue.Length == 2 && keyValue[0] == parameterName)
                    return keyValue[1];
            }
            return "";
        }
    }
}
