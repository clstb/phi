CREATE TABLE IF NOT EXISTS accounts_users (
	id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
	account UUID NOT NULL REFERENCES accounts(id) ON DELETE CASCADE,
	"user" UUID NOT NULL
)
