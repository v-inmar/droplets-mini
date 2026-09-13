# droplets_mini

A small event-driven Todo system built with **Go**, **NATS JetStream**, and **PostgreSQL**. It demonstrates an event-driven microservice pattern: one service produces domain events, a message broker persists and delivers them, and a separate service consumes those events to build an audit history.

## Overview

- **App service** (`app/`) — writes todos into its own PostgreSQL database and publishes a `todo.created` fat event to NATS JetStream after each insert.
- **History service** (`history/`) — a durable JetStream consumer that reads todo events, deduplicates them by `event_id`, and stores a history row per processed event.
- **NATS JetStream** — acts as the event bus and persistent message store.
- **Two PostgreSQL databases** — one for the app's todos, one for the history service — plus a processed-events table used for exactly-once-style deduplication.

## Architecture

```
                 ┌─────────────────────────────┐
                 │         NATS JetStream       │
                 │      stream: "TODO"         │
                 │   subjects: todo.created,   │
                 │   todo.completed, ...       │
                 └───────────┬─────────────────┘
                             │
        publishes todo.created (fat event)
                             │
┌──────────────────┐         │         ┌──────────────────────┐
│   App service    │         │         │   History service    │
│  ┌────────────┐  │         │         │  ┌────────────────┐  │
│  │ todo_model │  │         └────────►│  │history_model   │  │
│  └────────────┘  │                   │  └────────────────┘  │
│   PostgreSQL     │                   │  │processed_event │  │
└──────────────────┘                   │  │    _model      │  │
                                       │  └────────────────┘  │
                                       │   PostgreSQL         │
                                       └──────────────────────┘
```

Flow:
1. The app service creates a todo in `todo_model` (in a transaction).
2. It publishes a `TodoCreatedEvent` (containing the full todo snapshot) to the `todo.created` subject on the `TODO` JetStream stream.
3. The history service's durable consumer receives the event.
4. It first inserts the `event_id` into `processed_event_model` (unique constraint) — a duplicate insert returns `ErrEventIDExist` and the message is simply acknowledged.
5. On a new event ID it inserts a row into `history_model` and commits — both inserts happen in the same transaction as the message acknowledgment.

## Project structure

```
droplets_mini/
├── app/                      # Todo producer service
│   ├── cmd/main.go           # entrypoint: connect DB/NATS, migrate, seed todos, publish events
│   ├── internal/
│   │   ├── models/todo.go    # Todo model
│   │   ├── repo/todo.go      # Todo repository (transactional insert)
│   │   └── service/          # DB service + JetStream event service
│   ├── migrations/           # SQL migrations (todo_model table)
│   ├── .env                  # local dev environment (host networking)
│   └── .docker.env           # container environment (enchanted via docker-compose)
│
├── history/                  # Todo history consumer service
│   ├── cmd/main.go           # entrypoint: durable consumer + history writes
│   ├── internal/
│   │   ├── models/           # History + ProcessedEvent models
│   │   └── repo/             # History + ProcessedEvent repositories
│   ├── migrations/           # SQL migrations (processed_event_model, history_model)
│   ├── .env                  # local dev environment (host networking)
│   └── .docker.env           # container environment
│
├── pkg/                      # Shared libraries
│   ├── database/             # Postgres connect-with-retry + golang-migrate runner
│   ├── events/               # Event models (e.g. TodoCreatedEvent)
│   ├── stream/               # Stream name constants (TODO, HISTORY)
│   ├── subjects/             # Subject constants (todo.created, todo.deleted, ...)
│   └── utilities/            # NATS connect-with-retry helper
│
├── docker/
│   ├── Dockerfile.app        # multi-stage build for the app service
│   └── Dockerfile.history    # multi-stage build for the history service
│
├── docker-compose.*.yml      # one compose file per component
├── Makefile                  # helper targets for networks/services
└── go.mod                    # module: droplets_mini (Go 1.26.3)
```

## Technologies

| Component       | Technology                          |
|-----------------|-------------------------------------|
| Language        | Go 1.26.3                           |
| Message broker  | NATS 2 with JetStream (`-js`)       |
| Datastore       | PostgreSQL 16 (two databases)       |
| Database access | `sqlx`                              |
| Migrations      | `golang-migrate` (file source)      |
| Env loading     | `godotenv`                          |
| Orchestration   | Docker Compose                      |

## Messaging contract

**Streams** (`pkg/stream`):
- `TODO` — created by the app service; configured with file storage and subject wildcard `todo.>`

**Subjects** (`pkg/subjects`):
- `todo.created`
- `todo.completed`
- `todo.updated`
- `todo.deleted`
- `history.requested` (defined but unused so far)

**Event** (`pkg/events`): `TodoCreatedEvent` is the fat event published by the app:

```json
{
  "event_id": 1324567890,
  "todo_id": 1,
  "todo_value": "Go to sleep",
  "todo_completed": false,
  "todo_created_at": "2026-09-13T..."
}
```

## Running the system

### Prerequisites

- Docker and Docker Compose
- Go 1.26+ (only if running locally instead of in containers)

### With Docker Compose (recommended)

The Makefile wires everything up. The services span three isolated networks (`event_network`, `app_network`, `history_network`), so start dependencies in this order:

```bash
# 1. create the shared networks
make network

# 2. start NATS with JetStream      -> terminal 1
make nats

# 3. start the app's PostgreSQL     -> terminal 2
make appdb

# 4. start the history's PostgreSQL -> terminal 3
make historydb

# 5. start the app service          -> terminal 4
make app

# 6. start the history service      -> terminal 5
make history
```

The app service seeds a few sample todos, inserts each into PostgreSQL, and publishes `todo.created` events to JetStream. The history service consumes those events, deduplicates them, and writes them into its own database.

Ports:
- NATS client: `4222`, monitoring UI: `8222`
- App PostgreSQL: `5432`
- History PostgreSQL: `5433` (host port remapped to avoid collision with the app DB)

### Tear down

```bash
make down    # stop all services and networks
make clean   # down + no-cache rebuild + remove networks
```

### Running locally (without Docker)

Each service loads its local `.env` file (host networking: `localhost` hostnames). You need a running NATS server (with JetStream enabled) and two PostgreSQL instances, then:

```bash
go run ./app/cmd
go run ./history/cmd
```

## Database schemas

**App DB** — `droplet_appdb`:

```sql
CREATE TABLE todo_model (
    id BIGSERIAL PRIMARY KEY,
    value TEXT NOT NULL,
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

**History DB** — `droplet_historydb`:

```sql
CREATE TABLE processed_event_model (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT processed_event_model_eventid_unique UNIQUE (event_id)
);

CREATE TABLE history_model (
    id BIGSERIAL PRIMARY KEY,
    event TEXT NOT NULL,
    todo_value TEXT NOT NULL,
    todo_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);
```

## Notes

- The history service currently implements full handling only for `todo.created`; `todo.completed`, `todo.deleted`, and `todo.updated` are acknowledged via a dummy handler.
- Event deduplication relies on the `UNIQUE (event_id)` constraint in `processed_event_model`, so replays or redeliveries (NATS `MaxDeliver: 5`) do not create duplicate history rows.