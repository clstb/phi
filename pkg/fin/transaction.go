package fin

import (
	"time"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

// Transaction is composed of a fixed date, an entity executing,
// a reference describing the entities intention, a hash as origin identifier.
// It should always be balanced as defined in double entry bookkeeping.
// A transaction defines movement of currency(ies) over multiple accounts.
type Transaction struct {
	db.Transaction
	Postings Postings
}

// TransactionFromPB creates a new transaction from it's protobuf representation.
func TransactionFromPB(t *pb.Transaction) (Transaction, error) {
	var (
		id  uuid.UUID
		err error
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

	postings, err := PostingsFromPB(&pb.Postings{
		Data: t.Postings,
	})
	if err != nil {
		return Transaction{}, err
	}

	transaction := Transaction{}
	transaction.ID = id
	transaction.Date = date
	transaction.Entity = t.Entity
	transaction.Reference = t.Reference
	transaction.Hash = t.Hash
	transaction.Postings = postings

	return transaction, nil
}

// TransactionFromDB creates a new transaction from it's database representation.
func TransactionFromDB(t db.Transaction) Transaction {
	return Transaction{Transaction: t}
}

// PB marshalls the transaction into it's protobuf representation.
func (t Transaction) PB() *pb.Transaction {
	return &pb.Transaction{
		Id:        t.ID.String(),
		Date:      t.Date.Format("2006-01-02"),
		Entity:    t.Entity,
		Reference: t.Reference,
		Hash:      t.Hash,
		Postings:  t.Postings.PB(),
	}
}

// Sum calculates the sum of this transaction grouped by account.
func (t Transaction) Sum() (map[string]Amounts, error) {
	sum := make(map[string]Amounts)
	for _, posting := range t.Postings {
		weight, err := posting.Weight()
		if err != nil {
			return nil, err
		}

		v, ok := sum[posting.Account.String()]
		if !ok {
			v = Amounts{weight}
		} else {
			v = append(v, weight)
		}
		sum[posting.Account.String()] = v
	}

	return sum, nil
}

// Balanced returns an error if the transaction is unbalanced or nil otherwise.
func (t Transaction) Balanced() error {
	sum, err := t.Sum()
	if err != nil {
		return err
	}

	var amounts Amounts
	for _, v := range sum {
		amounts = append(amounts, v...)
	}

	amounts, err = amounts.Sum()
	if err != nil {
		return err
	}

	for _, amount := range amounts {
		if !amount.IsZero() {
			return ErrUnbalanced
		}
	}

	return nil
}
