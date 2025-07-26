using ParliamentMonitor.Contracts.Model.Votes;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Text.Json;
using System.Text.Json.Serialization;
using System.Threading.Tasks;

namespace ParliamentMonitor.ServiceImplementation.Utils
{
    internal class VotingRoundJsonSerializer : JsonConverter<Round>
    {
        public override Round? Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            using var jsonDoc = JsonDocument.ParseValue(ref reader);
            var root = jsonDoc.RootElement;

            var round = new Round
            {
                Id = root.GetProperty("id").GetGuid(),
                Title = root.GetProperty("title").GetString() ?? "",
                Description = root.GetProperty("description").GetString() ?? "",
                VoteDate = root.GetProperty("voteDate").GetDateTime(),
                VoteId = root.GetProperty("voteId").GetInt32()
            };

            if (root.TryGetProperty("voteResultIds", out var voteIdsProp))
            {
                var voteIds = voteIdsProp.EnumerateArray().Select(e => e.GetGuid());
                if(voteIds!= null)
                    round.VoteResultIds = voteIds.ToHashSet();
            }

            return round;
        }

        public override void Write(Utf8JsonWriter writer, Round value, JsonSerializerOptions options)
        {
            writer.WriteStartObject();

            writer.WriteString("id", value.Id.ToString());
            writer.WriteString("title", value.Title);
            writer.WriteString("description", value.Description);
            writer.WriteString("voteDate", value.VoteDate);
            writer.WriteNumber("voteId", value.VoteId);

            writer.WritePropertyName("voteResultIds");
            writer.WriteStartArray();
            foreach (var vote in value.VoteResults)
            {
                writer.WriteStringValue(vote.Id.ToString());
            }
            writer.WriteEndArray();

            writer.WriteEndObject();
        }
    }
}