package tinkgw

import (
	"net/http"

	"github.com/clstb/phi/pkg/pb"
	"github.com/clstb/phi/pkg/tink"
	"github.com/go-chi/chi"
	"github.com/gofrs/uuid"
	"github.com/jackc/pgx/v4/pgxpool"
)

type Server struct {
	pb.UnimplementedTinkGWServer
	clientID     string
	clientSecret string
	callbackURL  string
	db           *pgxpool.Pool
	tink         *tink.Client
	core         pb.CoreClient
	r            *chi.Mux
	states       map[string]uuid.UUID
}

func New(
	clientID string,
	clientSecret string,
	callbackURL string,
	db *pgxpool.Pool,
	core pb.CoreClient,
) (*Server, error) {
	tink, err := tink.NewClient(
		clientID,
		clientSecret,
		"user:create,authorization:grant",
	)
	if err != nil {
		return nil, err
	}

	s := &Server{
		clientID:     clientID,
		clientSecret: clientSecret,
		callbackURL:  callbackURL,
		db:           db,
		core:         core,
		tink:         tink,
		r:            chi.NewMux(),
		states:       make(map[string]uuid.UUID),
	}

	s.routes()

	return s, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.r.ServeHTTP(w, r)
}
