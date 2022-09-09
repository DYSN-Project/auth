-- +goose Up
CREATE TABLE IF NOT EXISTS recovery_password
(
    id              UUID PRIMARY KEY,
    email           varchar(255) UNIQUE NOT NULL REFERENCES public.users (email) ON DELETE CASCADE,
    confirm_code    varchar(255),
    status          int NOT NULL DEFAULT 0,
    created_at      TIMESTAMP WITH TIME ZONE,
    updated_at      TIMESTAMP WITH TIME ZONE,
    deleted_at      TIMESTAMP WITH TIME ZONE
    );

CREATE INDEX IF NOT EXISTS recovery_idx ON "recovery_password"("email",status);

-- +goose Down
DROP TABLE IF EXISTS recovery_password
DROP INDEX IF EXISTS recovery_idx

