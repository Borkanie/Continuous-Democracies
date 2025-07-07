using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;
using System.Xml.Linq;

namespace DataImporter
{
    internal class VotingDataimporter
    {
        private readonly IVotingService<Vote,Round> votingService;
        private readonly IPoliticianService<Politician> politicianService;
        private readonly IPartyService<Party> partyService;

        public VotingDataimporter( IVotingService<Vote, Round> votingService, IPoliticianService<Politician> politicianService, IPartyService<Party> partyService)
        {
            this.politicianService = politicianService;
            this.partyService = partyService;
            this.votingService = votingService;
        }

        public void ImportVotingRoundFromXml(string xmlFilePath)
        {

            using (var reader = new StreamReader(xmlFilePath, Encoding.UTF8))
            {
                var doc = XDocument.Load(reader);

                var votElement = doc.Descendants("ROWSET");

                var Id = votElement.Descendants("VOTID").First().Value;
                int id;
                if (int.TryParse(Id, out id))
                {
                    Console.WriteLine($"Vote Id:{id}");
                    var round = votingService.GetVotingRound(id);
                    var date = DateTime.UtcNow;
                    if (round != null)
                    {
                        Console.WriteLine($"Updatate voting round:{round.VoteId}");
                    }else
                    {
                        round = votingService.CreateVotingRound($"Vot electronic:{id}", date, id: id);
                    }
                    
                    var votesXML = votElement.Descendants("ROW");
                    List<Vote> votes = new List<Vote>();
                    Console.WriteLine($"Number of votes:{votesXML.Count()}");
                    foreach (var voteXML in votesXML)
                    {
                        CastVote(votElement, round, date, voteXML);
                    }
                }
                else
                {
                    Console.WriteLine($"Cannot parse {Id}, execution will be stopped");
                    return;
                }
            }

            
        }

        private void CastVote(IEnumerable<XElement> votElement, Round? round, DateTime date, XElement voteXML)
        {
            Console.WriteLine($"Processing :{voteXML} \n\n");
            var vote = new Vote();
            vote.Position = ConvertStringToVotePositon(votElement.Descendants("VOT").First().Value);
            vote.Politician = getPoliticianAndPartyFromVote(voteXML);
            vote.Round = round;
            var current = round.VoteResults.FirstOrDefault(x => x == vote);
            if (round.VoteDate != date || current == null)
            {
                Console.WriteLine($"Politician: {vote.Politician.Name} casted a vote");
                votingService.CastVote(round, vote.Politician, vote.Position);
            }
            else
            {
                Console.WriteLine($"Politician: {vote.Politician.Name} has changed it's vote");
                votingService.UpdateCastedVote(round, vote.Politician, vote.Position);
            }
        }

        private Politician getPoliticianAndPartyFromVote(XElement element)
        {
            var prename = element.Element("PRENUME")?.Value;
            var name = element.Element("NUME")?.Value;
            var partyName = element.Element("GRUP")?.Value;
            
            var party = partyService.GetParty(acronym: partyName);
            if (party != null)
            {
                Console.WriteLine($"Already Existing party {partyName} has voted");
            }
            else
            {
                party = partyService.CreateParty("partid", acronym: partyName);
            }

            var politican = politicianService.GetPolitician(name + " " + prename);
            if (politican != null)
            {
                Console.WriteLine($"Politican :{politican.Name} has voted");
            }
            else
            {
                politican = politicianService.CreatePolitican(prename + " " + name, party, WorkLocation.Parliament, Gender.Other, true)!;
            }

            return politican;
        }

        private VotePosition ConvertStringToVotePositon(string vote)
        {
            if(vote == "DA")
            {
                return VotePosition.Yes;
            }
            if (vote == "NU")
            {
                return VotePosition.No;
            }
            if (vote == "-")
            {
                return VotePosition.Abstain;
            }
            return VotePosition.Absent;
        }
    }
}
