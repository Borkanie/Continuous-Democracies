# brand_new_day

From-scratch Go + MongoDB backend for Continuous-Democracies, isolated from `Backend/` and `new-backend/`. Read-only REST API, contract-first via OpenAPI, served with a Swagger UI. Full architecture rationale: [`ARCHITECTURE.md`](./ARCHITECTURE.md).

## Start it

```bash
cd brand_new_day/api
docker compose up -d mongo   # starts MongoDB only (see note below on the api container)
make seed                    # loads the mock dataset: 5 parties, 330 politicians, 10 laws, random votes
make run                     # starts the API on :8090
```

- Swagger UI: **http://localhost:8090/swagger/**
- Raw OpenAPI spec: http://localhost:8090/openapi.yaml
- Health check: http://localhost:8090/health

> **Note:** `docker compose up -d` (bringing up the `api` container too) currently fails to *build* on networks with TLS-intercepting proxies (`go mod download` hits an unknown certificate authority inside the container). Running Mongo via Docker and the API directly on the host (`make run`) is the verified path — see the comment in `api/Dockerfile` for the two fix options (vendor the module graph, or inject the intercepting proxy's root CA into the image).

Stop everything: `docker compose down` (from `brand_new_day/api`), then Ctrl-C the `make run` process.

## Snapshot / restore MongoDB

```bash
make snapshot   # mongodump --archive=snapshot.gz --gzip
make restore    # mongorestore --archive=snapshot.gz --gzip --drop
```

## Database structure (MongoDB, database `brandnewdaydb`)

| Collection | Key fields | Notes |
|---|---|---|
| `parties` | `id, name, acronym, logoUrl, color(hex), active` | sparse unique index on `acronym` |
| `politicians` | `id, name, gender, imageUrl, partyId, active, workLocation` | 330 at steady state |
| `lawBuckets` | `_id (int = cdep's idp), version, plNumber, title, description, initiationDate, status` | the real-world law/bill; no embedded normatives array |
| `normatives` | `_id: {id, version}, lawBucketId, type (article\|amendment\|referencedAct\|wholeBillFinal), label, text, effectiveDate` | separate collection, 1:many child of `lawBuckets`; versioned as immutable docs — a law's text changing writes a new `version`, old `votingRounds` keep pointing at the exact version they voted on; index on `lawBucketId` |
| `votingRounds` | `_id (int = cdep's idv), title, description, voteDate, normativeId, normativeVersion, votes: [{politicianId, partyId, value}]` | one atomic vote event; `votes` is embedded (bounded ≤330, never outlives its round) with a multikey index on `votes.politicianId`; also text index on `{title,description}` and index on `(normativeId,normativeVersion)` |

There is no separate `votes` collection — votes live embedded on their `votingRound`. There is no CRUD; every endpoint below is a `GET`.

## Endpoints

```
GET /health
GET /parties                              GET /parties/{partyId}
GET /politicians                          GET /politicians/{politicianId}
GET /politicians/{politicianId}/votes     — every vote by this politician, enriched with round + law
GET /votingRounds                         GET /votingRounds/{roundId}
GET /votingRounds/{roundId}/votes         — every vote in this round, enriched with politician + party
GET /lawBuckets                           GET /lawBuckets/{lawBucketId}   — embeds the law's current normatives
GET /lawBuckets/{lawBucketId}/votingRounds
```

## Run the tests

```bash
cd brand_new_day/api
make test-unit   # go test ./tests/unit/...
make test        # unit + testcontainers integration tests (Docker must be running)
```

## Agent roles used to build this

`.claude/agents/` holds the role definitions used to build and evolve this backend: `coordinator` (Opus) dispatches `mongo-setup`, `services`, `openapi` (Sonnet) in parallel, then `unit-test-writer` and `integration-test-writer` (Haiku).
