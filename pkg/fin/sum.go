package fin

type Sum map[string]map[string]*Amount

func (s Sum) ByCurrency() map[string]*Amount {
	m := make(map[string]*Amount)

	for _, v := range s {
		for currency, a := range v {
			amount, ok := m[currency]
			if !ok {
				m[currency] = a
			} else {
				m[currency] = amount.Add(a)
			}
		}
	}

	return m

}
