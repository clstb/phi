package middleware

import (
	"context"
	"net/http"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	ory "github.com/ory/kratos-client-go"
	"go.uber.org/zap"
)

func Auth(
	ctx context.Context,
	logger *zap.Logger,
	jwksURL string,
) func(http.Handler) http.Handler {
	ar := jwk.NewAutoRefresh(ctx)
	ar.Configure(jwksURL)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			ks, err := ar.Fetch(r.Context(), jwksURL)
			if err != nil {
				logger.Error("fetching jwks", zap.Error(err))
				http.Error(rw, "fetching jwks", http.StatusFailedDependency)
				return
			}

			jwt, err := jwt.ParseHeader(
				r.Header,
				"Authorization",
				jwt.WithKeySet(ks),
				jwt.WithTypedClaim("session", ory.Session{}),
			)
			if err != nil {
				logger.Error("invalid token", zap.Error(err))
				http.Error(rw, "invalid token", http.StatusUnauthorized)
				return
			}

			session, ok := jwt.Get("session")
			if !ok || session == nil {
				logger.Error("invalid session", zap.Error(err))
				http.Error(rw, "invalid session", http.StatusUnauthorized)
				return
			}

			ctx := context.WithValue(r.Context(), "session", session)
			next.ServeHTTP(rw, r.WithContext(ctx))
		})
	}
}
