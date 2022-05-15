package server

import (
	"context"
	pb "github.com/clstb/phi/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) ProvisionTinkUser(ctx context.Context, in *emptypb.Empty) (*pb.ProvisionTinkUserResponse, error) {
	createdUser, err := s.tinkClient.CreateUserWithDefaults()
	if err != nil {
		s.Logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}

	s.Logger.Info("OK ---> ", createdUser.UserID)
	return &pb.ProvisionTinkUserResponse{
		ExternalUserId: createdUser.ExternalUserID,
		TinkId:         createdUser.UserID,
	}, nil
}
