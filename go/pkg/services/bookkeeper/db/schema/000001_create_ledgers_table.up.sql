CREATE TABLE IF NOT EXISTS ledgers (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	user_id UUID NOT NULL,
	data BYTEA NOT NULL,
	dk BYTEA NOT NULL
);
