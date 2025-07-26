using ParliamentMonitor.Contracts.Model.Votes;
using System;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace ParliamentMonitor.ServiceImplementation.Utils
{
    internal class VotingJsonSerializer : JsonConverter<Vote>
    {
        public override Vote? Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            using var jsonDoc = JsonDocument.ParseValue(ref reader);
            var root = jsonDoc.RootElement;

            var vote = new Vote
            {
                Id = root.GetProperty("id").GetGuid(),
                Position = Enum.TryParse<VotePosition>(root.GetProperty("position").GetString(), out var pos) ? pos : VotePosition.Absent,
                PoliticianId = root.TryGetProperty("politicianId", out var polIdProp) ? polIdProp.GetGuid() : null,
                RoundId = root.TryGetProperty("roundId", out var roundIdProp) ? roundIdProp.GetGuid() : null
            };

            return vote;
        }

        public override void Write(Utf8JsonWriter writer, Vote value, JsonSerializerOptions options)
        {
            writer.WriteStartObject();

            writer.WriteString("id", value.Id.ToString());
            writer.WriteString("position", value.Position.ToString());
            writer.WriteString("politicianId", value.Politician.Id.ToString());
            writer.WriteString("roundId", value.Round.Id.ToString());
            writer.WriteEndObject();
        }
    }
}