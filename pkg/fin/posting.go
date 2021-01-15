package fin

import (
	"github.com/clstb/phi/pkg/db"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

type Posting struct {
	db.Posting
}

func NewPosting(p db.Posting) Posting {
	return Posting{p}
}

func (p Posting) Weight() db.Amount {
	if !p.Cost.IsZero() {
		return p.Units.Mul(p.Cost)
	}
	if !p.Price.IsZero() {
		return p.Units.Mul(p.Price)
	}

	return p.Units
}

func (p Posting) PB() *pb.Posting {
	return &pb.Posting{
		Id:          p.ID.String(),
		Account:     p.Account.String(),
		Transaction: p.Transaction.String(),
		Units:       p.Units.String(),
		Cost:        p.Cost.String(),
		Price:       p.Price.String(),
	}
}

func PostingFromPB(pb *pb.Posting) (Posting, error) {
	id, err := uuid.FromString(pb.Id)
	if err != nil {
		return Posting{}, err
	}
	account, err := uuid.FromString(pb.Account)
	if err != nil {
		return Posting{}, err
	}
	transaction, err := uuid.FromString(pb.Transaction)
	if err != nil {
		return Posting{}, err
	}
	units, err := db.AmountFromString(pb.Units)
	if err != nil {
		return Posting{}, err
	}
	cost, err := db.AmountFromString(pb.Cost)
	if err != nil {
		return Posting{}, err
	}
	price, err := db.AmountFromString(pb.Price)
	if err != nil {
		return Posting{}, err
	}

	posting := db.Posting{
		ID:          id,
		Account:     account,
		Transaction: transaction,
		Units:       units,
		Cost:        cost,
		Price:       price,
	}

	return NewPosting(posting), nil
}
