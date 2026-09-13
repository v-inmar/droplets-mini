CREATE TABLE IF NOT EXISTS processed_event_model (
    id BIGSERIAL PRIMARY KEY,
    event_id BIGINT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT processed_event_model_eventid_unique UNIQUE (event_id)
);