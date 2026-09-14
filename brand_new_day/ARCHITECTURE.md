# brand_new_day — Architecture

Fresh rewrite of the Continuous-Democracies backend, isolated from every other stack in this repo (`Backend/`, `new-backend/`). Read-only Go monolith, MongoDB storage, contract-first OpenAPI.

## Why

- `Backend/` (C#/EF-Core/Postgres, abandoned Python/FastAPI spike) and `new-backend/` (Go+MongoDB, currently deployed) coexist today. Neither models a "law" as anything more than a single atomic vote event (`idv`) with a title/description filled in from whichever single PDF the scraper happened to find.
- A cached cdep.ro page (`Backend/ParliamentMonitor/ContinousDemocracyAPI/api.html`) shows cdep.ro itself groups multiple `idv` vote events under one legislative-project id (`idp`) — that grouping is never persisted anywhere in this repo today.
- This rewrite introduces that grouping as a first-class, versioned concept (`LawBucket` + `Normative`), while keeping the parts of the old schema (Party, Politician) that already work.

## Database: MongoDB, not CouchDB

CouchDB's native `_rev` MVCC versioning was considered (the schema needs versioned law text) but rejected: `_rev` targets replication conflict resolution, not durable audit history — compaction can purge old revisions, so durable versioning still needs our own `id`/`version` fields regardless of database. Switching would also cost the aggregation-pipeline join model (`$lookup`/`$unwind`) this app leans on constantly, and Go driver support for CouchDB is a community library, not the official driver. **MongoDB, with versioning handled explicitly via plain `id`+`version` fields and compound `_id`s.**

## Data model

- **Party**: `id, name, acronym, logoUrl, color(hex), active`
- **Politician**: `id, name, gender(int), imageUrl, partyId, active, workLocation(int)` — 330 at steady state
- **LawBucket**: `_id (int, = cdep's idp), version(int), plNumber, title, description, initiationDate, status` — the real-world law/bill
- **Normative** (separate collection, 1:many child of LawBucket): one votable unit — an article, amendment, referenced act, or the whole-bill final-adoption text (`type: article | amendment | referencedAct | wholeBillFinal`). Versioned via plain `id` + `version` fields; `_id` is the compound `{id, version}` — Mongo's native pattern for versioned documents, no string-concatenation tricks. Fields: `_id:{id,version}, lawBucketId, type, label, text, effectiveDate`. A law's text changing after a vote means a new immutable document at `version+1`; old `VotingRound`s keep pointing at the exact version they voted on.
  - Separate collection rather than an embedded `LawBucket.normatives[]` array because `VotingRound` needs to reference one exact version of one normative directly. `GET` on a law always embeds its current normatives via `$lookup` — there are no standalone normative-only endpoints.
- **VotingRound**: `_id (int, = cdep's idv, no redundant synthetic UUID), title, description, voteDate, normativeId, normativeVersion (FK pair → Normative._id)` — one atomic vote event, resolving to "the law" via `Normative → LawBucket` whether it's an article-level vote or a whole-bill final vote.
  - `votes: [{politicianId, partyId, value}]` — embedded array, not a separate collection. Votes cannot outlive their round and the array is bounded (≤330), so this is the idiomatic Mongo shape; a multikey index on `votes.politicianId` serves "all votes by this politician." `partyId` is snapshotted at vote time so a later party switch doesn't rewrite history.

Indexes: `voting_rounds` — multikey on `votes.politicianId`, text on `{title,description}`, index on `(normativeId,normativeVersion)`. `normatives` — index on `lawBucketId`. `parties` — sparse unique on `acronym`.

## API surface (read-only)

```
GET /health
GET /parties                              GET /parties/{partyId}
GET /politicians                          GET /politicians/{politicianId}
GET /politicians/{politicianId}/votes     — aggregation over voting_rounds.votes, enriched with round+law
GET /votingRounds                         GET /votingRounds/{roundId}
GET /votingRounds/{roundId}/votes         — reads embedded votes[], hydrated with politician+party
GET /lawBuckets                           GET /lawBuckets/{lawBucketId}   — embeds current normatives
GET /lawBuckets/{lawBucketId}/votingRounds
```

Contract-first: `api/openapi.yaml` is the source of truth; `oapi-codegen` generates `internal/generated/{types,server,spec}.go`; `Controller` implements the generated `ServerInterface`. Verbose operationIds throughout (`getVotesByPolitician`, not `getVotes`).

## Layers

`Controller` (one, implements generated interface) → `Service` (`PoliticianService`, `VotingService`, `LawService`) → `Repository` (one per collection: party, politician, lawbucket, normative, votinground — no separate vote repository, votes live on votinground) → MongoDB.

## Reused from `new-backend/api/`

- Mongo connection/config pattern: `internal/db/mongo.go`, `internal/config/config.go`
- Aggregation-pipeline joins: `internal/repository/*.go`
- Integration test harness: `tests/integration/` (testcontainers-go + JSON fixtures)
- Verbose naming / explicit receivers convention (no `h`, `r`, single-letter vars)

## Snapshot/startup

`docker compose up` (mongo:7 + api, bind-mounted `./data/mongo`) is the entire startup. Snapshot via `mongodump --archive=snapshot.gz --gzip` (portable, one file); tarring the bind-mounted directory is a faster local fallback.

## Build process for this session

Built via a coordinator (Opus) dispatching three parallel workers (Sonnet) — `mongo-setup`, `services`, `openapi` — plus two test-writer agents (Haiku) for unit and integration tests. See `.claude/agents/` for role definitions.

Full detail and rationale for every decision above: `/Users/bogdanioan.boboc/.claude/plans/ok-i-want-to-steady-boot.md`.
