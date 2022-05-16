package server

import (
	"github.com/clstb/phi/pkg"
	pb "github.com/clstb/phi/proto"
	"github.com/clstb/phi/tinkgw/internal/client/tink"
	"go.uber.org/zap"
)

type Server struct {
	pb.TinkGWServiceServer
	pb.TransactionGWServiceServer
	Logger           *zap.SugaredLogger
	tinkClientId     string
	tinkClientSecret string
	callbackURL      string
	basicClient      *tink.Client
}

func NewServer(tinkClientId, tinkClientSecret, callbackURL string) *Server {

	log := pkg.CreateLogger()
	s := &Server{
		Logger:           log.Named("tinkGW"),
		tinkClientId:     tinkClientId,
		tinkClientSecret: tinkClientSecret,
		callbackURL:      callbackURL,
		basicClient:      tink.NewClient(log),
	}
	return s
}

/*
func (s *Server) getUser(id string) (tink.User, error) {
	token, err := s.getToken(id)
	if err != nil {
		return tink.User{}, err
	}
	client := tink.NewClient(config.TinkUri)
	client.SetBearerToken(token.AccessToken)
	return client.GetUser()
}
*/

/*
func (s *Server) getToken(id string) (tink.Token, error) {
	code, err := s.tinkClient.GetAuthorizeGrantCode(
		"",
		id,
		config.GetAuthorizeGrantCodeRoles,
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
		config.GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: oauth token: %w", err)
	}
	return token, nil
}
*/
