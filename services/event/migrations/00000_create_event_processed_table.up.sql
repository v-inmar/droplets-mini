CREATE TABLE IF NOT EXISTS event_processed_model (
    id BIGSERIAL PRIMARY KEY,
    event_pid TEXT NOT NULL UNIQUE,
    event_retry INT,
    event_processed BOOLEAN DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);