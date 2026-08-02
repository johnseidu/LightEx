CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,

    email VARCHAR(254) NOT NULL UNIQUE,

    password_hash TEXT NOT NULL,

    email_verified BOOLEAN NOT NULL DEFAULT FALSE,

    status VARCHAR(20) NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,

    updated_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX idx_users_email
ON users(email);