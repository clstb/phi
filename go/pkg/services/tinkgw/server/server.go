package server

import (
	"net/http"
	"sync"

	"github.com/clstb/phi/go/pkg/services/tinkgw/pb"
	"github.com/clstb/phi/go/pkg/services/tinkgw/tink"
	"github.com/go-chi/chi"
	"go.uber.org/zap"
)

type PendingCallback struct {
	del   bool
	ch    chan<- string
	state string
}

type Server struct {
	pb.UnimplementedTinkGWServer
	clientID     string
	clientSecret string
	callbackURL  string
	callbacks    chan<- PendingCallback
	tink         *tink.Client
	r            *chi.Mux
	logger       *zap.Logger
	sync.RWMutex
}

func New(
	logger *zap.Logger,
	clientID string,
	clientSecret string,
	callbackURL string,
) (*Server, error) {
	tink, err := tink.NewClient(
		tink.WithAuth(
			clientID,
			clientSecret,
			"authorization:grant,user:create",
		),
	)

	if err != nil {
		return nil, err
	}

	s := &Server{
		logger:       logger,
		clientID:     clientID,
		clientSecret: clientSecret,
		callbackURL:  callbackURL,
		tink:         tink,
		r:            chi.NewMux(),
	}

	s.routes()

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
