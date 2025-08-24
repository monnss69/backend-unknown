CREATE TABLE IF NOT EXISTS components (
    id TEXT PRIMARY KEY,
    name TEXT UNIQUE NOT NULL,
    code TEXT NOT NULL,
    props_schema JSONB NOT NULL
);

