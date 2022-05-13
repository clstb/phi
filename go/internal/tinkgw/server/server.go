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
		TokenURL:     TinkTokenUri,
		Scopes:       []string{TinkAdminRoles},
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	tinkClient := tink.NewClient(
		TinkUri,
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

	s.routes(oryToken)

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
