# Viger AI Agent Instructions

This file is the primary project-specific instruction source for every AI agent working in this repository.

## Rule priority

1. The current user request, when safe and authorized.
2. This `AGENTS.md` and accepted Viger documentation.
3. The AI workspace foundation in `/home/septi/.ai` when available.
4. Relevant capabilities, playbooks, knowledge, and general engineering practices.

Explain material conflicts instead of resolving them silently.

## Language

- Use English for source code, identifiers, UI copy, API responses, documentation, tests, GitHub metadata, and commit messages.
- Keep technical writing direct and understandable to an engineer reviewing a take-home exercise.

## Product scope

Viger is a full-stack game review application. A visitor can browse, search, filter, sort, and page through seeded games; inspect a game and its paginated reviews; and submit a review containing a reviewer name, rating, and text. New reviews are persisted for the lifetime of the backend process and announced to connected clients through WebSocket events.

The following are intentionally out of scope:

- authentication and user accounts;
- game administration;
- editing or deleting reviews;
- external databases, Redis, and message brokers;
- CQRS, event sourcing, distributed tracing infrastructure, and Kubernetes;
- generated generic repository frameworks.

Do not expand this scope without explicit user approval.

## Architecture

- Keep the frontend and backend as separate components in one monorepo.
- Organize the Go backend using grouped hexagonal architecture under `internal/core`, `internal/adapters`, and `internal/platform`, following the dependency direction demonstrated by Nestori.
- Keep `Game` and `Review` as separate domain models.
- Domain and application code must not depend on HTTP, WebSocket, or the in-memory adapter.
- REST is the source of truth. WebSocket is a notification channel for newly created reviews.
- Keep frontend HTTP access, runtime validation, query state, form state, and presentation responsibilities explicit and separate.
- Prefer explicit wiring and small interfaces over reflection or dependency-injection frameworks.

## Engineering standards

- Optimize for maintainability, correctness, security, and clarity.
- Keep abstractions purposeful and functions focused.
- Validate all untrusted input on the backend; frontend validation exists for user experience only.
- Bound request bodies, pagination, search input, WebSocket messages, connections, and write rates.
- Do not log review content, credentials, secrets, or unnecessary personal data.
- Preserve accessibility, responsive behavior, graceful failure states, and reduced-motion preferences.
- Keep seed generation deterministic and route seed records through domain validation.
- Every behavior change requires relevant tests and documentation updates.

## Verification

Before declaring work complete, run the repository-prescribed formatting, linting, type checking, unit, integration, race, build, Docker, and end-to-end checks. Verify the built application through Docker Compose and record honest limitations in the README.

## Git workflow

- Preserve user changes and never rewrite published history or force-push.
- Use small, logical, reviewable commits with concise English intent-based messages.
- Do not commit secrets, local configuration, build output, test artifacts, or unrelated changes.
- Review each diff and verification result before committing.

