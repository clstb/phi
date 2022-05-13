package server

import (
	"context"
	"errors"
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
			GetAuthorizeGrantDelegateCodeRoles,
		)
		if err != nil {
			logger.Error("tink: authorize grant delegate", zap.Error(err))
			http.Error(rw, "tink: authorize grant delegate", http.StatusFailedDependency)
			return
		}

		link := fmt.Sprintf(
			LinkBankAccountUriFormat,
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

func (s *Server) getToken(id string) (tink.Token, error) {
	return GetToken(id, s.tinkClient, s.tinkClientId, s.tinkClientSecret)
}

func (s *Server) Token() http.HandlerFunc {
	logger := s.logger.With(
		zap.String("handler", "token"),
	)

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
		if err := json.NewEncoder(rw).Encode(&Response{Header: http.Header{
			"Authorization": []string{"Bearer " + token.AccessToken},
		}}); err != nil {
			logger.Error("marshalling response", zap.Error(err))
			http.Error(rw, "marshalling response", http.StatusInternalServerError)
		}
	}
}

func GetToken(id string, tinkClient *tink.Client, tinkClientId string, tinkClientSecret string) (tink.Token, error) {
	code, err := tinkClient.GetAuthorizeGrantCode(
		"",
		id,
		GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: authorize grant: %w", err)
	}
	token, err := tinkClient.GetToken(
		code,
		"",
		tinkClientId,
		tinkClientSecret,
		"authorization_code",
		GetAuthorizeGrantCodeRoles,
	)
	if err != nil {
		return tink.Token{}, fmt.Errorf("tink: oauth token: %w", err)
	}
	return token, nil
}

func (s *Server) RegisterTinkUser(oryToken string) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		err := createTinkClient(r.Context(), s.tinkClient, oryToken, s.tinkClientId, s.tinkClientSecret)
		if err != nil {
			http.Error(rw, err.description, err.httpCode)
		}
	}
}

func createTinkClient(ctx context.Context, tinkClient *tink.Client, oryToken string, tinkClientId string, tinkClientSecret string) *HttpError {
	oryConf := ory.NewConfiguration()
	oryConf.Servers = ory.ServerConfigurations{{URL: OriUrl}}
	oryConf.AddDefaultHeader("Authorization", "Bearer "+oryToken)
	oryConf.HTTPClient = &http.Client{}
	oryClient := ory.NewAPIClient(oryConf)

	session, ok := ctx.Value("session").(ory.Session)
	if !ok {
		return &HttpError{http.StatusUnauthorized, "missing session"}
	}
	createdUser, err := tinkClient.CreateUser(
		session.Identity.Id,
		"DE",
		"de_DE",
	)
	if err != nil {
		if !errors.Is(err, tink.ErrUserExists) {
			return &HttpError{http.StatusFailedDependency, "tink: creating user"}
		}
		user, err := getUser(session.Identity.Id, tinkClient, tinkClientId, tinkClientSecret)
		if err != nil {
			return &HttpError{http.StatusFailedDependency, "tink: getting user"}
		}
		createdUser.UserID = user.Id
	}

	traits := session.Identity.Traits.(map[string]interface{})
	traits["tink_id"] = createdUser.UserID
	identity, _, err := oryClient.V0alpha2Api.AdminUpdateIdentity(context.Background(), session.Identity.Id).AdminUpdateIdentityBody(ory.AdminUpdateIdentityBody{
		State:  *session.Identity.State,
		Traits: traits,
	}).Execute()
	if err != nil {
		return &HttpError{http.StatusFailedDependency, "ory: updating identity"}
	}
	session.Identity = *identity
	return nil
}
