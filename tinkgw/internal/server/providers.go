package server

import (
	pb "github.com/clstb/phi/proto"
	"github.com/clstb/phi/tinkgw/internal/client/tink"
	"github.com/clstb/phi/tinkgw/internal/config"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetProviders(in *pb.StringMessage, server pb.TransactionGWService_GetProvidersServer) error {

	s.tinkClient.SetBearerToken(in.Value)
	res, err := s.tinkClient.GetProviders(config.DefaultMarket)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, i := range res {
		p := mapProvider(i)
		err = server.Send(&p)
		if err != nil {
			return status.Error(codes.Aborted, err.Error())
		}
	}
	return nil
}

func mapProvider(p tink.Provider) pb.ProviderMessage {
	return pb.ProviderMessage{
		FinancialInstitutionId: p.FinancialInstitutionID,
		DisplayName:            p.DisplayName,
	}
}
