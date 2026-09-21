# droplets_mini

A microservices showcase built with **Go**, **Kafka**, and **PostgreSQL** following an **event driven architecture**. It demonstrates a small todo system where a REST API service produces domain events, Kafka moves them across the system, and a consumer service turns them into an queryable history.

This is a **showcase**, not a production-grade system. Each component is intentionally kept small so the design is easy to read and extend.

## Architecture

Services talk to each other only through **REST APIs** and **Kafka events** — no shared language or runtime coupling. Each service is **language independent** and can be rewritten in any other language (Go is used here by convenience, not constraint).

```
                        ┌────────────────────────────────────────┐
                        │              gateway :5500             │
                        │        (REST reverse proxy, chi)       │
                        └───────┬────────────────────┬───────────┘
                    /v1/todo ───┘                    └─── /v1/history
                                │                          │
                    ┌───────────▼─────────┐        ┌────────▼──────────┐
                    │     todo :5501      │        │   history :5502   │
                    │   (REST API, CRUD)  │        │   (REST API, GET) │
                    └───────────┬─────────┘        └────────┬──────────┘
                                │                          │
                       publish task.created /              reads
                       task.updated / task.deleted         history_db (Postgres)
                                │
                    ┌───────────▼──────────────────────────┐
                    │              Kafka                   │
                    │        topic: tasks / tasks.dlq      │
                    └───────────┬──────────────────────────┘
                                │ consume (group: task-service)
                    ┌───────────▼─────────┐
                    │     event (CLI)     │
                    │  dedup / retry / DLQ│
                    └───────────┬─────────┘
                                │ writes history
                                ▼
                     history_db (Postgres)
```

1. **todo** — REST CRUD on `/tasks`; writes to its own Postgres DB and publishes a fat `EventTask` to the Kafka `tasks` topic after each create/update/delete.
2. **event** — Kafka consumer. Deduplicates by unique `event_pid`, retries failed messages up to `MaxRetry=5`, then routes them to the `tasks.dlq` DLQ topic. Writes each processed event into the history Postgres DB.
3. **history** — REST `GET /items` over the built-up history.
4. **gateway** — REST reverse proxy exposing the todo and history services behind one entry point.

## Services

| Service   | Port | Type      | Role                                             | Depends on                             |
|-----------|------|-----------|--------------------------------------------------|----------------------------------------|
| `todo`    | 5501 | REST API  | CRUD tasks, publishes domain events to Kafka     | Postgres (`todo_db`), Kafka            |
| `event`   | —    | CLI       | Consumes events, dedup/retry/DLQ, writes history | Postgres (`event_db`, `history_db`), Kafka |
| `history` | 5502 | REST API  | Read-only history endpoint                       | Postgres (`history_db`)                |
| `gateway` | 5500 | REST proxy| Reverse-proxies `/v1/todo` and `/v1/history`     | todo, history                          |
| —         | 9092 | —         | Apache Kafka, topics `tasks` and `tasks.dlq`     | —                                      |
| —         | —    | —         | Three Postgres 16 databases (`todo_db`, `history_db`, `event_db`) | — |

## Repo layout

```
droplets/
├── docker-compose.storage.yml   # Kafka + three Postgres databases
├── Makefile                     # one-command orchestration
├── services/
│   ├── todo/                    # REST API producer  (own Go module)
│   ├── event/                   # Kafka consumer     (own Go module)
│   ├── history/                 # REST API reader    (own Go module)
│   └── gateway/                 # reverse proxy      (own Go module)
└── notes.txt
```

A **monorepo is used for convenience only**. Each `services/*` directory is its own Go module with its own `Dockerfile` and migrations, so any service can be split into a separate repository without refactoring.

## Tech stack

| Component       | Technology                       |
|-----------------|----------------------------------|
| Language        | Go 1.26.3                        |
| HTTP router     | `go-chi/chi`                     |
| Message broker  | Apache Kafka (`segmentio/kafka-go`) |
| Datastore       | PostgreSQL 16 (three databases)  |
| Database access | `sqlx`                           |
| Migrations      | `golang-migrate`                 |
| Env loading     | `godotenv`                       |
| Testing         | `go-sqlmock`                     |
| Orchestration   | Docker + Docker Compose          |

## Quickstart

**Prerequisites:** Docker and Docker Compose. Go 1.26+ only if running outside containers.

### Full stack — one terminal

```bash
make webup
```

This creates the Docker networks, starts storage (Kafka + Postgres), then builds and starts all four services.

### Per-component (multi-terminal)

```bash
make webup-network      # create the shared networks
make webup-storage      # Kafka + Postgres (terminal 1)
make webup-todo         # terminal 2
make webup-history      # terminal 3
make webup-event        # terminal 4
make webup-gateway      # terminal 5
```

### Manage

```bash
make webstatus   # status of all compose projects
make webdown     # stop everything
make webclean    # down + remove local images + networks
```

**Ports:** gateway `5500`, todo `5501`, history `5502`, Kafka `9092`. Only the gateway and Kafka are exposed on the host; the rest communicate over internal Docker networks.

Services are isolated across five external networks (`gateway_network`, `event_network`, `tododb_network`, `historydb_network`, `eventdb_network`) so each one only reaches what it needs.

### Running locally (without Docker)

Each service loads its local `.env` (host networking). With Kafka and Postgres running on `localhost`:

```bash
go run ./cmd/web     # todo, history, gateway
go run ./cmd/cli     # event
```

## Usage

```bash
# create a task
curl -X POST http://localhost:5500/v1/todo/tasks \
  -d '{"value":"write a readme"}' -H 'Content-Type: application/json'

# list tasks
curl http://localhost:5500/v1/todo/tasks

# read the event history
curl http://localhost:5500/v1/history/items
```

Create a task, then `curl /v1/history/items` shortly after — the `task.created` event should appear once, even if the message was redelivered.

## Event contract

```json
{
  "event_type": "task.created",
  "event_pid": 1767381204128391000,
  "event_created_at": "2026-09-21T...",
  "task_id": 1,
  "task_value": "write a readme",
  "task_pid": 1817643920171900000,
  "task_completed": false,
  "task_created_at": "2026-09-21T..."
}
```

- Event types: `task.created`, `task.updated`, `task.deleted` (fat events carrying the full task snapshot).
- Each event carries a unique `event_pid`; the consumer deduplicates on it (`event_processed_model` unique constraint).
- Failed messages retry up to `MaxRetry=5`, then go to the `tasks.dlq` topic.

## Limitations

- **Showcase, not production:** IDs are generated with `time.Now().UnixNano()`, so pid collision is possible at scale.
- **Outbox gap:** the todo service publishes to Kafka before committing its DB transaction, so a crash in between leaves an event without a DB row.
- **Offset-commit timing:** DLQ routing commits offsets before the dedup edge cases are fully resolved, which can duplicate under rare failure windows.

## Testing

```bash
cd services/todo && go test ./...
cd services/history && go test ./...
```