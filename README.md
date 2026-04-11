# Continuous-Democracies
This project will allow us to see what the fuck the parlament is doing on our money at all times.

## Architecture

### Deployment Pipeline

The CI/CD pipeline automates the build, containerization, and deployment of all application components using GitHub Actions workflows that trigger on code changes to the main branch. Three parallel workflows handle the backend API (.NET 8), frontend (React/Vite), and data ingestion scripts (Python), building Docker images and deploying them to a self-hosted runner environment. The Python orchestrator runs continuously in a daemon container, executing 12-hour cycles to scrape new voting data from the Romanian Parliament website (cdep.ro), import laws and votes, and update politician active status. All components are containerized with specific images (`borkanie/continousdemocracyapi:latest`, `borkanie/parliament-web:latest`, `borkanie/populate-laws:latest+SHA`) and deployed alongside a Kubernetes-managed PostgreSQL cluster with 3 replicas and persistent storage.

**Key Components:**
- **GitHub Actions Workflows**: `backend-container.yml` (API on ports 8080/8081), `docker-image-web.yml` (Frontend on port 5173), `populate-laws-build-deploy.yml` (Python daemon)
- **Build Stages**: Code checkout → Docker build → Push to Docker Hub → Deploy to self-hosted runner with automatic container restart
- **Python Orchestrator**: 12-hour loop executing `import_law.goFromLastforward(100)` to fetch new voting rounds, then `activate_all_groups()` to update politician status
- **Database Infrastructure**: PostgreSQL 14 with 3 replicas, 10Gi persistent volume (`/d/postgre`), exposed on port 5432 via Kubernetes NodePort service
- **Deployment Target**: Self-hosted GitHub Actions runner with Docker runtime, automatic container lifecycle management (stop old → pull latest → run new)
- **Image Versioning**: Backend and frontend use `latest` tag only; populate-laws uses both `latest` and commit SHA for rollback capabilities
- **Environment Configuration**: Database credentials (`DBSERVERADDRESS`, `DBNAME`, `DBUSER`, `DBPASSWORD`) and `LOG_LEVEL` injected via GitHub Secrets
- **Container Restart Policy**: `unless-stopped` for API and frontend ensures automatic recovery; populate-laws runs with `--rm -d` as ephemeral daemon
- **Kubernetes Resources**: ConfigMap for PostgreSQL credentials, PersistentVolume/PersistentVolumeClaim for data persistence, Service definitions for network access

```mermaid
graph TB
    subgraph "Code Repository"
        A[Developer Commits to main]
    end
    
    subgraph "GitHub Actions CI/CD"
        B1[Backend Workflow<br/>backend-container.yml]
        B2[Frontend Workflow<br/>docker-image-web.yml]
        B3[Populate-Laws Workflow<br/>populate-laws-build-deploy.yml]
        
        A --> B1
        A --> B2
        A --> B3
    end
    
    subgraph "Build Stages - Backend API"
        C1[Checkout Code]
        C2[Docker Build<br/>.NET 8.0 Multi-stage]
        C3[Push to Docker Hub<br/>borkanie/continousdemocracyapi:latest]
        C4[Deploy to Self-Hosted Runner]
        
        B1 --> C1 --> C2 --> C3 --> C4
    end
    
    subgraph "Build Stages - Frontend"
        D1[Checkout Code]
        D2[Docker Build<br/>Node 22 + pnpm]
        D3[Push to Docker Hub<br/>borkanie/parliament-web:latest]
        D4[Deploy to Self-Hosted Runner]
        
        B2 --> D1 --> D2 --> D3 --> D4
    end
    
    subgraph "Build Stages - Python Scripts"
        E1[Checkout Code]
        E2[Docker Build<br/>Python 3.11 + Tesseract OCR]
        E3[Push to Docker Hub<br/>borkanie/populate-laws:latest+SHA]
        E4[Deploy to Self-Hosted Runner]
        
        B3 --> E1 --> E2 --> E3 --> E4
    end
    
    subgraph "Self-Hosted Runner Environment"
        F1[API Container<br/>Ports: 8080 HTTP, 8081 HTTPS<br/>Restart: unless-stopped]
        F2[Frontend Container<br/>Port: 5173 Vite Dev<br/>Restart: unless-stopped]
        F3[Populate-Laws Container<br/>Daemon Mode<br/>12-hour cycles]
        
        C4 --> F1
        D4 --> F2
        E4 --> F3
    end
    
    subgraph "Python Orchestrator Loop"
        G1[Every 12 Hours]
        G2[import_law.py<br/>goFromLastforward 100]
        G3[activate_politicians.py<br/>activate_all_groups]
        G4[Scrape cdep.ro<br/>Parse XML votes]
        
        F3 --> G1
        G1 --> G2
        G2 --> G4
        G4 --> G3
        G3 --> G1
    end
    
    subgraph "Kubernetes Infrastructure"
        H1[PostgreSQL Deployment<br/>3 Replicas<br/>Port: 5432]
        H2[PersistentVolume<br/>10Gi /d/postgre]
        H3[ConfigMap<br/>DB Credentials]
        H4[NodePort Service<br/>30080 HTTP, 30081 HTTPS]
        
        H3 --> H1
        H2 --> H1
        H1 --> H4
    end
    
    F1 -.->|PostgreSQL Protocol| H1
    F3 -.->|psycopg2| H1
    F2 -.->|HTTP API Calls| F1
    
    style A fill:#e1f5ff
    style F1 fill:#c8e6c9
    style F2 fill:#c8e6c9
    style F3 fill:#c8e6c9
    style H1 fill:#ffccbc
    style G1 fill:#fff9c4
```

### Application Architecture

The application follows a multi-tier architecture with a React frontend communicating with a .NET 8 REST API that implements a service layer pattern for business logic and Entity Framework Core for data access to a PostgreSQL database. The backend exposes RESTful endpoints for politicians, parties, voting rounds, and individual votes, with comprehensive filtering capabilities documented via OpenAPI/Swagger specification. The frontend uses TanStack Router for navigation and React Query for state management, displaying vote breakdowns with Chart.js pie charts and enabling users to drill down from overall vote results to party-level and individual politician voting records. All data flows through the API layer, which queries the database via Entity Framework Core and returns JSON responses to the React application.

**Key Components:**
- **Frontend Stack**: React 19.1 + TypeScript 5.8 + Vite 7.0 build tool, TanStack Router for routing, TanStack Query for API state management with 5-minute cache
- **API Layer**: ASP.NET Core 8.0 controllers (`/api/Politicians`, `/api/Voting`, `/api/Party`) with Swagger/OpenAPI documentation, CORS middleware, dependency injection for services
- **Service Layer**: `PoliticianService`, `PartyService`, `VotingService`, `VotingRoundService` implementing business logic with filtering by date ranges, keywords, party affiliation, gender, work location
- **Data Access**: Entity Framework Core 8 with Npgsql provider, AppDBContext managing DbSets (Politicians, Parties, VotingRounds, Votes), custom Color value converter for hex storage
- **Database Schema**: PostgreSQL 14 with four core tables—Politicians (UUID, name, party, gender, location, active status, image URL), Parties (UUID, name, acronym, logo, color), VotingRounds (vote ID, title, description, date), Votes (linking politician, round, position: Yes/No/Abstain/Absent)
- **API Endpoints**: GET `/api/Voting/getAllRounds` (with date/keyword filters), `/api/Voting/GetResultForVote` (votes by round), `/api/Politicians/getAllPoliticians` (with party/gender/location filters), `/api/Party/all`
- **Frontend Routes**: `/` (redirects to latest round) → `/round/:roundId` (vote breakdown) → `/section/:sectionId` (party breakdown) → `/party/:partyId` (individual voters)
- **Data Flow**: User interaction → React components → TanStack Query → API HTTP calls (Vite proxy) → ASP.NET Core controllers → Service layer (LINQ queries) → Entity Framework Core → PostgreSQL database → JSON response → Frontend visualization
- **Component Architecture**: React functional components with hooks, Material UI design system, Chart.js for data visualization, TypeScript for type safety across API contracts

```mermaid
graph TB
    subgraph "Frontend Tier - React Application"
        A1[User Browser]
        A2[React 19.1 + TypeScript 5.8]
        A3[TanStack Router<br/>Routes: /, /round/:id, /section/:id, /party/:id]
        A4[TanStack Query<br/>5-min cache, API state]
        A5[Components<br/>PieChart, RoundsList, PoliticiansList<br/>Material UI Design System]
        A6[API Utils<br/>rounds.ts, politicians.ts]
        
        A1 --> A2
        A2 --> A3
        A3 --> A4
        A4 --> A5
        A4 --> A6
    end
    
    subgraph "Vite Build System"
        B1[Vite 7.0 Dev Server<br/>Port 5173]
        B2[Proxy: /api → backend.democratiacontinua.eu]
        
        A6 --> B1
        B1 --> B2
    end
    
    subgraph "API Tier - ASP.NET Core 8.0"
        C1[ASP.NET Core Web API<br/>Ports 8080 HTTP / 8081 HTTPS]
        C2[CORS Middleware]
        C3[PoliticiansController<br/>/api/Politicians/*<br/>GetById, GetByName, GetAll]
        C4[VotingController<br/>/api/Voting/*<br/>GetAllRounds, GetResultForVote]
        C5[PartyController<br/>/api/Party/*<br/>GetAll, GetById, Query]
        C6[Swagger/OpenAPI Docs<br/>SwaggerAPI.json]
        
        B2 -->|HTTP/JSON| C1
        C1 --> C2
        C2 --> C3
        C2 --> C4
        C2 --> C5
        C1 --> C6
    end
    
    subgraph "Service Layer - Business Logic"
        D1[PoliticianService<br/>IPoliticianService<br/>Filter by party/gender/location/active]
        D2[VotingService<br/>IVotingService<br/>Get votes by round/party<br/>CRUD operations]
        D3[VotingRoundService<br/>IVotingRoundService<br/>Filter by date range/keywords]
        D4[PartyService<br/>IPartyService<br/>Party CRUD operations]
        
        C3 --> D1
        C4 --> D2
        C4 --> D3
        C5 --> D4
    end
    
    subgraph "Data Access Layer - Entity Framework Core 8"
        E1[AppDBContext<br/>Npgsql Provider]
        E2[DbSet<Politicians>]
        E3[DbSet<Parties>]
        E4[DbSet<VotingRounds>]
        E5[DbSet<Votes>]
        E6[Value Converters<br/>Color ↔ HTML Hex]
        
        D1 --> E1
        D2 --> E1
        D3 --> E1
        D4 --> E1
        E1 --> E2
        E1 --> E3
        E1 --> E4
        E1 --> E5
        E1 --> E6
    end
    
    subgraph "Database - PostgreSQL 14"
        F1[(Politicians Table<br/>UUID, Name, PartyId, Gender,<br/>WorkLocation, Active, ImageURL)]
        F2[(Parties Table<br/>UUID, Name, Acronym,<br/>LogoURL, Color, Active)]
        F3[(VotingRounds Table<br/>UUID, VoteId, Title,<br/>Description, VoteDate)]
        F4[(Votes Table<br/>UUID, PoliticianId, RoundId,<br/>PartyId, Position 0-3)]
        
        E2 -.->|Npgsql Protocol| F1
        E3 -.->|Npgsql Protocol| F2
        E4 -.->|Npgsql Protocol| F3
        E5 -.->|Npgsql Protocol| F4
        
        F4 -->|Foreign Key| F1
        F4 -->|Foreign Key| F3
        F4 -->|Foreign Key| F2
        F1 -->|Foreign Key| F2
    end
    
    style A1 fill:#e3f2fd
    style A2 fill:#e3f2fd
    style C1 fill:#c8e6c9
    style D1 fill:#fff9c4
    style D2 fill:#fff9c4
    style D3 fill:#fff9c4
    style D4 fill:#fff9c4
    style E1 fill:#f0f4c3
    style F1 fill:#ffccbc
    style F2 fill:#ffccbc
    style F3 fill:#ffccbc
    style F4 fill:#ffccbc
```

## DataBase

It will use PostGRSQL on a kubernetes cluster to make it easy to move into the cloud and have a fast reliable data structure.
The porject will use EntityFramework to run the db. This feature also comes with a caching system that will be sufficient for our usecases.
In order to deploy changes done to model project to db isntance we have to run the follwoing commands in developer console in the folder of the DBManager (ParlimentMonitor.DataBaseConnector):
```
dotnet ef migrations add SomeName
dotnet ef database update
```
This will basically create a new commit and deploy it on the database. The commit can be seen in the intermediate file in the 'Migrations' direcotry inside the project. There the framework converts the changes to C# statements and than we run that on our DB after connecting to it.
Connection string will be saved in "Secrets.cs" file.
