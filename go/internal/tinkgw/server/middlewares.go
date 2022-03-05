package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/clstb/phi/go/pkg/client/tink"
	ory "github.com/ory/kratos-client-go"
	"go.uber.org/zap"
)

func (s *Server) getUser(id string) (tink.User, error) {
	token, err := s.getToken(id)
	if err != nil {
		return tink.User{}, err
	}

	client := tink.NewClient("https://api.tink.com")
	client.SetBearerToken(token.AccessToken)
	return client.GetUser()
}

func (s *Server) provisionTinkUser(oryToken string) func(http.Handler) http.Handler {
	oryConf := ory.NewConfiguration()
	oryConf.Servers = ory.ServerConfigurations{
		{
			URL: "https://romantic-kapitsa-wjt1qzo59j.projects.oryapis.com/api/kratos/admin",
		},
	}
	oryConf.AddDefaultHeader("Authorization", "Bearer "+oryToken)
	oryConf.HTTPClient = &http.Client{}
	oryClient := ory.NewAPIClient(oryConf)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			session, ok := r.Context().Value("session").(ory.Session)
			if !ok {
				s.logger.Error("missing session")
				http.Error(rw, "missing session", http.StatusUnauthorized)
				return
			}
			traits := session.Identity.Traits.(map[string]interface{})

			tinkId, ok := traits["tink_id"]
			if !ok || tinkId == "" {
				createdUser, err := s.tinkClient.CreateUser(
					session.Identity.Id,
					"DE",
					"de_DE",
				)
				if err != nil {
					if !errors.Is(err, tink.ErrUserExists) {
						s.logger.Error("tink: creating user", zap.Error(err))
						http.Error(rw, "tink: creating user", http.StatusFailedDependency)
						return
					}
					user, err := s.getUser(session.Identity.Id)
					if err != nil {
						s.logger.Error("tink: getting user", zap.Error(err))
						http.Error(rw, "tink: getting user", http.StatusFailedDependency)
						return
					}
					createdUser.UserID = user.Id
				}

				traits["tink_id"] = createdUser.UserID
				identity, _, err := oryClient.V0alpha2Api.AdminUpdateIdentity(
					context.Background(),
					session.Identity.Id,
				).AdminUpdateIdentityBody(ory.AdminUpdateIdentityBody{
					State:  *session.Identity.State,
					Traits: traits,
				}).Execute()
				if err != nil {
					s.logger.Error("ory: updating identity", zap.Error(err))
					http.Error(rw, "ory: updating identity", http.StatusFailedDependency)
					return
				}

				session.Identity = *identity
			}

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
