using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Globalization;
using System.Text.RegularExpressions;

namespace DataImporter
{
    internal class ImageImporter
    {

        IPoliticianService<Politician> politicianService;

        private static readonly HttpClient client = new HttpClient();

        public static async Task DownloadFileAsync(string fileUrl, string destinationPath)
        {
            Console.WriteLine($"Downloading image from {fileUrl} and saivng it to {destinationPath}");
            using (var response = await client.GetAsync(fileUrl))
            {
                response.EnsureSuccessStatusCode();

                using (var fileStream = new FileStream(destinationPath, FileMode.Create))
                {
                    await response.Content.CopyToAsync(fileStream);
                }
            }
        }

        public void SearchForNewUrls(bool overrideExisting)
        {
            var politicians = politicianService.GetAllPoliticiansAsync(number:10000).Result;
            var baseUrl = "https://www.cdep.ro/parlamentari/l2024/";
            foreach (var politician in politicians)
            {
                if(overrideExisting || politician.ImageUrl == null)
                {
                    LoadDataForPolitician(politician, baseUrl);
                }
            }
        }

        public static string MovelastNameAtTheBeginning(string fullName)
        {
            var names = fullName.Split(' ');
            var result = names.Last();
            for(int i=0;i<names.Length -1; i++)
            {
                result += " " + names[i];
            }
            return result;
        }

        public static string CleanName(string input)
        {
            if (string.IsNullOrEmpty(input))
                return string.Empty;
            Console.WriteLine($"Cleaning {input}");
            input = MovelastNameAtTheBeginning(input);
            // 1. Replace Romanian & Hungarian diacritics with ASCII equivalents
            string normalized = input.Normalize(NormalizationForm.FormD);
            var sb = new StringBuilder();

            foreach (var c in normalized)
            {
                var unicodeCategory = CharUnicodeInfo.GetUnicodeCategory(c);
                if (unicodeCategory != UnicodeCategory.NonSpacingMark)
                {
                    sb.Append(c);
                }
            }

            string asciiOnly = sb.ToString().Normalize(NormalizationForm.FormC);

            // 2. Remove hyphens
            asciiOnly = asciiOnly.Replace("-", "");

            // 3. Remove any other characters not allowed in URLs (keep letters, numbers, underscore, dot and space)
            string urlSafe = Regex.Replace(asciiOnly, @"[^a-zA-Z0-9 _\.]", "");

            // 4. Optionally, replace spaces with underscores or hyphens
            urlSafe = urlSafe.Replace(" ", "");

            Console.WriteLine($"Result: {urlSafe}");
            return urlSafe;
        }

        public ImageImporter( IPoliticianService<Politician> politicianService)
        {
            this.politicianService = politicianService;
        }


        public void LoadDataForPolitician(Politician target, string baseSourceUrl)
        {
            string targetLocationRelative = CleanName(target.Name.Trim());
            string targetUrl = baseSourceUrl + targetLocationRelative + ".JPG";
            var file = "D:/Continuous-Democracies/Resources/DeputyPortraits/" + targetLocationRelative + ".jpg";
            Console.WriteLine($"Image Path: {file}");
            if(!File.Exists(file))
            {
                try
                {
                    DownloadFileAsync(targetUrl, file).Wait();
                    politicianService.UpdatePoliticianAsync(target.Id, imageUrl: file);
                }catch(Exception ex)
                {
                    Console.WriteLine($" Error when lading image for: {target.Name}. \n Message:{ex.Message}");
                }
            }
            
        }
    }
}
