package server

import (
	"context"
	"fmt"
	"github.com/clstb/phi/go/ledger/config"
	"github.com/clstb/phi/go/ledger/internal"
	pb "github.com/clstb/phi/go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
)

func (s *LedgerServer) SyncLedger(ctx context.Context, in *pb.UserNameMessage) (*emptypb.Empty, error) {

	file, err := os.Open(fmt.Sprintf("%s/%s", config.DataDirPath, in.Username))
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	userLedger := internal.NewLedger(file)
	err = s.Sync(userLedger)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *LedgerServer) Sync(ledger internal.Ledger) error {

	providers, err := GetProvidersRPC()
	if err != nil {
		return err
	}

	accounts, err := GetAccountsRPC()
	if err != nil {
		return err
	}

	transactions, err := GetTransactionRPC()
	if err != nil {
		return err
	}

	var filteredTransactions []internal.TinkTransaction
	for _, transaction := range transactions {
		if transaction.Status != "BOOKED" {
			continue
		}
		filteredTransactions = append(filteredTransactions, transaction)
	}

	ledger.UpdateLedger(providers, accounts, filteredTransactions)
	return nil
}

func GetProvidersRPC() ([]internal.Provider, error) {
	return nil, nil
}

func GetTransactionRPC() ([]internal.TinkTransaction, error) {
	return nil, nil
}

func GetAccountsRPC() ([]internal.Account, error) {
	return nil, nil
}
