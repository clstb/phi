CREATE TABLE IF NOT EXISTS postings (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  account UUID NOT NULL REFERENCES accounts(id),
  transaction UUID NOT NULL REFERENCES transactions(id),
  units STRING NOT NULL,
  units_cur STRING NOT NULL,
  cost STRING,
  cost_cur STRING,
  price STRING,
  price_cur STRING
);
CREATE VIEW postings_joined (
  id,
  account,
  account_name,
  transaction,
  date,
  units,
  units_cur,
  cost,
  cost_cur,
  price,
  price_cur
) AS
SELECT
  postings.id,
  postings.account,
  accounts.name,
  postings.transaction,
  transactions.date,
  postings.units,
  postings.units_cur,
  postings.cost,
  postings.cost_cur,
  postings.price,
  postings.price_cur
FROM
  postings
  JOIN accounts ON accounts.id = postings.account
  JOIN transactions ON transactions.id = postings.transaction;
