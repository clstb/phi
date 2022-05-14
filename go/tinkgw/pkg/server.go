package pkg

import (
	"fmt"
	pb "github.com/clstb/phi/go/proto"
	"github.com/clstb/phi/go/tinkgw/pkg/client/tink"
	"go.uber.org/zap"
)

type Server struct {
	pb.TinkGWServiceServer
	Logger           *zap.SugaredLogger
	tinkClientId     string
	tinkClientSecret string
	tinkClient       *tink.Client
	callbackURL      string
}

func NewServer(tinkClientId, tinkClientSecret, callbackURL string) *Server {

	tinkClient := tink.NewClient(TinkUri)

	s := &Server{
		Logger:           Sugar,
		tinkClientId:     tinkClientId,
		tinkClientSecret: tinkClientSecret,
		tinkClient:       tinkClient,
		callbackURL:      callbackURL,
	}

	return s
}

func (s *Server) getUser(id string) (tink.User, error) {
	token, err := s.getToken(id)
	if err != nil {
		return tink.User{}, err
	}
	client := tink.NewClient(TinkUri)
	client.SetBearerToken(token.AccessToken)
	return client.GetUser()
}

func (s *Server) getToken(id string) (tink.Token, error) {
	code, err := s.tinkClient.GetAuthorizeGrantCode(
		"",
		id,
		GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: authorize grant: %w", err)
	}
	token, err := s.tinkClient.GetToken(
		code,
		"",
		s.tinkClientId,
		s.tinkClientSecret,
		"authorization_code",
		GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: oauth token: %w", err)
	}
	return token, nil
}
