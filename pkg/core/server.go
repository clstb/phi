package core

import (
	"database/sql"

	"github.com/clstb/phi/pkg/pb"
)

type Server struct {
	pb.UnimplementedCoreServer
	db *sql.DB
}

func New(opts ...Opt) *Server {
	s := &Server{}

	for _, opt := range opts {
		opt(s)
	}

	return s
}

type Opt func(s *Server)

func WithDB(db *sql.DB) Opt {
	return func(s *Server) {
		s.db = db
	}
}
