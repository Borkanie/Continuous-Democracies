using ParliamentMonitor.Contracts.Model;
using StackExchange.Redis;
using System.Text.Json;
using System.Text.Json.Serialization;

namespace ParliamentMonitor.Contracts.Services
{
    public class RedisService
    {
        protected string serviceKey;

        protected readonly IConnectionMultiplexer _redis;
        protected readonly IDatabase _cache;
        protected readonly JsonConverter? _customConverter;

        public virtual string MakeKey(string id)
        {
            return $"{serviceKey}:{id}";
        }

        public RedisService(IConnectionMultiplexer redis, string serviceKey, JsonConverter? customConverter = null)
        {
            _redis = redis;
            _cache = redis.GetDatabase();
            this.serviceKey = serviceKey;
            _customConverter = customConverter;
        }

        public async Task SetAsync<T>(string key, T value, TimeSpan? expiry = null)
        {
            var options = new JsonSerializerOptions();
            if (_customConverter != null)
                options.Converters.Add(_customConverter);
            var json = JsonSerializer.Serialize(value, options);
            await _cache.StringSetAsync(key, json, expiry);
            await _cache.SetAddAsync(serviceKey, key);
        }

        public async Task<T?> GetAsync<T>(string key)
        {
            var value = await _cache.StringGetAsync(key);
            if (!value.HasValue)
                return default;

            var options = new JsonSerializerOptions();
            if (_customConverter != null)
                options.Converters.Add(_customConverter);
            return JsonSerializer.Deserialize<T>(value, options);
        }

        public async Task<bool> ExistsAsync(string key) =>
            await _cache.KeyExistsAsync(key);

        public async Task RemoveAsync(string key) =>
            await _cache.KeyDeleteAsync(key);
    }
}