CREATE TABLE IF NOT EXISTS postings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account UUID NOT NULL REFERENCES accounts(id),
  transaction UUID NOT NULL REFERENCES transactions(id),
  units STRING NOT NULL,
  cost STRING,
  price STRING
);
