package server

import (
	"context"
	"fmt"
	pb "github.com/clstb/phi/go/proto"
	"github.com/clstb/phi/go/tinkgw/pkg/config"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"os"
)

func (s *Server) CreateTinkLink(ctx context.Context, in *pb.UserIdMessage) (*pb.TinkLinkMessage, error) {

	tinkClientId := os.Getenv("TINK_CLIENT_ID")
	code, err := s.tinkClient.GetDelegatedAutorizationCode(tinkClientId, in.UserId)
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
	return &pb.TinkLinkMessage{TinkLink: link}, nil

}
