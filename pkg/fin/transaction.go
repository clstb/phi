package fin

import (
	"database/sql"
	"time"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
	"github.com/shopspring/decimal"
)

// Transaction is composed of a fixed date, an entity executing,
// a reference describing the entities intention, a hash as origin identifier.
// It should always be balanced as defined in double entry bookkeeping.
// A transaction defines movement of currency(ies) over multiple accounts.
type Transaction struct {
	db.Transaction
	Units Amount
	Cost  Amount
	Price Amount
}

// TransactionFromPB creates a new transaction from it's protobuf representation.
func TransactionFromPB(t *pb.Transaction) (Transaction, error) {
	var (
		id   uuid.UUID
		user uuid.UUID
		from uuid.UUID
		to   uuid.UUID
		err  error
	)

	if t.Id == "" {
		id = uuid.Nil
	} else {
		id, err = uuid.FromString(t.Id)
		if err != nil {
			return Transaction{}, err
		}
	}

	date, err := time.Parse("2006-01-02", t.Date)
	if err != nil {
		return Transaction{}, err
	}

	user, err = uuid.FromString(t.User)
	if err != nil {
		return Transaction{}, err
	}

	from, err = uuid.FromString(t.From)
	if err != nil {
		return Transaction{}, err
	}

	to, err = uuid.FromString(t.To)
	if err != nil {
		return Transaction{}, err
	}

	units, err := decimal.NewFromString(t.Units)
	if err != nil {
		return Transaction{}, err
	}

	cost, err := decimal.NewFromString(t.Cost)
	if err != nil {
		return Transaction{}, err
	}

	price, err := decimal.NewFromString(t.Price)
	if err != nil {
		return Transaction{}, err
	}

	transaction := Transaction{}
	transaction.ID = id
	transaction.Date = date
	transaction.Entity = t.Entity
	if t.Reference != "" {
		transaction.Reference = sql.NullString{
			String: t.Reference,
			Valid:  true,
		}
	}
	transaction.User = user
	transaction.From = from
	transaction.To = to
	transaction.Units = Amount{Decimal: units, Currency: t.UnitsCur}
	transaction.Cost = Amount{Decimal: cost, Currency: t.CostCur}
	transaction.Price = Amount{Decimal: price, Currency: t.PriceCur}
	if t.TinkId != "" {
		transaction.TinkID = sql.NullString{
			String: t.TinkId,
			Valid:  true,
		}
	}
	transaction.Debit = t.Debit

	return transaction, nil
}

// TransactionFromDB creates a new transaction from it's database representation.
func TransactionFromDB(t db.Transaction) Transaction {
	return Transaction{
		Transaction: t,
		Units:       Amount{Decimal: t.Units, Currency: t.Unitscur},
		Cost:        Amount{Decimal: t.Cost, Currency: t.Costcur},
		Price:       Amount{Decimal: t.Price, Currency: t.Pricecur},
	}
}

// PB marshalls the transaction into it's protobuf representation.
func (t Transaction) PB() *pb.Transaction {
	return &pb.Transaction{
		Id:        t.ID.String(),
		Date:      t.Date.Format("2006-01-02"),
		Entity:    t.Entity,
		Reference: t.Reference.String,
		User:      t.User.String(),
		From:      t.From.String(),
		To:        t.To.String(),
		Units:     t.Units.Decimal.String(),
		UnitsCur:  t.Unitscur,
		Cost:      t.Cost.Decimal.String(),
		CostCur:   t.Costcur,
		Price:     t.Price.Decimal.String(),
		PriceCur:  t.Pricecur,
		TinkId:    t.TinkID.String,
		Debit:     t.Debit,
	}
}

// Weight calculates a currency correct balancing value from units, cost and
// price of this transaction. It only returns positive values. The caller is
// supposed to used ".Neg()"" depending on "from" and "to" fields of this
// transaction.
func (t Transaction) Weight() (Amount, error) {
	if !t.Cost.IsZero() {
		t.Units.Currency = t.Cost.Currency
		return t.Units.Abs().Mul(t.Cost.Abs())
	}
	if !t.Price.IsZero() {
		t.Units.Currency = t.Cost.Currency
		return t.Units.Abs().Mul(t.Price.Abs())
	}

	return t.Units.Abs(), nil
}
