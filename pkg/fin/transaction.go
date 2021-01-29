package fin

import (
	"fmt"
	"time"

	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

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

func (t Transaction) Balanced() bool {
	var amounts Amounts
	for _, v := range t.Postings.Sum() {
		amounts = append(amounts, v...)
	}

	amounts = amounts.Sum()

	for _, amount := range amounts {
		if !amount.IsZero() {
			return false
		}
	}

	return true
}

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

func TransactionFromDB(t db.Transaction) Transaction {
	return Transaction{
		Transaction: t,
	}
}

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
