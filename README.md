# Viger

Viger is a full-stack game review application. Visitors can explore a curated game catalog, search and filter it, read player reviews, publish their own review, and see new reviews appear live in another browser.

The project was built as a take-home exercise with an emphasis on maintainability, explicit boundaries, input safety, useful automated tests, and a reproducible Docker workflow.

## Features

- 48 deterministic example games and 278 deterministic example reviews.
- Search, genre and platform filters, sorting, and server-side pagination.
- Game details, rating aggregates, rating distribution, and paginated reviews.
- Review submission with backend-authoritative validation.
- Live `review.created` notifications over WebSocket.
- Responsive, accessible editorial UI without external artwork or runtime services.
- Structured logs, request IDs, health checks, and Prometheus-compatible metrics.
- In-memory, concurrency-safe storage as required by the exercise.

## Run with Docker

Requirements:

- Docker Engine with Docker Compose, or Docker Desktop.
- Ports `3000` and `8080` available on the loopback interface.

Build and start both components with one command:

```bash
docker compose up --build
```

Open:

- Web application: <http://localhost:3000>
- REST API: <http://localhost:8080/v1/games>
- Readiness check: <http://localhost:8080/health/ready>
- Metrics: <http://localhost:8080/metrics>

Stop the application:

```bash
docker compose down
```

The containers run as non-root users with read-only filesystems, dropped Linux capabilities, health checks, and explicit resource limits.

## Run from source

Requirements:

- Go 1.25
- Node.js 24.18
- Corepack and pnpm 11.16

Start the backend:

```bash
cd backend
go run ./cmd/server
```

In a second terminal, start the frontend:

```bash
cd frontend
corepack enable
pnpm install --frozen-lockfile
pnpm dev
```

Default local configuration:

| Variable | Default | Purpose |
| --- | --- | --- |
| `VIGER_HTTP_ADDRESS` | `:8080` | Backend listen address |
| `VIGER_ALLOWED_ORIGINS` | `http://localhost:3000` | Comma-separated browser and WebSocket origin allowlist |
| `VIGER_REVIEW_RATE_LIMIT` | `10` | Review writes allowed per client during the rate window |
| `VIGER_MAX_WS_CONNECTIONS` | `500` | Maximum concurrent WebSocket connections |
| `NEXT_PUBLIC_API_URL` | `http://localhost:8080` | REST URL embedded in the frontend build |
| `NEXT_PUBLIC_WS_URL` | `ws://localhost:8080/v1/ws` | WebSocket URL embedded in the frontend build |

`NEXT_PUBLIC_*` values are build-time values. Pass them as Docker build arguments when building for a different browser-visible host.

## Tests and quality checks

Run backend tests, including the race detector:

```bash
cd backend
go test -race ./...
go vet ./...
```

Run frontend checks:

```bash
cd frontend
pnpm lint
pnpm typecheck
pnpm test
pnpm build
```

Run browser tests against the Docker deployment:

```bash
docker compose up -d --build
cd frontend
pnpm exec playwright install --with-deps chromium
pnpm test:e2e
```

The E2E suite opens two isolated browser contexts and verifies that a review submitted in one appears in the other through WebSocket notification and REST cache refresh.

Convenience commands are also available from the repository root:

```bash
make check
make test
make build
make up
make down
```

## Architecture

The repository contains two independently built components:

- `frontend`: Next.js, React, TypeScript, TanStack Query, React Hook Form, and Zod.
- `backend`: Go standard-library HTTP server with grouped hexagonal architecture.

The backend groups each core model into domain, ports, and services. Inbound HTTP handlers and the outbound memory store depend on those ports; core code does not depend on transport or storage details. REST is the source of truth. WebSocket only announces that a successfully stored review was created, and clients refresh authoritative REST data after receiving the event.

See [Architecture](docs/architecture.md), [OpenAPI](docs/openapi.yaml), and the decision records under [`docs/decisions`](docs/decisions/) for details.

## API overview

| Method | Endpoint | Purpose |
| --- | --- | --- |
| `GET` | `/health/live` | Process liveness |
| `GET` | `/health/ready` | Catalog readiness |
| `GET` | `/metrics` | Prometheus-compatible metrics |
| `GET` | `/v1/games` | Search, filter, sort, and page games |
| `GET` | `/v1/games/{gameID}` | Read game details and rating summary |
| `GET` | `/v1/games/{gameID}/reviews` | Sort and page reviews |
| `POST` | `/v1/games/{gameID}/reviews` | Validate and create a review |
| `GET` | `/v1/ws` | Upgrade to the review event WebSocket |

## Key choices

- **Hexagonal backend:** keeps domain rules, use cases, storage, HTTP, and realtime concerns independently understandable and testable.
- **In-memory storage:** directly follows the exercise constraint. A mutex protects concurrent reads and writes, and seed data passes through the same domain constructors used by runtime data.
- **REST plus WebSocket:** mutations and reads remain predictable REST operations; WebSocket provides low-latency cross-browser notification without becoming a second source of truth.
- **Page-based pagination:** easy to understand and sufficient for the bounded in-memory catalog. A production data store could introduce cursor pagination without changing the HTTP handlers' responsibility.
- **Generated covers:** creates a distinctive interface without external requests, image licensing issues, or nondeterministic builds.
- **Small dependency surface:** libraries are used for specific responsibilities; the Go backend otherwise favors explicit standard-library wiring.

## Limitations and next steps

- Data resets to the deterministic seed whenever the backend restarts.
- The rate limiter and WebSocket hub are process-local. Multiple backend replicas would require shared persistence and pub/sub coordination.
- The client IP limiter intentionally uses the direct peer address. A deployment behind a proxy would need an explicit trusted-proxy policy before accepting forwarded headers.
- Reviews cannot be edited or deleted, games cannot be administered, and there are no accounts; these are deliberate exercise boundaries.
- With more time, I would add persistent storage, cursor pagination for a large catalog, moderation workflows, authenticated authorship, more browser/device coverage, and production telemetry export.

## License

This project is available under the [MIT License](LICENSE).

