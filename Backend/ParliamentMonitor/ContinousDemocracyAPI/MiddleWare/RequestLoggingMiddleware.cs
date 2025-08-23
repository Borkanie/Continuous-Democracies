namespace ContinousDemocracyAPI.MiddleWare
{
    public class RequestLoggingMiddleware
    {
        private readonly RequestDelegate _next;
        private readonly ILogger<RequestLoggingMiddleware> _logger;

        public RequestLoggingMiddleware(RequestDelegate next, ILogger<RequestLoggingMiddleware> logger)
        {
            _next = next;
            _logger = logger;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            var request = context.Request;
            var ip = context.Connection.RemoteIpAddress?.ToString();

            // Log request
            _logger.LogInformation("Incoming {Method} {Path} from {IP}",
                request.Method, request.Path, ip);

            // Capture response
            var originalBodyStream = context.Response.Body;
            using var responseBody = new MemoryStream();
            context.Response.Body = responseBody;

            await _next(context);

            context.Response.Body.Seek(0, SeekOrigin.Begin);
            var responseText = await new StreamReader(context.Response.Body).ReadToEndAsync();
            context.Response.Body.Seek(0, SeekOrigin.Begin);

            _logger.LogInformation("Response {StatusCode} for {Path}",
                context.Response.StatusCode, request.Path);

            await responseBody.CopyToAsync(originalBodyStream);
        }
    }

}
