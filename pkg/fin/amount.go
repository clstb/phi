package fin

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/shopspring/decimal"
)

// Amount holds an decimal and currency.
type Amount struct {
	Decimal  decimal.Decimal
	Currency string
}

// StringRaw returns the amount decimal as string.
func (a Amount) StringRaw() string {
	return a.Decimal.String()
}

// String is like StringRaw but includes the currency.
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

// ColorRaw returns the amount decimal in red when negative and green when positive.
// The color choice can be swapped with invert
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

// Color is like ColorRaw but includes the currency.
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

// IsZero return true when the amounts decimal is zero or false otherwise.
func (a Amount) IsZero() bool {
	return a.Decimal.IsZero()
}

// Abs returns a new amount with absolute decimal.
func (a Amount) Abs() Amount {
	amount := Amount{}
	amount.Decimal = a.Decimal.Abs()
	amount.Currency = a.Currency

	return amount
}

// Neg returns a new amount with negated decimal.
func (a Amount) Neg() Amount {
	amount := Amount{}
	amount.Decimal = a.Decimal.Neg()
	amount.Currency = a.Currency

	return amount
}

// Add adds two amounts. An error is returned when currencies don't match.
func (a Amount) Add(a2 Amount) (Amount, error) {
	if a.Currency != a2.Currency {
		return Amount{}, ErrMismatchedCurrency
	}

	amount := Amount{}
	amount.Decimal = a.Decimal.Add(a2.Decimal)
	amount.Currency = a.Currency

	return amount, nil
}

// Mul multiplies two amounts. An error is returned when currencies don't match.
func (a Amount) Mul(a2 Amount) (Amount, error) {
	if a.Currency != a2.Currency {
		return Amount{}, ErrMismatchedCurrency
	}

	amount := Amount{}
	amount.Decimal = a.Decimal.Mul(a2.Decimal)
	amount.Currency = a.Currency

	return amount, nil
}

// AmountRE defines valid string representations of an amount.
var AmountRE = regexp.MustCompile(`(-?\d*\.?\d*) ([A-Z]+)`)

// AmountFromString parses a string into an Amount.
// The input string must match AmountRE.
// Formatters can modify the input string before parsing.
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

// AmountFormatter modifies an amount string.
type AmountFormatter func(string) string

// AmountEU changes the input string from EU currency standard to US currency standard.
func AmountEU(s string) string {
	s = strings.ReplaceAll(s, ".", "")
	return strings.ReplaceAll(s, ",", ".")
}
