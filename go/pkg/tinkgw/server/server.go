package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/clstb/phi/go/pkg/tink"
	"github.com/clstb/phi/go/pkg/tinkgw/pb"
	"github.com/go-chi/chi"
	"github.com/jackc/pgx/v4/pgxpool"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
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
	db *pgxpool.Pool
}

func New(
	ctx context.Context,
	logger *zap.Logger,
	clientID string,
	clientSecret string,
	callbackURL string,
	db *pgxpool.Pool,
) *Server {
	config := &clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     "https://api.tink.com/api/v1/oauth/token",
		Scopes:       []string{"authorization:grant,user:create"},
		AuthStyle:    oauth2.AuthStyleInParams,
	}
	tinkClient := tink.NewClient(config.Client(ctx))

	s := &Server{
		logger:       logger,
		clientID:     clientID,
		clientSecret: clientSecret,
		callbackURL:  callbackURL,
		tink:         tinkClient,
		r:            chi.NewMux(),
		db:           db,
	}

	s.routes()

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
