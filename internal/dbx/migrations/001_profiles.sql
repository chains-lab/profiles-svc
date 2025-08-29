-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE profiles (
    user_id     UUID PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    pseudonym   VARCHAR(128),
    description VARCHAR(255),
    avatar      TEXT,
    sex         VARCHAR(16),
    birth_date  TIMESTAMP,

    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    created_at  TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS biographies CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
