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

func TestAmountAdd(t *testing.T) {
	is := is.New(t)

	type test struct {
		do    func() (fin.Amount, error)
		check func(fin.Amount, error)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}
	fromString := func(s string) fin.Amount {
		amount, err := fin.AmountFromString(s)
		is.NoErr(err)
		return amount
	}

	add(test{
		do: func() (fin.Amount, error) {
			return fromString("1 EUR").Add(fromString("1 EUR"))
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.String(), fromString("2 EUR").String())
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fromString("0.1 EUR").Add(fromString("0.1 EUR"))
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.String(), fromString("0.2 EUR").String())
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fromString("0.1 EUR").Add(fromString("0 EUR"))
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.String(), fromString("0.1 EUR").String())
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fromString("1 EUR").Add(fromString("1 CAD"))
		},
		check: func(a fin.Amount, e error) {
			is.Equal(e, fin.ErrMismatchedCurrency)
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}

func TestAmountMul(t *testing.T) {
	is := is.New(t)

	type test struct {
		do    func() (fin.Amount, error)
		check func(fin.Amount, error)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}
	fromString := func(s string) fin.Amount {
		amount, err := fin.AmountFromString(s)
		is.NoErr(err)
		return amount
	}

	add(test{
		do: func() (fin.Amount, error) {
			return fromString("2 EUR").Mul(fromString("5 EUR"))
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.String(), fromString("10 EUR").String())
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fromString("0.1 EUR").Mul(fromString("10 EUR"))
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.String(), fromString("1 EUR").String())
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fromString("0.1 EUR").Mul(fromString("0 EUR"))
		},
		check: func(a fin.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.String(), fromString("0 EUR").String())
		},
	})
	add(test{
		do: func() (fin.Amount, error) {
			return fromString("2 EUR").Mul(fromString("5 CAD"))
		},
		check: func(a fin.Amount, e error) {
			is.Equal(e, fin.ErrMismatchedCurrency)
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}
