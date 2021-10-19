package server

import pb "github.com/clstb/phi/go/pkg/services/auth/pb"

type Server struct {
	pb.UnimplementedAuthServer
	signingSecret []byte
}

func New(signingSecret []byte) *Server {
	s := &Server{
		signingSecret: signingSecret,
	}

	return s
}
