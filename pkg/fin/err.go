package fin

import (
	"errors"
	"fmt"
)

var (
	// ErrMismatchedCurrency occurs when two amounts are added and currencies don't match.
	ErrMismatchedCurrency = errors.New("currencies don't match")
	// ErrUnbalanced occurs when a transaction is unbalanced.
	ErrUnbalanced = errors.New("transaction is unbalanced")
)

// ErrNotFound occurs when a resource that was searched for is not found
type ErrNotFound struct {
	kind string
	name string
}

// Error returns the error as string.
func (e ErrNotFound) Error() string {
	return fmt.Sprintf(
		"%s not found: %s",
		e.kind,
		e.name,
	)
}
