CREATE TABLE IF NOT EXISTS history_model (
    id BIGSERIAL PRIMARY KEY,
    event_pid BIGINT NOT NULL UNIQUE,
    event_happened TEXT NOT NULL,
    event_happened_at TIMESTAMPTZ NOT NULL,
    task_id BIGINT NOT NULL,
    task_value TEXT NOT NULL,
    task_pid BIGINT NOT NULL,
    task_completed BOOLEAN NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);


