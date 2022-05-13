package server

import (
	"github.com/clstb/phi/go/pkg/client/tink"
	"go.uber.org/zap"
	"net/http"
)

type Response struct {
	Header http.Header `json:"header"`
}

type Server struct {
	Logger *zap.Logger

	tinkClientId     string
	tinkClientSecret string
	tinkClient       *tink.Client
	callbackURL      string
	OryToken         string
}
