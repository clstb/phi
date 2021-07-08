CREATE TABLE IF NOT EXISTS accounts (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  name string NOT NULL,
  "user" UUID NOT NULL
);
