-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS collections (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT ''
);

CREATE TABLE IF NOT EXISTS requests (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    url TEXT NOT NULL,
    method TEXT NOT NULL DEFAULT 'GET',
    headers TEXT DEFAULT '{}',
    body TEXT DEFAULT '',
    collection_id INTEGER REFERENCES collections(id) ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS history (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL,
    method TEXT NOT NULL,
    status_code INTEGER,
    duration_ms BIGINT,
    timestamp BIGINT NOT NULL
);

-- Insert default collection if not exists
INSERT INTO collections (id, name, description) 
VALUES (1, 'Default', 'Default collection')
ON CONFLICT (id) DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS history;
DROP TABLE IF EXISTS requests;
DROP TABLE IF EXISTS collections;
-- +goose StatementEnd
