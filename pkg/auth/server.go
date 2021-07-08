package auth

import (
	"github.com/clstb/phi/pkg/pb"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	pb.UnimplementedAuthServer
	db            *pgxpool.Pool
	signingSecret []byte
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

func WithSigningSecret(signingSecret []byte) Opt {
	return func(s *Server) {
		s.signingSecret = signingSecret
	}
}
