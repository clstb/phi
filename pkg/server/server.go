package server

import (
	"database/sql"
	"log"
	"os"

	"github.com/doug-martin/goqu/v9"
	_ "github.com/doug-martin/goqu/v9/dialect/postgres"
)

type Server struct {
	Core *core
	db   *goqu.Database
}

func NewServer(opts ...ServerOpt) *Server {
	s := &Server{}
	internalOpts := []ServerOpt{
		withCore(),
	}

	for _, opt := range opts {
		opt(s)
	}

	for _, opt := range internalOpts {
		opt(s)
	}

	return s
}

type ServerOpt func(*Server)

// public opts

func WithDB(db *sql.DB) ServerOpt {
	return func(s *Server) {
		s.db = goqu.New("postgres", db)
		s.db.Logger(log.New(os.Stdout, "DB:", 0))
	}
}

// internal opts

func withCore() ServerOpt {
	return func(s *Server) {
		s.Core = &core{
			db: s.db,
		}
	}
}
