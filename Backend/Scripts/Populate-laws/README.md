# Orchestrator — Temporal Diagram

This README documents the current runtime data flow for the orchestrator and scraping logic used to import politicians, laws and votes.

## Overview

- The `orchestrator.py` script runs hourly and coordinates three main actions:
  - `import_politicians.importPoliticians(...)` — imports politicians for the current year.
  - `import_law.goFromLastforward(n)` — imports laws (up to `n`) starting from the last processed item.
  - `activate_politicians.activate_all_groups()` — runs every 12 hours.

## Temporal / Sequence Diagram

Below is a small sequence diagram showing the order of events when a cycle runs.

```mermaid
sequenceDiagram
    participant Orch as Orchestrator (orchestrator.py)
    participant DB as Database (shared connection)
    participant ImportPol as import_politicians
    participant ImportLaw as import_law.goFromLastforward
    participant Scraper as Scraper / Downloader modules

    Orch->>DB: init_db() (on start)
    Orch->>ImportPol: importPoliticians(1..333, year)
    ImportPol->>Scraper: fetch politician pages/data
    Scraper->>DB: save politician records

    Orch->>ImportLaw: goFromLastforward(100)
    ImportLaw->>Scraper: request next law(s) since last processed
    Scraper->>ImportLaw: return law data + vote links
    ImportLaw->>Scraper: fetch votes for each law
    Scraper->>DB: save laws and votes

    Orch->>Orch: sleep(3600) (loop hourly)
    Orch->>ImportPol: (every hour loop continues)

    Note over Orch: every 12h: activate_all_groups() -> updates DB
    Orch->>DB: close_db() (on shutdown)
```


