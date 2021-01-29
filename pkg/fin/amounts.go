package fin

// Amounts is a slice of amount.
type Amounts []Amount

// Sum calculates the sum of all amounts grouped by the currency of each amount.
func (a Amounts) Sum() (Amounts, error) {
	sums := make(map[string]Amount)
	for _, amount := range a {
		sum, ok := sums[amount.Currency]
		if !ok {
			sum = amount
		} else {
			v, err := sum.Add(amount)
			if err != nil {
				return Amounts{}, err
			}
			sum = v
		}
		sums[amount.Currency] = sum
	}

	var amounts Amounts
	for _, v := range sums {
		amounts = append(amounts, v)
	}

	return amounts, nil
}

// ByCurrency returns the first amount matching given currency.
// An empty amount is returned when no amount is found.
// This is usually used after summing amounts.
func (a Amounts) ByCurrency(currency string) Amount {
	for _, amount := range a {
		if amount.Currency == currency {
			return amount
		}
	}

	return Amount{}
}

// Currencies returns a slice of all currencies in amounts.
func (a Amounts) Currencies() []string {
	var currencies []string
	for _, amount := range a {
		currencies = append(currencies, amount.Currency)
	}

	return currencies
}
