package core

import (
	"github.com/clstb/phi/pkg/pb"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	pb.UnimplementedCoreServer
	db *pgxpool.Pool
}

func New(opts ...Opt) *Server {
	s := &Server{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Opt func(s *Server)

func WithDB(db *pgxpool.Pool) Opt {
	return func(s *Server) {
		s.db = db
	}
}
