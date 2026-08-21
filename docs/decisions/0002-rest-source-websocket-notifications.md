# 0002 — Keep REST authoritative and WebSocket event-only

Status: Accepted

## Context

New reviews must appear immediately for anyone viewing a game. The exercise also requires frontend and backend communication through REST.

## Decision

Create and read reviews through REST. After a review is validated and stored, broadcast a `review.created` WebSocket event. Connected clients use the event to invalidate and refresh their REST cache. Reconnects do not replay events.

## Consequences

- REST contracts remain independently usable and testable.
- Other open browsers receive low-latency updates without polling.
- Temporary WebSocket failure does not prevent normal reads or writes.
- Multi-instance deployment would require shared pub/sub and durable event coordination.

