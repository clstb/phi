package core_test

import (
	"database/sql"
	"log"
	"net"
	"os"
	"testing"
	"time"

	"github.com/clstb/phi/pkg/core"
	coredb "github.com/clstb/phi/pkg/db/core"
	"github.com/clstb/phi/pkg/pb"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/test/bufconn"
)

const bufSize = 1024 * 1024

var conn *grpc.ClientConn
var db *coredb.Queries

func client() pb.CoreClient {
	return pb.NewCoreClient(conn)
}

func TestMain(m *testing.M) {
	sqlDB, err := sql.Open(
		"postgres",
		"postgres://phi@127.0.0.1:26257/phi?sslmode=disable",
	)
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()
	if err := sqlDB.Ping(); err != nil {
		log.Fatal(err)
	}
	db = coredb.New(sqlDB)

	lis := bufconn.Listen(bufSize)
	s := grpc.NewServer()
	server := core.New(
		core.WithDB(sqlDB),
	)
	if err != nil {
		log.Fatal(err)
	}
	pb.RegisterCoreServer(s, server)

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
