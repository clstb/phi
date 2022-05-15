package server

import (
	"context"
	"fmt"
	pb "github.com/clstb/phi/proto"
	"github.com/clstb/phi/tinkgw/internal/config"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

func (s *Server) CreateTinkLink(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {

	tinkClientId := os.Getenv("TINK_CLIENT_ID")
	code, err := s.tinkClient.GetDelegatedAuthorizationCode(tinkClientId, in.Value)
	if err != nil {
		s.Logger.Error("tink: authorize grant delegate", zap.Error(err))
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}

	link := fmt.Sprintf(
		config.LinkBankAccountUriFormat,
		s.tinkClientId,
		s.callbackURL,
		config.DefaultMarket,
		config.DefaultLocale,
		code,
	)
	return &pb.StringMessage{Value: link}, nil

}
