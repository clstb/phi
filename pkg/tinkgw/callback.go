package tinkgw

import (
	"net/http"

	db "github.com/clstb/phi/pkg/db/tinkgw"
	"github.com/jackc/pgx/v4"
)

func (s *Server) Callback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		state := r.URL.Query().Get("state")
		if state == "" {
			http.Error(w, "missing query parameter: state", http.StatusBadRequest)
			return
		}

		user, ok := s.states[state]
		if !ok {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}

		credentialsID := r.URL.Query().Get("credentialsId")
		if credentialsID == "" {
			http.Error(w, "missing query parameter: credentialsId", http.StatusBadRequest)
			return
		}

		tx, ok := ctx.Value("tx").(pgx.Tx)
		if !ok {
			http.Error(w, "context: missing transaction", http.StatusInternalServerError)
			return
		}
		q := db.New(tx)

		_, err := q.CreateCredential(ctx, db.CreateCredentialParams{
			User:          user,
			CredentialsID: credentialsID,
		})
		if err != nil {
			http.Error(w, "storing credentialID failed", http.StatusInternalServerError)
			return
		}
	}
}
