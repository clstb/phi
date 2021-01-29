package fin

import "errors"

var (
	ErrMismatchedCurrency = errors.New("currencies don't match")
	ErrUnbalanced         = errors.New("transaction is unbalanced")
)
