-- +goose Up
CREATE TABLE IF NOT EXISTS users
(
    id              UUID PRIMARY KEY,
    email           VARCHAR(255) NOT NULL,
    password        VARCHAR(255) NOT NULL,
    is_confirmed    BOOL DEFAULT(false),
    confirm_code    VARCHAR(255),
    lang            VARCHAR(255) DEFAULT('en'),
    created_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE
);
CREATE UNIQUE INDEX IF NOT EXISTS users_email_uniq_idx ON "users"("email");

-- +goose Down
DROP INDEX IF EXISTS users_email_uniq_idx;
DROP TABLE IF EXISTS users;
