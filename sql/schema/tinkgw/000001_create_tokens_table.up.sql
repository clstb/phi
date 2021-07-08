CREATE TABLE IF NOT EXISTS tokens (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  access_token STRING NOT NULL,
  refresh_token STRING NOT NULL,
  token_type STRING NOT NULL,
  expires_at TIMESTAMP NOT NULL,
  scope STRING NOT NULL,
  "user" UUID NOT NULL
);
