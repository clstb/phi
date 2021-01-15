package server

import (
	"database/sql"
	"fmt"
)

type Server struct {
	Core *core
	db   *sql.DB
}

func NewServer(opts ...ServerOpt) (*Server, error) {
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

	if s.db == nil {
		return nil, fmt.Errorf("no database configured")
	}

	return s, nil
}

type ServerOpt func(*Server)

// public opts

func WithDB(db *sql.DB) ServerOpt {
	return func(s *Server) {
		s.db = db
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
