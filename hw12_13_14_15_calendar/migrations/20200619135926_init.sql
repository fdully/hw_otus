-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE events (
    id uuid primary key,
    subject text,
    start_time timestamptz NOT NULL,
    end_time timestamptz NOT NULL,
    description text,
    user_id VARCHAR(100) NOT NULL,
    notify_time int DEFAULT 900
);

-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE events;
