package fin

import (
	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

type Posting struct {
	db.Posting
	Units Amount
	Cost  Amount
	Price Amount
}

func (p Posting) Weight() Amount {
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

func PostingFromDB(p db.Posting) (Posting, error) {
	posting := Posting{
		Posting: p,
	}

	units, err := AmountFromString(posting.UnitsStr)
	if err != nil {
		return Posting{}, err
	}
	cost, err := AmountFromString(posting.CostStr)
	if err != nil {
		return Posting{}, err
	}
	price, err := AmountFromString(posting.PriceStr)
	if err != nil {
		return Posting{}, err
	}

	posting.Units = units
	posting.Cost = cost
	posting.Price = price

	return posting, nil
}

func PostingFromPB(pb *pb.Posting) (Posting, error) {
	units, err := AmountFromString(pb.Units)
	if err != nil {
		return Posting{}, err
	}
	cost, err := AmountFromString(pb.Cost)
	if err != nil {
		return Posting{}, err
	}
	price, err := AmountFromString(pb.Price)
	if err != nil {
		return Posting{}, err
	}

	posting := Posting{}
	posting.ID = uuid.FromStringOrNil(pb.Id)
	posting.Account = uuid.FromStringOrNil(pb.Account)
	posting.Transaction = uuid.FromStringOrNil(pb.Transaction)
	posting.Units = units
	posting.Cost = cost
	posting.Price = price

	return posting, nil
}
