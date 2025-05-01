CREATE TABLE IF NOT EXISTS clients (
    id TEXT PRIMARY KEY,
    capacity INTEGER NOT NULL,
    rate_per_sec INTEGER NOT NULL
);