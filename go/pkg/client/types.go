package client

import (
	"fmt"
)

type HttpError struct {
	HttpCode    int
	Description string
}

func (h HttpError) Error() string {
	return fmt.Sprintf("http %d: %s: syntax error", h.HttpCode, h.Description)
}
