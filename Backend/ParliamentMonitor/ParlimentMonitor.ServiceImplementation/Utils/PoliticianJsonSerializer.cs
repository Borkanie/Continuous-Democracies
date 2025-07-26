using ParliamentMonitor.Contracts.Model;
using System;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace ParliamentMonitor.ServiceImplementation.Utils
{
    internal class PoliticianJsonSerializer : JsonConverter<Politician>
    {
        public override Politician? Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            using var jsonDoc = JsonDocument.ParseValue(ref reader);
            var root = jsonDoc.RootElement;

            var politician = new Politician
            {
                Id = root.GetProperty("id").GetGuid(),
                Name = root.GetProperty("name").GetString() ?? string.Empty,
                Gender = Enum.TryParse<Gender>(root.GetProperty("gender").GetString(), out var gender) ? gender : Gender.Male,
                ImageUrl = root.TryGetProperty("imageUrl", out var imgProp) ? imgProp.GetString() : null,
                PartyId = root.TryGetProperty("partyId", out var partyIdProp) ? partyIdProp.GetGuid() : null,
                Active = root.TryGetProperty("active", out var activeProp) && activeProp.GetBoolean(),
                WorkLocation = Enum.TryParse<WorkLocation>(root.GetProperty("workLocation").GetString(), out var loc) ? loc : WorkLocation.Parliament
            };

            return politician;
        }

        public override void Write(Utf8JsonWriter writer, Politician value, JsonSerializerOptions options)
        {
            writer.WriteStartObject();
            writer.WriteString("id", value.Id.ToString());
            writer.WriteString("name", value.Name);
            writer.WriteString("gender", value.Gender.ToString());
            if (value.ImageUrl != null)
                writer.WriteString("imageUrl", value.ImageUrl);
            writer.WriteString("partyId", value.Party.Id.ToString());
            writer.WriteBoolean("active", value.Active);
            writer.WriteString("workLocation", value.WorkLocation.ToString());
            writer.WriteEndObject();
        }
    }
}