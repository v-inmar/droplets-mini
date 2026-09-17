CREATE TABLE IF NOT EXISTS event_processed_model (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT event_processed_model_eventid_unique UNIQUE (event_id)
);