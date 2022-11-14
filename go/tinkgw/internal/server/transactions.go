package server

import (
	pb "github.com/clstb/phi/proto"
	"github.com/clstb/phi/tinkgw/internal/client/tink"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Server) GetAccounts(in *pb.StringMessage, server pb.TransactionGWService_GetAccountsServer) error {

	authorizedClient := tink.NewAuthorizedClient(in.Value, s.Logger)
	res, err := authorizedClient.GetAccounts()

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

func (s *Server) GetProviders(in *pb.StringMessage, server pb.TransactionGWService_GetProvidersServer) error {

	authorizedClient := tink.NewAuthorizedClient(in.Value, s.Logger)
	res, err := authorizedClient.GetProviders()

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

func (s *Server) GetTransactions(in *pb.StringMessage, server pb.TransactionGWService_GetTransactionsServer) error {

	authorizedClient := tink.NewAuthorizedClient(in.Value, s.Logger)
	res, err := authorizedClient.GetTransactions()

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
