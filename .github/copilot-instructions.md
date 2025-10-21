## Continuous-Democracies — Agent instructions (concise)

Target: help an AI coding agent become productive quickly in this repository.

1) Big picture (how the pieces fit)
- This repo is a .NET 8 solution (Backend/ParliamentMonitor and FrontEnd/ParlimentMonitor.WebInterface).
- Major components:
  - Backend/ParliamentMonitor/ContinousDemocracyAPI — ASP.NET Core WebAPI exposing controllers under `api/*`.
  - Backend/ParliamentMonitor/ParliamentMonitor.DataBaseConnector — EF Core DbContext and models (PostgreSQL).
  - Backend/ParliamentMonitor/ParlimentMonitor.ServiceImplementation — concrete service implementations (IPartyService, IPoliticianService, IVotingService, etc.).
  - Backend/ParliamentMonitor/ParliamentMonitor.Contracts — shared DTOs/interfaces used across API, services and frontend.
    - FrontEnd/web — SPA (Vite + React + TypeScript) used as the primary frontend client. It lives under `FrontEnd/web` and uses pnpm (see `pnpm-lock.yaml`).
  - Backend/ParliamentMonitor/DataImporter — scrapers and importers (ParliamentScraper, VoteDownloader) which populate the DB.

2) How data flows
- DataImporter scrapes external sources and writes to the DB via `AppDBContext`.
- ServiceImplementation classes use `AppDBContext` (EF Core) and are registered as scoped in the API and frontend.
- Controllers call the generic service interfaces (e.g., `IPartyService<Party>`) and return results to clients.

3) How to build, run and test locally (concrete)
- General restore + build (from repo root):
  - dotnet restore
  - dotnet build --configuration Release
- Run the API locally:
  - dotnet run --project Backend/ParliamentMonitor/ContinousDemocracyAPI/ContinousDemocracyAPI.csproj
 - Run the Frontend (Vite + React SPA):
   - cd FrontEnd/web
   - pnpm install
   - pnpm dev
   - Build for production: `pnpm build` (output in `dist/`)
- Run tests (xUnit):
  - dotnet test Backend/ParliamentMonitor/ParliamentMonitor.Tests
- Migrations / database updates (in DB connector folder):
  - cd Backend/ParliamentMonitor/ParliamentMonitor.DataBaseConnector
  - dotnet ef migrations add <Name>
  - dotnet ef database update
  (The top-level README contains the same commands — keep migrations in the DataBaseConnector project.)

4) CI / Docker / Deployment hints
- buildspec.yaml (root) is used by CI (CodeBuild style). It expects .NET 8.0 and publishes the API to `publish/`.
- Dockerfile for the API produces a runtime image and exposes ports 8080/8081. Entry point: `dotnet ContinousDemocracyAPI.dll`.
- Kubernetes manifests for Postgres live under `Backend/` (ps-deployment.yaml, ps-service.yaml, psql-claim.yaml, psql-pv.yaml).

5) Project-specific patterns & conventions an agent should follow
- Dependency registration:
  - API: services are registered in `ContinousDemocracyAPI/Program.cs` using AddScoped for service implementations and AddSingleton for `AppDBContext` in some places.
  - Frontend registers a singleton `AppDBContext` (see FrontEnd/ParlimentMonitor.WebInterface/Program.cs).
- Generic service interfaces: most services implement generic interfaces like `IPartyService<Party>`. When editing services prefer keeping the generic signatures.
- Controllers frequently call `.Result` on async service methods (synchronous waiting). Be careful when refactoring to avoid deadlocks — prefer returning async Task<ActionResult<T>> and awaiting where appropriate.
- EF Core model nuance: `AppDBContext.OnModelCreating` converts System.Drawing.Color to a string via a ValueConverter — preserve this when changing party color handling.

6) Integration & secrets
- The DbContext uses PostgreSQL (RDS). `Program.cs` for the API reads configuration key `RDS` for the connection string: `builder.Configuration.GetConnectionString("RDS")`.
- NOTE: There is a hard-coded connection string present in `AppDBContext`'s parameterless constructor — treat this as sensitive. Prefer using environment variables or appsettings and do not commit secrets.

7) Tests & test helpers
- Tests use xUnit and `Microsoft.EntityFrameworkCore.InMemory` to create in-memory DBs for service tests. Look at `ParliamentMonitor.Tests` for examples of how services are instantiated for tests.
- Coverage collector: `coverlet.collector` is referenced in the test csproj.

8) Useful files to open when working on a change
- Backend/ParliamentMonitor/ContinousDemocracyAPI/Program.cs (DI, middleware)
- Backend/ParliamentMonitor/ContinousDemocracyAPI/Controllers/* (PartyController, VotingController, PoliticiansController)
- Backend/ParliamentMonitor/ParliamentMonitor.DataBaseConnector/AppDBContext.cs (models, converters, connection patterns)
- Backend/ParliamentMonitor/ParlimentMonitor.ServiceImplementation/* (PartyService.cs, VotingService.cs, PoliticianService.cs)
 - FrontEnd/web/package.json (Vite React client): check `scripts` for `dev`/`build` commands and `pnpm-lock.yaml` to confirm pnpm usage.
- Backend/ParliamentMonitor/DataImporter/* (scrapers and importers)
- buildspec.yaml and Dockerfile (CI and containerization guidance)

9) Security & safety hints for code changes
- Avoid leaving credentials in code. If you find secrets in source, flag them and prefer env/config-based secrets. Replace hard-coded strings with `builder.Configuration.GetConnectionString("RDS")` or secured secret store.

10) When proposing edits, include
- Specific files changed and why (one-line).
- A short developer test plan: how to run locally and which test(s) to run.
- If changing DB schema, include EF migration commands and where to run them.

If any section is unclear or you want deeper detail (example calls, model shapes, or tests), tell me which area to expand and I will iterate.
