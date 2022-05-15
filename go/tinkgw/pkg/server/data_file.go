package server

import (
	"context"
	"fmt"
	pb "github.com/clstb/phi/go/proto"
	"github.com/clstb/phi/go/tinkgw/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"os"
	"runtime/debug"
)

func (s *Server) ProvisionFSStructure(ctx context.Context, in *pb.UserNameMessage, opts ...grpc.CallOption) (*emptypb.Empty, error) {
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
