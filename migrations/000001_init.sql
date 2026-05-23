-- Oplati: users (internal/storages/inmemory/oplati) + auth (internal/storages/inmemory/auth)

CREATE TABLE users (
    id      UUID    PRIMARY KEY,
    balance INTEGER NOT NULL DEFAULT 0,

    CONSTRAINT users_balance_non_negative CHECK (balance >= 0)
);

CREATE TABLE accounts (
    login         TEXT PRIMARY KEY,
    password_hash TEXT NOT NULL,
    user_id       UUID NOT NULL
);
