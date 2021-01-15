package server_test

import (
	"database/sql"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/server"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var conn *grpc.ClientConn
var db *sql.DB

func coreClient() pb.CoreClient {
	return pb.NewCoreClient(conn)
}

func TestMain(m *testing.M) {
	db, err := sql.Open(
		"postgres",
		"postgres://phi@127.0.0.1:26257/phi?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	server, err := server.NewServer(
		server.WithDB(db),
	)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterCoreServer(s, server.Core)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatal(err)
		}
	}()

	conn, err = grpc.Dial(
		"bufnet",
		grpc.WithInsecure(),
		grpc.WithDialer(func(s string, d time.Duration) (net.Conn, error) {
			return lis.Dial()
		}),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	status := m.Run()
	s.GracefulStop()
	os.Exit(status)
}
