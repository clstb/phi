package server

func (s *Server) routes() {
	s.r.Get("/link", s.Link())
	s.r.Post("/token", s.Token())
}
