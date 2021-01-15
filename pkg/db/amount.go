package db

import (
	"database/sql/driver"
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Amount struct {
	Decimal  decimal.Decimal
	Currency string
}

func (a Amount) Value() (driver.Value, error) {
	return a.String(), nil
}

func (a *Amount) Scan(src interface{}) error {
	s, ok := src.(string)
	if !ok {
		return fmt.Errorf("incompatible type for Amount")
	}

	amount, err := AmountFromString(s)
	if err != nil {
		return fmt.Errorf("parsing string for Amount failed")
	}

	a.Decimal = amount.Decimal
	a.Currency = amount.Currency
	return nil
}

func (a Amount) StringRaw() string {
	return a.Decimal.String()
}

func (a Amount) String() string {
	return fmt.Sprintf(
		"%s %s",
		a.Decimal.String(),
		a.Currency,
	)
}

func (a Amount) IsZero() bool {
	return a.Decimal.IsZero()
}

func (a Amount) Abs() Amount {
	return Amount{
		Decimal:  a.Decimal.Abs(),
		Currency: a.Currency,
	}
}

func (a Amount) Neg() Amount {
	return Amount{
		Decimal:  a.Decimal.Neg(),
		Currency: a.Currency,
	}
}

func (a Amount) Add(amount Amount) Amount {
	value := a.Decimal.Add(amount.Decimal)

	return Amount{
		Decimal:  value,
		Currency: amount.Currency,
	}
}

func (a Amount) Mul(amount Amount) Amount {
	value := a.Decimal.Mul(amount.Decimal)

	return Amount{
		Decimal:  value,
		Currency: amount.Currency,
	}
}

func AmountFromString(s string, fmts ...AmountFormatter) (Amount, error) {
	if s == "" {
		return Amount{}, nil
	}

	blocks := strings.Split(s, " ")
	if len(blocks) != 2 {
		return Amount{}, fmt.Errorf("invalid format: expected \"<decimal> <currency>\"")
	}

	for _, fmt := range fmts {
		blocks[0] = fmt(blocks[0])
	}

	value, err := decimal.NewFromString(blocks[0])
	if err != nil {
		return Amount{}, fmt.Errorf("parsing decimal: %w", err)
	}

	return Amount{
		Decimal:  value,
		Currency: blocks[1],
	}, nil
}

type AmountFormatter func(string) string

func AmountEU(s string) string {
	s = strings.ReplaceAll(s, ".", "")
	return strings.ReplaceAll(s, ",", ".")
}
