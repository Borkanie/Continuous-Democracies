using ParliamentMonitor.Contracts.Model;
using System;
using System.Drawing;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace ParliamentMonitor.ServiceImplementation.Utils
{
    internal class PartyJsonSerializer : JsonConverter<Party>
    {
        public override Party? Read(ref Utf8JsonReader reader, Type typeToConvert, JsonSerializerOptions options)
        {
            using var jsonDoc = JsonDocument.ParseValue(ref reader);
            var root = jsonDoc.RootElement;

            var party = new Party
            {
                Id = root.GetProperty("id").GetGuid(),
                Name = root.GetProperty("name").GetString() ?? string.Empty,
                Acronym = root.TryGetProperty("acronym", out var acrProp) ? acrProp.GetString() : null,
                LogoUrl = root.TryGetProperty("logoUrl", out var logoProp) ? logoProp.GetString() : null,
                Color = root.TryGetProperty("color", out var colorProp) ? ColorTranslator.FromHtml(colorProp.GetString() ?? "#D3D3D3") : Color.LightGray,
                Active = root.TryGetProperty("active", out var activeProp) && activeProp.GetBoolean()
            };

            if (root.TryGetProperty("politicianIds", out var polIdsProp))
            {
                var ids = new HashSet<Guid>();
                foreach (var idElem in polIdsProp.EnumerateArray())
                {
                    if (idElem.TryGetGuid(out var id))
                        ids.Add(id);
                }
                party.PoliticianIds = ids;
            }

            return party;
        }

        public override void Write(Utf8JsonWriter writer, Party value, JsonSerializerOptions options)
        {
            writer.WriteStartObject();
            writer.WriteString("id", value.Id.ToString());
            writer.WriteString("name", value.Name);
            if (value.Acronym != null)
                writer.WriteString("acronym", value.Acronym);
            if (value.LogoUrl != null)
                writer.WriteString("logoUrl", value.LogoUrl);
            writer.WriteString("color", ColorTranslator.ToHtml(value.Color));
            writer.WriteBoolean("active", value.Active);
            writer.WritePropertyName("politicianIds");
            writer.WriteStartArray();
            foreach (var politician in value.Politicians)
                writer.WriteStringValue(politician.Id.ToString());
            writer.WriteEndArray();
            writer.WriteEndObject();
        }
    }
}