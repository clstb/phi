package fin

import (
	"fmt"

	"github.com/clstb/phi/pkg/pb"
)

type Posting struct {
	Id          string
	Account     string
	Transaction string
	Units       *Amount
	Cost        *Amount
	Price       *Amount
}

func NewPosting() *Posting {
	return &Posting{}
}

func (p *Posting) Weight() *Amount {
	if p.Cost != nil {
		return p.Units.Mul(p.Cost)
	}
	if p.Price != nil {
		return p.Units.Mul(p.Price)
	}

	return p.Units
}

func (p *Posting) FromPB(pb *pb.Posting) error {
	units := NewAmount()
	if err := units.FromPB(pb.Units); err != nil {
		return fmt.Errorf("units: %w", err)
	}
	p.Units = units

	if pb.Cost != nil {
		cost := NewAmount()
		if err := cost.FromPB(pb.Cost); err != nil {
			return fmt.Errorf("cost: %w", err)
		}
		p.Cost = cost
	}

	if pb.Price != nil {
		price := NewAmount()
		if err := price.FromPB(pb.Price); err != nil {
			return fmt.Errorf("price: %w", err)

		}
		p.Price = price
	}

	p.Id = pb.Id
	p.Account = pb.Account
	p.Transaction = pb.Transaction

	return nil
}

func (p *Posting) PB() (*pb.Posting, error) {
	pb := &pb.Posting{
		Id:          p.Id,
		Account:     p.Account,
		Transaction: p.Transaction,
		Units:       p.Units.PB(),
	}

	if p.Cost != nil {
		pb.Cost = p.Cost.PB()
	}

	if p.Price != nil {
		pb.Price = p.Price.PB()
	}

	return pb, nil
}
