using System;
using System.IO;
using System.Net.Http;
using System.Text;
using System.Threading.Tasks;

namespace ParliamentDownloader
{
    public class VoteDownloader
    {
        private readonly HttpClient _httpClient;
        private readonly string _outputFolder;

        public VoteDownloader(string outputFolder = "Votes")
        {
            _httpClient = new HttpClient();
            _outputFolder = outputFolder;

            // Ensure output folder exists
            Directory.CreateDirectory(_outputFolder);
        }

        public async Task DownloadVotesAsync(int par1, int startVoteIndex, int endVoteIndex, int delaySeconds = 5)
        {
            for (int voteIndex = startVoteIndex; voteIndex > endVoteIndex; voteIndex--)
            {
                if(File.Exists($"{_outputFolder}\\vote_{voteIndex}.xml"))
                {
                    Console.WriteLine($"vote_{voteIndex}.xml already exists");
                    continue;
                }
                string url = $"https://www.cdep.ro/pls/steno/evot2015.xml?par1={par1}&par2={voteIndex}";
                Console.WriteLine($"[{DateTime.Now:HH:mm:ss}] Downloading vote {voteIndex}...");

                try
                {
                    var response = await _httpClient.GetAsync(url);
                    response.EnsureSuccessStatusCode();

                    // Decode ISO-8859-2 content and re-encode to UTF-8
                    byte[] rawBytes = await response.Content.ReadAsByteArrayAsync();
                    string content = Encoding.GetEncoding("ISO-8859-2").GetString(rawBytes);

                    string filename = Path.Combine(_outputFolder, $"vote_{voteIndex}.xml");
                    await File.WriteAllTextAsync(filename, content, Encoding.UTF8);

                    Console.WriteLine($"Saved to {filename}");
                }
                catch (Exception ex)
                {
                    Console.WriteLine($"Error fetching vote {voteIndex}: {ex.Message}");
                }

                // Delay between requests
                await Task.Delay(delaySeconds * 1000);
            }

            Console.WriteLine("Finished downloading all votes.");
        }
    }
}
