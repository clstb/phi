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
	Postings  *Postings
}

func NewTransaction() *Transaction {
	return &Transaction{}
}

func (t *Transaction) Balanced() bool {
	sum := t.Postings.Sum().ByCurrency()
	for _, amount := range sum {
		if !amount.IsZero() {
			return false
		}
	}

	return true
}

func (t *Transaction) FromPB(pb *pb.Transaction) error {
	postings := NewPostings()
	if err := postings.FromPB(pb.Postings); err != nil {
		return fmt.Errorf("postings: %w", err)
	}

	t.Id = pb.Id
	t.Date = pb.Date
	t.Entity = pb.Entity
	t.Reference = pb.Reference
	t.Postings = postings

	return nil
}

func (t *Transaction) PB() (*pb.Transaction, error) {
	postings, err := t.Postings.PB()
	if err != nil {
		return nil, err
	}

	return &pb.Transaction{
		Id:        t.Id,
		Date:      t.Date,
		Entity:    t.Entity,
		Reference: t.Reference,
		Postings:  postings,
	}, nil
}
