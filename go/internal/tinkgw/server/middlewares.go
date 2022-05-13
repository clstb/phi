package server

import (
	"github.com/clstb/phi/go/pkg/client/tink"
)

func (s *Server) getUser(id string) (tink.User, error) {
	return getUser(id, s.tinkClient, s.tinkClientId, s.tinkClientSecret)
}

func getUser(id string, tinkClient *tink.Client, tinkClientId string, tinkClientSecret string) (tink.User, error) {
	token, err := GetToken(id, tinkClient, tinkClientId, tinkClientSecret)
	if err != nil {
		return tink.User{}, err
	}
	client := tink.NewClient(TinkUri)
	client.SetBearerToken(token.AccessToken)
	return client.GetUser()
}

/*
func (s *Server) provisionTinkUser(oryToken string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			session, ok := r.Context().Value("session").(ory.Session)
			if !ok {
				s.Logger.Error("missing session")
				http.Error(rw, "missing session", http.StatusUnauthorized)
				return
			}
			traits := session.Identity.Traits.(map[string]interface{})

			tinkId, ok := traits["tink_id"]
			if !ok || tinkId == "" {
				 err := createTinkClient(r.Context(), s.tinkClient, oryToken, s.tinkClientId, s.tinkClientSecret)
				if err != nil {
					http.Error(rw, err.Description, err.HttpCode)
				}
			}
			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
*/
