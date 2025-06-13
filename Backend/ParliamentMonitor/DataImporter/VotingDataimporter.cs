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

            using (var reader = new StreamReader(xmlFilePath, Encoding.GetEncoding("ISO-8859-2")))
            {
                var doc = XDocument.Load(reader);

                var votElement = doc.Descendants("ROWSET");

                var Id = votElement.Descendants("VOTID").First().Value;
                int id;
                if (int.TryParse(Id, out id))
                {
                    Console.WriteLine($"Vote Id:{id}");
                    var round = votingService.GetVotingRound(id) ?? votingService.CreateVotingRound($"Vot electronic:{id}", DateTime.UtcNow);
                    var votesXML = votElement.Descendants("ROW");
                    List<Vote> votes = new List<Vote>();
                    Console.WriteLine($"Number of votes:{votesXML.Count()}");
                    foreach (var voteXML in votesXML)
                    {
                        Console.WriteLine($"Processing :{voteXML} \n\n");
                        var vote = new Vote();
                        vote.Position = ConvertStringToVotePositon(votElement.Descendants("VOT").First().Value);
                        vote.Politician = getPoliticianFromVote(voteXML);
                        vote.Round = round;
                        votes.Add(vote);
                    }
                    votingService.UpdateVoteResult(round.Id, votes: votes);
                }
                else
                {
                    Console.WriteLine($"Cannot parse {Id}, execution will be stopped");
                    return;
                }
            }

            
        }

        private Politician getPoliticianFromVote(XElement element)
        {
            var prename = element.Element("PRENUME")?.Value;
            var name = element.Element("NUME")?.Value;
            var prenume = element.Element("PRENUME")?.Value;
            var partyName = element.Element("GRUP")?.Value;

            var party = partyService.GetParty(acronym: partyName) ?? partyService.CreateParty("partid", acronym: partyName);

            return politicianService.GetPolitician(prename + " " + name) ?? 
                politicianService.CreatePolitican(prename + " " + name, party, WorkLocation.Parliament, Gender.Other, true)!;
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
