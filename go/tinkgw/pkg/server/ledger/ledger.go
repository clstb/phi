package ledger

import (
	"bufio"
	"fmt"
	"io"
)

type stringer string

func (s stringer) String() string {
	return string(s)
}

type Ledger []fmt.Stringer

func NewLedger(r io.Reader) Ledger {
	s := bufio.NewScanner(r)

	var ledger Ledger
	parse := func(t string) {
		switch {
		case TransactionRE.MatchString(t):
			ledger = append(ledger, parseTransaction(t))
		case OpenRE.MatchString(t):
			ledger = append(ledger, parseOpen(t))
		default:
			i := len(ledger) - 1
			if i < 0 {
				ledger = append(ledger, stringer(t))
				return
			}

			switch v := ledger[i].(type) {
			case Transaction:
				switch {
				case PostingRE.MatchString(t):
					v.Postings = append(v.Postings, parsePosting(t))
					ledger[i] = v
				case MetadataRE.MatchString(t):
					v.Metadata = append(v.Metadata, parseMetadata(t))
					ledger[i] = v
				default:
					ledger = append(ledger, stringer(t))
				}
			case Open:
				switch {
				case MetadataRE.MatchString(t):
					v.Metadata = append(v.Metadata, parseMetadata(t))
					ledger[i] = v
				default:
					ledger = append(ledger, stringer(t))
				}
			default:
				ledger = append(ledger, stringer(t))
			}
		}
	}

	for s.Scan() {
		parse(s.Text())
	}

	return ledger
}

func (l Ledger) Transactions() Transactions {
	var transactions Transactions
	for _, v := range l {
		transaction, ok := v.(Transaction)
		if !ok {
			continue
		}
		transactions = append(transactions, transaction)
	}
	return transactions
}

func (l Ledger) Opens() Opens {
	var opens Opens
	for _, v := range l {
		open, ok := v.(Open)
		if !ok {
			continue
		}
		opens = append(opens, open)
	}
	return opens
}
