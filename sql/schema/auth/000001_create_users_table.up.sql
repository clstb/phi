CREATE TABLE IF NOT EXISTS users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	name string NOT NULL UNIQUE,
	password BYTEA NOT NULL
);
