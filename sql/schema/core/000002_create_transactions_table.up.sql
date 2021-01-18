CREATE TABLE IF NOT EXISTS transactions (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  date DATE NOT NULL,
  entity STRING NOT NULL,
  reference STRING NOT NULL,
  hash STRING NOT NULL
);
