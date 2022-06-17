package main

import (
	"fmt"
	"github.com/aaarrti/RL-2048/2048/env"
	"github.com/aaarrti/RL-2048/2048/util"
	pb "github.com/aaarrti/RL-2048/proto/go/proto"
	grpczap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

func main() {

	addr := fmt.Sprintf("0.0.0.0:%v", util.MainPort)
	listener, err := net.Listen("tcp", addr)
	util.Must(err)

	logger := util.CreateLogger()

	log.Printf("----> GRPC listeninng on %v\n\n", addr)

	_server := grpc.NewServer(
		grpc.UnaryInterceptor(grpczap.UnaryServerInterceptor(logger)),
	)
	pb.RegisterEnvServiceServer(_server, &env.GameServer{})
	reflection.Register(_server)
	err = _server.Serve(listener)
	util.Must(err)
}
