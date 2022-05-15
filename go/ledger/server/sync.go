package server

import (
	"context"
	"fmt"
	"github.com/clstb/phi/go/ledger/beanacount"
	"github.com/clstb/phi/go/ledger/config"
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

	userLedger := beanacount.NewLedger(file)
	err = s.Sync(userLedger)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *LedgerServer) Sync(ledger beanacount.Ledger) error {

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

	var filteredTransactions []beanacount.TinkTransaction
	for _, transaction := range transactions {
		if transaction.Status != "BOOKED" {
			continue
		}
		filteredTransactions = append(filteredTransactions, transaction)
	}

	ledger.UpdateLedger(providers, accounts, filteredTransactions)
	return nil
}

func GetProvidersRPC() ([]beanacount.Provider, error) {
	return nil, nil
}

func GetTransactionRPC() ([]beanacount.TinkTransaction, error) {
	return nil, nil
}

func GetAccountsRPC() ([]beanacount.Account, error) {
	return nil, nil
}
