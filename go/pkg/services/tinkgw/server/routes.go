package server

import "github.com/go-chi/chi/middleware"

func (s *Server) routes() {
	s.r.Use(
		middleware.RequestID,
		middleware.RealIP,
		middleware.Logger,
	)

	callbacks := make(chan PendingCallback)
	s.callbacks = callbacks
	s.r.Get("/callback", s.Callback(callbacks))
}
