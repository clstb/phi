package server

import "github.com/clstb/phi/go/pkg/services/bookkeeper/pb"

type Server struct {
	pb.UnimplementedbookkeeperServer
	signingSecret []byte
}

func New(signingSecret []byte) *Server {
	s := &Server{
		signingSecret: signingSecret,
	}

	return s
}
