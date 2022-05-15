package server

import (
	"context"
	"fmt"
	pb "github.com/clstb/phi/go/proto"
	"github.com/clstb/phi/go/tinkgw/pkg/client/tink"
	"github.com/clstb/phi/go/tinkgw/pkg/config"
	"github.com/clstb/phi/go/tinkgw/pkg/server/ledger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
)

func (s *Server) SyncLedger(ctx context.Context, in *pb.UserNameMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {

	file, err := os.Open(fmt.Sprintf("%s/%s", config.DataDirPath, in.Username))
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}

	userLedger := ledger.NewLedger(file)
	err = s.Sync(userLedger)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

func (s *Server) Sync(ledger ledger.Ledger) error {

	providers, err := s.tinkClient.GetProviders("DE")
	if err != nil {
		return err
	}

	accounts, err := s.tinkClient.GetAccounts("")
	if err != nil {
		return err
	}

	transactions, err := s.tinkClient.GetTransactions("")
	if err != nil {
		return err
	}

	var filteredTransactions []tink.Transaction
	for _, transaction := range transactions {
		if transaction.Status != "BOOKED" {
			continue
		}
		filteredTransactions = append(filteredTransactions, transaction)
	}

	ledger.UpdateLedger(providers, accounts, filteredTransactions)
	return nil
}
