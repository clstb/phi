package fin_test

import (
	"testing"

	"github.com/clstb/phi/pkg/fin"
	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestAmountFromString(t *testing.T) {
	is := is.New(t)

	type test struct {
		do    func() (fin.Amount, error)
		check func(fin.Amount, error)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}

	add(test{
		do: func() (fin.Amount, error) {
			return fin.AmountFromString("0.1 EUR")
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.Currency, "EUR")
			is.Equal(a.Decimal, decimal.NewFromFloat(0.1))
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fin.AmountFromString("0,1 EUR", fin.AmountEU)
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.Currency, "EUR")
			is.Equal(a.Decimal, decimal.NewFromFloat(0.1))
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fin.AmountFromString("")
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.Currency, "")
			is.Equal(a.Decimal, decimal.Zero)
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fin.AmountFromString("1")
		},
		check: func(a fin.Amount, e error) {
			is.True(e != nil)
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fin.AmountFromString("characters")
		},
		check: func(a fin.Amount, e error) {
			is.True(e != nil)
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fin.AmountFromString("characters EUR")
		},
		check: func(a fin.Amount, e error) {
			is.True(e != nil)
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}
