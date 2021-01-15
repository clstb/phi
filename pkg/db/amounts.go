package db

type Amounts []Amount

func (a Amounts) Sum() Amounts {
	sums := make(map[string]Amount)
	for _, amount := range a {
		sum, ok := sums[amount.Currency]
		if !ok {
			sum = amount
		} else {
			sum = sum.Add(amount)
		}
		sums[amount.Currency] = sum
	}

	var amounts Amounts
	for _, v := range sums {
		amounts = append(amounts, v)
	}

	return amounts
}

func (a Amounts) ByCurrency(currency string) Amount {
	for _, amount := range a {
		if amount.Currency == currency {
			return amount
		}
	}

	return Amount{}
}

func (a Amounts) Currencies() []string {
	var currencies []string
	for _, amount := range a {
		currencies = append(currencies, amount.Currency)
	}

	return currencies
}
