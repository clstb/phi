package db_test

import (
	"testing"

	"github.com/clstb/phi/pkg/core/db"
	"github.com/matryer/is"
	"github.com/shopspring/decimal"
)

func TestAmounts(t *testing.T) {
	is := is.New(t)

	type test struct {
		do    func() db.Amounts
		check func(db.Amounts)
	}
	var tests []test
	add := func(t test) {
		tests = append(tests, t)
	}
	add(test{
		do: func() db.Amounts {
			return db.Amounts{
				{Currency: "EUR", Decimal: decimal.NewFromFloat(0.1)},
				{Currency: "EUR", Decimal: decimal.NewFromFloat(0.1)},
				{Currency: "EUR", Decimal: decimal.NewFromFloat(0.1)},
				{Currency: "EUR", Decimal: decimal.NewFromFloat(0.1)},
				{Currency: "EUR", Decimal: decimal.NewFromFloat(0.1)},
				{Currency: "USD", Decimal: decimal.NewFromFloat(0.3)},
				{Currency: "USD", Decimal: decimal.NewFromFloat(0.3)},
				{Currency: "USD", Decimal: decimal.NewFromFloat(0.3)},
				{Currency: "USD", Decimal: decimal.NewFromFloat(0.3)},
				{Currency: "USD", Decimal: decimal.NewFromFloat(0.3)},
				{Currency: "CAD", Decimal: decimal.NewFromFloat(1)},
				{Currency: "CAD", Decimal: decimal.NewFromFloat(1)},
				{Currency: "CAD", Decimal: decimal.NewFromFloat(1)},
				{Currency: "CAD", Decimal: decimal.NewFromFloat(1)},
				{Currency: "CAD", Decimal: decimal.NewFromFloat(1)},
			}.Sum()
		},
		check: func(a db.Amounts) {
			is.Equal(a.ByCurrency("EUR"), db.Amount{Currency: "EUR", Decimal: decimal.NewFromFloat(0.5)})
			is.Equal(a.ByCurrency("USD"), db.Amount{Currency: "USD", Decimal: decimal.NewFromFloat(1.5)})
			is.Equal(a.ByCurrency("CAD"), db.Amount{Currency: "CAD", Decimal: decimal.NewFromFloat(5)})
			is.Equal(a.Currencies(), []string{"EUR", "USD", "CAD"})
		},
	})

	for _, t := range tests {
		t.check(t.do())
	}
}
