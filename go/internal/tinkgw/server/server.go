package server

import (
	"context"
	"github.com/clstb/phi/go/pkg/client/tink"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
)

func NewServer(
	tinkClientId,
	tinkClientSecret,
	oryToken,
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
		Logger: logger,

		tinkClientId:     tinkClientId,
		tinkClientSecret: tinkClientSecret,
		tinkClient:       tinkClient,
		callbackURL:      callbackURL,
		OryToken:         oryToken,
	}

	return s, nil
}
