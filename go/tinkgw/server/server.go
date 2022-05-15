package server

import (
	"context"
	"fmt"
	pb "github.com/clstb/phi/go/proto"
	tink2 "github.com/clstb/phi/go/tinkgw/client/tink"
	"github.com/clstb/phi/go/tinkgw/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"log"
	"time"
)

type Server struct {
	pb.TinkGWServiceServer
	pb.TransactionGWServiceServer
	Logger           *zap.SugaredLogger
	tinkClientId     string
	tinkClientSecret string
	tinkClient       *tink2.Client
	callbackURL      string
}

var _, sugar = func() (*zap.Logger, *zap.SugaredLogger) {
	loggerConfig := zap.NewProductionConfig()
	loggerConfig.EncoderConfig.TimeKey = "timestamp"
	loggerConfig.EncoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout(time.RFC3339)

	_logger, err := loggerConfig.Build()
	if err != nil {
		log.Fatal(err)
	}
	_sugar := _logger.Sugar()
	return _logger, _sugar
}()

func NewServer(tinkClientId, tinkClientSecret, callbackURL string) *Server {

	oauthConfig := &clientcredentials.Config{
		ClientID:     tinkClientId,
		ClientSecret: tinkClientSecret,
		TokenURL:     config.TinkTokenUri,
		Scopes:       []string{config.TinkAdminRoles},
		AuthStyle:    oauth2.AuthStyleInParams,
	}

	ctx := context.Background()

	tinkClient := tink2.NewClient(config.TinkUri, tink2.WithHTTPClient(oauthConfig.Client(ctx)))

	s := &Server{
		Logger:           sugar,
		tinkClientId:     tinkClientId,
		tinkClientSecret: tinkClientSecret,
		tinkClient:       tinkClient,
		callbackURL:      callbackURL,
	}

	return s
}

func (s *Server) getUser(id string) (tink2.User, error) {
	token, err := s.getToken(id)
	if err != nil {
		return tink2.User{}, err
	}
	client := tink2.NewClient(config.TinkUri)
	client.SetBearerToken(token.AccessToken)
	return client.GetUser()
}

func (s *Server) getToken(id string) (tink2.Token, error) {
	code, err := s.tinkClient.GetAuthorizeGrantCode(
		"",
		id,
		config.GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink2.Token{}, fmt.Errorf("tink: authorize grant: %w", err)
	}
	token, err := s.tinkClient.GetToken(
		code,
		"",
		s.tinkClientId,
		s.tinkClientSecret,
		"authorization_code",
		config.GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink2.Token{}, fmt.Errorf("tink: oauth token: %w", err)
	}
	return token, nil
}
