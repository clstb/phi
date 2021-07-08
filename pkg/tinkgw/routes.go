package tinkgw

import "github.com/go-chi/chi/middleware"

func (s *Server) routes() {
	s.r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
		TX(s.db),
	)

	s.r.Get("/callback", s.Callback())
}
