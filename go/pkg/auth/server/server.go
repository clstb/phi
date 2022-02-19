package server

import "github.com/clstb/phi/go/pkg/auth/pb"

type Server struct {
	pb.UnimplementedAuthServer
	signingSecret []byte
}

func New(signingSecret []byte) *Server {
	return &Server{
		signingSecret: signingSecret,
	}
}
