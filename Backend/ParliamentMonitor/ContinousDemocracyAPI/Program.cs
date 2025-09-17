using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using ParliamentMonitor.ServiceImplementation;
using ContinousDemocracyAPI.MiddleWare;

var builder = WebApplication.CreateBuilder(args);
Console.WriteLine("Starting ContinousDemocracyAPI...");

builder.Services.AddSwaggerGen(c =>
{
    c.SwaggerDoc("v1",
        new Microsoft.OpenApi.Models.OpenApiInfo
        {
            Title = "API democratia contiua",
            Version = "v1"
        }
     );

    var filePath = Path.Combine(System.AppContext.BaseDirectory, "MyApi.xml");
    c.IncludeXmlComments(filePath);
}
);

// Add services to the container.
builder.Services.AddSingleton(new AppDBContext(builder.Configuration.GetConnectionString("RDS")!));
builder.Services.AddScoped<IPartyService<Party>, PartyService>();
builder.Services.AddScoped<IPoliticianService<Politician>, PoliticianService>();
builder.Services.AddScoped<IVotingService<Vote>, VotingService>();
builder.Services.AddScoped<IVotingRoundService<Round>, VotingRoundService>();
builder.Services.AddControllers();

// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

app.UseMiddleware<RequestLoggingMiddleware>();

// Configure the HTTP request pipeline.
app.UseSwagger();
app.UseSwaggerUI();

app.UseHttpsRedirection();

app.UseAuthorization();

app.MapControllers();

app.MapGet("/health", () => Results.Ok("Healthy"));

Console.WriteLine("ContinousDemocracyAPI started successfully.");


app.Run();