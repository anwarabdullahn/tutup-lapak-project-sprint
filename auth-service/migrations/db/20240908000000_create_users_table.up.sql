-- Create schema for auth service (optional but recommended)
CREATE SCHEMA IF NOT EXISTS auth;

-- Users table without FK to other services
CREATE TABLE IF NOT EXISTS auth.users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) NULL DEFAULT '',
    phone VARCHAR(50) NULL DEFAULT '',
    password VARCHAR(255) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- Helpful indexes
CREATE UNIQUE INDEX IF NOT EXISTS idx_users_email_unique ON auth.users (email) WHERE email <> '';
CREATE INDEX IF NOT EXISTS idx_users_phone ON auth.users (phone);

