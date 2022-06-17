package env

import (
	"context"
	"github.com/aaarrti/RL-2048/2048/game"
	"github.com/aaarrti/RL-2048/2048/util"
	pb "github.com/aaarrti/RL-2048/proto/go/proto"
	"google.golang.org/protobuf/types/known/emptypb"
)

type GameServer struct {
	pb.EnvServiceServer
}

var portGenerator = util.MainPort

func (c *GameServer) ProvisionEnvironment(context.Context, *emptypb.Empty) (*pb.IntMessage, error) {
	portGenerator++
	go func() {
		s := game.NewServerGame()
		s.ServerGame(portGenerator)
	}()
	return &pb.IntMessage{Value: int32(portGenerator)}, nil
}
