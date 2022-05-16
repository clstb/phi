package server

import (
	"context"
	"fmt"
	pb "github.com/clstb/phi/proto"
	"github.com/clstb/phi/tinkgw/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
)

func (s *Server) ProvisionTinkUser(ctx context.Context, in *emptypb.Empty) (*pb.ProvisionTinkUserResponse, error) {
	return nil, status.Error(codes.Unimplemented, "Tink Admin account needed!")
}

func buildLink(clientId string, test bool) string {
	link := fmt.Sprintf(
		"https://link.tink.com/1.0/transactions/connect-accounts/?client_id=%s&redirect_uri=https",
		clientId) + "%3A%2F%2Fconsole.tink.com%2Fcallback" +
		fmt.Sprintf("&market=%s&locale=%s", config.DefaultMarket, config.DefaultLocale)
	if !test {
		return link
	}
	return link + "&test=true"
}
func (s *Server) GetTinkAuthLink(ctx context.Context, in *emptypb.Empty) (*pb.BytesMessage, error) {

	clientId := os.Getenv("TINK_CLIENT_ID")
	link := buildLink(clientId, false)
	s.Logger.Info(link)
	return &pb.BytesMessage{Arr: []byte(link)}, nil

}

func (s *Server) GetTestAuthLink(ctx context.Context, in *emptypb.Empty) (*pb.BytesMessage, error) {

	clientId := os.Getenv("TINK_CLIENT_ID")
	link := buildLink(clientId, true)
	s.Logger.Info(link)
	return &pb.BytesMessage{Arr: []byte(link)}, nil

}

func (s *Server) ExchangeAuthCodeToToken(ctx context.Context, in *pb.StringMessage) (*pb.StringMessage, error) {
	res, err := s.basicClient.ExchangeAccessCodeForToken(in.Value)
	if err != nil {
		s.Logger.Error(err)
		return nil, status.Error(codes.FailedPrecondition, err.Error())
	}
	s.Logger.Info(res)
	return &pb.StringMessage{Value: res}, nil
}
