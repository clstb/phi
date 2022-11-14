package beanacount

import (
	"fmt"
)

func (o Open) String() string {
	s := fmt.Sprintf("%s open %s\n", o.Date, o.Account.String())

	for _, md := range o.Metadata {
		s += md.String()
	}
	return s
}

type Opens []Open

func (t Opens) ByTinkId() map[string]Open {
	m := map[string]Open{}
	for _, open := range t {
		for _, md := range open.Metadata {
			if md.Key == "tink_id" {
				m[md.Value[1:len(md.Value)-1]] = open
			}
		}
	}
	return m
}

func (t Opens) Filter(f func(open Open) bool) []Open {
	var filtered Opens
	for _, open := range t {
		if f(open) {
			filtered = append(filtered, open)
		}
	}

	return filtered
}
