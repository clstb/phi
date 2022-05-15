package server

import (
	pb "github.com/clstb/phi/go/proto"
	"github.com/clstb/phi/go/tinkgw/client/tink"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAccounts(in *pb.StringMessage, server pb.TransactionGWService_GetAccountsServer) error {

	s.tinkClient.SetBearerToken(in.Value)
	res, err := s.tinkClient.GetAccounts()

	if err != nil {
		s.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}
	for _, i := range res {
		m := mapAccount(i)
		err = server.Send(&m)
		if err != nil {
			s.Logger.Error(err)
			return err
		}
	}
	return nil
}

func mapAccount(acc tink.Account) pb.AccountMessage {

	return pb.AccountMessage{
		FinancialInstitutionId: acc.FinancialInstitutionID,
		ID:                     acc.ID,
		Name:                   acc.Name,
	}

}
