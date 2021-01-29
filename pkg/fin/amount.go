package fin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/shopspring/decimal"
)

type Amount struct {
	Decimal  decimal.Decimal
	Currency string
}

func (a Amount) StringRaw() string {
	return a.Decimal.String()
}

func (a Amount) String() string {
	if a.IsZero() {
		return ""
	}

	return fmt.Sprintf(
		"%s %s",
		a.StringRaw(),
		a.Currency,
	)
}

func (a Amount) ColorRaw(invert bool) string {
	red := false
	if a.Decimal.LessThan(decimal.Zero) {
		red = !red
	}
	if invert {
		red = !red
	}

	s := a.StringRaw()
	if red {
		return color.RedString(s)
	} else {
		return color.GreenString(s)
	}
}

func (a Amount) Color(invert bool) string {
	if a.IsZero() {
		return ""
	}

	return fmt.Sprintf(
		"%s %s",
		a.ColorRaw(invert),
		a.Currency,
	)
}

func (a Amount) IsZero() bool {
	return a.Decimal.IsZero()
}

func (a Amount) Abs() Amount {
	amount := Amount{}
	amount.Decimal = a.Decimal.Abs()
	amount.Currency = a.Currency

	return amount
}

func (a Amount) Neg() Amount {
	amount := Amount{}
	amount.Decimal = a.Decimal.Neg()
	amount.Currency = a.Currency

	return amount
}

func (a Amount) Add(a2 Amount) (Amount, error) {
	if a.Currency != a2.Currency {
		return Amount{}, ErrMismatchedCurrency
	}

	amount := Amount{}
	amount.Decimal = a.Decimal.Add(a2.Decimal)
	amount.Currency = a.Currency

	return amount, nil
}

func (a Amount) Mul(a2 Amount) (Amount, error) {
	if a.Currency != a2.Currency {
		return Amount{}, ErrMismatchedCurrency
	}

	amount := Amount{}
	amount.Decimal = a.Decimal.Mul(a2.Decimal)
	amount.Currency = a.Currency

	return amount, nil
}

var AmountRE = regexp.MustCompile(`(-?\d*\.?\d*) ([A-Z]+)`)

func AmountFromString(s string, fmts ...AmountFormatter) (Amount, error) {
	if s == "" {
		return Amount{
			Decimal: decimal.Zero,
		}, nil
	}

	for _, fmt := range fmts {
		s = fmt(s)
	}

	if !AmountRE.MatchString(s) {
		return Amount{}, fmt.Errorf(
			"invalid format %s: expected %s",
			s,
			AmountRE.String(),
		)
	}

	matches := AmountRE.FindAllStringSubmatch(s, -1)

	decimal, err := decimal.NewFromString(matches[0][1])
	if err != nil {
		return Amount{}, fmt.Errorf("parsing decimal: %w", err)
	}

	currency := matches[0][2]

	amount := Amount{}
	amount.Decimal = decimal
	amount.Currency = currency

	return amount, nil
}

type AmountFormatter func(string) string

func AmountEU(s string) string {
	s = strings.ReplaceAll(s, ".", "")
	return strings.ReplaceAll(s, ",", ".")
}
