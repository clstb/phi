package fin

import "errors"

var (
	// ErrMismatchedCurrency occurs when two amounts are added and currencies don't match.
	ErrMismatchedCurrency = errors.New("currencies don't match")
	// ErrUnbalanced occurs when a transaction is unbalanced.
	ErrUnbalanced = errors.New("transaction is unbalanced")
)
