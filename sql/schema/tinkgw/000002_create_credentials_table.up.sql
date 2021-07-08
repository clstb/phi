CREATE TABLE IF NOT EXISTS credentials (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  "user" UUID NOT NULL,
  credentials_id STRING NOT NULL
)
