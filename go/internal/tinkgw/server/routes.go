package server

func (s *Server) routes(oryToken string) {
	s.r.Get("/api/link", s.Link())
	s.r.Post("/api/token", s.Token())
	s.r.Post("/api/tink-user", s.RegisterTinkUser(oryToken))
}
