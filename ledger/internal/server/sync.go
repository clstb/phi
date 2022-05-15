package server

import (
	"context"
	"fmt"
	"github.com/clstb/phi/ledger/internal/beanacount"
	"github.com/clstb/phi/ledger/internal/config"
	pb "github.com/clstb/phi/proto"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"io"
	"os"
)

func (s *LedgerServer) SyncLedger(ctx context.Context, in *pb.SyncMessage) (*emptypb.Empty, error) {
	file, err := os.Open(fmt.Sprintf("%s/%s", config.DataDirPath, in.Username))
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	userLedger := beanacount.NewLedger(file)
	err = s.Sync(userLedger, in.Token, in.Username)
	if err != nil {
		s.Logger.Error(err)
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *LedgerServer) Sync(ledger beanacount.Ledger, token string, username string) error {
	connection, err := grpc.Dial(config.TinkGwAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithStreamInterceptor(grpczap.StreamClientInterceptor(s.Logger.Desugar())),
	)
	if err != nil {
		return err
	}
	defer connection.Close()
	gwServiceClient := pb.NewTransactionGWServiceClient(connection)

	providers, err := GetMockProvidersRPC()
	if err != nil {
		return err
	}

	accounts, err := GetAccountsRPC(gwServiceClient, token)
	if err != nil {
		return err
	}

	transactions, err := GetTransactionRPC(gwServiceClient, token)
	if err != nil {
		return err
	}

	var filteredTransactions []beanacount.TinkTransaction
	for _, transaction := range transactions {
		if transaction.Status != "BOOKED" {
			continue
		}
		filteredTransactions = append(filteredTransactions, transaction)
	}

	ledger.UpdateLedger(providers, accounts, filteredTransactions)
	err = ledger.PersistLedger(username)
	if err != nil {
		s.Logger.Error(err)
		return err
	}
	return nil
}

func GetMockProvidersRPC() ([]beanacount.Provider, error) {
	var slice []beanacount.Provider
	return slice, nil
}

// GetProvidersRPC Doesn't work with TINK admin account :(
func GetProvidersRPC(client pb.TransactionGWServiceClient, token string) ([]beanacount.Provider, error) {
	stream, err := client.GetProviders(context.Background(), &pb.StringMessage{Value: token})
	if err != nil {
		return nil, err
	}

	var providers []beanacount.Provider
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		a := mapProvider(in)
		providers = append(providers, a)
	}
	return providers, nil
}

func mapProvider(pr *pb.ProviderMessage) beanacount.Provider {
	return beanacount.Provider{
		FinancialInstitutionId: pr.FinancialInstitutionId,
		DisplayName:            pr.DisplayName,
	}
}

func GetTransactionRPC(client pb.TransactionGWServiceClient, token string) ([]beanacount.TinkTransaction, error) {
	stream, err := client.GetTransactions(context.Background(), &pb.StringMessage{Value: token})
	if err != nil {
		return nil, err
	}

	var transaction []beanacount.TinkTransaction
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		a := mapTransaction(in)
		transaction = append(transaction, a)
	}
	return transaction, nil
}

func mapTransaction(tr *pb.TinkTransactionMessage) beanacount.TinkTransaction {
	return beanacount.TinkTransaction{
		Status:    tr.Status,
		AccountID: tr.AccountID,
		ID:        tr.ID,
		Amount: beanacount.Amount{
			CurrencyCode: tr.Amount.CurrencyCode,
			Value: beanacount.Value{
				Scale:         tr.Amount.Value.Scale,
				UnscaledValue: tr.Amount.Value.UnscaledValue,
			},
		},
		Dates: beanacount.Dates{
			Booked: tr.Dates.Booked,
			Value:  tr.Dates.Value,
		},
		Reference:    tr.Reference,
		Descriptions: tr.Description,
	}

}

func GetAccountsRPC(client pb.TransactionGWServiceClient, token string) ([]beanacount.Account, error) {
	stream, err := client.GetAccounts(context.Background(), &pb.StringMessage{Value: token})
	if err != nil {
		return nil, err
	}

	var accounts []beanacount.Account
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			break
		}
		a := mapAccount(in)
		accounts = append(accounts, a)
	}
	return accounts, nil
}

func mapAccount(acc *pb.AccountMessage) beanacount.Account {
	return beanacount.Account{
		FinancialInstitutionId: acc.FinancialInstitutionId,
		ID:                     acc.ID,
		Name:                   acc.Name,
	}
}
