package fin

import (
	"fmt"

	"github.com/clstb/phi/pkg/db"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
	"github.com/golang/protobuf/ptypes"
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

func (t Transaction) Balanced(postings Postings) bool {
	var amounts db.Amounts
	for _, v := range postings.Sum() {
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

func (t Transaction) PB() (*pb.Transaction, error) {
	date, err := ptypes.TimestampProto(t.Date)
	if err != nil {

	}

	return &pb.Transaction{
		Id:        t.ID.String(),
		Date:      date,
		Entity:    t.Entity,
		Reference: t.Reference,
		Hash:      t.Hash,
		Postings:  t.Postings.PB(),
	}, nil
}

func TransactionFromPB(t *pb.Transaction) (Transaction, error) {
	postings, err := PostingsFromPB(&pb.Postings{
		Data: t.Postings,
	})
	if err != nil {
		return Transaction{}, fmt.Errorf("transaction: %w", err)
	}

	date, err := ptypes.Timestamp(t.Date)
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
