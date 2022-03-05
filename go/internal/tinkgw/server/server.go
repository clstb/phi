package server

import (
	"context"
	"net/http"

	"github.com/clstb/phi/go/pkg/client/tink"
	"github.com/clstb/phi/go/pkg/middleware"
	chi "github.com/go-chi/chi/v5"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

type Server struct {
	r      *chi.Mux
	logger *zap.Logger

	tinkClientId     string
	tinkClientSecret string
	tinkClient       *tink.Client
	callbackURL      string
}

func NewServer(
	tinkClientId,
	tinkClientSecret,
	oryToken,
	jwksURL,
	callbackURL string,
) (*Server, error) {
	ctx := context.Background()

	loggerConfig := zap.NewProductionConfig()
	loggerConfig.OutputPaths = []string{"stdout"}
	logger, err := loggerConfig.Build()
	if err != nil {
		return nil, err
	}

	oauthConfig := &clientcredentials.Config{
		ClientID:     tinkClientId,
		ClientSecret: tinkClientSecret,
		TokenURL:     "https://api.tink.com/api/v1/oauth/token",
		Scopes:       []string{"authorization:grant,user:create"},
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	tinkClient := tink.NewClient(
		"https://api.tink.com",
		tink.WithHTTPClient(oauthConfig.Client(ctx)),
	)

	s := &Server{
		r:      chi.NewRouter(),
		logger: logger,

		tinkClientId:     tinkClientId,
		tinkClientSecret: tinkClientSecret,
		tinkClient:       tinkClient,
		callbackURL:      callbackURL,
	}

	s.r.Use(
		middleware.Auth(ctx, logger, jwksURL),
		s.provisionTinkUser(oryToken),
	)

	s.routes()

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
