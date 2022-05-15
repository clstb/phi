package server

import (
	pb "github.com/clstb/phi/go/proto"
	"github.com/clstb/phi/go/tinkgw/client/tink"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetTransactions(in *pb.StringMessage, server pb.TransactionGWService_GetTransactionsServer) error {
	s.tinkClient.SetBearerToken(in.Value)
	res, err := s.tinkClient.GetTransactions()
	if err != nil {
		s.Logger.Error(err)
		return status.Error(codes.Internal, err.Error())
	}

	for _, i := range res {
		m := mapTransaction(i)
		err = server.Send(&m)
		if err != nil {
			s.Logger.Error(err)
			return err
		}
	}
	return nil
}

func mapTransaction(tr tink.Transaction) pb.TinkTransactionMessage {

	datesMessage := pb.DatesMessage{
		Booked: tr.Dates.Booked,
		Value:  tr.Dates.Value,
	}
	valueMessage := pb.ValueMessage{
		Scale:         tr.Amount.Value.Scale,
		UnscaledValue: tr.Amount.Value.UnscaledValue,
	}
	amountMessage := pb.AmountMessage{
		CurrencyCode: tr.Amount.CurrencyCode,
		Value:        &valueMessage,
	}
	return pb.TinkTransactionMessage{
		AccountID:   tr.AccountID,
		ID:          tr.ID,
		Amount:      &amountMessage,
		Dates:       &datesMessage,
		Reference:   tr.Reference,
		Description: tr.Descriptions.Display,
		Status:      tr.Status,
	}
}
