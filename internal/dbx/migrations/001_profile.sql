-- +migrate Up
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE user_profiles (
    user_id     UUID PRIMARY KEY,
    username    VARCHAR(32) NOT NULL UNIQUE,
    pseudonym   VARCHAR(128),
    description VARCHAR(1024),
    avatar      TEXT,
    official    BOOLEAN NOT NULL DEFAULT FALSE,
    updated_at  TIMESTAMP NOT NULL DEFAULT now(),
    created_at  TIMESTAMP NOT NULL DEFAULT now()
);

CREATE TYPE sex_enum AS ENUM ('female','male','other');
CREATE TABLE users_biographic (
    user_id UUID PRIMARY KEY REFERENCES user_profiles(user_id) ON DELETE CASCADE,

    sex              sex_enum,
    birthday         timestamp,
    citizenship      VARCHAR(255),
    nationality      VARCHAR(255),
    primary_language VARCHAR(255),

    sex_updated_at              timestamp,
    citizenship_updated_at      timestamp,
    nationality_updated_at      timestamp,
    primary_language_updated_at timestamp
)

CREATE TABLE users_residence (
    user_id UUID PRIMARY KEY REFERENCES user_profiles(user_id) ON DELETE CASCADE,

    country_residence VARCHAR(255),
    city_residence    VARCHAR(255),

    country_updated_at timestamp,
    city_updated_at    timestamp
)

CREATE TYPE degrees_enum AS ENUM ('no degree', 'middle', 'incomplete higher', 'higher', 'candidate/doctor of sciences')
CREATE TYPE job_industry_enum AS ENUM ('IT', 'trade', 'service sector', 'hard physical work')
CREATE TYPE income_range_enum AS ENUM ('-200', '200-500', '500-1000', '1000-2000', '2000-5000', '5000-10000', '10000-50000', '50000+')
CREATE TABLE users_job (
    user_id  UUID PRIMARY KEY REFERENCES user_profiles(user_id) ON DELETE CASCADE,

    degree degrees_enum,
    industry job_industry_enum,
    income income_range_enum,

    degree_updated_at    timestamp,
    job_scope_updated_at timestamp,
    income_updated_at    timestamp,
)

-- +migrate Down
DROP TABLE IF EXISTS users_job CASCADE;
DROP TABLE IF EXISTS users_residence CASCADE;
DROP TABLE IF EXISTS users_personalities_data CASCADE;
DROP TABLE IF EXISTS profiles CASCADE;
DROP TYPE IF EXISTS sex_enums;
DROP TYPE IF EXISTS degrees_enum;
DROP TYPE IF EXISTS job_scope_enum;
DROP TYPE IF EXISTS income_enum;
