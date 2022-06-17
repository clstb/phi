package game

import (
	"context"
	"fmt"
	"github.com/aaarrti/RL-2048/2048/util"
	pb "github.com/aaarrti/RL-2048/proto/go/proto"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type ServerGameType struct {
	pb.GameServiceServer
	game     *IBoard
	maxScore int
}

func NewServerGame() ServerGameType {
	b := New()
	b.AddElement()
	b.AddElement()
	return ServerGameType{game: &b, maxScore: 0}
}

func (s *ServerGameType) ServerGame(port int) {
	addr := fmt.Sprintf("0.0.0.0:%v", port)
	listener, err := net.Listen("tcp", addr)
	util.Must(err)

	logger := util.CreateLogger()

	log.Printf("----> GRPC serving Game on %v\n\n", addr)

	_server := grpc.NewServer(
		grpc.UnaryInterceptor(grpczap.UnaryServerInterceptor(logger)),
	)
	pb.RegisterGameServiceServer(_server, s)
	reflection.Register(_server)
	err = _server.Serve(listener)
	util.Must(err)
}

func (s *ServerGameType) DoMove(ctx context.Context, in *pb.MoveMessage) (*pb.GameState, error) {
	fmt.Printf("Received Move: %v\n", in.Value)
	board := *s.game
	move := mapMove(in.Value)
	board.Move(move)

	res := pb.GameState{Value: flattenMatrix(board.(*SBoard).Matrix)}
	return &res, nil
}

func mapMove(enum pb.MoveEnum) Dir {
	switch enum {
	case pb.MoveEnum_UP:
		return UP
	case pb.MoveEnum_DOWN:
		return DOWN
	case pb.MoveEnum_LEFT:
		return LEFT
	case pb.MoveEnum_RIGHT:
		return RIGHT
	default:
		return NO_DIR
	}
}

func flattenMatrix(matrix [][]int) []int32 {
	var flatArr []int32
	for _, row := range matrix {
		for _, i := range row {
			flatArr = append(flatArr, int32(i))
		}
	}
	return flatArr
}
