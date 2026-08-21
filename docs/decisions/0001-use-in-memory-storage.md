# 0001 — Use deterministic in-memory storage

Status: Accepted

## Context

The exercise explicitly excludes an external database and says database design is not under evaluation. The application still needs shared mutable data during a process lifetime and enough records to demonstrate search and pagination.

## Decision

Use a concurrency-safe in-memory adapter populated with 48 games and 278 deterministic reviews. Construct seed records through domain validation and protect all access with a read/write mutex.

## Consequences

- Startup and Docker operation require no data service or migration.
- Tests are fast and deterministic.
- Data resets on restart and cannot be shared by backend replicas.
- A durable adapter can later implement the existing ports without changing domain code or HTTP ownership.

