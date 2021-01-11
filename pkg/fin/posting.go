package fin

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
)

type Posting struct {
	Id          string
	Account     string
	Transaction string
	Units       Amount
	Cost        Amount
	Price       Amount
}

func (p *Posting) Weight() Amount {
	if !p.Cost.IsZero() {
		return p.Units.Mul(p.Cost)
	}
	if !p.Price.IsZero() {
		return p.Units.Mul(p.Price)
	}

	return p.Units
}

func PostingFromPB(pb *pb.Posting) (Posting, error) {
	units, err := AmountFromString(pb.Units)
	if err != nil {
		return Posting{}, fmt.Errorf("units: %w", err)
	}

	cost, err := AmountFromString(pb.Cost)
	if err != nil {
		return Posting{}, fmt.Errorf("cost: %w", err)
	}

	price, err := AmountFromString(pb.Price)
	if err != nil {
		return Posting{}, fmt.Errorf("price: %w", err)
	}

	return Posting{
		Id:          pb.Id,
		Account:     pb.Account,
		Transaction: pb.Transaction,
		Units:       units,
		Cost:        cost,
		Price:       price,
	}, nil
}

func (p Posting) PB() (*pb.Posting, error) {
	pb := &pb.Posting{
		Id:          p.Id,
		Account:     p.Account,
		Transaction: p.Transaction,
		Units:       p.Units.String(),
	}

	if !p.Cost.IsZero() {
		pb.Cost = p.Cost.String()
	}

	if !p.Price.IsZero() {
		pb.Price = p.Price.String()
	}

	return pb, nil
}
