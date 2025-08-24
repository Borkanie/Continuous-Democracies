namespace ContinousDemocracyAPI.MiddleWare
{
    public class RequestLoggingMiddleware
    {
        private readonly RequestDelegate _requestDelegate;
        private readonly ILogger<RequestLoggingMiddleware> _logger;

        public RequestLoggingMiddleware(RequestDelegate next, ILogger<RequestLoggingMiddleware> logger)
        {
            _requestDelegate = next;
            _logger = logger;
        }

        public async Task InvokeAsync(HttpContext context)
        {
            var request = context.Request;
            var ip = context.Connection.RemoteIpAddress?.ToString();

            // Log request
            _logger.LogInformation($"Incoming {request.Method} {request.Path} from {ip}");

            // Capture response
            var originalBodyStream = context.Response.Body;
            using var responseBody = new MemoryStream();
            context.Response.Body = responseBody;

            // pass forward the request to the next middleware
            await _requestDelegate(context);

            context.Response.Body.Seek(0, SeekOrigin.Begin);
            var responseText = await new StreamReader(context.Response.Body).ReadToEndAsync();
            context.Response.Body.Seek(0, SeekOrigin.Begin);

            _logger.LogInformation($"Response {context.Response.StatusCode} for {request.Path} to {ip} with body: {responseText}");

            await responseBody.CopyToAsync(originalBodyStream);
        }
    }

}
