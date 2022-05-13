package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"

	"github.com/lestrrat-go/jwx/jwk"
	"github.com/lestrrat-go/jwx/jwt"
	ory "github.com/ory/kratos-client-go"
	"go.uber.org/zap"
)

func Auth(logger *zap.Logger, jwksURL string) func(ctx *gin.Context) {
	return func(ctx *gin.Context) {
		ar := jwk.NewAutoRefresh(ctx)
		ar.Configure(jwksURL)

		ks, err := ar.Fetch(ctx, jwksURL)
		if err != nil {
			logger.Error("fetching jwks", zap.Error(err))
			ctx.AbortWithError(http.StatusFailedDependency, err)
			return
		}

		jwt, err := jwt.ParseHeader(
			ctx.Request.Header,
			"Authorization",
			jwt.WithKeySet(ks),
			jwt.WithTypedClaim("session", ory.Session{}),
		)
		if err != nil {
			logger.Error("invalid token", zap.Error(err))
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}

		session, ok := jwt.Get("session")
		if !ok || session == nil {
			logger.Error("invalid session", zap.Error(err))
			ctx.AbortWithError(http.StatusUnauthorized, err)
			return
		}
		ctx.Set("session", session)
		ctx.Next()
	}
}
