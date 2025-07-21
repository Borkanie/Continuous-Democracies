using ParliamentMonitor.Contracts.Model;
using StackExchange.Redis;
using System.Text.Json;

namespace ParliamentMonitor.Contracts.Services
{
    public class RedisService
    {
        protected string serviceKey;

        protected readonly IConnectionMultiplexer _redis;
        protected readonly IDatabase _cache;

        protected string MakeKey(String id)
        {
            return $"{serviceKey}:{id}";
        }

        public RedisService(IConnectionMultiplexer redis,string serviceKey)
        {
            _redis = redis;
            _cache = redis.GetDatabase();
            this.serviceKey = serviceKey;
        }

        
        public async Task SetAsync<T>(string key, T value, TimeSpan? expiry = null)
        {
            var json = JsonSerializer.Serialize(value);
            await _cache.StringSetAsync(key, json, expiry);
            await _cache.SetAddAsync(serviceKey, key);
        }

        public async Task<T?> GetAsync<T>(string key)
        {
            var value = await _cache.StringGetAsync(key);
            return value.HasValue
                ? JsonSerializer.Deserialize<T>(value)
                : default;
        }

        public async Task<bool> ExistsAsync(string key) =>
            await _cache.KeyExistsAsync(key);

        public async Task RemoveAsync(string key) =>
            await _cache.KeyDeleteAsync(key);
    }
}
