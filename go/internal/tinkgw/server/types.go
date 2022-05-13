package server

import (
	"fmt"
	"github.com/clstb/phi/go/pkg/client/tink"
	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"net/http"
)

type HttpError struct {
	httpCode    int
	description string
}

func (h HttpError) Error() string {
	return fmt.Sprintf("http %d: %s: syntax error", h.httpCode, h.description)
}

type Response struct {
	Header http.Header `json:"header"`
}

type Server struct {
	r      *chi.Mux
	logger *zap.Logger

	tinkClientId     string
	tinkClientSecret string
	tinkClient       *tink.Client
	callbackURL      string
}
