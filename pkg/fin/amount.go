package fin

import (
	"fmt"
	"strings"

	"github.com/shopspring/decimal"
)

type Amount struct {
	Value    decimal.Decimal
	Currency string
}

func (a Amount) StringRaw() string {
	return a.Value.String()
}

func (a Amount) String() string {
	return fmt.Sprintf(
		"%s %s",
		a.Value.String(),
		a.Currency,
	)
}

func (a Amount) IsZero() bool {
	return a.Value.IsZero()
}

func (a Amount) Abs() Amount {
	return Amount{
		Value:    a.Value.Abs(),
		Currency: a.Currency,
	}
}

func (a Amount) Neg() Amount {
	return Amount{
		Value:    a.Value.Neg(),
		Currency: a.Currency,
	}
}

func (a Amount) Add(amount Amount) Amount {
	value := a.Value.Add(amount.Value)

	return Amount{
		Value:    value,
		Currency: amount.Currency,
	}
}

func (a Amount) Mul(amount Amount) Amount {
	value := a.Value.Mul(amount.Value)

	return Amount{
		Value:    value,
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
		Value:    value,
		Currency: blocks[1],
	}, nil
}

type AmountFormatter func(string) string

func AmountEU(s string) string {
	s = strings.ReplaceAll(s, ".", "")
	return strings.ReplaceAll(s, ",", ".")
}
