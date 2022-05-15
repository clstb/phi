package server

import (
	"context"
	"fmt"
	"github.com/clstb/phi/go/ledger/config"
	pb "github.com/clstb/phi/go/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"runtime/debug"
)

func (s *LedgerServer) ProvisionFSStructure(ctx context.Context, in *pb.UserNameMessage) (*emptypb.Empty, error) {
	err := createUserDir(in.Username)
	if err != nil {
		return &emptypb.Empty{}, status.Error(codes.Internal, err.Error())
	}
	return &emptypb.Empty{}, nil
}

//Each user needs
//.data/username/accounts.bean
//.data/username/transactions/

func createUserDir(username string) error {
	err := os.MkdirAll(fmt.Sprintf("%s/%s", config.DataDirPath, username), os.ModePerm)
	if err != nil {
		debug.PrintStack()
		return err
	}
	err = os.MkdirAll(fmt.Sprintf("%s/%s/transactions", config.DataDirPath, username), os.ModePerm)
	if err != nil {
		debug.PrintStack()
		return err
	}
	_, err = os.Create(fmt.Sprintf("%s/%s/accounts.bean", config.DataDirPath, username))
	if err != nil {
		debug.PrintStack()
		return err
	}
	return nil
}
