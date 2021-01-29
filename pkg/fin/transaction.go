package fin

import (
	"fmt"
	"time"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

// Transaction is composed of
// * a fixed date
// * an entity executing it
// * a reference describing the entities intention
// * a hash as origin identifier
// It should always be balanced as defined in double entry bookkeeping.
// A transaction defines movement of currency(ies) over multiple accounts.
type Transaction struct {
	db.Transaction
	Postings Postings
}

func NewTransaction(
	t db.Transaction,
	postings Postings,
) Transaction {
	return Transaction{t, postings}
}

// Balanced returns an error if the transaction is unbalanced or nil otherwise.
func (t Transaction) Balanced() error {
	var amounts Amounts
	sum, err := t.Postings.Sum()
	if err != nil {
		return err
	}
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

// TransactionFromDB creates a new transaction from it's database representation.
func TransactionFromDB(t db.Transaction) Transaction {
	return Transaction{
		Transaction: t,
	}
}

// TransactionFromPB creates a new transaction from it's protobuf representation.
func TransactionFromPB(t *pb.Transaction) (Transaction, error) {
	postings, err := PostingsFromPB(&pb.Postings{
		Data: t.Postings,
	})
	if err != nil {
		return Transaction{}, fmt.Errorf("transaction: %w", err)
	}

	date, err := time.Parse("2006-01-02", t.Date)
	if err != nil {
		return Transaction{}, err
	}

	transaction := db.Transaction{
		ID:        uuid.FromStringOrNil(t.Id),
		Date:      date,
		Entity:    t.Entity,
		Reference: t.Reference,
		Hash:      t.Hash,
	}

	return NewTransaction(transaction, postings), nil
}
