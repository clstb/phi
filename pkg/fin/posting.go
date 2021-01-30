package fin

import (
	db "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	"github.com/gofrs/uuid"
)

// Posting is part of an transaction and associated to an account.
// It defines a movement of currency on an account.
type Posting struct {
	db.Posting
	Units Amount
	Cost  Amount
	Price Amount
}

// PostingFromPB creates a new posting from it's protobuf representation.
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
	posting.ID = id
	posting.Account = account
	posting.Transaction = transaction
	posting.Units = units
	posting.Cost = cost
	posting.Price = price

	return posting, nil
}

// PostingFromDB creates a new posting from it's database representation.
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

// PB marshalls the posting into it's protobuf representation.
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

// Weight calculates the weight used for balancing this posting against postings.
func (p Posting) Weight() (Amount, error) {
	if !p.Cost.IsZero() {
		p.Units.Currency = p.Cost.Currency
		return p.Units.Mul(p.Cost)
	}
	if !p.Price.IsZero() {
		p.Units.Currency = p.Cost.Currency
		return p.Units.Mul(p.Price)
	}

	return p.Units, nil
}
