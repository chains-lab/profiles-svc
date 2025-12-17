-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE profiles (
    account_id  UUID PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    pseudonym   VARCHAR(128),
    description VARCHAR(255),
    avatar      TEXT,

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

CREATE TYPE outbox_event_status AS ENUM (
    'pending',
    'sent',
    'failed'
);

CREATE TABLE outbox_events (
    id            UUID  NOT NULL,
    topic         TEXT  NOT NULL,
    event_type    TEXT  NOT NULL,
    event_version INT   NOT NULL,
    key           TEXT  NOT NULL,
    payload       JSONB NOT NULL,

    status        outbox_event_status NOT NULL DEFAULT 'pending', -- pending | sent | failed
    attempts      INT         NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    sent_at       TIMESTAMPTZ NULL,

    PRIMARY KEY (id, topic)
);

CREATE TYPE inbox_event_status AS ENUM (
    'pending',
    'processed',
    'failed'
);

CREATE TABLE inbox_events (
    id            UUID  NOT NULL,
    topic         TEXT  NOT NULL,
    event_type    TEXT  NOT NULL,
    event_version INT   NOT NULL,
    key           TEXT  NOT NULL,
    payload       JSONB NOT NULL,

    status        inbox_event_status NOT NULL DEFAULT 'pending', -- pending | processed | failed
    attempts      INT         NOT NULL DEFAULT 0,
    next_retry_at TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),

    created_at    TIMESTAMPTZ NOT NULL DEFAULT (now() AT TIME ZONE 'UTC'),
    processed_at  TIMESTAMPTZ NULL,

    PRIMARY KEY (id, topic)
);


-- +migrate Down
DROP TABLE IF EXISTS profiles CASCADE;
DROP TABLE IF EXISTS outbox_events CASCADE;
DROP TABLE IF EXISTS inbox_events CASCADE;

DROP TYPE IF EXISTS outbox_event_status;
DROP TYPE IF EXISTS inbox_event_status;
