package db_test

import (
	"testing"

	"github.com/clstb/phi/pkg/core/db"
	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestAmountFromString(t *testing.T) {
	is := is.New(t)

	type test struct {
		do    func() (db.Amount, error)
		check func(db.Amount, error)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}

	add(test{
		do: func() (db.Amount, error) {
			return db.AmountFromString("0.1 EUR")
		},
		check: func(a db.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.Currency, "EUR")
			is.Equal(a.Decimal, decimal.NewFromFloat(0.1))
		},
	})
	add(test{
		do: func() (db.Amount, error) {
			return db.AmountFromString("0,1 EUR", db.AmountEU)
		},
		check: func(a db.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.Currency, "EUR")
			is.Equal(a.Decimal, decimal.NewFromFloat(0.1))
		},
	})
	add(test{
		do: func() (db.Amount, error) {
			return db.AmountFromString("")
		},
		check: func(a db.Amount, e error) {
			is.NoErr(e)
			is.Equal(a.Currency, "")
			is.Equal(a.Decimal, decimal.Zero)
		},
	})
	add(test{
		do: func() (db.Amount, error) {
			return db.AmountFromString("1")
		},
		check: func(a db.Amount, e error) {
			is.True(e != nil)
		},
	})
	add(test{
		do: func() (db.Amount, error) {
			return db.AmountFromString("characters")
		},
		check: func(a db.Amount, e error) {
			is.True(e != nil)
		},
	})
	add(test{
		do: func() (db.Amount, error) {
			return db.AmountFromString("characters EUR")
		},
		check: func(a db.Amount, e error) {
			is.True(e != nil)
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}
