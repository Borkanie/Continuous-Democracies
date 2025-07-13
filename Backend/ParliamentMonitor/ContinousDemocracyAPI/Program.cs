using ParliamentMonitor.Contracts.Model.Votes;
using ParliamentMonitor.Contracts.Model;
using ParliamentMonitor.Contracts.Services;
using ParliamentMonitor.DataBaseConnector;
using ParliamentMonitor.ServiceImplementation;
using StackExchange.Redis;

var builder = WebApplication.CreateBuilder(args);

// Add Redis singleton
builder.Services.AddSingleton<IConnectionMultiplexer>(sp =>
{
    var configuration = builder.Configuration.GetConnectionString("Redis");
    return ConnectionMultiplexer.Connect(configuration);
});

// Add services to the container.
builder.Services.AddSingleton(new AppDBContext());
builder.Services.AddScoped<IVotingService<Vote, Round>, VotingService>();
builder.Services.AddScoped<IPartyService<Party>, PartyService>();
builder.Services.AddScoped<IPoliticianService<Politician>, PoliticianService>();
builder.Services.AddControllers();

// Learn more about configuring Swagger/OpenAPI at https://aka.ms/aspnetcore/swashbuckle
builder.Services.AddEndpointsApiExplorer();
builder.Services.AddSwaggerGen();

var app = builder.Build();

// Configure the HTTP request pipeline.
if (app.Environment.IsDevelopment())
{
    app.UseSwagger();
    app.UseSwaggerUI();
}

app.UseHttpsRedirection();

app.UseAuthorization();

app.MapControllers();

app.Run();
