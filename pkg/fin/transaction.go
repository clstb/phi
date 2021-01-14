package fin

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
)

type Transaction struct {
	Id        string
	Date      string
	Entity    string
	Reference string
	Hash      string
	Postings  Postings
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (t Transaction) Sum() Sum {
	return t.Postings.Sum()
}

func (t Transaction) Balanced() bool {
	sum := t.Postings.Sum().ByCurrency()
	for _, amount := range sum {
		if !amount.IsZero() {
			return false
		}
	}

	return true
}

func TransactionFromPB(pb *pb.Transaction) (Transaction, error) {
	postings, err := PostingsFromPB(pb.Postings)
	if err != nil {
		return Transaction{}, fmt.Errorf("transaction: %w", err)
	}

	return Transaction{
		Id:        pb.Id,
		Date:      pb.Date,
		Entity:    pb.Entity,
		Reference: pb.Reference,
		Hash:      pb.Hash,
		Postings:  postings,
	}, nil
}

func (t Transaction) PB() (*pb.Transaction, error) {
	postings, err := t.Postings.PB()
	if err != nil {
		return nil, err
	}

	return &pb.Transaction{
		Id:        t.Id,
		Date:      t.Date,
		Entity:    t.Entity,
		Reference: t.Reference,
		Hash:      t.Hash,
		Postings:  postings,
	}, nil
}
