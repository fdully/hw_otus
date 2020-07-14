-- +goose Up
-- SQL in this section is executed when the migration is applied.

CREATE TABLE events (
    id uuid primary key,
    subject text,
    start_time timestamptz NOT NULL,
    end_time timestamptz NOT NULL,
    description text,
    owner_id VARCHAR(100) NOT NULL,
    notify_period bigint DEFAULT 900
);

COMMENT on column events.notify_period is 'notify period in seconds. 0 to disable notification';
COMMENT on column events.owner_id is 'owner of the event';


-- +goose Down
-- SQL in this section is executed when the migration is rolled back.
DROP TABLE events;
