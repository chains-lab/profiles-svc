-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TYPE sex_enum AS ENUM ('female','male','other');

CREATE TABLE profiles (
    id          UUID PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    alias       VARCHAR(64),
    description VARCHAR(1024),
    avatar      TEXT,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    created_at  TIMESTAMP NOT NULL DEFAULT now()
);

-- +migrate Down
DROP TABLE IF EXISTS profiles CASCADE;
DROP TYPE IF EXISTS sex_enums;