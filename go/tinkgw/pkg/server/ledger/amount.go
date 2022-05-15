package ledger

import (
	"fmt"

	"github.com/shopspring/decimal"
)

type Amount struct {
	Decimal  decimal.Decimal
	Currency string
}

func (a Amount) String() string {
	return fmt.Sprintf(
		"%s %s",
		a.Decimal.String(),
		a.Currency,
	)
}

func (a Amount) Neg() Amount {
	return Amount{
		Decimal:  a.Decimal.Neg(),
		Currency: a.Currency,
	}
}

func (a Amount) IsZero() bool {
	return a.Decimal.IsZero()
}

func (a Amount) IsNegative() bool {
	return a.Decimal.IsNegative()
}

func (a Amount) Add(v Amount) Amount {
	return Amount{
		Decimal:  a.Decimal.Add(v.Decimal),
		Currency: v.Currency,
	}
}

func (a Amount) Mul(v Amount) Amount {
	return Amount{
		Decimal:  a.Decimal.Mul(v.Decimal),
		Currency: v.Currency,
	}
}
