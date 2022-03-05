package server

import (
	"fmt"
	"net/http"

	"github.com/clstb/phi/go/pkg/client/tink"
	"github.com/goccy/go-json"
	ory "github.com/ory/kratos-client-go"
	"go.uber.org/zap"
)

func (s *Server) Link() http.HandlerFunc {
	logger := s.logger.With(
		zap.String("handler", "link"),
	)

	return func(rw http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value("session").(ory.Session)
		if !ok {
			s.logger.Error("missing session")
			http.Error(rw, "missing session", http.StatusUnauthorized)
			return
		}

		code, err := s.tinkClient.GetAuthorizeGrantDelegateCode(
			"code",
			"",
			session.Identity.Id,
			session.Identity.Id,
			"authorization:read,authorization:grant,credentials:refresh,credentials:read,credentials:write,providers:read,user:read",
		)
		if err != nil {
			logger.Error("tink: authorize grant delegate", zap.Error(err))
			http.Error(rw, "tink: authorize grant delegate", http.StatusFailedDependency)
			return
		}

		link := fmt.Sprintf(
			"https://link.tink.com/1.0/transactions/connect-accounts?client_id=%s&redirect_uri=%s&market=%s&locale=%s&authorization_code=%s",
			s.tinkClientId,
			s.callbackURL,
			"DE",    // req.Market TODO,
			"de_DE", // req.Locale TODO,
			code,
		)

		rw.Header().Set("Content-Type", "text/plain")
		if _, err := rw.Write([]byte(link)); err != nil {
			logger.Error("writing repsonse", zap.Error(err))
		}
	}
}

func (s *Server) getToken(
	id string,
) (tink.Token, error) {
	code, err := s.tinkClient.GetAuthorizeGrantCode(
		"",
		id,
		"transactions:read,accounts:read,provider-consents:read,user:read",
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: authorize grant: %w", err)
	}
	token, err := s.tinkClient.GetToken(
		code,
		"",
		s.tinkClientId,
		s.tinkClientSecret,
		"authorization_code",
		"transactions:read,accounts:read,provider-consents:read,user:read",
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: oauth token: %w", err)
	}

	return token, nil
}

func (s *Server) Token() http.HandlerFunc {
	logger := s.logger.With(
		zap.String("handler", "token"),
	)
	type response struct {
		Header http.Header `json:"header"`
	}

	return func(rw http.ResponseWriter, r *http.Request) {
		session, ok := r.Context().Value("session").(ory.Session)
		if !ok {
			s.logger.Error("missing session")
			http.Error(rw, "missing session", http.StatusUnauthorized)
			return
		}

		token, err := s.getToken(session.Identity.Id)
		if err != nil {
			logger.Error("getting token", zap.Error(err))
			http.Error(rw, "getting token", http.StatusFailedDependency)
			return
		}

		if err := json.NewEncoder(rw).Encode(&response{Header: http.Header{
			"Authorization": []string{"Bearer " + token.AccessToken},
		}}); err != nil {
			logger.Error("marshalling response", zap.Error(err))
			http.Error(rw, "marshalling response", http.StatusInternalServerError)
		}
	}
}
