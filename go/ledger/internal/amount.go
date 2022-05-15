package internal

import (
	"fmt"
)

func (a AmountType) String() string {
	return fmt.Sprintf(
		"%s %s",
		a.Decimal.String(),
		a.Currency,
	)
}

func (a AmountType) Neg() AmountType {
	return AmountType{
		Decimal:  a.Decimal.Neg(),
		Currency: a.Currency,
	}
}

func (a AmountType) IsZero() bool {
	return a.Decimal.IsZero()
}

func (a AmountType) IsNegative() bool {
	return a.Decimal.IsNegative()
}

func (a AmountType) Add(v AmountType) AmountType {
	return AmountType{
		Decimal:  a.Decimal.Add(v.Decimal),
		Currency: v.Currency,
	}
}

func (a AmountType) Mul(v AmountType) AmountType {
	return AmountType{
		Decimal:  a.Decimal.Mul(v.Decimal),
		Currency: v.Currency,
	}
}
