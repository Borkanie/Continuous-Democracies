using System.Text;

namespace DataImporter
{
    internal static class UrlBuilder
    {
        private static readonly Dictionary<char, char> charMap = new()
        {
        // Romanian
        { 'ă', 'a' }, { 'â', 'a' }, { 'î', 'i' }, { 'ș', 's' }, { 'ş', 's' }, { 'ț', 't' }, { 'ţ', 't' },
        { 'Ă', 'A' }, { 'Â', 'A' }, { 'Î', 'I' }, { 'Ș', 'S' }, { 'Ş', 'S' }, { 'Ț', 'T' }, { 'Ţ', 'T' },

        // Hungarian
        { 'á', 'a' }, { 'é', 'e' }, { 'í', 'i' }, { 'ó', 'o' }, { 'ö', 'o' }, { 'ő', 'o' },
        { 'ú', 'u' }, { 'ü', 'u' }, { 'ű', 'u' },
        { 'Á', 'A' }, { 'É', 'E' }, { 'Í', 'I' }, { 'Ó', 'O' }, { 'Ö', 'O' }, { 'Ő', 'O' },
        { 'Ú', 'U' }, { 'Ü', 'U' }, { 'Ű', 'U' }
    };

        public static string CleanName(string input)
        {
            if (string.IsNullOrWhiteSpace(input))
                return string.Empty;

            // Split by whitespace (last word = last name)
            var parts = input.Split(' ', StringSplitOptions.RemoveEmptyEntries);
            if (parts.Length == 0)
                return string.Empty;

            var lastName = parts[^1]; // last word
            var rest = parts.Take(parts.Length - 1);

            // Reorder: last name first, then rest
            var reordered = string.Join("", new[] { lastName }.Concat(rest));

            var sb = new StringBuilder(reordered.Length);

            foreach (var c in reordered)
            {
                if (c == ' ' || c == '-') // remove spaces and dashes
                    continue;

                if (charMap.ContainsKey(c))
                    _ = sb.Append(charMap[c]); // replace special chars
                else
                    _ = sb.Append(c); // keep as is
            }

            return sb.ToString() + ".jpg";
        }
    }
}

