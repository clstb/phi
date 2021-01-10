package fin

import (
	"fmt"
	"strings"

	"github.com/clstb/phi/pkg/pb"
	"github.com/shopspring/decimal"
)

type Amount struct {
	Value    decimal.Decimal
	Currency string
}

func NewAmount() *Amount {
	return &Amount{}
}

func (a *Amount) String() string {
	return a.Value.String()
}

func (a *Amount) IsZero() bool {
	return a.Value.IsZero()
}

func (a *Amount) Abs() *Amount {
	return &Amount{
		Value:    a.Value.Abs(),
		Currency: a.Currency,
	}
}

func (a *Amount) Neg() *Amount {
	return &Amount{
		Value:    a.Value.Neg(),
		Currency: a.Currency,
	}
}

func (a *Amount) Add(amount *Amount) *Amount {
	value := a.Value.Add(amount.Value)

	return &Amount{
		Value:    value,
		Currency: amount.Currency,
	}
}

func (a *Amount) Mul(amount *Amount) *Amount {
	value := a.Value.Mul(amount.Value)

	return &Amount{
		Value:    value,
		Currency: amount.Currency,
	}
}

func (a *Amount) FromString(s string, fmts ...AmountFormatter) error {
	if s == "" {
		return nil
	}

	for _, fmt := range fmts {
		s = fmt(s)
	}

	blocks := strings.Split(s, " ")
	if len(blocks) != 2 {
		return fmt.Errorf("invalid format: expected \"<decimal> <currency>\"")
	}

	value, err := decimal.NewFromString(blocks[0])
	if err != nil {
		return fmt.Errorf("parsing decimal: %w", err)
	}

	a.Value = value
	a.Currency = blocks[1]

	return nil
}

func (a *Amount) FromPB(pb *pb.Amount) error {
	value := decimal.Zero
	value, err := decimal.NewFromString(pb.Value)
	if err != nil {
		return fmt.Errorf("unmarshalling decimal from bytes: %w", err)
	}

	a.Value = value
	a.Currency = pb.Currency

	return nil
}

func (a *Amount) PB() *pb.Amount {
	return &pb.Amount{
		Value:    a.Value.String(),
		Currency: a.Currency,
	}
}

type AmountFormatter func(string) string

func AmountEU(s string) string {
	s = strings.ReplaceAll(s, ".", "")
	return strings.ReplaceAll(s, ",", ".")
}
